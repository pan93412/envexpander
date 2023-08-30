package envexpander_test

import (
	"fmt"
	"testing"

	"github.com/pan93412/envexpander"
	"github.com/stretchr/testify/assert"
)

func TestReplacerIntegrate(t *testing.T) {
	testmap := []struct {
		Value     string
		Variables map[string]string
		Result    string
	}{
		{
			Value:     "A=${B}",
			Variables: map[string]string{"B": "1234"},
			Result:    "A=1234",
		},
		{
			Value:     "A=${B}${C}",
			Variables: map[string]string{"B": "1234", "C": "5678"},
			Result:    "A=12345678",
		},
		{
			Value:     "A=1234",
			Variables: map[string]string{"B": "1234", "C": "5678"},
			Result:    "A=1234",
		},
		{
			Value:     "A=${B}${C}",
			Variables: map[string]string{"B": "1234"},
			Result:    "A=1234${C}",
		},
		{
			Value:     "A=${B}${C}",
			Variables: map[string]string{"C": "5678"},
			Result:    "A=${B}5678",
		},
	}

	for _, v := range testmap {
		v := v
		t.Run(fmt.Sprintf("%s_to_%s", v.Value, v.Result), func(t *testing.T) {
			t.Parallel()

			r := envexpander.Replacer{
				Value:     v.Value,
				Variables: make(map[string]string, len(v.Variables)),
			}

			for key, val := range v.Variables {
				val := val

				r.Variables[key] = val
			}

			assert.Equal(t, v.Result, r.Integrate())
		})
	}
}

func BenchmarkReplacerIntegrate(b *testing.B) {
	r := envexpander.Replacer{
		Value: "A=${B}${C}${D}${D}",
		Variables: map[string]string{
			"B": "1234",
			"C": "5678",
			"D": "9012",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Integrate()
	}
}

func BenchmarkReplacerIntegrateWithCache(b *testing.B) {
	cvp := envexpander.NewCachedVariablePos()
	r := envexpander.Replacer{
		Value: "A=${B}${C}${D}${D}",
		Variables: map[string]string{
			"B": "1234",
			"C": "5678",
			"D": "9012",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.IntegrateWithCache(cvp)
	}
}
