package util

import (
	"fmt"
	"testing"
)

func TestGetLinkFromDeleteDescription(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
		isErr    bool
	}{
		{
			name:     "valid link",
			input:    "[testing](https://discord.com/channels/530946869023604746/530949686962421770/755612707146825839)",
			expected: "https://discord.com/channels/530946869023604746/530949686962421770/755612707146825839",
		},
		{
			name:  "invalid link",
			input: "https://google.com",
			isErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := GetLinkFromDeleteDescription(tc.input)
			if err == nil && tc.isErr {
				t.Fatalf("expected err: %v", tc.name)
			}
			if got != tc.expected {
				t.Fatalf("expected: %v, got: %v", tc.expected, got)
			}
		})
	}
}

func TestParseFieldHeadCount(t *testing.T) {
	cases := []struct {
		name          string
		input         string
		expectedCount int
		expectedLimit int
	}{
		{
			name:          "no limit with count",
			input:         "✅ Accepted (17)",
			expectedCount: 17,
			expectedLimit: -1,
		},
		{
			name:          "limit",
			input:         "✅ Accepted (1/10)",
			expectedCount: 1,
			expectedLimit: 10,
		},
		{
			name:          "no limit without count",
			input:         "❔ Tentative",
			expectedCount: -1,
			expectedLimit: -1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotCount, gotLimit, err := ParseFieldHeadCount(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotCount != tc.expectedCount && gotLimit != tc.expectedLimit {
				t.Fatalf("expected count %d, limit %d: got %d, %d", tc.expectedCount, tc.expectedLimit, gotCount, gotLimit)
			}
		})
	}
}

func TestIncrementFieldCounter(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
		isErr    bool
	}{
		{
			name:     "no limit with count",
			input:    "✅ Accepted (100)",
			expected: "✅ Accepted (101)",
		},
		{
			name:     "limit",
			input:    "✅ Accepted (0/3)",
			expected: "✅ Accepted (1/3)",
		},
		{
			name:     "no limit without count",
			input:    "✅ Accepted",
			expected: "✅ Accepted (1)",
		},
		{
			name:  "limit exceeded",
			input: "✅ Accepted (3/3)",
			isErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := IncrementFieldCounter(tc.input)
			if tc.isErr && err == nil {
				fmt.Println(got)
				t.Fatalf("expected err: %v", err)
			}
			if got != tc.expected {
				t.Fatalf("expected: %s, got: %s", tc.expected, got)
			}
		})
	}
}

func TestDecrementFieldCounter(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
		isErr    bool
	}{
		{
			name:     "no limit with count",
			input:    "✅ Accepted (98)",
			expected: "✅ Accepted (97)",
		},
		{
			name:     "limit",
			input:    "✅ Accepted (1/3)",
			expected: "✅ Accepted (0/3)",
		},
		{
			name:     "no limit without count",
			input:    "✅ Accepted (1)",
			expected: "✅ Accepted",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DecrementFieldCounter(tc.input)
			if tc.isErr && err == nil {
				fmt.Println(got)
				t.Fatalf("expected err: %v", err)
			}
			if got != tc.expected {
				t.Fatalf("expected: %s, got: %s", tc.expected, got)
			}
		})
	}
}

func TestAddUserToField(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		userName string
		expected string
	}{
		{
			name:     "empty",
			value:    "-",
			userName: "foo",
			expected: "> foo",
		},
		{
			name:     "add user",
			value:    "> foo",
			userName: "bar",
			expected: "> foo\n> bar",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := AddUserToField(tc.value, tc.userName)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got != tc.expected {
				t.Fatalf("expected: %s, got: %s", tc.expected, got)
			}
		})
	}
}

func TestRemoveUserFromField(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		userName string
		expected string
	}{
		{
			name:     "empty",
			value:    "-",
			userName: "foo",
			expected: "-",
		},
		{
			name:     "remove user",
			value:    "> foo",
			userName: "foo",
			expected: "-",
		},
		{
			name:     "remove user with remaining list",
			value:    "> foo\n> bar\n> baz",
			userName: "bar",
			expected: "> foo\n> baz",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := RemoveUserFromField(tc.value, tc.userName)
			if err != nil {
				t.Fatalf("unexpected err: %v", err)
			}
			if got != tc.expected {
				t.Fatalf("expected: %s, got: %s", tc.expected, got)
			}
		})
	}
}
