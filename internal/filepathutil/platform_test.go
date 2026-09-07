package filepathutil_test

import (
	"testing"

	"github.com/maramkhaledn/protolint/internal/filepathutil"
)

func TestToUnixPath(t *testing.T) {
	for _, test := range []struct {
		name                        string
		inputCrossPlatformPath      string
		inputIsWindowsPathSeparator bool
		want                        string
	}{
		{
			name:                   "returns an unix path as is",
			inputCrossPlatformPath: "proto/generated/search.proto",
			want:                   "proto/generated/search.proto",
		},
		{
			name:                   "returns a separator-free path as is",
			inputCrossPlatformPath: "search.proto",
			want:                   "search.proto",
		},
		{
			name:                   "returns an empty path as is",
			inputCrossPlatformPath: "",
			want:                   "",
		},
		{
			name:                        "converts all the windows separators to unix ones",
			inputCrossPlatformPath:      `proto\generated\search.proto`,
			inputIsWindowsPathSeparator: true,
			want:                        "proto/generated/search.proto",
		},
		{
			name:                        "returns a separator-free path as is on windows",
			inputCrossPlatformPath:      "search.proto",
			inputIsWindowsPathSeparator: true,
			want:                        "search.proto",
		},
		{
			name:                        "keeps the unix separators already contained in a windows path",
			inputCrossPlatformPath:      `proto\generated/search.proto`,
			inputIsWindowsPathSeparator: true,
			want:                        "proto/generated/search.proto",
		},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			defer setOSPathSeparator(t, test.inputIsWindowsPathSeparator)()

			got := filepathutil.ToUnixPath(test.inputCrossPlatformPath)
			if got != test.want {
				t.Errorf("got %s, but want %s", got, test.want)
			}
		})
	}
}

// setOSPathSeparator overrides the package level separator and returns the func to restore it.
func setOSPathSeparator(t *testing.T, isWindows bool) func() {
	t.Helper()

	sep := '/'
	if isWindows {
		sep = '\\'
	}
	prev := filepathutil.OSPathSeparator
	filepathutil.OSPathSeparator = sep
	return func() {
		filepathutil.OSPathSeparator = prev
	}
}
