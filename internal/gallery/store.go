package gallery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	_ "modernc.org/sqlite"

	"opengraphy/internal/model"
	"opengraphy/internal/security"
)

type Store struct {
	db  *sql.DB
	rdb *redis.Client
	ttl time.Duration
}

type Health struct {
	RedisRecent int64 `json:"redis_recent"`
	SQLiteWorks int64 `json:"sqlite_works"`
	Total       int64 `json:"total"`
}

func NewStore(path string, rdb *redis.Client, ttl time.Duration) (*Store, error) {
	if dir := filepath.Dir(path); dir != "." && dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, rdb: rdb, ttl: ttl}
	if err := s.migrate(); err != nil {
		return nil, err
	}
	_ = s.Hydrate(context.Background())
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	_, err := s.db.Exec(`CREATE TABLE IF NOT EXISTS works (
		id TEXT PRIMARY KEY,
		url TEXT NOT NULL,
		final_url TEXT NOT NULL,
		domain TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		image TEXT NOT NULL,
		score INTEGER NOT NULL,
		rating INTEGER NOT NULL DEFAULT 1400,
		wins INTEGER NOT NULL DEFAULT 0,
		losses INTEGER NOT NULL DEFAULT 0,
		first_seen TEXT NOT NULL,
		last_seen TEXT NOT NULL
	);`)
	return err
}

func (s *Store) Qualifies(result model.CheckResult) bool {
	return result.Score >= 70 && result.Preview.Title != "" && result.Preview.Image != "" && result.Image.Reachable
}

func (s *Store) Upsert(ctx context.Context, result model.CheckResult) error {
	if !s.Qualifies(result) {
		return nil
	}
	id := security.Hash(result.Preview.URL)
	now := time.Now().UTC()
	w := model.Work{
		ID:          id,
		URL:         result.URL,
		FinalURL:    result.FinalURL,
		Domain:      domain(result.Preview.URL),
		Title:       result.Preview.Title,
		Description: result.Preview.Description,
		Image:       result.Preview.Image,
		Score:       result.Score,
		Rating:      1400,
		FirstSeen:   now,
		LastSeen:    now,
	}
	existing, ok := s.Get(ctx, id)
	if ok {
		w.Rating = existing.Rating
		w.Wins = existing.Wins
		w.Losses = existing.Losses
		w.FirstSeen = existing.FirstSeen
	}
	raw, _ := json.Marshal(w)
	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, "og:work:"+id, raw, s.ttl)
	pipe.ZAdd(ctx, "og:works:recent", redis.Z{Score: float64(now.Unix()), Member: id})
	pipe.ZAdd(ctx, "og:works:rating", redis.Z{Score: float64(w.Rating), Member: id})
	pipe.Expire(ctx, "og:works:recent", s.ttl)
	pipe.Expire(ctx, "og:works:rating", s.ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	return s.persistWork(ctx, w)
}

func (s *Store) Get(ctx context.Context, id string) (model.Work, bool) {
	raw, err := s.rdb.Get(ctx, "og:work:"+id).Bytes()
	if err == nil {
		var w model.Work
		if json.Unmarshal(raw, &w) == nil {
			return w, true
		}
	}
	row := s.db.QueryRowContext(ctx, `SELECT id,url,final_url,domain,title,description,image,score,rating,wins,losses,first_seen,last_seen FROM works WHERE id=?`, id)
	w, err := scanWork(row)
	if err != nil {
		return model.Work{}, false
	}
	_ = s.cacheWork(ctx, w)
	return w, true
}

func (s *Store) Recent(ctx context.Context, limit int) []model.Work {
	ids := s.rdb.ZRevRange(ctx, "og:works:recent", 0, int64(limit-1)).Val()
	return mergeWorks(s.worksByIDs(ctx, ids, limit), s.recentFromDB(ctx, limit), limit)
}

func (s *Store) Leaderboard(ctx context.Context, limit int) []model.Work {
	ids := s.rdb.ZRevRange(ctx, "og:works:rating", 0, int64(limit-1)).Val()
	return mergeWorks(s.worksByIDs(ctx, ids, limit), s.leaderboardFromDB(ctx, limit), limit)
}

