package repository

import (
	"reflect"
	"testing"
)

// memoryRepositoryTests contains a table of test cases
var memoryRepositoryTests = []struct {
	name string
	data map[int]string
}{
	{"valid", map[int]string{1: "CREATE TABLE migration (version text NOT NULL PRIMARY KEY);\n"}},
}

// TestMemoryRepository iterates over all entries in the memoryRepositoryTests table and executes the tests
func TestMemoryRepository(t *testing.T) {
	for _, tc := range memoryRepositoryTests {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize a new memory based repository
			r := FromMemory(tc.data)

			// Load the repository
			m, _ := r.Load()

			// Make sure that the results are equal
			if !reflect.DeepEqual(tc.data, m) {
				t.Errorf("%s: repo did not match the expected value", tc.name)
				t.Logf("output: %#v", m)
				t.Logf("expected: %#v", tc.data)
			}
		})
	}
}
