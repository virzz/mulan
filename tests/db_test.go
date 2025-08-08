package tests

import (
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/virzz/mulan/db"
)

func TestNewDB_MySQL(t *testing.T) {
	cfg := db.Config{
		DSN:  "mysql://root:@127.0.0.1/mysql",
		Pass: "123456",
		Args: map[string]string{
			"parseTime": "true",
		},
	}
	dialector := mysql.New(mysql.Config{
		DSN:                    cfg.String(),
		DefaultStringSize:      255,
		DontSupportRenameIndex: true,
	})
	_, err := db.New(dialector, nil, &gorm.Config{Logger: db.NewLogger(false)})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDB_PgSQL(t *testing.T) {
	cfg := db.Config{
		DSN: "postgres://mozhu:@127.0.0.1/postgres",
		Args: map[string]string{
			"sslmode": "disable",
		},
	}
	dialector := postgres.New(postgres.Config{
		DSN:                  cfg.String(),
		PreferSimpleProtocol: true,
	})
	_, err := db.New(dialector, nil, &gorm.Config{Logger: db.NewLogger(true)})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDB_SQLiteMemory(t *testing.T) {
	cfg := db.Config{
		DSN: "sqlite3://:memory:",
		Args: map[string]string{
			"cache": "shared",
		},
	}
	dsn := cfg.String()
	t.Log(dsn)
	dialector := sqlite.Open(dsn)
	_, err := db.New(dialector, nil, &gorm.Config{Logger: db.NewLogger(true)})
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDB_SQLiteFile(t *testing.T) {
	dsns := []string{
		"sqlite3://:memory:",
		"sqlite3://test.db",
		"sqlite3://./test.db",
		"sqlite3:///test.db",
	}
	defer func() { os.Remove("./test.db") }()
	for _, dsn := range dsns {
		t.Run(dsn, func(t *testing.T) {
			cfg := db.Config{DSN: dsn}
			_, err := db.New(sqlite.Open(cfg.String()), nil)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