func (s *Store) Random(ctx context.Context, limit int) []model.Work {
	ids := s.rdb.ZRange(ctx, "og:works:recent", 0, -1).Val()
	rand.Shuffle(len(ids), func(i, j int) { ids[i], ids[j] = ids[j], ids[i] })
	return mergeWorks(s.worksByIDs(ctx, ids, limit), s.randomFromDB(ctx, limit), limit)
}

func (s *Store) PKPair(ctx context.Context) (model.Work, model.Work, error) {
	items := s.Random(ctx, 2)
	if len(items) < 2 {
		return model.Work{}, model.Work{}, errors.New("not enough works")
	}
	return items[0], items[1], nil
}

func (s *Store) Health(ctx context.Context) Health {
	sqliteWorks := s.countDB(ctx)
	redisRecent := s.rdb.ZCard(ctx, "og:works:recent").Val()
	total := sqliteWorks
	if redisRecent > total {
		total = redisRecent
	}
	return Health{RedisRecent: redisRecent, SQLiteWorks: sqliteWorks, Total: total}
}

func (s *Store) Count(ctx context.Context) int64 {
	return s.Health(ctx).Total
}

func (s *Store) SeedDefaultsIfEmpty(ctx context.Context) error {
	if s.Count(ctx) > 0 {
		return nil
	}
	now := time.Now().UTC()
	for i, w := range defaultWorks(now) {
		w.LastSeen = now.Add(time.Duration(i) * time.Second)
		if err := s.persistWork(ctx, w); err != nil {
			return err
		}
		if err := s.cacheWork(ctx, w); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) Vote(ctx context.Context, winnerID, loserID string) error {
	winner, ok := s.Get(ctx, winnerID)
	if !ok {
		return errors.New("winner not found")
	}
	loser, ok := s.Get(ctx, loserID)
	if !ok {
		return errors.New("loser not found")
	}
	winner.Rating++
	winner.Wins++
	winner.LastSeen = time.Now().UTC()
	loser.Rating--
	loser.Losses++
	loser.LastSeen = time.Now().UTC()
	if err := s.cacheWork(ctx, winner); err != nil {
		return err
	}
	if err := s.cacheWork(ctx, loser); err != nil {
		return err
	}
	pipe := s.rdb.Pipeline()
	pipe.ZAdd(ctx, "og:works:rating", redis.Z{Score: float64(winner.Rating), Member: winner.ID})
	pipe.ZAdd(ctx, "og:works:rating", redis.Z{Score: float64(loser.Rating), Member: loser.ID})
	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}
	if err := s.persistWork(ctx, winner); err != nil {
		return err
	}
	return s.persistWork(ctx, loser)
}

func (s *Store) SyncEvery(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			_ = s.Sync(ctx)
			return
		case <-ticker.C:
			_ = s.Sync(ctx)
		}
	}
}

func (s *Store) Sync(ctx context.Context) error {
	ids := s.rdb.ZRange(ctx, "og:works:recent", 0, -1).Val()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO works (id,url,final_url,domain,title,description,image,score,rating,wins,losses,first_seen,last_seen)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
	ON CONFLICT(id) DO UPDATE SET url=excluded.url, final_url=excluded.final_url, domain=excluded.domain, title=excluded.title,
	description=excluded.description, image=excluded.image, score=excluded.score, rating=excluded.rating, wins=excluded.wins,
	losses=excluded.losses, last_seen=excluded.last_seen`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, id := range ids {
		w, ok := s.Get(ctx, id)
		if !ok {
			continue
		}
		_, _ = stmt.ExecContext(ctx, w.ID, w.URL, w.FinalURL, w.Domain, w.Title, w.Description, w.Image, w.Score, w.Rating, w.Wins, w.Losses, w.FirstSeen.Format(time.RFC3339), w.LastSeen.Format(time.RFC3339))
	}
	return tx.Commit()
}

func (s *Store) Hydrate(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `SELECT id,url,final_url,domain,title,description,image,score,rating,wins,losses,first_seen,last_seen FROM works ORDER BY last_seen DESC LIMIT 1000`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		w, err := scanWork(rows)
		if err == nil {
			_ = s.cacheWork(ctx, w)
		}
	}
	return rows.Err()
}

func (s *Store) cacheWork(ctx context.Context, w model.Work) error {
	raw, _ := json.Marshal(w)
	pipe := s.rdb.Pipeline()
	pipe.Set(ctx, "og:work:"+w.ID, raw, s.ttl)
	pipe.ZAdd(ctx, "og:works:recent", redis.Z{Score: float64(w.LastSeen.Unix()), Member: w.ID})
	pipe.ZAdd(ctx, "og:works:rating", redis.Z{Score: float64(w.Rating), Member: w.ID})
	_, err := pipe.Exec(ctx)
	return err
}

func (s *Store) persistWork(ctx context.Context, w model.Work) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO works (id,url,final_url,domain,title,description,image,score,rating,wins,losses,first_seen,last_seen)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
	ON CONFLICT(id) DO UPDATE SET url=excluded.url, final_url=excluded.final_url, domain=excluded.domain, title=excluded.title,
	description=excluded.description, image=excluded.image, score=excluded.score, rating=excluded.rating, wins=excluded.wins,
	losses=excluded.losses, last_seen=excluded.last_seen`,
		w.ID, w.URL, w.FinalURL, w.Domain, w.Title, w.Description, w.Image, w.Score, w.Rating, w.Wins, w.Losses, w.FirstSeen.Format(time.RFC3339), w.LastSeen.Format(time.RFC3339))
	return err
}

