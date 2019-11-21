package migrator

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/osm/migrator/repository"
)

// testDb is the test database
const testDb string = "./test.db"

// compareVersion compares the current version of the database
func compareVersion(t *testing.T, db *sql.DB, version int) {
	// Get current version from database
	var v int
	err := db.QueryRow("SELECT version FROM migration ORDER BY cast(version AS int) DESC").Scan(&v)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Make sure it is what we expect
	if v != version {
		t.Errorf("should have version %d, got version %d", version, v)
	}
}

func TestMigrator(t *testing.T) {
	// Remove the test db for each run
	os.Remove(testDb)

	// Initialize a new test db
	db, err := sql.Open("sqlite3", testDb)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Create a new mem repo
	repo := repository.FromMemory(map[int]string{
		1: "CREATE TABLE migration (version text NOT NULL PRIMARY KEY);\n",
		2: "CREATE TABLE foo (version text NOT NULL PRIMARY KEY);\n",
		3: "INSERT INTO foo VALUES(123);\n",
	})

	// Migrate to version 1
	err = ToVersion(db, repo, 1)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify version
	compareVersion(t, db, 1)

	// Migrate to latest version
	err = ToLatest(db, repo)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify version
	compareVersion(t, db, 3)
}

func TestMigratorWithManyMigrations(t *testing.T) {
	// Remove the test db for each run
	os.Remove(testDb)

	// Initialize a new test db
	db, err := sql.Open("sqlite3", testDb)
	if err != nil {
		t.Errorf(err.Error())
	}

	// Create a new mem repo
	repo := repository.FromMemory(map[int]string{
		1:  "CREATE TABLE migration (version text NOT NULL PRIMARY KEY);\n",
		2:  "CREATE TABLE foo_1 (version text NOT NULL PRIMARY KEY);\n",
		3:  "CREATE TABLE foo_2 (version text NOT NULL PRIMARY KEY);\n",
		4:  "CREATE TABLE foo_3 (version text NOT NULL PRIMARY KEY);\n",
		5:  "CREATE TABLE foo_4 (version text NOT NULL PRIMARY KEY);\n",
		6:  "CREATE TABLE foo_5 (version text NOT NULL PRIMARY KEY);\n",
		7:  "CREATE TABLE foo_6 (version text NOT NULL PRIMARY KEY);\n",
		8:  "CREATE TABLE foo_7 (version text NOT NULL PRIMARY KEY);\n",
		9:  "CREATE TABLE foo_8 (version text NOT NULL PRIMARY KEY);\n",
		10: "CREATE TABLE foo_9 (version text NOT NULL PRIMARY KEY);\n",
	})

	// Migrate to version 1
	err = ToVersion(db, repo, 1)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify version
	compareVersion(t, db, 1)

	// Migrate to latest version
	err = ToLatest(db, repo)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Verify version
	compareVersion(t, db, 10)

	// Run the migration again, we don't expect any action to be taken now
	err = ToLatest(db, repo)
	if err != nil {
		t.Errorf("%v", err)
	}
}
