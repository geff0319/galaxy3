package bridge

// Process
const ProcessTable = `CREATE TABLE IF NOT EXISTS "process" (
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
const ProcessAll = `SELECT id,pid,url, params,info,progress,output,biliMeta FROM process where is_delete = 0 ORDER BY create_time DESC`

// yaml config
const ConfigTable = `CREATE TABLE IF NOT EXISTS "config" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "config_name" TEXT,
  "config_value" TEXT
);`

const ConfigInsert = `INSERT INTO config ("config_name", "config_value") VALUES (?, ?);`
const ConfigUpdate = `UPDATE config SET "config_value"= ? WHERE "config_name" = ?;`
const ConfigAll = `SELECT id, config_name,config_value FROM config`
const ConfigUser = `SELECT config_value FROM config WHERE config_name = 'user' limit 1;`
const ConfigYtdlp = `SELECT config_value FROM config WHERE config_name = 'ytdlp' limit 1;`
