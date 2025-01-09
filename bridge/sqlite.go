package bridge

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ge-fei-fan/gefflog"
	"github.com/wailsapp/wails/v3/pkg/application"
	_ "modernc.org/sqlite"
)

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
	err := s.Init()
	if err != nil {
		gefflog.Err(err.Error())
		return err
	}
	// 读取user配置
	Config.Unmarshal()
	gefflog.ChangeLogger(0, Config.LogPath)

	// ytdlp配置
	InitYtDlpConfig()

	//替换视频状态
	res, err := s.Select(`SELECT id FROM process WHERE JSON_EXTRACT(progress, '$.process_status') != 2 AND is_delete = 0;`)
	if err != nil {
		gefflog.Err("查询视频状态error: " + err.Error())
	}
	for _, i := range res {
		v, ok := i["id"].(int64)
		if ok {
			_, err := s.Execute(`UPDATE process SET progress = '{"process_status":3,"percentage":"","speed":0,"eta":0}' WHERE id = ?`, v)
			if err != nil {
				gefflog.Err("替换视频状态error: " + err.Error())
			}
		}
	}
	//失败视频放进列表
	res, err = SqliteS.Select(`SELECT id,pid,url, params,info,progress,output,biliMeta FROM process where is_delete = 0 AND json_extract(progress, '$.process_status') != 2 ORDER BY create_time DESC`)
	if err != nil {
		gefflog.Err("获取失败视频error: " + err.Error())
	}
	//解析失败时pid不存在
	for _, r := range res {
		var p Process
		p.Unmarshal(r)
		if p.Pid == "" {
			YdpConfig.Mq.Publish(&p)
		} else {
			YdpConfig.Mdb.Set(&p)
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
	_, err = s.conn.Exec("PRAGMA journal_mode=WAL;")
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
		return errors.New("create process err: " + err.Error())
	}
	_, err = s.Execute(ConfigTable)
	if err != nil {
		return errors.New("create config err: " + err.Error())
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
