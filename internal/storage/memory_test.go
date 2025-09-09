package storage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateIndexAndSelectWhere(t *testing.T) {
	dataDir := t.TempDir()
	db := NewDatabase(dataDir)

	if msg := db.CreateTable("users", []string{"id", "name", "email"}); !strings.Contains(msg, "created") {
		t.Fatalf("create table failed: %s", msg)
	}

	_ = db.Insert("users", []string{"1", "Hareesh", "hareesh@example.com"})
	_ = db.Insert("users", []string{"2", "Bhittam", "bhittam@example.com"})

	if msg := db.CreateIndex("users", "email"); !strings.Contains(msg, "Index created") {
		t.Fatalf("create index failed: %s", msg)
	}

	out := db.SelectWhere("users", "email", "bhittam@example.com")
	if !strings.Contains(out, "2 | Bhittam | bhittam@example.com") {
		t.Fatalf("expected indexed lookup to find Bhittam, got:\n%s", out)
	}

	// Restart database to ensure indexes are rebuilt from metadata
	db = NewDatabase(dataDir)
	out = db.SelectWhere("users", "email", "hareesh@example.com")
	if !strings.Contains(out, "1 | Hareesh | hareesh@example.com") {
		t.Fatalf("expected rebuilt index/fullscan to find Hareesh, got:\n%s", out)
	}
}

func TestIndexMaintenanceOnUpdateDelete(t *testing.T) {
	dataDir := t.TempDir()
	db := NewDatabase(dataDir)

	_ = db.CreateTable("t", []string{"k", "v"})
	_ = db.Insert("t", []string{"a", "1"})
	_ = db.Insert("t", []string{"b", "2"})
	_ = db.CreateIndex("t", "k")

	// Update row 1 from b->c
	msg := db.Update("t", 1, []string{"c", "2"})
	if !strings.Contains(msg, "updated") {
		t.Fatalf("update failed: %s", msg)
	}
	// old key should not be found
	if out := db.SelectWhere("t", "k", "b"); !strings.Contains(out, "(no rows)") {
		t.Fatalf("expected no rows for key b after update, got:\n%s", out)
	}
	// new key should be found
	if out := db.SelectWhere("t", "k", "c"); !strings.Contains(out, "c | 2") {
		t.Fatalf("expected row for key c after update, got:\n%s", out)
	}

	// Delete first row and ensure index rebuild removes it
	msg = db.Delete("t", 0)
	if !strings.Contains(msg, "deleted") {
		t.Fatalf("delete failed: %s", msg)
	}
	if out := db.SelectWhere("t", "k", "a"); !strings.Contains(out, "(no rows)") {
		t.Fatalf("expected no rows for key a after delete, got:\n%s", out)
	}

	// Sanity: tables persisted with index metadata
	if _, err := os.Stat(filepath.Join(dataDir, "t.harudb")); err != nil {
		t.Fatalf("expected persisted table file: %v", err)
	}
}
