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
	merge(b.data["ja"], jaExtra)
	merge(b.data["ko"], koExtra)
	merge(b.data["es"], esExtra)
	merge(b.data["fr"], frExtra)
	merge(b.data["de"], deExtra)
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
	"apply_language":   "Apply",
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
	"apply_language":   "套用語言",
	"consent_full":     "我們只使用語言偏好、同意狀態與匿名在線人數所需的必要 cookie，不使用追蹤 cookie 或廣告 cookie。",
	"reject_necessary": "僅必要",
	"score_label":      "分數", "fetch_label": "抓取", "source_label": "來源", "cached": "快取", "fresh": "即時",
	"title_label": "標題", "description_label": "描述", "image_label": "圖片", "final_url_label": "最終 URL", "copy": "複製",
	"all_good_title": "狀態良好", "all_good_body": "沒有偵測到主要 metadata 問題。",
	"privacy_title": "隱私政策", "privacy_body_1": "OpenGraphy 不出售資料，也不使用廣告 cookie。提交的 URL 只會用於產生 metadata 診斷，並在 Redis 中快取最多 4 小時。符合條件的公開網站預覽可能會作為公開範例顯示在作品牆。", "privacy_body_2": "在線人數使用匿名隨機 session cookie。速率限制使用短期加鹽雜湊，不保存明文 IP。",
	"terms_title": "服務條款", "terms_body": "請只檢測你有權公開訪問的網站。服務會阻擋私有網段，並保留限制濫用流量的權利。作品牆內容是公開網站預覽，未來可由管理員刷新或移除。",
}

var zhCNExtra = map[string]string{}

