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

var responses = []string{
	`[[["Привет","hi",null,null,1]
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
]`,
	`[[["Привет","hello",null,null,10]
,[null,null,"Privet","həˈlō"]
]
,[["глагол",["здороваться","звать","окликать"]
,[["здороваться",["greet","hello","salute","hallo","hullo","bow"]
,null,0.0028530264]
,["звать",["call","invite","shout","hail","hallo","hello"]
,null,2.753645E-5]
,["окликать",["hail","holler","call","challenge","speak","hello"]
,null,2.753645E-5]
]
,"hello",2]
,["имя существительное",["приветствие","приветственный возглас","возглас удивления"]
,[["приветствие",["greeting","welcome","salute","salutation","hello","hail"]
,null,0.0014801305,null,3]
,["приветственный возглас",["hallo","halloa","hello","viva"]
,null,2.753645E-5,null,1]
,["возглас удивления",["hallo","halloa","hello"]
,null,2.753645E-5]
]
,"hello",1]
]
,"en",null,null,[["hello",null,[["Привет",1000,true,false,[10,3,0]
]
,["приветствовать",1000,true,false,[10]
]
,["Здравствуйте",1000,true,false,[10]
]
]
,[[0,5]
]
,"hello",0,0]
]
,1,[]
,[["en"]
,null,[1]
,["en"]
]
,null,null,null,[["восклицание",[["used as a greeting or to begin a phone conversation.","m_en_gbus0460730.012","hello there, Katie!"]
]
,"hello"]
,["имя существительное",[["an utterance of “hello”; a greeting.","m_en_gbus0460730.025","she was getting polite nods and hellos from people"]
]
,"hello"]
,["глагол",[["say or shout “hello”; greet someone.","m_en_gbus0460730.034","I pressed the phone button and helloed"]
]
,"hello"]
]
,[[["\u003cb\u003ehello\u003c/b\u003e there, Katie!",null,null,null,3,"m_en_gbus0460730.012"]
]
]
]
`,
}

func TestTranslateAdvanced(t *testing.T) {
	for _, tc := range []struct {
		name string
		text string
		body string
		want Translated
	}{
		{
			name: "hi",
			text: "hi",
			body: responses[0],
			want: Translated{
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
			},
		},
		{
			name: "hello",
			text: "hello",
			body: responses[1],
			want: Translated{
				Original:              "hello",
				OriginalLanguage:      "en",
				OriginalPronunciation: "həˈlō",
				Text:                  "Привет",
				Pronunciation:         "Privet",
				Definitions: []string{
					"used as a greeting or to begin a phone conversation.",
					"an utterance of “hello”; a greeting.",
					"say or shout “hello”; greet someone.",
				},
				Examples: []string{
					"hello there, Katie!",
					"she was getting polite nods and hellos from people",
					"I pressed the phone button and helloed",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			WithTestServer(t, 200, []byte(tc.body), func(url string, requestBodyCh <-chan []byte) {
				translated, err := TranslateAdvanced(tc.text, language.English, language.Russian, url)
				require.NoError(t, err)
				assert.Equal(t, tc.want, translated)
			})
		})
	}
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
