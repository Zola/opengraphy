package i18n

import (
	"net/http"
	"sort"
	"strings"
)

type Bundle struct {
	data map[string]map[string]string
}

func New() *Bundle {
	b := &Bundle{data: map[string]map[string]string{
		"en": {
			"brand": "OpenGraphy", "home": "Home", "gallery": "Gallery", "privacy": "Privacy", "terms": "Terms",
			"hero_title": "Trustworthy social previews for every website", "hero_subtitle": "Check Open Graph and Twitter Card metadata, preview social cards, and get repair suggestions before your links go public.",
			"url_placeholder": "URL to check", "check": "Check Preview", "stats_online": "online", "stats_visits": "visits", "stats_checks": "checks",
			"recent_title": "Newest qualified previews", "random_title": "Discover and vote", "leaderboard_title": "Most loved previews",
			"pk_title": "Which preview wins?", "vote_left": "Left wins", "vote_right": "Right wins", "diagnostics": "Diagnostics", "metadata": "Metadata", "code": "Code",
			"consent": "We use necessary cookies for language, consent, and anonymous presence. No tracking cookies.", "accept": "Accept", "reject": "Reject optional",
		},
		"zh-TW": {
			"brand": "OpenGraphy", "home": "首頁", "gallery": "作品牆", "privacy": "隱私", "terms": "條款",
			"hero_title": "讓每個網站都有可信的社群預覽", "hero_subtitle": "偵測 Open Graph 與 Twitter Card metadata，預覽社群卡片，並在連結公開前取得修復建議。",
			"url_placeholder": "輸入要檢測的網址", "check": "檢測預覽", "stats_online": "在線", "stats_visits": "訪問", "stats_checks": "檢測",
			"recent_title": "最新合格作品", "random_title": "探索並投票", "leaderboard_title": "最受歡迎作品榜",
			"pk_title": "哪個預覽更吸引你？", "vote_left": "左邊勝出", "vote_right": "右邊勝出", "diagnostics": "診斷報告", "metadata": "元資料", "code": "代碼",
			"consent": "我們只使用語言、同意狀態與匿名在線人數所需 cookie，不使用追蹤 cookie。", "accept": "接受", "reject": "拒絕可選",
		},
		"zh-CN": {"brand": "OpenGraphy", "home": "首页", "gallery": "作品墙", "privacy": "隐私", "terms": "条款", "hero_title": "让每个网站都有可信的社交预览", "hero_subtitle": "检测 Open Graph 与 Twitter Card metadata，预览社交卡片，并获得修复建议。", "url_placeholder": "输入要检测的网址", "check": "检测预览", "stats_online": "在线", "stats_visits": "访问", "stats_checks": "检测", "recent_title": "最新合格作品", "random_title": "探索并投票", "leaderboard_title": "最受欢迎作品榜", "pk_title": "哪个预览更吸引你？", "vote_left": "左边胜出", "vote_right": "右边胜出", "diagnostics": "诊断报告", "metadata": "元数据", "code": "代码", "consent": "我们只使用必要 cookie，不使用追踪 cookie。", "accept": "接受", "reject": "拒绝可选"},
		"ja":    {"brand": "OpenGraphy", "home": "ホーム", "gallery": "ギャラリー", "privacy": "プライバシー", "terms": "規約", "hero_title": "信頼できるソーシャルプレビューをすべてのサイトに", "hero_subtitle": "Open Graph と Twitter Card を検査し、修正提案を表示します。", "url_placeholder": "チェックするURL", "check": "プレビュー確認", "stats_online": "オンライン", "stats_visits": "訪問", "stats_checks": "検査", "recent_title": "新着プレビュー", "random_title": "発見して投票", "leaderboard_title": "人気ランキング", "pk_title": "どちらが良いですか？", "vote_left": "左", "vote_right": "右", "diagnostics": "診断", "metadata": "Metadata", "code": "コード", "consent": "必要なCookieのみ使用します。", "accept": "同意", "reject": "拒否"},
		"ko":    {"brand": "OpenGraphy", "home": "홈", "gallery": "갤러리", "privacy": "개인정보", "terms": "약관", "hero_title": "모든 웹사이트를 위한 신뢰 가능한 소셜 미리보기", "hero_subtitle": "Open Graph와 Twitter Card를 검사하고 개선 제안을 제공합니다.", "url_placeholder": "검사할 URL", "check": "검사", "stats_online": "온라인", "stats_visits": "방문", "stats_checks": "검사", "recent_title": "최신 미리보기", "random_title": "둘러보고 투표", "leaderboard_title": "인기 순위", "pk_title": "어느 쪽이 더 좋나요?", "vote_left": "왼쪽", "vote_right": "오른쪽", "diagnostics": "진단", "metadata": "Metadata", "code": "코드", "consent": "필수 쿠키만 사용합니다.", "accept": "동의", "reject": "거절"},
		"es":    {"brand": "OpenGraphy", "home": "Inicio", "gallery": "Galería", "privacy": "Privacidad", "terms": "Términos", "hero_title": "Previsualizaciones sociales confiables", "hero_subtitle": "Comprueba Open Graph y Twitter Card con sugerencias de mejora.", "url_placeholder": "URL para comprobar", "check": "Comprobar", "stats_online": "en línea", "stats_visits": "visitas", "stats_checks": "checks", "recent_title": "Últimos previews", "random_title": "Descubre y vota", "leaderboard_title": "Ranking", "pk_title": "¿Cuál gana?", "vote_left": "Izquierda", "vote_right": "Derecha", "diagnostics": "Diagnóstico", "metadata": "Metadata", "code": "Código", "consent": "Usamos solo cookies necesarias.", "accept": "Aceptar", "reject": "Rechazar"},
		"fr":    {"brand": "OpenGraphy", "home": "Accueil", "gallery": "Galerie", "privacy": "Confidentialité", "terms": "Conditions", "hero_title": "Aperçus sociaux fiables", "hero_subtitle": "Vérifiez Open Graph et Twitter Card avec des conseils de correction.", "url_placeholder": "URL à vérifier", "check": "Vérifier", "stats_online": "en ligne", "stats_visits": "visites", "stats_checks": "tests", "recent_title": "Aperçus récents", "random_title": "Découvrir et voter", "leaderboard_title": "Classement", "pk_title": "Lequel gagne ?", "vote_left": "Gauche", "vote_right": "Droite", "diagnostics": "Diagnostic", "metadata": "Metadata", "code": "Code", "consent": "Nous utilisons uniquement des cookies nécessaires.", "accept": "Accepter", "reject": "Refuser"},
		"de":    {"brand": "OpenGraphy", "home": "Start", "gallery": "Galerie", "privacy": "Datenschutz", "terms": "Bedingungen", "hero_title": "Verlässliche Social Previews", "hero_subtitle": "Prüfen Sie Open Graph und Twitter Cards mit konkreten Empfehlungen.", "url_placeholder": "URL prüfen", "check": "Prüfen", "stats_online": "online", "stats_visits": "Besuche", "stats_checks": "Checks", "recent_title": "Neue Previews", "random_title": "Entdecken und bewerten", "leaderboard_title": "Bestenliste", "pk_title": "Welches gewinnt?", "vote_left": "Links", "vote_right": "Rechts", "diagnostics": "Diagnose", "metadata": "Metadata", "code": "Code", "consent": "Wir verwenden nur notwendige Cookies.", "accept": "Akzeptieren", "reject": "Ablehnen"},
	}}
	merge(b.data["en"], enExtra)
	merge(b.data["zh-TW"], zhTWExtra)
	merge(b.data["zh-CN"], zhCNExtra)
	return b
}