var jaExtra = map[string]string{
	"meta_description":     "OpenGraphy は Open Graph、Twitter Card、ソーシャルプレビューを無料で検査し、metadata 解析、各プラットフォームのプレビュー、診断、修正提案を提供します。",
	"hero_eyebrow":         "Open Graph / Twitter Card / ソーシャルプレビュー検査",
	"hero_title_full":      "Open Graph meta tags を生成、検査、プレビュー",
	"hero_subtitle_full":   "URL を入力すると、OpenGraphy がサーバー側で HTML を取得し、Open Graph と Twitter Card の metadata を解析し、各プラットフォームのカードと修正案を表示します。",
	"url_placeholder_full": "検査する URL を入力。例: https://example.com",
	"stats_online_full":    "オンライン",
	"stats_visits_full":    "訪問",
	"stats_checks_full":    "検査",
	"mock_title":           "公開前に共有プレビューを確認",
	"mock_body":            "Facebook、LinkedIn、Slack、メッセージアプリで一貫した分かりやすいリンクカードを表示します。",
	"preview_score":        "プレビュースコア",
	"checks_passed":        "項目合格",
	"platforms_title":      "リンクが実際に共有される場所に対応",
	"secure_title":         "安全な取得、正確な解析、信頼できるプレビュー",
	"secure_body":          "OpenGraphy は利用者のブラウザではなくバックエンドで crawler に近い動作を行い、DNS、リダイレクト、プライベートネットワーク、HTML、画像の到達性を確認します。",
	"feature_ssrf_title":   "SSRF 保護",
	"feature_ssrf_body":    "取得ごとに DNS、IP 範囲、リダイレクトを検証し、localhost、内部ネットワーク、予約アドレスへのアクセスを防ぎます。",
	"feature_cache_title":  "4 時間キャッシュ",
	"feature_cache_body":   "metadata とプレビュー結果は Redis に一時保存され、URL 履歴や個人情報を長期保存しません。",
	"feature_diag_title":   "診断提案",
	"feature_diag_body":    "不足タグ、相対画像 URL、画像サイズ、HTTP 状態、タイトル長、説明文長を確認します。",
	"knowledge_title":      "共有リンクを設計された広告枠のように見せる",
	"knowledge_body":       "Open Graph は、ソーシャルプラットフォームやメッセージアプリがページプレビューを読むための共通言語です。",
	"what_og_title":        "Open Graph とは？",
	"what_og_body":         "OG タグは crawler 向けの metadata です。多くのサービスは og:title、og:description、og:image を読んでリンクプレビューを生成します。",
	"image_size_title":     "推奨画像サイズ",
	"image_size_body":      "安全なプレビュー画像サイズは 1200 x 630、比率は約 1.91:1 です。軽量な画像ほどメッセージアプリで素早く表示されます。",
	"backend_fetch_title":  "なぜバックエンドで取得するのか？",
	"backend_fetch_body":   "多くの crawler は JavaScript を実行せず、SPA の描画を待ちません。サーバー側取得のほうが実際の crawler の見え方に近くなります。",
	"minimal_tags_title":   "最小限必要な Open Graph タグ",
	"minimal_tags_body":    "まず title、description、image、canonical URL を補うと、多くのプレビュー問題を解消できます。",
	"platform_rules_title": "プラットフォームごとに metadata の解釈は異なります",
	"framework_title":      "主要な技術スタックでの導入",
	"faq_title":            "よくある質問",
	"gallery_empty":        "検査済みで条件を満たした Web サイトのプレビューがここに表示されます。",
	"gallery_intro":        "OpenGraphy が見つけた公開プレビューです。最初の行には新しい投稿を優先表示します。",
	"no_works":             "まだ対象作品はありません。",
	"ranking_empty":        "投票後にランキングが表示されます。",
	"load_pair":            "作品を読み込み",
	"footer_body":          "安全な Open Graph 診断、ソーシャルプレビュー検査、metadata 修正提案。",
	"language_label":       "言語",
	"apply_language":       "適用",
	"facebook_debugger":    "Facebook デバッガー",
	"linkedin_inspector":   "LinkedIn Post Inspector",
	"consent_full":         "言語設定、同意状態、匿名オンライン表示に必要な cookie のみ使用します。追跡 cookie や広告 cookie は使用しません。",
	"reject_necessary":     "必要なもののみ",
	"score_label":          "スコア",
	"fetch_label":          "取得",
	"source_label":         "ソース",
	"cached":               "キャッシュ",
	"fresh":                "新規",
	"title_label":          "タイトル",
	"description_label":    "説明",
	"image_label":          "画像",
	"final_url_label":      "最終 URL",
	"copy":                 "コピー",
	"all_good_title":       "問題ありません",
	"all_good_body":        "主要な metadata 問題は検出されませんでした。",
	"privacy_title":        "プライバシーポリシー",
	"privacy_body_1":       "OpenGraphy はデータを販売せず、広告 cookie を使用しません。送信された URL は診断のために取得され、Redis に最大 4 時間キャッシュされます。",
	"privacy_body_2":       "オンライン表示には匿名セッション cookie を使用します。レート制限では短期のソルト付きハッシュを使い、IP を平文保存しません。",
	"terms_title":          "利用規約",
	"terms_body":           "公開アクセスが許可されている Web サイトのみ検査してください。本サービスはプライベートネットワークをブロックし、不正利用を制限する場合があります。",
}

