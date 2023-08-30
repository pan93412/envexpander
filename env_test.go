package envexpander_test

import (
	"fmt"
	"testing"

	"github.com/pan93412/envexpander"
	"github.com/stretchr/testify/assert"
)

func TestFindVariableReferenceMap(t *testing.T) {
	cvp := envexpander.NewCachedVariablePos()
	testmap := map[string]map[string]struct{}{
		"A=${B}":                 {"B": struct{}{}},
		"A=${B}${C}":             {"B": struct{}{}, "C": struct{}{}},
		"A=${B}${C}${D}":         {"B": struct{}{}, "C": struct{}{}, "D": struct{}{}},
		"A=${B}${C}${D}${D}":     {"B": struct{}{}, "C": struct{}{}, "D": struct{}{}},
		"A=${B}${C}${C}${D}${D}": {"B": struct{}{}, "C": struct{}{}, "D": struct{}{}},
		"A=1234":                 {},
		"A=${B":                  {},
	}

	for k, v := range testmap {
		k := k
		v := v

		t.Run(k, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, envexpander.FindVariableReferenceMap(cvp, k), v)
		})
	}
}

func BenchmarkFindVariableReferenceMap(b *testing.B) {
	cvp := envexpander.NewCachedVariablePos()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		envexpander.FindVariableReferenceMap(cvp, "A=${B}${C}${D}${D}")
	}
}

func TestResolveEnvVariable(t *testing.T) {
	testmap := []struct {
		Raw      map[string]string
		Resolved map[string]string
	}{
		{
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
			Raw: map[string]string{
				"A": "1234",
				"B": "C${A}",
			},
			Resolved: map[string]string{
				"A": "1234",
				"B": "C1234",
			},
		},
		{
			Raw: map[string]string{
				"A": "C${B}",
				"B": "C${A}",
			},
			Resolved: map[string]string{
				"A": "C",
				"B": "C",
			},
		},
		{
			Raw: map[string]string{
				"A": "C${B}",
				"B": "C${A}",
				"C": "C${B}",
			},
			Resolved: map[string]string{
				"A": "C",
				"B": "C",
				"C": "CC",
			},
		},
		{
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
			Raw: map[string]string{
				"A": "${CCC}",
			},
			Resolved: map[string]string{
				"A": "${CCC}",
			},
		},
		{
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
			Raw: map[string]string{
				"A": "${A}",
			},
			Resolved: map[string]string{
				"A": "",
			},
		},
	}

	for _, v := range testmap {
		v := v
		t.Run(fmt.Sprintf("%v", v.Raw), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, v.Resolved, envexpander.ResolveEnvVariable(v.Raw))
		})
	}
}

// fixme: the iteration order of `Raw` is not fixed;
// therefore, the test execution time is not deterministic.
func BenchmarkResolveEnvVariable_Basic(b *testing.B) {
	test := map[string]string{
		"CONTACT_MAIL":      "foo@bar.tld",
		"LISTEN_HOST":       "http://0.0.0.0:${PORT}",
		"DATABASE_URI":      "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable",
		"PORT":              "8080",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_PASSWORD": "postgres",
		"POSTGRES_HOST":     "internal.postgres.host",
		"POSTGRES_PORT":     "1145",
		"POSTGRES_DB":       "testdb",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		envexpander.ResolveEnvVariable(test)
	}
}

// fixme: the iteration order of `Raw` is not fixed;
// therefore, the test execution time is not deterministic.
func BenchmarkResolveEnvVariable_Complex(b *testing.B) {
	test := map[string]string{
		"A": "1",
		"B": "${A}2",
		"C": "${A}${B}3",
		"D": "${B}${C}4",
		"E": "${A}${C}${D}5",
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		envexpander.ResolveEnvVariable(test)
	}
}