var enExtra = map[string]string{
	"meta_description":     "OpenGraphy is a free Open Graph, Twitter Card, and social preview checker with metadata parsing, platform previews, diagnostics, and repair suggestions.",
	"hero_eyebrow":         "Open Graph / Twitter Card / Social Preview Checker",
	"hero_title_full":      "Generate, inspect, and preview Open Graph meta tags in one place",
	"hero_subtitle_full":   "Enter a URL and OpenGraphy fetches the original HTML from the server, parses Open Graph and Twitter Card metadata, previews social cards, and returns practical repair suggestions.",
	"url_placeholder_full": "Enter a URL to check, for example https://example.com",
	"stats_online_full":    "online", "stats_visits_full": "visits", "stats_checks_full": "checks",
	"mock_title":    "Check your share preview before launch",
	"mock_body":     "Make your links look consistent, clear, and compelling across Facebook, LinkedIn, Slack, and messaging apps.",
	"preview_score": "Preview score", "checks_passed": "checks passed",
	"platforms_title":    "Built for the places links actually travel",
	"secure_title":       "Secure fetching, accurate parsing, trustworthy previews",
	"secure_body":        "OpenGraphy does not fetch target sites from the user's browser. The backend simulates crawler behavior: DNS and redirect checks, private network blocking, initial HTML parsing, image reachability checks, and clear diagnostics.",
	"feature_ssrf_title": "SSRF protection", "feature_ssrf_body": "Every fetch validates DNS, IP ranges, and each redirect step to avoid localhost, private networks, and reserved addresses.",
	"feature_cache_title": "4-hour cache", "feature_cache_body": "Metadata and preview results are temporarily cached in Redis without long-term URL history or personal data.",
	"feature_diag_title": "Diagnostics", "feature_diag_body": "Missing tags, relative image URLs, image size, HTTP status, title length, and description length are checked.",
	"knowledge_title": "Make every share look like a carefully designed ad placement",
	"knowledge_body":  "Open Graph is the shared language used by social platforms, messengers, and content tools to read page previews. Put the right tags in your HTML head and you control the title, description, image, URL, and content type.",
	"what_og_title":   "What is Open Graph?", "what_og_body": "OG tags are metadata for crawlers. Facebook, LinkedIn, Slack, Discord, WhatsApp, Telegram, and X usually read og:title, og:description, and og:image when generating link previews.",
	"image_size_title": "Recommended image size", "image_size_body": "The safest preview image size is 1200 x 630, about 1.91:1. Keep images lightweight so messaging apps can generate previews quickly.",
	"backend_fetch_title": "Why fetch from the backend?", "backend_fetch_body": "Most platform crawlers do not run browser JavaScript or wait for SPA rendering. Server-side fetching is closer to how real crawlers see your page.",
	"minimal_tags_title": "Minimal useful Open Graph tags", "minimal_tags_body": "If you can only fix one thing first, add title, description, image, and canonical URL. These solve most broken preview problems.",
	"platform_rules_title": "Different platforms interpret the same metadata differently",
	"framework_title":      "How common stacks integrate",
	"faq_title":            "FAQ",
	"gallery_empty":        "Qualified website previews will appear here after checks.",
	"gallery_intro":        "Public, qualified previews discovered by OpenGraphy. The first row favors the newest submissions so contributors can see their work immediately.",
	"no_works":             "No qualified works yet.", "ranking_empty": "Ranking will appear after votes.", "load_pair": "Load pair",
	"footer_body":    "Secure Open Graph diagnostics, social preview testing, and metadata repair suggestions.",
	"language_label": "Language", "facebook_debugger": "Facebook Debugger", "linkedin_inspector": "LinkedIn Post Inspector",
	"consent_full":     "We only use necessary cookies for language preference, consent status, and anonymous online presence. No tracking or advertising cookies.",
	"reject_necessary": "Necessary only",
	"score_label":      "score", "fetch_label": "fetch", "source_label": "source", "cached": "cached", "fresh": "fresh",
	"title_label": "Title", "description_label": "Description", "image_label": "Image", "final_url_label": "Final URL", "copy": "Copy",
	"all_good_title": "All good", "all_good_body": "No major metadata problems detected.",
	"privacy_title": "Privacy Policy", "privacy_body_1": "OpenGraphy does not sell data and does not use advertising cookies. Submitted URLs are fetched to produce metadata diagnostics and cached in Redis for up to 4 hours. Qualified public website previews may appear in the gallery as public examples.", "privacy_body_2": "Presence uses an anonymous random session cookie. Rate limiting uses salted short-lived hashes and does not store plain IP addresses.",
	"terms_title": "Terms", "terms_body": "Use OpenGraphy only for public websites you are allowed to inspect. The service blocks private networks and reserves the right to rate limit abusive traffic. Gallery entries are public website previews and can be refreshed or removed by administrators later.",
}

