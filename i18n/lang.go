package i18n

type Lang string

const (
	ZH = "zh"
	EN = "en"
	JP = "jp"
	KR = "kr"
)

var defaultLang Lang = ZH

func SetDefaultLang(lang Lang) {
	defaultLang = lang
}
