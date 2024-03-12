package i18n

type I18nText interface {
	String() string
	Translate(lang Lang) string
}

type Text map[Lang]string

func (t Text) String() string {
	return t[defaultLang]
}

func (t Text) Translate(lang Lang) string {
	if v, ok := t[lang]; ok {
		return v
	}
	return t.String()
}