var koExtra = map[string]string{
	"meta_description":     "OpenGraphy는 Open Graph, Twitter Card, 소셜 미리보기를 무료로 검사하고 metadata 분석, 플랫폼 미리보기, 진단, 수정 제안을 제공합니다.",
	"hero_eyebrow":         "Open Graph / Twitter Card / 소셜 미리보기 검사",
	"hero_title_full":      "Open Graph meta tags를 한곳에서 생성, 검사, 미리보기",
	"hero_subtitle_full":   "URL을 입력하면 OpenGraphy가 서버에서 원본 HTML을 가져와 Open Graph와 Twitter Card metadata를 분석하고 플랫폼별 카드와 수정 제안을 보여줍니다.",
	"url_placeholder_full": "검사할 URL 입력, 예: https://example.com",
	"stats_online_full":    "온라인",
	"stats_visits_full":    "방문",
	"stats_checks_full":    "검사",
	"mock_title":           "공개 전에 공유 미리보기 확인",
	"mock_body":            "Facebook, LinkedIn, Slack, 메시지 앱에서 일관되고 매력적인 링크 카드를 보여줍니다.",
	"preview_score":        "미리보기 점수",
	"checks_passed":        "개 검사 통과",
	"platforms_title":      "링크가 실제로 공유되는 플랫폼을 위해 설계",
	"secure_title":         "안전한 가져오기, 정확한 분석, 신뢰 가능한 미리보기",
	"secure_body":          "OpenGraphy는 사용자의 브라우저가 아니라 백엔드에서 crawler와 비슷하게 DNS, 리디렉션, 사설망 차단, HTML 분석, 이미지 접근성을 확인합니다.",
	"feature_ssrf_title":   "SSRF 보호",
	"feature_ssrf_body":    "모든 요청에서 DNS, IP 범위, 리디렉션을 확인해 localhost, 내부망, 예약 주소 접근을 막습니다.",
	"feature_cache_title":  "4시간 캐시",
	"feature_cache_body":   "metadata와 미리보기 결과는 Redis에 임시 저장되며 URL 기록이나 개인정보를 장기 저장하지 않습니다.",
	"feature_diag_title":   "진단 제안",
	"feature_diag_body":    "누락 태그, 상대 이미지 URL, 이미지 크기, HTTP 상태, 제목과 설명 길이를 확인합니다.",
	"knowledge_title":      "모든 공유를 잘 설계된 광고 영역처럼 보이게",
	"knowledge_body":       "Open Graph는 소셜 플랫폼과 메신저가 페이지 미리보기를 읽는 공통 언어입니다.",
	"what_og_title":        "Open Graph란?",
	"what_og_body":         "OG 태그는 crawler를 위한 metadata입니다. 많은 플랫폼이 og:title, og:description, og:image를 읽어 링크 미리보기를 만듭니다.",
	"image_size_title":     "권장 이미지 크기",
	"image_size_body":      "가장 안정적인 미리보기 이미지는 1200 x 630, 약 1.91:1 비율입니다. 가벼운 이미지일수록 메신저에서 빠르게 표시됩니다.",
	"backend_fetch_title":  "왜 백엔드에서 가져오나요?",
	"backend_fetch_body":   "대부분의 crawler는 JavaScript를 실행하거나 SPA 렌더링을 기다리지 않습니다. 서버 측 가져오기가 실제 crawler 동작에 더 가깝습니다.",
	"minimal_tags_title":   "최소 Open Graph 태그",
	"minimal_tags_body":    "먼저 title, description, image, canonical URL을 채우면 대부분의 미리보기 문제를 줄일 수 있습니다.",
	"platform_rules_title": "플랫폼마다 metadata 해석 방식이 다릅니다",
	"framework_title":      "주요 기술 스택 연동",
	"faq_title":            "자주 묻는 질문",
	"gallery_empty":        "검사를 통과한 웹사이트 미리보기가 여기에 표시됩니다.",
	"gallery_intro":        "OpenGraphy가 발견한 공개 미리보기입니다. 첫 줄은 최신 제출을 우선 표시합니다.",
	"no_works":             "아직 적합한 작품이 없습니다.",
	"ranking_empty":        "투표 후 순위가 표시됩니다.",
	"load_pair":            "작품 불러오기",
	"footer_body":          "안전한 Open Graph 진단, 소셜 미리보기 테스트, metadata 수정 제안.",
	"language_label":       "언어",
	"apply_language":       "적용",
	"facebook_debugger":    "Facebook 디버거",
	"linkedin_inspector":   "LinkedIn Post Inspector",
	"consent_full":         "언어 설정, 동의 상태, 익명 온라인 표시를 위한 필수 cookie만 사용합니다. 추적 또는 광고 cookie는 사용하지 않습니다.",
	"reject_necessary":     "필수만",
	"score_label":          "점수",
	"fetch_label":          "가져오기",
	"source_label":         "소스",
	"cached":               "캐시",
	"fresh":                "새로 가져옴",
	"title_label":          "제목",
	"description_label":    "설명",
	"image_label":          "이미지",
	"final_url_label":      "최종 URL",
	"copy":                 "복사",
	"all_good_title":       "문제 없음",
	"all_good_body":        "주요 metadata 문제가 감지되지 않았습니다.",
	"privacy_title":        "개인정보 처리방침",
	"privacy_body_1":       "OpenGraphy는 데이터를 판매하지 않으며 광고 cookie를 사용하지 않습니다. 제출된 URL은 진단 생성을 위해 가져오고 Redis에 최대 4시간 캐시합니다.",
	"privacy_body_2":       "온라인 표시는 익명 세션 cookie를 사용합니다. 속도 제한은 단기 salted hash를 사용하며 IP를 평문 저장하지 않습니다.",
	"terms_title":          "이용 약관",
	"terms_body":           "공개적으로 검사할 권한이 있는 웹사이트만 사용하세요. 서비스는 사설 네트워크를 차단하며 남용 트래픽을 제한할 수 있습니다.",
}

