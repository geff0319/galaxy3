package bridge

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ge-fei-fan/gefflog"
	"github.com/geff0319/galaxy3/bridge/ytdlp"
	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

const ProcessTable = `CREATE TABLE "process" (
  "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
  "pid" text,
  "url" text,
  "params" text,
  "info" text,
  "progress" text,
  "output" text,
  "biliMeta" text,
  "is_delete" integer DEFAULT 0,
  "create_time" DATE DEFAULT CURRENT_TIMESTAMP
);`

const ProcessInsert = `INSERT INTO process ("url", "params", "info", "progress", "output", "biliMeta") VALUES (?, ?, ?, ?, ?, ?);`
const ProcessUpdate = `UPDATE process SET "pid" = ?,"info" = ?, "progress" = ?, "output" = ?, "biliMeta" = ? WHERE "id" = ?;`
const ProcessUpdateProgress = `UPDATE process SET "progress" = ? WHERE "id" = ?;`
const ProcessDelete = `UPDATE process SET "is_delete" = ? WHERE "id" = ?;`
const ProcessAll = `SELECT id,pid,url, params,info,progress,output FROM process where is_delete = 0 ORDER BY create_time DESC`

type SqliteService struct {
	DBFile string
	conn   *sql.DB
}

var SqliteS *SqliteService

func SqliteNew(config string) *SqliteService {
	SqliteS = &SqliteService{
		DBFile: config,
	}
	return SqliteS
}

// OnShutdown is called when the app is shutting down
// You can use this to clean up any resources you have allocated
func (s *SqliteService) OnShutdown() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// Name returns the name of the plugin.
// You should use the go module format e.g. github.com/myuser/myplugin
func (s *SqliteService) Name() string {
	return "github.com/wailsapp/wails/v3/plugins/sqlite"
}

// OnStartup is called when the app is starting up. You can use this to
// initialise any resources you need.
func (s *SqliteService) OnStartup(ctx context.Context, options application.ServiceOptions) error {
	if s.DBFile == "" {
		gefflog.Err(`no database file specified. Please set DBFile in the config to either a filename or use ":memory:" to use an in-memory database`)
		return errors.New(`no database file specified. Please set DBFile in the config to either a filename or use ":memory:" to use an in-memory database`)
	}
	gefflog.Info("连接数据库·····")
	if !ytdlp.IsFileExist(s.DBFile) {
		err := s.Init()
		if err != nil {
			gefflog.Err(err.Error())
			return err
		}
	} else {
		_, err := s.Open(s.DBFile)
		if err != nil {
			gefflog.Err(err.Error())
			return err
		}
	}

	return nil
}

func (s *SqliteService) Open(dbPath string) (string, error) {
	var err error
	s.conn, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return "", err
	}
	return "Database connection opened", nil
}
func (s *SqliteService) Init() error {
	_, err := s.Open(s.DBFile)
	if err != nil {
		return err
	}
	_, err = s.Execute(ProcessTable)
	if err != nil {
		return err
	}

	return nil
}
func (s *SqliteService) Execute(query string, args ...any) (sql.Result, error) {
	if s.conn == nil {
		return nil, errors.New("no open database connection")
	}

	return s.conn.Exec(query, args...)
}

func (s *SqliteService) Select(query string, args ...any) ([]map[string]any, error) {
	if s.conn == nil {
		return nil, errors.New("no open database connection")
	}

	rows, err := s.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	var results []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))

		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, err
		}

		row := make(map[string]any, len(columns))
		for i, column := range columns {
			row[column] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}
func (s *SqliteService) Close() (string, error) {
	if s.conn == nil {
		return "", errors.New("no open database connection")
	}

	err := s.conn.Close()
	if err != nil {
		return "", err
	}
	s.conn = nil
	return "Database connection closed", nil
}
