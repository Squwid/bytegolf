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
	LanguageKey("node", "22"): {
		Image:     "node:22-alpine3.19",
		Extension: "js",
		Command:   "node",
	},
	LanguageKey("php", "8.2.18"): {
		Image:     "php:8.2.18-fpm-alpine3.19",
		Extension: "php",
		Command:   "php",
	},
	LanguageKey("go", "1.22.2"): {
		Image:     "golang:1.22.2-alpine3.19",
		Extension: "go",
		Command:   "go run",
	},
	LanguageKey("bash", "5.2.26"): {
		Image:     "bash:5.2.26",
		Extension: "sh",
		Command:   "bash",
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