var esExtra = map[string]string{
	"meta_description":     "OpenGraphy es una herramienta gratuita para revisar Open Graph, Twitter Card y vistas previas sociales con análisis de metadata, diagnósticos y sugerencias de reparación.",
	"hero_eyebrow":         "Open Graph / Twitter Card / Comprobador de vista previa social",
	"hero_title_full":      "Genera, inspecciona y previsualiza meta tags Open Graph en un solo lugar",
	"hero_subtitle_full":   "Introduce una URL y OpenGraphy obtiene el HTML original desde el servidor, analiza Open Graph y Twitter Card, muestra tarjetas sociales y propone mejoras concretas.",
	"url_placeholder_full": "Introduce una URL, por ejemplo https://example.com",
	"stats_online_full":    "en línea", "stats_visits_full": "visitas", "stats_checks_full": "comprobaciones",
	"mock_title":    "Comprueba la vista previa antes de publicar",
	"mock_body":     "Haz que tus enlaces se vean claros, coherentes y atractivos en Facebook, LinkedIn, Slack y apps de mensajería.",
	"preview_score": "Puntuación", "checks_passed": "comprobaciones superadas",
	"platforms_title":    "Preparado para los lugares donde realmente viajan los enlaces",
	"secure_title":       "Captura segura, análisis preciso y vistas previas confiables",
	"secure_body":        "OpenGraphy no obtiene el sitio desde el navegador del usuario. El backend simula el comportamiento de un crawler y revisa DNS, redirecciones, redes privadas, HTML inicial e imágenes.",
	"feature_ssrf_title": "Protección SSRF", "feature_ssrf_body": "Cada petición valida DNS, rangos IP y redirecciones para evitar localhost, redes privadas y direcciones reservadas.",
	"feature_cache_title": "Caché de 4 horas", "feature_cache_body": "La metadata y los resultados se guardan temporalmente en Redis sin historial largo de URLs ni datos personales.",
	"feature_diag_title": "Diagnósticos", "feature_diag_body": "Comprueba etiquetas faltantes, imágenes relativas, tamaño de imagen, estado HTTP y longitud de título y descripción.",
	"knowledge_title": "Haz que cada enlace compartido parezca un anuncio bien diseñado",
	"knowledge_body":  "Open Graph es el lenguaje común que usan plataformas sociales y mensajeros para leer vistas previas de páginas.",
	"what_og_title":   "¿Qué es Open Graph?", "what_og_body": "Las etiquetas OG son metadata para crawlers. Muchas plataformas leen og:title, og:description y og:image para crear vistas previas.",
	"image_size_title": "Tamaño recomendado de imagen", "image_size_body": "El tamaño más seguro es 1200 x 630, con proporción aproximada 1.91:1. Las imágenes ligeras se cargan mejor en mensajería.",
	"backend_fetch_title": "¿Por qué capturar desde el backend?", "backend_fetch_body": "La mayoría de crawlers no ejecutan JavaScript ni esperan a que renderice una SPA. La captura del servidor se parece más a lo que ve el crawler real.",
	"minimal_tags_title": "Etiquetas Open Graph mínimas", "minimal_tags_body": "Empieza por título, descripción, imagen y URL canónica. Eso soluciona la mayoría de problemas de vista previa.",
	"platform_rules_title": "Cada plataforma interpreta la misma metadata de forma distinta",
	"framework_title":      "Integración en stacks comunes",
	"faq_title":            "Preguntas frecuentes",
	"gallery_empty":        "Las vistas previas calificadas aparecerán aquí después de las comprobaciones.",
	"gallery_intro":        "Vistas previas públicas calificadas por OpenGraphy. La primera fila favorece los envíos más recientes.",
	"no_works":             "Aún no hay trabajos calificados.", "ranking_empty": "El ranking aparecerá después de los votos.", "load_pair": "Cargar par",
	"footer_body":    "Diagnóstico seguro de Open Graph, pruebas de vista previa social y sugerencias de metadata.",
	"language_label": "Idioma", "apply_language": "Aplicar", "facebook_debugger": "Depurador de Facebook", "linkedin_inspector": "LinkedIn Post Inspector",
	"consent_full":     "Solo usamos cookies necesarias para idioma, consentimiento y presencia anónima. Sin cookies de rastreo ni publicidad.",
	"reject_necessary": "Solo necesarias",
	"score_label":      "puntuación", "fetch_label": "captura", "source_label": "origen", "cached": "caché", "fresh": "nuevo",
	"title_label": "Título", "description_label": "Descripción", "image_label": "Imagen", "final_url_label": "URL final", "copy": "Copiar",
	"all_good_title": "Todo bien", "all_good_body": "No se detectaron problemas importantes de metadata.",
	"privacy_title": "Política de privacidad", "privacy_body_1": "OpenGraphy no vende datos ni usa cookies publicitarias. Las URLs enviadas se usan para generar diagnósticos y se cachean en Redis hasta 4 horas.", "privacy_body_2": "La presencia usa una cookie anónima. Los límites de uso emplean hashes temporales con sal y no guardan IP en texto claro.",
	"terms_title": "Términos", "terms_body": "Usa OpenGraphy solo con sitios públicos que puedas inspeccionar. El servicio bloquea redes privadas y puede limitar tráfico abusivo.",
}

