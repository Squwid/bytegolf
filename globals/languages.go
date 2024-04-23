package globals

import "fmt"

type Languages map[string]JdoodleLang

type BGCompilerLanguage struct {
	Image     string
	Extension string
	Command   string
}

type JdoodleLang struct {
	JdoodleLang    string
	JdoodleVersion string
}

var languages = Languages{
	LanguageKey("python2", "2.7.16"): {
		JdoodleLang:    "python2",
		JdoodleVersion: "2",
	},

	LanguageKey("python3", "3.7.4"): {
		JdoodleLang:    "python3",
		JdoodleVersion: "3",
	},

	LanguageKey("java", "JDK 11.0.4"): {
		JdoodleLang:    "java",
		JdoodleVersion: "3",
	},

	LanguageKey("javascript", "12.11.1"): {
		JdoodleLang:    "nodejs",
		JdoodleVersion: "3",
	},

	LanguageKey("c++", "g++ 17 GCC 9.10"): {
		JdoodleLang:    "cpp17",
		JdoodleVersion: "0",
	},

	LanguageKey("php", "7.3.10"): {
		JdoodleLang:    "php",
		JdoodleVersion: "3",
	},

	LanguageKey("rust", "1.38.0"): {
		JdoodleLang:    "rust",
		JdoodleVersion: "3",
	},

	LanguageKey("go", "1.13.1"): {
		JdoodleLang:    "go",
		JdoodleVersion: "3",
	},

	LanguageKey("bash", "5.0.011"): {
		JdoodleLang:    "bash",
		JdoodleVersion: "3",
	},
}

func (langs Languages) get(lang, version string) *JdoodleLang {
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
func GetLanguage(lang, version string) *JdoodleLang {
	return languages.get(lang, version)
}