var zhTWExtra = map[string]string{
	"meta_description":     "OpenGraphy 是免費的 Open Graph、Twitter Card 與社群分享預覽檢測工具，提供 metadata 解析、平台預覽、診斷報告與修復建議。",
	"hero_eyebrow":         "Open Graph / Twitter Card / 社群預覽檢測",
	"hero_title_full":      "一站式生成、檢測與預覽 Open Graph Meta Tags",
	"hero_subtitle_full":   "輸入網址，OpenGraphy 會從伺服器端抓取網頁原始 HTML，解析 Open Graph 與 Twitter Card metadata，展示各平台分享卡片，並給出可執行的修復建議。",
	"url_placeholder_full": "輸入要檢測的網址，例如 https://example.com",
	"stats_online_full":    "人在線", "stats_visits_full": "次訪問", "stats_checks_full": "次檢測",
	"mock_title": "發布前就確認分享預覽", "mock_body": "讓客戶在 Facebook、LinkedIn、Slack 或訊息 App 裡看到一致、清楚、有吸引力的連結卡片。",
	"preview_score": "預覽分數", "checks_passed": "項檢查通過",
	"platforms_title":    "支援連結真正會被分享出去的平台",
	"secure_title":       "安全抓取、準確解析、可信預覽",
	"secure_body":        "OpenGraphy 不在使用者瀏覽器裡抓取目標網站，而是由後端模擬 crawler 行為：檢查 DNS 與重定向、阻擋私有網段、解析初始 HTML、驗證圖片是否可訪問，最後把結果整理成診斷報告與修復建議。",
	"feature_ssrf_title": "SSRF 防護", "feature_ssrf_body": "每次抓取都檢查 DNS、IP 範圍與每一步重定向，避免訪問 localhost、內網與保留位址。",
	"feature_cache_title": "4 小時快取", "feature_cache_body": "metadata 與預覽結果暫存在 Redis，不長期保存 URL 歷史，也不保存使用者個資。",
	"feature_diag_title": "診斷建議", "feature_diag_body": "檢查缺失標籤、相對圖片 URL、圖片尺寸、HTTP 狀態、標題長度與描述長度。",
	"knowledge_title": "用 Open Graph 讓每一次分享都像精心設計的廣告位",
	"knowledge_body":  "Open Graph 是社群平台、即時通訊軟體和內容工具讀取網頁預覽的共同語言。只要在 HTML 的 head 裡放對標籤，你就能控制分享卡片的標題、描述、圖片、網址和內容類型。",
	"what_og_title":   "Open Graph 是什麼？", "what_og_body": "OG 標籤是一組給 crawler 看的 metadata。Facebook、LinkedIn、Slack、Discord、WhatsApp、Telegram、X 等平台在產生連結預覽時，通常會先讀取 og:title、og:description 和 og:image。",
	"image_size_title": "建議圖片尺寸", "image_size_body": "最穩妥的預覽圖尺寸是 1200 x 630，比例約為 1.91:1。圖片請盡量壓小，讓通訊工具更快生成預覽。",
	"backend_fetch_title": "為什麼要從後端抓取？", "backend_fetch_body": "大多數平台 crawler 不會執行瀏覽器裡的 JavaScript，也不會等待 SPA 完成渲染。後端抓取更接近平台實際看到的頁面。",
	"minimal_tags_title": "最小可用的 Open Graph 標籤", "minimal_tags_body": "如果你只能先做一件事，請先把標題、描述、圖片和 canonical URL 補齊。這幾個標籤能消除大部分預覽問題。",
	"platform_rules_title": "不同平台會用不同規則解讀同一組 metadata",
	"framework_title":      "常見技術棧如何接入",
	"faq_title":            "常見問題",
	"gallery_empty":        "完成檢測並符合展示條件的網站預覽，會出現在這裡。",
	"gallery_intro":        "OpenGraphy 發現的公開合格預覽。第一排優先展示最新提交，讓第一次提交的網站也能立即被看到。",
	"no_works":             "目前還沒有合格作品。", "ranking_empty": "投票後會出現排行榜。", "load_pair": "載入作品",
	"footer_body":    "安全的 Open Graph 診斷、社群預覽測試與 metadata 修復建議。",
	"language_label": "語言", "facebook_debugger": "Facebook 官方除錯工具", "linkedin_inspector": "LinkedIn Post Inspector",
	"consent_full":     "我們只使用語言偏好、同意狀態與匿名在線人數所需的必要 cookie，不使用追蹤 cookie 或廣告 cookie。",
	"reject_necessary": "僅必要",
	"score_label":      "分數", "fetch_label": "抓取", "source_label": "來源", "cached": "快取", "fresh": "即時",
	"title_label": "標題", "description_label": "描述", "image_label": "圖片", "final_url_label": "最終 URL", "copy": "複製",
	"all_good_title": "狀態良好", "all_good_body": "沒有偵測到主要 metadata 問題。",
	"privacy_title": "隱私政策", "privacy_body_1": "OpenGraphy 不出售資料，也不使用廣告 cookie。提交的 URL 只會用於產生 metadata 診斷，並在 Redis 中快取最多 4 小時。符合條件的公開網站預覽可能會作為公開範例顯示在作品牆。", "privacy_body_2": "在線人數使用匿名隨機 session cookie。速率限制使用短期加鹽雜湊，不保存明文 IP。",
	"terms_title": "服務條款", "terms_body": "請只檢測你有權公開訪問的網站。服務會阻擋私有網段，並保留限制濫用流量的權利。作品牆內容是公開網站預覽，未來可由管理員刷新或移除。",
}