var frExtra = map[string]string{
	"meta_description":     "OpenGraphy est un outil gratuit pour vérifier Open Graph, Twitter Card et les aperçus sociaux avec analyse de metadata, diagnostics et conseils de correction.",
	"hero_eyebrow":         "Open Graph / Twitter Card / Vérificateur d’aperçu social",
	"hero_title_full":      "Générez, inspectez et prévisualisez les meta tags Open Graph au même endroit",
	"hero_subtitle_full":   "Saisissez une URL et OpenGraphy récupère le HTML côté serveur, analyse Open Graph et Twitter Card, affiche les cartes sociales et propose des corrections.",
	"url_placeholder_full": "Saisissez une URL, par exemple https://example.com",
	"stats_online_full":    "en ligne", "stats_visits_full": "visites", "stats_checks_full": "vérifications",
	"mock_title":    "Vérifiez l’aperçu avant publication",
	"mock_body":     "Rendez vos liens cohérents, clairs et attractifs sur Facebook, LinkedIn, Slack et les applications de messagerie.",
	"preview_score": "Score d’aperçu", "checks_passed": "vérifications réussies",
	"platforms_title":    "Conçu pour les endroits où les liens circulent réellement",
	"secure_title":       "Récupération sûre, analyse précise, aperçus fiables",
	"secure_body":        "OpenGraphy ne récupère pas les sites depuis le navigateur. Le backend simule un crawler et vérifie DNS, redirections, réseaux privés, HTML initial et images.",
	"feature_ssrf_title": "Protection SSRF", "feature_ssrf_body": "Chaque requête valide DNS, plages IP et redirections afin d’éviter localhost, réseaux privés et adresses réservées.",
	"feature_cache_title": "Cache de 4 heures", "feature_cache_body": "La metadata et les aperçus sont temporairement stockés dans Redis sans historique long ni données personnelles.",
	"feature_diag_title": "Diagnostics", "feature_diag_body": "Balises manquantes, images relatives, taille d’image, statut HTTP, longueur du titre et de la description sont vérifiés.",
	"knowledge_title": "Transformez chaque partage en emplacement marketing soigné",
	"knowledge_body":  "Open Graph est le langage commun utilisé par les plateformes sociales et messageries pour lire les aperçus de page.",
	"what_og_title":   "Qu’est-ce qu’Open Graph ?", "what_og_body": "Les balises OG sont de la metadata pour crawlers. Beaucoup de plateformes lisent og:title, og:description et og:image pour créer un aperçu.",
	"image_size_title": "Taille d’image recommandée", "image_size_body": "La taille la plus sûre est 1200 x 630, ratio environ 1.91:1. Les images légères s’affichent plus vite dans les messageries.",
	"backend_fetch_title": "Pourquoi récupérer côté backend ?", "backend_fetch_body": "La plupart des crawlers n’exécutent pas JavaScript et n’attendent pas les SPA. La récupération serveur reflète mieux ce qu’ils voient.",
	"minimal_tags_title": "Balises Open Graph minimales utiles", "minimal_tags_body": "Commencez par title, description, image et URL canonique. Cela résout la plupart des aperçus cassés.",
	"platform_rules_title": "Chaque plateforme interprète la même metadata différemment",
	"framework_title":      "Intégration dans les stacks courantes",
	"faq_title":            "FAQ",
	"gallery_empty":        "Les aperçus qualifiés apparaîtront ici après vérification.",
	"gallery_intro":        "Aperçus publics qualifiés trouvés par OpenGraphy. La première ligne privilégie les soumissions récentes.",
	"no_works":             "Aucun aperçu qualifié pour le moment.", "ranking_empty": "Le classement apparaîtra après les votes.", "load_pair": "Charger une paire",
	"footer_body":    "Diagnostics Open Graph sécurisés, tests d’aperçus sociaux et conseils de metadata.",
	"language_label": "Langue", "apply_language": "Appliquer", "facebook_debugger": "Débogueur Facebook", "linkedin_inspector": "LinkedIn Post Inspector",
	"consent_full":     "Nous utilisons uniquement les cookies nécessaires pour la langue, le consentement et la présence anonyme. Aucun cookie publicitaire ou de suivi.",
	"reject_necessary": "Nécessaires seulement",
	"score_label":      "score", "fetch_label": "récupération", "source_label": "source", "cached": "cache", "fresh": "récent",
	"title_label": "Titre", "description_label": "Description", "image_label": "Image", "final_url_label": "URL finale", "copy": "Copier",
	"all_good_title": "Tout va bien", "all_good_body": "Aucun problème majeur de metadata détecté.",
	"privacy_title": "Politique de confidentialité", "privacy_body_1": "OpenGraphy ne vend pas de données et n’utilise pas de cookies publicitaires. Les URLs soumises servent au diagnostic et sont mises en cache dans Redis jusqu’à 4 heures.", "privacy_body_2": "La présence utilise une session anonyme. Les limites de débit utilisent des hashes temporaires salés sans stocker les IP en clair.",
	"terms_title": "Conditions", "terms_body": "Utilisez OpenGraphy uniquement pour des sites publics que vous êtes autorisé à inspecter. Le service bloque les réseaux privés et peut limiter les abus.",
}

