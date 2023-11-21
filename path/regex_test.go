package path

import (
	"testing"
)

type argument struct {
	in  string
	out string
}

func TestGetPathRegex(t *testing.T) {
	paths := []argument{
		{"/path/to/example*", "example.txt"},
		{"/path/to/example*", "example1234123412.txt"},
		{"/path/to/example*", "example1$^&@1A.txt"},
		{"/path/to/example.txt", "example.txt"},
		{"/path/to/example", "example"},
		{"/path/to/(example){2}", "exampleexample"},
		{"/path/to/*.mp4", "example.mp4"},
	}

	for _, p := range paths {
		t.Run(p.in, func(t *testing.T) {
			res, err := GetPathRegex(p.in)
			if err != nil {
				t.Fatalf("Failed to extract regex from path. %s", err)
			}

			if !res.MatchString(p.out) {
				t.Fatalf("Failed to match example.txt to regex")
			}
		})
	}
}
