package service

import (
  "yx-api/util"
  "database/sql"
  "fmt"
  "github.com/Masterminds/squirrel"
  _ "github.com/go-sql-driver/mysql"
  jsoniter "github.com/json-iterator/go"
  "time"
)

var (
  db      *sql.DB
  dbCache squirrel.DBProxyBeginner
  dbLog   = util.AppLog.With("file", "service.db.go")
  json    = jsoniter.ConfigCompatibleWithStandardLibrary
  conf    = util.Settings{}
)

func OpenDB() squirrel.DBProxyBeginner {
  log := dbLog.With("func", "OpenDB")
  var err error
  conf = util.Get()
  log.Infof("Connecting to %s ... \n", conf.DBUrl)
  db, err = sql.Open(conf.DriverName, conf.DBUrl)
  if err != nil {
    log.Errorf("DB conection set up failed, %s\n", err.Error())
    panic(err)
  }
  db.SetMaxOpenConns(300)
  db.SetMaxIdleConns(0)
  db.SetConnMaxLifetime(10 * time.Second)
  err = db.Ping()
  if err != nil {
    log.Errorf("DB conection set up failed, %s\n", err.Error())
    panic(err)
  }
  log.Infof("DB conection set up successfully \n")
  dbCache = squirrel.NewStmtCacheProxy(db)
  return dbCache
}

func CloseDB() {
  db.Close()
}

type txDBProxyBeginnerWrapper struct {
  tx *sql.Tx
}

func (w txDBProxyBeginnerWrapper) Begin() (*sql.Tx, error) {
  return w.tx, nil
}
func (w txDBProxyBeginnerWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
  return w.tx.Exec(query, args...)
}

func (w txDBProxyBeginnerWrapper) Query(query string, args ...interface{}) (*sql.Rows, error) {
  return w.tx.Query(query, args...)
}
func (w txDBProxyBeginnerWrapper) QueryRow(query string, args ...interface{}) squirrel.RowScanner {
  return w.tx.QueryRow(query, args...)
}
func (w txDBProxyBeginnerWrapper) Prepare(query string) (*sql.Stmt, error) {
  return w.tx.Prepare(query)
}
func txToBeginner(tx *sql.Tx) squirrel.DBProxyBeginner {
  return txDBProxyBeginnerWrapper{tx: tx}
}

func inTx(f func(squirrel.DBProxyBeginner) error) error {
  tx, err := dbCache.Begin()

  if err != nil {
    fmt.Printf(err.Error())
    return err
  }
  defer tx.Rollback()

  dbProxyBeginner := txToBeginner(tx)
  if err := f(dbProxyBeginner); err != nil {
    return err
  }
  tx.Commit()
  return nil
}
