package users

import (
	"testing"

	db "github.com/Squwid/bytegolf/database"
)

func TestCreateTable(t *testing.T) {
	// t.Skip("Should only be run by hand")
	stmt, err := db.DB.Prepare(`
	CREATE TABLE IF NOT EXISTS "user" (
		id INT PRIMARY KEY,
		username VARCHAR(150) NULL,
		picture_uri VARCHAR(300) NULL,
		github_uri VARCHAR(300) NULL,
		name VARCHAR(200) NULL,
		total_submissions INT NOT NULL,
		total_correct INT NOT NULL,
		total_bytes INT NOT NULL
	);`)
	if err != nil {
		t.Errorf("could not prepare: %v\n", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		t.Errorf("could not query: %v\n", err)
		return
	}
}
