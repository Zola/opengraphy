package gallery

import (
	"time"

	"opengraphy/internal/model"
	"opengraphy/internal/security"
)

func defaultWorks(now time.Time) []model.Work {
	items := []struct {
		url         string
		title       string
		description string
		image       string
		score       int
		rating      int
	}{
		{
			url:         "https://opengraph.aston.tw",
			title:       "可信 Open Graph 預覽診斷工具",
			description: "一鍵檢測 Open Graph 與 Twitter Card metadata，預覽社群卡片並取得修復建議。",
			image:       "https://picsum.photos/seed/opengraphy-expanded-0/1200/628",
			score:       96,
			rating:      1530,
		},
		{
			url:         "https://aston.tw",
			title:       "Aston Technology Studio",
			description: "展示技術服務、產品能力和品牌案例的公司網站。",
			image:       "https://picsum.photos/seed/opengraphy-expanded-1/1200/628",
			score:       94,
			rating:      1513,
		},
		{
			url:         "https://huangwenjue.com",
			title:       "設計師作品與品牌案例",
			description: "以清楚的社群預覽展示公開作品，讓更多訪客點進原站。",
			image:       "https://picsum.photos/seed/opengraphy-expanded-2/1200/628",
			score:       92,
			rating:      1496,
		},
		{
			url:         "https://voicematch.aston.tw",
			title:       "VoiceMatch 語音工具",
			description: "可視化語音資料與應用能力的產品展示頁。",
			image:       "https://picsum.photos/seed/opengraphy-expanded-3/1200/628",
			score:       91,
			rating:      1479,
		},
		{
			url:         "https://asr.aston.tw",
			title:       "ASR 語音辨識服務",
			description: "即時語音轉文字服務與模型能力展示。",
			image:       "https://picsum.photos/seed/opengraphy-expanded-4/1200/628",
			score:       90,
			rating:      1462,
		},
		{
			url:         "https://stickers.aston.tw",
			title:       "貼圖作品展示",
			description: "適合社群分享的圖像與作品集入口。",
			image:       "https://picsum.photos/seed/opengraphy-expanded-5/1200/628",
			score:       88,
			rating:      1445,
		},
	}
	works := make([]model.Work, 0, len(items))
	for i, item := range items {
		seen := now.Add(-time.Duration(len(items)-i) * time.Minute)
		works = append(works, model.Work{
			ID:          security.Hash(item.url),
			URL:         item.url,
			FinalURL:    item.url,
			Domain:      domain(item.url),
			Title:       item.title,
			Description: item.description,
			Image:       item.image,
			Score:       item.score,
			Rating:      item.rating,
			FirstSeen:   seen,
			LastSeen:    seen,
		})
	}
	return works
}
