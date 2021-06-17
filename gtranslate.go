package gtranslate

import (
	"time"

	"golang.org/x/text/language"
)

var GoogleHost = "https://translate.google.com"

// TranslationParams is a util struct to pass as parameter to indicate how to translate
type TranslationParams struct {
	From       string
	To         string
	Tries      int
	Delay      time.Duration
	GoogleHost string
}

// Translate translate a text using native tags offer by go language
func Translate(text string, from language.Tag, to language.Tag, googleHost ...string) (string, error) {
	host := GoogleHost
	if len(googleHost) != 0 && googleHost[0] != "" {
		host = googleHost[0]
	}
	translated, err := translate(text, from.String(), to.String(), false, 2, 0, host)
	if err != nil {
		return "", err
	}

	return translated, nil
}

// TranslateWithParams translate a text with simple params as string
func TranslateWithParams(text string, params TranslationParams) (string, error) {
	if params.GoogleHost == "" {
		params.GoogleHost = GoogleHost
	}
	translated, err := translate(text, params.From, params.To, true, params.Tries, params.Delay, params.GoogleHost)
	if err != nil {
		return "", err
	}
	return translated, nil
}

type Translated struct {
	OriginalLanguage      string
	Original              string
	OriginalPronunciation string
	TextLanguage          string
	Text                  string
	Pronunciation         string
	Definitions           []string
	Examples              []string
}

func TranslateAdvanced(text string, from language.Tag, to language.Tag, googleHost ...string) (Translated, error) {
	host := GoogleHost
	if len(googleHost) != 0 && googleHost[0] != "" {
		host = googleHost[0]
	}

	translated, err := translateAdvanced(text, from.String(), to.String(), false, 2, 0, host)
	if err != nil {
		return Translated{}, err
	}

	return translated, nil
}
