package gtranslate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/text/language"

	"github.com/robertkrimen/otto"
)

var ttk otto.Value

func init() {
	ttk, _ = otto.ToValue("0")
}

const (
	defaultNumberOfRetries = 2
)

func translate(text, from, to string, withVerification bool, tries int, delay time.Duration, host string) (string, error) {
	if tries == 0 {
		tries = defaultNumberOfRetries
	}

	if withVerification {
		if _, err := language.Parse(from); err != nil && from != "auto" {
			log.Println("[WARNING], '" + from + "' is a invalid language, switching to 'auto'")
			from = "auto"
		}
		if _, err := language.Parse(to); err != nil {
			log.Println("[WARNING], '" + to + "' is a invalid language, switching to 'en'")
			to = "en"
		}
	}

	t, _ := otto.ToValue(text)

	urll := fmt.Sprintf("https://translate.%s/translate_a/single", host)

	token := get(t, ttk)

	data := map[string]string{
		"client": "gtx",
		"sl":     from,
		"tl":     to,
		"hl":     to,
		// "dt":     []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"},
		"ie":   "UTF-8",
		"oe":   "UTF-8",
		"otf":  "1",
		"ssel": "0",
		"tsel": "0",
		"kc":   "7",
		"q":    text,
	}

	u, err := url.Parse(urll)
	if err != nil {
		return "", nil
	}

	parameters := url.Values{}

	for k, v := range data {
		parameters.Add(k, v)
	}
	for _, v := range []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"} {
		parameters.Add("dt", v)
	}

	parameters.Add("tk", token)
	u.RawQuery = parameters.Encode()

	var r *http.Response

	for tries > 0 {
		r, err = http.Get(u.String())
		if err != nil {
			if err == http.ErrHandlerTimeout {
				return "", errBadNetwork
			}
			return "", err
		}

		if r.StatusCode == http.StatusOK {
			break
		}

		if r.StatusCode == http.StatusForbidden {
			tries--
			time.Sleep(delay)
		}
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}

	var resp []interface{}

	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return "", err
	}

	responseText := ""
	for _, obj := range resp[0].([]interface{}) {
		if len(obj.([]interface{})) == 0 {
			break
		}

		t, ok := obj.([]interface{})[0].(string)
		if ok {
			responseText += t
		}
	}

	return responseText, nil
}

func translateAdvanced(
	text,
	from,
	to string,
	withVerification bool,
	tries int,
	delay time.Duration,
	host string,
) (Translated, error) {
	if tries == 0 {
		tries = defaultNumberOfRetries
	}

	if withVerification {
		if _, err := language.Parse(from); err != nil && from != "auto" {
			log.Println("[WARNING], '" + from + "' is a invalid language, switching to 'auto'")
			from = "auto"
		}
		if _, err := language.Parse(to); err != nil {
			log.Println("[WARNING], '" + to + "' is a invalid language, switching to 'en'")
			to = "en"
		}
	}

	t, err := otto.ToValue(text)
	if err != nil {
		return Translated{}, err
	}

	urll := fmt.Sprintf("%s/translate_a/single", host)

	token := get(t, ttk)

	data := map[string]string{
		"client": "gtx",
		"sl":     from,
		"tl":     to,
		"hl":     to,
		// "dt":     []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"},
		"ie":   "UTF-8",
		"oe":   "UTF-8",
		"otf":  "1",
		"ssel": "0",
		"tsel": "0",
		"kc":   "7",
		"q":    text,
	}

	u, err := url.Parse(urll)
	if err != nil {
		return Translated{}, nil
	}

	parameters := url.Values{}

	for k, v := range data {
		parameters.Add(k, v)
	}
	for _, v := range []string{"at", "bd", "ex", "ld", "md", "qca", "rw", "rm", "ss", "t"} {
		parameters.Add("dt", v)
	}

	parameters.Add("tk", token)
	u.RawQuery = parameters.Encode()

	var r *http.Response

	for tries > 0 {
		r, err = http.Get(u.String())
		if err != nil {
			if err == http.ErrHandlerTimeout {
				return Translated{}, errBadNetwork
			}
			return Translated{}, err
		}

		if r.StatusCode == http.StatusOK {
			break
		}

		if r.StatusCode == http.StatusForbidden {
			tries--
			time.Sleep(delay)
		}
	}

	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return Translated{}, err
	}

	var resp []interface{}

	err = json.Unmarshal(raw, &resp)
	if err != nil {
		return Translated{}, err
	}

	result := Translated{
		Original: text,
	}

	for i, obj := range resp {
		part, ok := partsMapping[i]
		if !ok {
			continue
		}
		switch part {
		case dataTypeTranslation:
			var parts []string
			v, ok := obj.([]interface{})
			if !ok {
				continue
			}
			for i, v1 := range v {
				v2, ok := v1.([]interface{})
				if !ok || len(v2) == 0 {
					continue
				}

				if part, ok := v2[0].(string); ok {
					parts = append(parts, part)
				}

				// if not last element
				if i != len(v)-1 && len(v2) < 2 {
					continue
				}

				if originalPronunciation, ok := v2[len(v2)-1].(string); ok {
					result.OriginalPronunciation = originalPronunciation
				}

				if pronunciation, ok := v2[len(v2)-2].(string); ok {
					result.Pronunciation = pronunciation
				}
			}
			result.Text = strings.Join(parts, "")
		case dataTypeOriginalLanguage:
			l, ok := obj.(string)
			if !ok {
				continue
			}
			result.OriginalLanguage = l
		case dataTypeDefinitions:
			v, ok := obj.([]interface{})
			if !ok {
				continue
			}
			for _, v1 := range v {
				v2, ok := v1.([]interface{})
				if !ok || len(v2) < 2 {
					continue
				}

				v3, ok := v2[1].([]interface{})
				if !ok || len(v3) == 0 {
					continue
				}

				v4, ok := v3[0].([]interface{})
				if !ok || len(v4) == 0 {
					continue
				}

				definition, ok := v4[0].(string)
				if !ok {
					continue
				}
				result.Definitions = append(result.Definitions, definition)

				if len(v4) < 3 {
					continue
				}
				example, ok := v4[2].(string)
				if !ok {
					continue
				}
				result.Examples = append(result.Examples, example)
			}
		}
	}

	return result, nil
}

type dataType int

const (
	dataTypeTranslation dataType = iota
	dataTypeAllTranslation
	dataTypeOriginalLanguage
	dataTypePossibleTranslations
	dataTypeConfidence
	dataTypePossibleMistakes
	dataTypeLanguage
	dataTypeSynonyms
	dataTypeDefinitions
	dataTypeExamples
	dataTypeSeeAlso
)

var partsMapping = map[int]dataType{
	0:  dataTypeTranslation,
	1:  dataTypeAllTranslation,
	2:  dataTypeOriginalLanguage,
	5:  dataTypePossibleTranslations,
	6:  dataTypeConfidence,
	7:  dataTypePossibleMistakes,
	8:  dataTypeLanguage,
	11: dataTypeSynonyms,
	12: dataTypeDefinitions,
	13: dataTypeExamples,
	14: dataTypeSeeAlso,
}