var zhCNExtra = map[string]string{}

func init() {
	zhCNExtra = make(map[string]string, len(zhTWExtra))
	for k, v := range zhTWExtra {
		zhCNExtra[k] = strings.NewReplacer("預", "预", "檢", "检", "測", "测", "與", "与", "隱", "隐", "條", "条", "網", "网", "準", "准", "確", "确", "覽", "览", "議", "议", "語", "语", "態", "态", "線", "线", "項", "项", "圖", "图", "後", "后", "這", "这", "會", "会", "體", "体", "讀", "读", "標", "标", "籤", "签", "過", "过", "據", "据", "暫", "暂", "歷", "历", "長", "长", "資訊", "信息", "發", "发", "現", "现", "優", "优", "顯", "显", "個", "个", "數", "数", "複", "复", "態", "态").Replace(v)
	}
	merge(zhCNExtra, map[string]string{
		"meta_description":     "OpenGraphy 是免费的 Open Graph、Twitter Card 与社交分享预览检测工具，提供 metadata 解析、平台预览、诊断报告与修复建议。",
		"hero_eyebrow":         "Open Graph / Twitter Card / 社交预览检测",
		"hero_title_full":      "一站式生成、检测与预览 Open Graph Meta Tags",
		"hero_subtitle_full":   "输入网址，OpenGraphy 会从服务器端抓取网页原始 HTML，解析 Open Graph 与 Twitter Card metadata，展示各平台分享卡片，并给出可执行的修复建议。",
		"url_placeholder_full": "输入要检测的网址，例如 https://example.com",
		"stats_online_full":    "人在线",
		"stats_visits_full":    "次访问",
		"stats_checks_full":    "次检测",
		"mock_title":           "发布前就确认分享预览",
		"mock_body":            "让客户在 Facebook、LinkedIn、Slack 或消息 App 里看到一致、清楚、有吸引力的链接卡片。",
		"preview_score":        "预览分数",
		"checks_passed":        "项检查通过",
		"platforms_title":      "支持链接真正会被分享出去的平台",
		"secure_title":         "安全抓取、准确解析、可信预览",
		"secure_body":          "OpenGraphy 不在用户浏览器里抓取目标网站，而是由后端模拟 crawler 行为：检查 DNS 与重定向、阻挡私有网段、解析初始 HTML、验证图片是否可访问，最后把结果整理成诊断报告与修复建议。",
		"feature_ssrf_title":   "SSRF 防护",
		"feature_ssrf_body":    "每次抓取都检查 DNS、IP 范围与每一步重定向，避免访问 localhost、内网与保留地址。",
		"feature_cache_title":  "4 小时缓存",
		"feature_cache_body":   "metadata 与预览结果暂存在 Redis，不长期保存 URL 历史，也不保存用户个人资料。",
		"feature_diag_title":   "诊断建议",
		"feature_diag_body":    "检查缺失标签、相对图片 URL、图片尺寸、HTTP 状态、标题长度与描述长度。",
		"knowledge_title":      "用 Open Graph 让每一次分享都像精心设计的广告位",
		"knowledge_body":       "Open Graph 是社交平台、即时通讯软件和内容工具读取网页预览的共同语言。只要在 HTML 的 head 里放对标签，你就能控制分享卡片的标题、描述、图片、网址和内容类型。",
		"what_og_title":        "Open Graph 是什么？",
		"what_og_body":         "OG 标签是一组给 crawler 看的 metadata。Facebook、LinkedIn、Slack、Discord、WhatsApp、Telegram、X 等平台在生成链接预览时，通常会先读取 og:title、og:description 和 og:image。",
		"image_size_title":     "建议图片尺寸",
		"image_size_body":      "最稳妥的预览图尺寸是 1200 x 630，比例约为 1.91:1。图片请尽量压小，让通讯工具更快生成预览。",
		"backend_fetch_title":  "为什么要从后端抓取？",
		"backend_fetch_body":   "大多数平台 crawler 不会执行浏览器里的 JavaScript，也不会等待 SPA 完成渲染。后端抓取更接近平台实际看到的页面。",
		"minimal_tags_title":   "最小可用的 Open Graph 标签",
		"minimal_tags_body":    "如果你只能先做一件事，请先把标题、描述、图片和 canonical URL 补齐。这几个标签能消除大部分预览问题。",
		"platform_rules_title": "不同平台会用不同规则解读同一组 metadata",
		"framework_title":      "常见技术栈如何接入",
		"faq_title":            "常见问题",
		"gallery_empty":        "完成检测并符合展示条件的网站预览，会出现在这里。",
		"gallery_intro":        "OpenGraphy 发现的公开合格预览。第一排优先展示最新提交，让第一次提交的网站也能立即被看到。",
		"no_works":             "目前还没有合格作品。",
		"ranking_empty":        "投票后会出现排行榜。",
		"load_pair":            "载入作品",
		"footer_body":          "安全的 Open Graph 诊断、社交预览测试与 metadata 修复建议。",
		"language_label":       "语言",
		"facebook_debugger":    "Facebook 官方调试工具",
		"consent_full":         "我们只使用语言偏好、同意状态与匿名在线人数所需的必要 cookie，不使用追踪 cookie 或广告 cookie。",
		"reject_necessary":     "仅必要",
		"score_label":          "分数",
		"fetch_label":          "抓取",
		"source_label":         "来源",
		"cached":               "缓存",
		"fresh":                "实时",
		"title_label":          "标题",
		"description_label":    "描述",
		"image_label":          "图片",
		"final_url_label":      "最终 URL",
		"copy":                 "复制",
		"all_good_title":       "状态良好",
		"all_good_body":        "没有检测到主要 metadata 问题。",
		"privacy_title":        "隐私政策",
		"privacy_body_1":       "OpenGraphy 不出售资料，也不使用广告 cookie。提交的 URL 只会用于生成 metadata 诊断，并在 Redis 中缓存最多 4 小时。符合条件的公开网站预览可能会作为公开示例显示在作品墙。",
		"privacy_body_2":       "在线人数使用匿名随机 session cookie。速率限制使用短期加盐哈希，不保存明文 IP。",
		"terms_title":          "服务条款",
		"terms_body":           "请只检测你有权公开访问的网站。服务会阻挡私有网段，并保留限制滥用流量的权利。作品墙内容是公开网站预览，未来可由管理员刷新或移除。",
	})
}

