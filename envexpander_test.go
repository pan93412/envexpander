package envexpander_test

import (
	"maps"
	"testing"

	"github.com/pan93412/envexpander/v3"
)

var testmap = []struct {
	Name     string
	Raw      map[string]string
	Resolved map[string]string
}{
	{
		Name: "simple",
		Raw: map[string]string{
			"A": "C${B}",
			"B": "1234",
		},
		Resolved: map[string]string{
			"A": "C1234",
			"B": "1234",
		},
	},
	{
		Name: "simple-2",
		Raw: map[string]string{
			"A": "1234",
			"B": "C${A}",
		},
		Resolved: map[string]string{
			"A": "1234",
			"B": "C1234",
		},
	},
	// FIXME: Currently we don't handle cyclic reference
	// well – we should break the cycle instead of expanding
	// in a random order.
	// {
	// 	Name: "cycle",
	// 	Raw: map[string]string{
	// 		"A": "C${B}",
	// 		"B": "C${A}",
	// 	},
	// 	Resolved: map[string]string{
	// 		"A": "C",
	// 		"B": "C",
	// 	},
	// },
	// {
	// 	Name: "cycle-2",
	// 	Raw: map[string]string{
	// 		"A": "C${B}",
	// 		"B": "C${A}",
	// 		"C": "C${B}",
	// 	},
	// 	Resolved: map[string]string{
	// 		"A": "C",
	// 		"B": "C",
	// 		"C": "C",
	// 	},
	// },
	{
		Name: "complex",
		Raw: map[string]string{
			"A": "1",
			"B": "${A}2",
			"C": "${A}${B}3",
			"D": "${B}${C}4",
			"E": "${A}${C}${D}5",
		},
		Resolved: map[string]string{
			"A": "1",
			"B": "12",
			"C": "1123",
			"D": "1211234",
			"E": "1112312112345",
		},
	},
	{
		Name: "very-complex",
		Raw: map[string]string{
			"A": "1",
			"B": "${A}2",
			"C": "${A}${B}3",
			"D": "${B}${C}4",
			"E": "${A}${C}${D}5",
			"F": "${A}${B}${C}${D}${E}${A}${A}${A}${B}${C}6",
		},
		Resolved: map[string]string{
			"A": "1",
			"B": "12",
			"C": "1123",
			"D": "1211234",
			"E": "1112312112345",
			"F": "1121123121123411123121123451111211236",
		},
	},
	{
		Name: "unknown-reference",
		Raw: map[string]string{
			"A": "C${B}",
		},
		Resolved: map[string]string{
			"A": "C${B}",
		},
	},
	{
		Name: "unknown-reference-2",
		Raw: map[string]string{
			"A": "C${D}",
			"B": "1234${A}",
		},
		Resolved: map[string]string{
			"A": "C${D}",
			"B": "1234C${D}",
		},
	},
	{
		Name: "realcase-1",
		Raw: map[string]string{
			"CONTACT_MAIL":      "foo@bar.tld",
			"LISTEN_HOST":       "http://0.0.0.0:${PORT}",
			"DATABASE_URI":      "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable",
			"PORT":              "8080",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_HOST":     "internal.postgres.host",
			"POSTGRES_PORT":     "1145",
			"POSTGRES_DB":       "testdb",
		},
		Resolved: map[string]string{
			"CONTACT_MAIL":      "foo@bar.tld",
			"LISTEN_HOST":       "http://0.0.0.0:8080",
			"DATABASE_URI":      "postgres://postgres:postgres@internal.postgres.host:1145/testdb?sslmode=disable",
			"PORT":              "8080",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_HOST":     "internal.postgres.host",
			"POSTGRES_PORT":     "1145",
			"POSTGRES_DB":       "testdb",
		},
	},
	{
		Name: "simple-3",
		Raw: map[string]string{
			"A": "${B}${B}${B}${B}${B}${B}",
			"B": "114514",
		},
		Resolved: map[string]string{
			"A": "114514114514114514114514114514114514",
			"B": "114514",
		},
	},
	{
		Name: "selfreference",
		Raw: map[string]string{
			"A": "${A}",
		},
		Resolved: map[string]string{
			"A": "${A}",
		},
	},
}

var reftestmap = []struct {
	Name string
	Raw  string
	Ref  []string
}{
	{
		Name: "simple",
		Raw:  "C${B}",
		Ref:  []string{"B"},
	},
	{
		Name: "simple-2",
		Raw:  "C${B}D",
		Ref:  []string{"B"},
	},
	{
		Name: "simple-3",
		Raw:  "C${B}D${E}",
		Ref:  []string{"B", "E"},
	},
	{
		Name: "simple-4",
		Raw:  "C${B}D${E}F",
		Ref:  []string{"B", "E"},
	},
	{
		Name: "invalid-syntax-1",
		Raw:  "C${B",
		Ref:  []string{},
	},
	{
		Name: "invalid-syntax-2",
		Raw:  "C$B}",
		Ref:  []string{},
	},
	{
		Name: "invalid-syntax-3",
		Raw:  "C${B}}}",
		Ref:  []string{"B"},
	},
	{
		Name: "invalid-syntax-4",
		Raw:  "C$B${C}",
		Ref:  []string{"C"},
	},
	{
		Name: "invalid-syntax-5",
		Raw:  "C${B${C}}",
		Ref:  []string{"B${C"}, // FIXME: ??
	},
	{
		Name: "escape",
		Raw:  "C\\${B}${C}",
		Ref:  []string{"C"},
	},
}

func TestEnvExpand(t *testing.T) {
	t.Parallel()

	for _, test := range testmap {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			resolved := envexpander.Expand(test.Raw)
			if !maps.Equal(test.Resolved, resolved) {
				t.Errorf("[%s] expected: %v, got: %v", test.Name, test.Resolved, resolved)
			}
		})
	}
}

func TestEnvExpandOne(t *testing.T) {
	t.Parallel()

	testmap := []struct {
		Name     string
		Value    string
		Variable map[string]string
		Resolved string
	}{
		{
			Name:  "simple",
			Value: "C${B}",
			Variable: map[string]string{
				"B": "1234",
			},
			Resolved: "C1234",
		},
		{
			Name:  "simple-2",
			Value: "1234",
			Variable: map[string]string{
				"A": "1234",
			},
			Resolved: "1234",
		},
		{
			Name:  "complex",
			Value: "1${A}2${B}3${C}4${D}5${E}",
			Variable: map[string]string{
				"A": "1",
				"B": "2",
				"C": "3",
				"D": "4",
				"E": "5",
			},
			Resolved: "1122334455",
		},
	}

	for _, test := range testmap {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			resolved := envexpander.ExpandOne(test.Value, test.Variable)
			if test.Resolved != resolved {
				t.Errorf("[%s] expected: %v, got: %v", test.Name, test.Resolved, resolved)
			}
		})
	}
}

func BenchmarkEnvExpand_V3(b *testing.B) {
	for _, test := range testmap {
		b.Run(test.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = envexpander.Expand(test.Raw)
			}
		})
	}
}

func TestRef(t *testing.T) {
	t.Parallel()

	for _, test := range reftestmap {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			ref := envexpander.Refs(test.Raw)
			for i, r := range ref {
				if test.Ref[i] != r.Variable(test.Raw) {
					t.Errorf("[%s] expected: %v, got: %v", test.Name, test.Ref[i], r.Variable(test.Raw))
				}
			}
		})
	}
}

func BenchmarkRefV3(b *testing.B) {
	for _, test := range reftestmap {
		b.Run(test.Name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = envexpander.Refs(test.Raw)
			}
		})
	}
}
