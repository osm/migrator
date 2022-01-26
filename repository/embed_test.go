package repository

import (
	"embed"
	"reflect"
	"testing"
)

//go:embed example/*
var embedFS embed.FS

// fileRepositoryTests contains a table of test cases
var embedFileRepositoryTests = []struct {
	name string
	dir  string
	err  string
	res  map[int]string
}{
	{"valid", "example", "", map[int]string{1: "CREATE TABLE migration (version text NOT NULL PRIMARY KEY);\n"}},
	{"no_migrations", ".", "no migrations found in .", nil},
	{"invalid_repository", "/this/dir/should/not/exist", "open /this/dir/should/not/exist: file does not exist", nil},
}

// TestFilesRepository iterates over all entries in the fileRepositoryTests table and executes the tests
func TestEmbedFilesRepository(t *testing.T) {
	for _, tc := range embedFileRepositoryTests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize a new embed file repository
			r := FromEmbedded(embedFS, tc.dir)

			// Load the repository
			m, err := r.Load()

			// There should be an error, let's make sure that either err is not nil or that the error messages match
			if tc.err != "" && (err == nil || tc.err != err.Error()) {
				t.Errorf("%s: an error was expected", tc.name)
				t.Logf("error: %#v", err)
				t.Logf("expected: %#v", tc.err)

			}

			// Make sure that the results are equal
			if tc.res != nil && !reflect.DeepEqual(tc.res, m) {
				t.Errorf("%s: repository did not match the expected value", tc.name)
				t.Logf("output: %#v", m)
				t.Logf("expected: %#v", tc.res)
			}
		})
	}
}