var deExtra = map[string]string{
	"meta_description":     "OpenGraphy ist ein kostenloser Prüfer für Open Graph, Twitter Card und Social Previews mit metadata Analyse, Plattformvorschau, Diagnosen und Reparaturhinweisen.",
	"hero_eyebrow":         "Open Graph / Twitter Card / Social Preview Prüfung",
	"hero_title_full":      "Open Graph meta tags an einem Ort erzeugen, prüfen und anzeigen",
	"hero_subtitle_full":   "Geben Sie eine URL ein. OpenGraphy lädt das ursprüngliche HTML serverseitig, analysiert Open Graph und Twitter Card, zeigt Social Cards und konkrete Verbesserungen.",
	"url_placeholder_full": "URL eingeben, zum Beispiel https://example.com",
	"stats_online_full":    "online", "stats_visits_full": "Besuche", "stats_checks_full": "Prüfungen",
	"mock_title":    "Vorschau vor der Veröffentlichung prüfen",
	"mock_body":     "Ihre Links wirken konsistent, klar und überzeugend auf Facebook, LinkedIn, Slack und in Messenger-Apps.",
	"preview_score": "Preview Score", "checks_passed": "Prüfungen bestanden",
	"platforms_title":    "Gebaut für die Orte, an denen Links wirklich geteilt werden",
	"secure_title":       "Sicheres Abrufen, präzise Analyse, verlässliche Vorschauen",
	"secure_body":        "OpenGraphy ruft Zielseiten nicht im Browser des Nutzers ab. Das Backend simuliert crawler Verhalten und prüft DNS, Weiterleitungen, private Netze, initiales HTML und Bilder.",
	"feature_ssrf_title": "SSRF-Schutz", "feature_ssrf_body": "Jeder Abruf validiert DNS, IP-Bereiche und Weiterleitungen, um localhost, private Netze und reservierte Adressen zu vermeiden.",
	"feature_cache_title": "4-Stunden-Cache", "feature_cache_body": "Metadata und Vorschauen werden temporär in Redis gespeichert, ohne langfristige URL-Historie oder personenbezogene Daten.",
	"feature_diag_title": "Diagnosen", "feature_diag_body": "Fehlende Tags, relative Bild-URLs, Bildgröße, HTTP-Status sowie Titel- und Beschreibungslänge werden geprüft.",
	"knowledge_title": "Jeder geteilte Link wie eine gestaltete Anzeige",
	"knowledge_body":  "Open Graph ist die gemeinsame Sprache, mit der soziale Plattformen und Messenger Seitenvorschauen lesen.",
	"what_og_title":   "Was ist Open Graph?", "what_og_body": "OG-Tags sind metadata für crawler. Viele Plattformen lesen og:title, og:description und og:image, um Linkvorschauen zu erzeugen.",
	"image_size_title": "Empfohlene Bildgröße", "image_size_body": "Die sicherste Größe ist 1200 x 630 mit etwa 1.91:1. Kleine Bilder werden in Messengern schneller angezeigt.",
	"backend_fetch_title": "Warum serverseitig abrufen?", "backend_fetch_body": "Die meisten crawler führen kein JavaScript aus und warten nicht auf SPAs. Serverseitiges Abrufen ist näher an ihrer tatsächlichen Sicht.",
	"minimal_tags_title": "Minimale nützliche Open Graph Tags", "minimal_tags_body": "Beginnen Sie mit title, description, image und canonical URL. Das behebt die meisten Preview-Probleme.",
	"platform_rules_title": "Plattformen interpretieren dieselbe metadata unterschiedlich",
	"framework_title":      "Integration in gängige Stacks",
	"faq_title":            "FAQ",
	"gallery_empty":        "Qualifizierte Website-Vorschauen erscheinen hier nach der Prüfung.",
	"gallery_intro":        "Öffentliche qualifizierte Vorschauen, gefunden von OpenGraphy. Die erste Reihe bevorzugt neue Einreichungen.",
	"no_works":             "Noch keine qualifizierten Einträge.", "ranking_empty": "Die Rangliste erscheint nach den Abstimmungen.", "load_pair": "Paar laden",
	"footer_body":    "Sichere Open Graph Diagnosen, Social Preview Tests und metadata Reparaturhinweise.",
	"language_label": "Sprache", "apply_language": "Anwenden", "facebook_debugger": "Facebook Debugger", "linkedin_inspector": "LinkedIn Post Inspector",
	"consent_full":     "Wir verwenden nur notwendige Cookies für Sprache, Zustimmung und anonyme Online-Präsenz. Keine Tracking- oder Werbe-Cookies.",
	"reject_necessary": "Nur notwendige",
	"score_label":      "Score", "fetch_label": "Abruf", "source_label": "Quelle", "cached": "Cache", "fresh": "neu",
	"title_label": "Titel", "description_label": "Beschreibung", "image_label": "Bild", "final_url_label": "Finale URL", "copy": "Kopieren",
	"all_good_title": "Alles gut", "all_good_body": "Keine größeren metadata Probleme erkannt.",
	"privacy_title": "Datenschutzerklärung", "privacy_body_1": "OpenGraphy verkauft keine Daten und verwendet keine Werbe-Cookies. Eingereichte URLs werden für Diagnosen genutzt und bis zu 4 Stunden in Redis gecacht.", "privacy_body_2": "Präsenz nutzt eine anonyme Session. Rate Limits verwenden kurzlebige gesalzene Hashes und speichern IPs nicht im Klartext.",
	"terms_title": "Bedingungen", "terms_body": "Nutzen Sie OpenGraphy nur für öffentliche Websites, die Sie prüfen dürfen. Der Dienst blockiert private Netze und kann missbräuchlichen Traffic begrenzen.",
}

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
		"apply_language":       "应用语言",
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
