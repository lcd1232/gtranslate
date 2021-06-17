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
	`[[["частный","private",null,null,10]
,[null,null,"chastnyy","ˈprīvit"]
]
,[["имя прилагательное",["частный","личный","конфиденциальный","рядовой","уединенный","негласный","неофициальный","тайный"]
,[["частный",["private","partial","particular","individual","proprietary","local"]
,null,0.08737902]
,["личный",["private","personal","individual","intimate","identity","direct"]
,null,0.016939197]
,["конфиденциальный",["confidential","private","privy","tete-a-tete"]
,null,9.2624204E-4]
,["рядовой",["private","common","common-or-garden","rank-and-file"]
,null,8.4335293E-4]
,["уединенный",["secluded","solitary","lonely","private","remote","cloistered"]
,null,1.8526005E-4]
,["негласный",["secret","private","backroom"]
,null,9.611165E-5]
,["неофициальный",["informal","unofficial","private","inofficial","officious","unceremonious"]
,null,1.7778868E-5]
,["тайный",["secret","covert","clandestine","arcane","undercover","private"]
,null,1.6964714E-5]
]
,"private",3]
,["имя существительное",["рядовой","половые органы"]
,[["рядовой",["soldier","private","common soldier","squaddie","Tommy Atkins","tommy"]
,null,8.4335293E-4,null,1]
,["половые органы",["genitals","genitalia","private parts","private","privates","privy parts"]
]
]
,"private",1]
]
,"en",null,null,[["private",null,[["частный",1000,true,false,[10,3,0]
]
]
,[[0,7]
]
,"private",0,0]
]
,0.9760956,[]
,[["en"]
,null,[0.9760956]
,["en"]
]
,null,null,[["имя прилагательное",[[["personal","one's own","individual","particular","special","exclusive","privately owned"]
,"m_en_gbus0816730.006"]
,[["confidential","strictly confidential","secret","top secret","classified","unofficial","off the record","not for publication","not to be made public","not to be disclosed","closet","backstage","offstage","privileged","one-on-one","tête-à-tête","covert","clandestine","surreptitious","in camera"]
,"m_en_gbus0816730.008"]
,[["hush-hush"]
,"m_en_gbus0816730.008",[["informal"]
]
]
,[["intimate","personal","secret","innermost","inward","unspoken","undeclared","undisclosed","unvoiced","sneaking","hidden"]
,"m_en_gbus0816730.009"]
,[["reserved","introvert","introverted","self-contained","reticent","discreet","uncommunicative","noncommunicative","media-shy","unforthcoming","secretive","retiring","ungregarious","unsocial","unsociable","withdrawn","solitary","insular","reclusive","hermitlike","hermitic"]
,"m_en_gbus0816730.010"]
,[["secluded","secret","quiet","undisturbed","concealed","hidden","remote","isolated","out of the way","sequestered"]
,"m_en_gbus0816730.012"]
,[["unofficial","personal","non-official","nonpublic"]
,"m_en_gbus0816730.016"]
,[["independent","non-state-controlled","non-state-run","privatized","denationalized","nonpublic","commercial","private-enterprise"]
,"m_en_gbus0816730.018"]
]
,"private"]
,["имя существительное",[[["swad","swaddy"]
,"m_en_gbus0816730.023",[["archaic"]
]
]
,[["private soldier","common soldier","infantryman","foot soldier","trooper","sapper","ranker","GI","enlisted man","poilu","jawan"]
,"m_en_gbus0816730.023"]
,[["Tommy","squaddie","Tommy Atkins","grunt","buck private","digger","troopie"]
,"m_en_gbus0816730.023",[["informal"]
]
]
]
,"private"]
]
,[["имя прилагательное",[["belonging to or for the use of one particular person or group of people only.","m_en_gbus0816730.006","all bedrooms have private facilities"]
,["(of a person) having no official or public role or position.","m_en_gbus0816730.015","the paintings were sold to a private collector"]
,["(of a service or industry) provided or owned by an individual or an independent, commercial company rather than by the government.","m_en_gbus0816730.018","research projects carried out by private industry"]
]
,"private"]
,["имя существительное",[["a soldier of the lowest rank, in particular an enlisted person in the US Army or Marine Corps ranking below private first class.","m_en_gbus0816730.023"]
,["short for private parts.","m_en_gbus0816730.027",null,[["informal"]
]
]
]
,"private"]
]
,[[["a small \u003cb\u003eprivate\u003c/b\u003e service in the chapel",null,null,null,3,"m_en_gbus0816730.008"]
,["his \u003cb\u003eprivate\u003c/b\u003e plane",null,null,null,3,"m_en_gbus0816730.006"]
,["he was a very \u003cb\u003eprivate\u003c/b\u003e man",null,null,null,3,"m_en_gbus0816730.010"]
,["all bedrooms have \u003cb\u003eprivate\u003c/b\u003e facilities",null,null,null,3,"m_en_gbus0816730.006"]
,["the paintings were sold to a \u003cb\u003eprivate\u003c/b\u003e collector",null,null,null,3,"m_en_gbus0816730.015"]
,["he would continue to represent her in a \u003cb\u003eprivate\u003c/b\u003e capacity as advisor and confidant",null,null,null,3,"m_en_gbus0816730.016"]
,["she felt awkward intruding on \u003cb\u003eprivate\u003c/b\u003e grief",null,null,null,3,"m_en_gbus0816730.009"]
,["more than 1,400 state enterprises that were about to go \u003cb\u003eprivate\u003c/b\u003e",null,null,null,3,"m_en_gbus0816730.018"]
,["if I could afford it I'd go \u003cb\u003eprivate\u003c/b\u003e",null,null,null,3,"m_en_gbus0816730.019"]
,["research projects carried out by \u003cb\u003eprivate\u003c/b\u003e industry",null,null,null,3,"m_en_gbus0816730.018"]
,["\u003cb\u003eprivate\u003c/b\u003e education",null,null,null,3,"m_en_gbus0816730.019"]
,["it was a \u003cb\u003eprivate\u003c/b\u003e sale—no agent's commission",null,null,null,3,"m_en_gbus0816730.020"]
,["can we go somewhere a little more \u003cb\u003eprivate\u003c/b\u003e?",null,null,null,3,"m_en_gbus0816730.012"]
,["this is a \u003cb\u003eprivate\u003c/b\u003e conversation",null,null,null,3,"m_en_gbus0816730.011"]
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
		{
			name: "private",
			text: "private",
			body: responses[2],
			want: Translated{
				Original:              "private",
				OriginalLanguage:      "en",
				OriginalPronunciation: "ˈprīvit",
				Text:                  "частный",
				Pronunciation:         "chastnyy",
				Definitions: []string{
					"belonging to or for the use of one particular person or group of people only.",
					"a soldier of the lowest rank, in particular an enlisted person in the US Army or Marine Corps ranking below private first class.",
				},
				Examples: []string{
					"all bedrooms have private facilities",
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