func merge(dst, src map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}

func (b *Bundle) T(locale, key string) string {
	if v := b.data[locale][key]; v != "" {
		return v
	}
	if v := b.data["en"][key]; v != "" {
		return v
	}
	return key
}

func (b *Bundle) Locale(r *http.Request) string {
	if c, err := r.Cookie("lang"); err == nil && b.Has(c.Value) {
		return c.Value
	}
	header := r.Header.Get("Accept-Language")
	for _, part := range strings.Split(header, ",") {
		code := strings.TrimSpace(strings.Split(part, ";")[0])
		if b.Has(code) {
			return code
		}
		if strings.HasPrefix(code, "zh") {
			if strings.Contains(strings.ToUpper(code), "TW") || strings.Contains(strings.ToUpper(code), "HK") {
				return "zh-TW"
			}
			return "zh-CN"
		}
		if len(code) >= 2 && b.Has(code[:2]) {
			return code[:2]
		}
	}
	return "en"
}

func (b *Bundle) Has(locale string) bool {
	_, ok := b.data[locale]
	return ok
}

func (b *Bundle) Locales() []string {
	out := make([]string, 0, len(b.data))
	for k := range b.data {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func (b *Bundle) LocaleLabels() map[string]string {
	return map[string]string{
		"de":    "Deutsch",
		"en":    "English",
		"es":    "Español",
		"fr":    "Français",
		"ja":    "日本語",
		"ko":    "한국어",
		"zh-CN": "简体中文",
		"zh-TW": "繁體中文",
	}
}
