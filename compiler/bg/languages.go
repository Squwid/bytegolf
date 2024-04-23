package bg

import "fmt"

type Languages map[string]Language

type Language struct {
	Image     string `json:"image"`
	Extension string `json:"extension"`
	Command   string `json:"cmd"`
}

var languages = Languages{
	LanguageKey("python3", "3.12.2"): {
		Image:     "python:3.12.2-alpine3.19",
		Extension: "py",
		Command:   "python3",
	},
}

func (langs Languages) get(lang, version string) *Language {
	l, ok := langs[LanguageKey(lang, version)]
	if !ok {
		return nil
	}
	return &l
}

// LanguageKey takes a language and version and crafts a lanugage key
func LanguageKey(lang, version string) string {
	return fmt.Sprintf("%s:%s", lang, version)
}

// GetLanguage gets an active language using the language and version. If no
// language exists, nil will be returned
func GetLanguage(lang, version string) *Language {
	return languages.get(lang, version)
}
