package gtranslate

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

func TestTranslateWithFromTo(t *testing.T) {
	for i := 0; i < 4; i++ {
		for _, ta := range testingTable {
			resp, err := TranslateWithParams(ta.inText, TranslationParams{
				From:       ta.langFrom,
				To:         ta.langTo,
				Tries:      5,
				Delay:      time.Second,
				GoogleHost: "google.cn",
			})
			if err != nil {
				t.Error(err, err.Error())
				t.Fail()
			}
			if resp != ta.outText {
				t.Error("translated text is not the expected", ta.outText, " != ", resp)
			}
		}
	}
}

func TestTranslateAdvanced(t *testing.T) {
	responseBody := `[[["Привет","hi",null,null,1]
,[null,null,"Privet","hī"]
]
,[["сокращение",["Гавайи"]
,[["Гавайи",["HI"]
,null,0.029729217]
]
,"HI",6]
]
,"en",null,null,[["hi",null,[["Привет",1000,true,false,[1]
]
,["Здравствуй",1000,true,false,[1]
]
]
,[[0,2]
]
,"hi",0,0]
]
,0,[]
,[["en"]
,null,[0]
,["en"]
]
,null,null,null,[["восклицание",[["used as a friendly greeting or to attract attention.","m_en_gbus0465760.005","“Hi there. How was the flight?”",[["informal"]
]
]
]
,"hi"]
,["сокращение",[["Hawaii (in official postal use).","m_en_gbus0465770.003"]
]
,"hi"]
]
]`

	WithTestServer(t, 200, []byte(responseBody), func(url string, requestBodyCh <-chan []byte) {
		translated, err := TranslateAdvanced("hi", language.English, language.Russian, url)
		require.NoError(t, err)
		assert.Equal(t, Translated{
			Original:              "hi",
			OriginalLanguage:      "en",
			OriginalPronunciation: "hī",
			Text:                  "Привет",
			Pronunciation:         "Privet",
			Definitions: []string{
				"used as a friendly greeting or to attract attention.",
				"Hawaii (in official postal use).",
			},
			Examples: []string{
				"“Hi there. How was the flight?”",
			},
		}, translated)
	})
}

func WithTestServer(t *testing.T, responseCode int, responseBody []byte, f func(url string, requestBodyCh <-chan []byte)) {
	ch := make(chan []byte, 1)
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		w.WriteHeader(responseCode)
		if len(responseBody) > 0 {
			_, err := w.Write(responseBody)
			assert.NoError(t, err)
		}
		ch <- b
	}))
	defer s.Close()

	f(s.URL+"/", ch)
}
