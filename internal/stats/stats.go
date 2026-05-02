package stats

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type Stats struct {
	rdb *redis.Client
}

type Snapshot struct {
	VisitsTotal int64 `json:"visits_total"`
	ChecksTotal int64 `json:"checks_total"`
	Online      int64 `json:"online"`
	CacheHits   int64 `json:"cache_hits"`
	CacheMisses int64 `json:"cache_misses"`
	WorksTotal  int64 `json:"works_total"`
}

func New(rdb *redis.Client) *Stats {
	return &Stats{rdb: rdb}
}

func (s *Stats) Visit(ctx context.Context, locale string) {
	pipe := s.rdb.Pipeline()
	pipe.Incr(ctx, "stats:visits:total")
	if locale != "" {
		pipe.Incr(ctx, "stats:locale:"+locale)
	}
	_, _ = pipe.Exec(ctx)
}

func (s *Stats) Check(ctx context.Context) {
	_ = s.rdb.Incr(ctx, "stats:checks:total").Err()
}

func (s *Stats) CacheHit(ctx context.Context, hit bool) {
	key := "stats:cache:misses"
	if hit {
		key = "stats:cache:hits"
	}
	_ = s.rdb.Incr(ctx, key).Err()
}

func (s *Stats) Presence(ctx context.Context, sessionID string) {
	now := time.Now().Unix()
	pipe := s.rdb.Pipeline()
	pipe.ZAdd(ctx, "stats:online", redis.Z{Score: float64(now), Member: sessionID})
	pipe.ZRemRangeByScore(ctx, "stats:online", "0", strconv.FormatInt(now-90, 10))
	pipe.Expire(ctx, "stats:online", 2*time.Minute)
	_, _ = pipe.Exec(ctx)
}

func (s *Stats) Snapshot(ctx context.Context) Snapshot {
	now := time.Now().Unix()
	_ = s.rdb.ZRemRangeByScore(ctx, "stats:online", "0", strconv.FormatInt(now-90, 10)).Err()
	vals, _ := s.rdb.MGet(ctx,
		"stats:visits:total",
		"stats:checks:total",
		"stats:cache:hits",
		"stats:cache:misses",
	).Result()
	return Snapshot{
		VisitsTotal: toInt(vals, 0),
		ChecksTotal: toInt(vals, 1),
		CacheHits:   toInt(vals, 2),
		CacheMisses: toInt(vals, 3),
		Online:      s.rdb.ZCard(ctx, "stats:online").Val(),
		WorksTotal:  s.rdb.ZCard(ctx, "og:works:recent").Val(),
	}
}

func toInt(vals []any, index int) int64 {
	if index >= len(vals) || vals[index] == nil {
		return 0
	}
	switch v := vals[index].(type) {
	case string:
		var n int64
		for _, r := range v {
			if r >= '0' && r <= '9' {
				n = n*10 + int64(r-'0')
			}
		}
		return n
	case int64:
		return v
	default:
		return 0
	}
}
