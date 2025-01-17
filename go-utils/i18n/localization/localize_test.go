package localization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalize(t *testing.T) {
	const (
		nl1 = "Hoe gaat het met u?"
		nl2 = "U heeft {games, plural, =0 {geen gratis spins} one {1 gratis spin} other {# gratis spins}} gewonnen."
		nl3 = "U heeft {games, plural, one {1 gratis spin} oops {# gratis spins}} gewonnen."
		nl4 = "Dit is de {times, selectordinal, x {#e}} gratis spin."
		nl5 = "U kunt {min|max, pluralrange, one {0-1 gratis spin} x {#-# gratis spins}} winnen."
		nl6 = "{index} van {count}"
		nl7 = "{a1, number} - {a2, number, integer} - {a3, number, percent} - {a4, number, currency} - {a5, number, currency:GBP}"
		en1 = "How are you?"
		en2 = "You have won {games, plural, =0 {no free spins} one {1 free spin} other {# free spins}}."
		en3 = "You have won {games, plural, one {one free spin} x {# free spins}}."
		en4 = "This is the {times, selectordinal, one {#st} two {#nd} few {#rd} other {#th}} free spin."
		en7 = "{a1, number} - {a2, number, integer} - {a3, number, percent} - {a4, number, currency} - {a5, number, currency:GBP}"
		it1 = "Come stai?"
		it2 = "{games, plural, =0 {Non hai vinto giri gratis} one {Hai vinto 1 giro gratuito} other {Hai vinto # giri gratis}}."
		it3 = "Hai vinto {games, plural, one {1 giro gratuito} x {# giri gratis}}."
		it4 = "Questa è {times, selectordinal, many {l'#°} other {la #°}} volta."
		it7 = "{a1, number} - {a2, number, integer} - {a3, number, percent} - {a4, number, currency} - {a5, number, currency:GBP}"
	)

	testCases := []struct {
		name   string
		locale string
		input  string
		params map[string]string
		want   string
	}{
		{
			name:   "NL - msg1",
			locale: "nl",
			input:  nl1,
			want:   nl1,
		},
		{
			name:   "NL - msg2 - 0",
			locale: "nl",
			input:  nl2,
			params: map[string]string{"games": "0"},
			want:   "U heeft geen gratis spins gewonnen.",
		},
		{
			name:   "NL - msg2 - 1",
			locale: "nl-NL",
			input:  nl2,
			params: map[string]string{"games": "1"},
			want:   "U heeft 1 gratis spin gewonnen.",
		},
		{
			name:   "NL - msg2 - 2",
			locale: "nl-BE",
			input:  nl2,
			params: map[string]string{"games": "2"},
			want:   "U heeft 2 gratis spins gewonnen.",
		},
		{
			name:   "NL - msg3 - 0",
			locale: "nl",
			input:  nl3,
			params: map[string]string{"games": "0"},
			want:   "U heeft 0 gratis spins gewonnen.",
		},
		{
			name:   "NL - msg3 - 1",
			locale: "nl-be",
			input:  nl3,
			params: map[string]string{"games": "1"},
			want:   "U heeft 1 gratis spin gewonnen.",
		},
		{
			name:   "NL - msg3 - 8",
			locale: "nl",
			input:  nl3,
			params: map[string]string{"games": "8"},
			want:   "U heeft 8 gratis spins gewonnen.",
		},
		{
			name:   "NL - msg4 - 1",
			locale: "nl-nl",
			input:  nl4,
			params: map[string]string{"times": "1"},
			want:   "Dit is de 1e gratis spin.",
		},
		{
			name:   "NL - msg4 - 2",
			locale: "nl",
			input:  nl4,
			params: map[string]string{"times": "2"},
			want:   "Dit is de 2e gratis spin.",
		},
		{
			name:   "NL - msg4 - 5",
			locale: "nl",
			input:  nl4,
			params: map[string]string{"times": "5"},
			want:   "Dit is de 5e gratis spin.",
		},
		{
			name:   "NL - msg5 - 0,1",
			locale: "nl",
			input:  nl5,
			params: map[string]string{"min": "0", "max": "1"},
			want:   "U kunt 0-1 gratis spin winnen.",
		},
		{
			name:   "NL - msg5 - 1,2",
			locale: "nl",
			input:  nl5,
			params: map[string]string{"min": "1", "max": "2"},
			want:   "U kunt 1-2 gratis spins winnen.",
		},
		{
			name:   "NL - msg5 - 0,2",
			locale: "nl",
			input:  nl5,
			params: map[string]string{"min": "0", "max": "2"},
			want:   "U kunt 0-2 gratis spins winnen.",
		},
		{
			name:   "NL - msg6 - 1,5",
			locale: "nl",
			input:  nl6,
			params: map[string]string{"index": "1", "count": "5"},
			want:   "1 van 5",
		},
		{
			name:   "NL - msg6 - 3,8",
			locale: "nl",
			input:  nl6,
			params: map[string]string{"index": "3", "count": "8"},
			want:   "3 van 8",
		},
		{
			name:   "NL - msg7 - 1,2,3,4,5",
			locale: "nl",
			input:  nl7,
			params: map[string]string{"a1": "1", "a2": "2", "a3": "3", "a4": "4", "a5": "5"},
			want:   "1 - 2 - 3% - € 4,00 - £ 5,00",
		},
		{
			name:   "NL - msg7 - 1.1,2.2,3.3,4.4,5.5",
			locale: "nl",
			input:  nl7,
			params: map[string]string{"a1": "1.1", "a2": "2.2", "a3": "3.3", "a4": "4.4", "a5": "5.5"},
			want:   "1,1 - 2 - 3,3% - € 4,40 - £ 5,50",
		},
		{
			name:   "EN - msg1",
			locale: "en",
			input:  en1,
			want:   en1,
		},
		{
			name:   "EN - msg2 - 0",
			locale: "en-us",
			input:  en2,
			params: map[string]string{"games": "0"},
			want:   "You have won no free spins.",
		},
		{
			name:   "EN - msg2 - 1",
			locale: "en-GB",
			input:  en2,
			params: map[string]string{"games": "1"},
			want:   "You have won 1 free spin.",
		},
		{
			name:   "EN - msg2 - 2",
			locale: "en-GB",
			input:  en2,
			params: map[string]string{"games": "2"},
			want:   "You have won 2 free spins.",
		},
		{
			name:   "EN - msg3 - 0",
			locale: "en",
			input:  en3,
			params: map[string]string{"games": "0"},
			want:   "You have won 0 free spins.",
		},
		{
			name:   "EN - msg3 - 1",
			locale: "en-US",
			input:  en3,
			params: map[string]string{"games": "1"},
			want:   "You have won one free spin.",
		},
		{
			name:   "EN - msg3 - 8",
			locale: "en",
			input:  en3,
			params: map[string]string{"games": "8"},
			want:   "You have won 8 free spins.",
		},
		{
			name:   "EN - msg4 - 1",
			locale: "en",
			input:  en4,
			params: map[string]string{"times": "1"},
			want:   "This is the 1st free spin.",
		},
		{
			name:   "EN - msg4 - 2",
			locale: "en",
			input:  en4,
			params: map[string]string{"times": "2"},
			want:   "This is the 2nd free spin.",
		},
		{
			name:   "EN - msg4 - 5",
			locale: "en",
			input:  en4,
			params: map[string]string{"times": "5"},
			want:   "This is the 5th free spin.",
		},
		{
			name:   "EN - msg4 - 33",
			locale: "en",
			input:  en4,
			params: map[string]string{"times": "33"},
			want:   "This is the 33rd free spin.",
		},
		{
			name:   "EN - msg7 - 1,2,3,4,5",
			locale: "en",
			input:  en7,
			params: map[string]string{"a1": "1", "a2": "2", "a3": "3", "a4": "4", "a5": "5"},
			want:   "1 - 2 - 3% - €4.00 - £5.00",
		},
		{
			name:   "EN - msg7 - 1.1,2.2,3.3,4.4,5.5",
			locale: "en",
			input:  en7,
			params: map[string]string{"a1": "1.1", "a2": "2.2", "a3": "3.3", "a4": "4.4", "a5": "5.5"},
			want:   "1.1 - 2 - 3.3% - €4.40 - £5.50",
		},
		{
			name:   "IT - msg1",
			locale: "it",
			input:  it1,
			want:   it1,
		},
		{
			name:   "IT - msg2 - 0",
			locale: "it",
			input:  it2,
			params: map[string]string{"games": "0"},
			want:   "Non hai vinto giri gratis.",
		},
		{
			name:   "IT - msg2 - 1",
			locale: "it-IT",
			input:  it2,
			params: map[string]string{"games": "1"},
			want:   "Hai vinto 1 giro gratuito.",
		},
		{
			name:   "IT - msg2 - 2",
			locale: "it",
			input:  it2,
			params: map[string]string{"games": "2"},
			want:   "Hai vinto 2 giri gratis.",
		},
		{
			name:   "IT - msg3 - 0",
			locale: "it-it",
			input:  it3,
			params: map[string]string{"games": "0"},
			want:   "Hai vinto 0 giri gratis.",
		},
		{
			name:   "IT - msg3 - 1",
			locale: "it",
			input:  it3,
			params: map[string]string{"games": "1"},
			want:   "Hai vinto 1 giro gratuito.",
		},
		{
			name:   "IT - msg3 - 8",
			locale: "it",
			input:  it3,
			params: map[string]string{"games": "8"},
			want:   "Hai vinto 8 giri gratis.",
		},
		{
			name:   "IT - msg4 - 1",
			locale: "it",
			input:  it4,
			params: map[string]string{"times": "1"},
			want:   "Questa è la 1° volta.",
		},
		{
			name:   "IT - msg4 - 2",
			locale: "it",
			input:  it4,
			params: map[string]string{"times": "2"},
			want:   "Questa è la 2° volta.",
		},
		{
			name:   "IT - msg4 - 5",
			locale: "it",
			input:  it4,
			params: map[string]string{"times": "5"},
			want:   "Questa è la 5° volta.",
		},
		{
			name:   "IT - msg4 - 11",
			locale: "it",
			input:  it4,
			params: map[string]string{"times": "11"},
			want:   "Questa è l'11° volta.",
		},
		{
			name:   "IT - msg4 - 80",
			locale: "it",
			input:  it4,
			params: map[string]string{"times": "80"},
			want:   "Questa è l'80° volta.",
		},
		{
			name:   "IT - msg7 - 1,2,3,4,5",
			locale: "it",
			input:  it7,
			params: map[string]string{"a1": "1", "a2": "2", "a3": "3", "a4": "4", "a5": "5"},
			want:   "1 - 2 - 3% - 4,00 € - 5,00 £",
		},
		{
			name:   "IT - msg7 - 1.1,2.2,3.3,4.4,5.5",
			locale: "it",
			input:  it7,
			params: map[string]string{"a1": "1.1", "a2": "2.2", "a3": "3.3", "a4": "4.4", "a5": "5.5"},
			want:   "1,1 - 2 - 3,3% - 4,40 € - 5,50 £",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := Localize(tc.locale, tc.input, tc.params)
			assert.Equal(t, tc.want, got)
		})
	}
}
