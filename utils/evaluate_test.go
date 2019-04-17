package utils

import "testing"

func TestEvaluate(t *testing.T) {
	items := []struct {
		str      string
		expected string
	}{
		{str: "test\\ntest", expected: "test\ntest"},
		{str: "test\ntest", expected: "test\ntest"},
		{str: "test\\\\ntest", expected: "test\\ntest"},
		{str: "\\n", expected: "\n"},
		{str: "\\\\n\\n", expected: "\\n\n"},
		{str: "\\n", expected: "\n"},
		{str: "\n", expected: "\n"},
		{str: "\\n\\thi", expected: "\n\thi"},
		{str: "\\n\t\\n\\t\n\n", expected: "\n\t\n\t\n\n"},
		{str: "\\\\n\\\\thallo\\\\t\\\\n", expected: "\\n\\thallo\\t\\n"},
	}
	for _, item := range items {
		t.Run(item.str, func(t *testing.T) {
			eval := Evaluate(item.str)
			if eval != item.expected {
				t.Fatalf("str '%s' should be evaluated to '%s' but was '%s'.", item.str, item.expected, eval)
			}
		})
	}
}