func (s *Store) worksByIDs(ctx context.Context, ids []string, limit int) []model.Work {
	out := make([]model.Work, 0, limit)
	for _, id := range ids {
		w, ok := s.Get(ctx, id)
		if ok {
			out = append(out, w)
		}
		if len(out) >= limit {
			break
		}
	}
	return out
}

func (s *Store) recentFromDB(ctx context.Context, limit int) []model.Work {
	return s.queryWorks(ctx, `SELECT id,url,final_url,domain,title,description,image,score,rating,wins,losses,first_seen,last_seen FROM works ORDER BY last_seen DESC LIMIT ?`, limit)
}

func (s *Store) leaderboardFromDB(ctx context.Context, limit int) []model.Work {
	return s.queryWorks(ctx, `SELECT id,url,final_url,domain,title,description,image,score,rating,wins,losses,first_seen,last_seen FROM works ORDER BY rating DESC, last_seen DESC LIMIT ?`, limit)
}

func (s *Store) randomFromDB(ctx context.Context, limit int) []model.Work {
	candidates := s.queryWorks(ctx, `SELECT id,url,final_url,domain,title,description,image,score,rating,wins,losses,first_seen,last_seen FROM works ORDER BY last_seen DESC LIMIT ?`, 200)
	rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
	if len(candidates) > limit {
		return candidates[:limit]
	}
	return candidates
}

func (s *Store) queryWorks(ctx context.Context, query string, limit int) []model.Work {
	if limit <= 0 {
		return nil
	}
	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()
	out := make([]model.Work, 0, limit)
	for rows.Next() {
		w, err := scanWork(rows)
		if err == nil {
			out = append(out, w)
			_ = s.cacheWork(ctx, w)
		}
	}
	return out
}

func (s *Store) countDB(ctx context.Context) int64 {
	var count int64
	_ = s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM works`).Scan(&count)
	return count
}

func mergeWorks(primary, fallback []model.Work, limit int) []model.Work {
	out := make([]model.Work, 0, limit)
	seen := make(map[string]struct{}, limit)
	for _, list := range [][]model.Work{primary, fallback} {
		for _, w := range list {
			if _, ok := seen[w.ID]; ok {
				continue
			}
			seen[w.ID] = struct{}{}
			out = append(out, w)
			if len(out) >= limit {
				return out
			}
		}
	}
	return out
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanWork(row rowScanner) (model.Work, error) {
	var w model.Work
	var first, last string
	err := row.Scan(&w.ID, &w.URL, &w.FinalURL, &w.Domain, &w.Title, &w.Description, &w.Image, &w.Score, &w.Rating, &w.Wins, &w.Losses, &first, &last)
	if err != nil {
		return w, err
	}
	w.FirstSeen, _ = time.Parse(time.RFC3339, first)
	w.LastSeen, _ = time.Parse(time.RFC3339, last)
	return w, nil
}

func domain(raw string) string {
	u, err := url.Parse(raw)
	if err == nil && u.Hostname() != "" {
		return strings.TrimPrefix(u.Hostname(), "www.")
	}
	return raw
}
