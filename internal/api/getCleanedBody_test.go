package api

import "testing"

func TestGetCleanedBody(t *testing.T) {
	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	tests := []struct {
		name     string
		body     string
		expected string
	}{
		{
			name:     "no profanity",
			body:     "I had something interesting for breakfast",
			expected: "I had something interesting for breakfast",
		},
		{
			name:     "profanity with lower case",
			body:     "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			name:     "profanity with upper case",
			body:     "I hear Mastodon is better than Chirpy. Sharbert I need to migrate",
			expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			name:     "profanity with camel case",
			body:     "I hear Mastodon is better than Chirpy. sharBert I need to migrate",
			expected: "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			name:     "profanity with lower case and upper case",
			body:     "I really need a kerfuffle to go to bed sooner, Fornax !",
			expected: "I really need a **** to go to bed sooner, **** !",
		},
		{
			name:     "profanity with punctuation attached (should NOT be replaced)",
			body:     "That's a sharbert! situation",
			expected: "That's a sharbert! situation",
		},
		{
			name:     "multiple profane words in one body",
			body:     "kerfuffle sharbert fornax",
			expected: "**** **** ****",
		},
		{
			name:     "empty body",
			body:     "",
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := getCleanedBody(test.body, badWords)
			if result != test.expected {
				t.Errorf("failed test case '%s':\n expected %q,\n got %q", test.name, test.expected, result)
			}
		})
	}

}
