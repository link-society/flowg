package models

import (
	"testing"
)

func TestConvert_FilterDSL_to_ExprLang(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{`foo="bar"`, `foo == "bar"`},
		{`foo!="bar"`, `foo != "bar"`},
		{`foo>10`, `foo > 10`},
		{`foo<20`, `foo < 20`},
		{`foo>=30`, `foo >= 30`},
		{`foo<=40`, `foo <= 40`},
		{`foo="bar" and baz>10`, `foo == "bar" and baz > 10`},
		{`foo="bar" or baz<20`, `foo == "bar" or baz < 20`},
		{`(foo="bar" and baz>10) or qux<=5`, `( foo == "bar" and baz > 10 ) or qux <= 5`},
		{`foo = "hello=world"`, `foo == "hello=world"`},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			got, err := convertFilterdslToExprlang(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != tc.expected {
				t.Fatalf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}
