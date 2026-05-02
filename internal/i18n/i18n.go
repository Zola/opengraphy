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
	return &Bundle{data: map[string]map[string]string{
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
