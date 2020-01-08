package util

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "github.com/yb7/alilog"
  "strings"
)

type (
  Pagination struct {
    PageNo   uint64 `json:"pageNo" desc:"第几页(1开始计数)"`
    PageSize uint64 `json:"pageSize" desc:"每页有几条"`
    RowsNo   uint64 `json:"rowsNo" desc:"总条数"`
    PagesNo  uint64 `json:"pagesNo" desc:"总页数"`
  }
)

func (p *Pagination) Start() uint64 {
  if p.PageNo > 0 && p.PageSize > 0 {
    page := p.PageNo
    perPage := p.PageSize
    return (page - 1) * perPage
  }
  return 0
}

func (p *Pagination) PageLimit(q squirrel.SelectBuilder) squirrel.SelectBuilder {
  if p.PageSize == 0 {
    p.PageSize = 10
  }
  //if p.PageNo > 0 && p.PageSize > 0 {
  //  return q.Offset(p.Start()).Limit(p.PageSize)
  //}
  if p.PageNo == 0 {
    p.PageNo = 1
  }
  return q.Offset(p.Start()).Limit(p.PageSize)
}

func (p *Pagination) BuildBy(totalRows uint64) {
  p.RowsNo = totalRows
  if p.PageSize == 0 {
    p.PageSize = 20
  }
  if totalRows%p.PageSize == 0 {
    p.PagesNo = totalRows / p.PageSize
  } else {
    p.PagesNo = totalRows/p.PageSize + 1
  }
}

func PrintQuery(log *alilog.SLog, query squirrel.SelectBuilder) {
  str, args, err := query.ToSql()
  if err != nil {
    log.Error(err)
  }
  var msg []string
  //msg = append(msg, fmt.Sprintf("START >> SQL in [%s.%s]", log.["file"], log["method"]))
  msg = append(msg, str)
  if len(args) > 0 {
    msg = append(msg, "args:")
    for i, v := range args {
      msg = append(msg, fmt.Sprintf("%d => %v", i, v))
    }
  }
  //msg = append(msg, "END  <<")
  log.Debugf(strings.Join(msg, "\n"))
}

func PrintUpdate(log *alilog.SLog, query squirrel.UpdateBuilder) {
  str, args, err := query.ToSql()
  if err != nil {
    log.Error(err)
  }
  var msg []string
  //msg = append(msg, fmt.Sprintf("START >> SQL in [%s.%s]", log["file"], log["method"]))
  msg = append(msg, str)
  if len(args) > 0 {
    msg = append(msg, "args:")
    for i, v := range args {
      msg = append(msg, fmt.Sprintf("%d => %v", i, v))
    }
  }
  //msg = append(msg, "END  <<")
  log.Debugf(strings.Join(msg, "\n"))
}

func PrintInsert(log *alilog.SLog, query squirrel.InsertBuilder) {
  str, args, err := query.ToSql()
  if err != nil {
    log.Error(err)
  }
  var msg []string
  //msg = append(msg, fmt.Sprintf("START >> SQL in [%s.%s]", log["file"], log["method"]))
  msg = append(msg, str)
  if len(args) > 0 {
    msg = append(msg, "args:")
    for i, v := range args {
      msg = append(msg, fmt.Sprintf("%d => %v", i, v))
    }
  }
  //msg = append(msg, "END  <<")
  log.Debugf(strings.Join(msg, "\n"))
}

func PrintDelete(log *alilog.SLog, query squirrel.DeleteBuilder) {
  str, args, err := query.ToSql()
  if err != nil {
    log.Error(err)
  }
  var msg []string
  //msg = append(msg, fmt.Sprintf("START >> SQL in [%s.%s]", log["file"], log["method"]))
  msg = append(msg, str)
  if len(args) > 0 {
    msg = append(msg, "args:")
    for i, v := range args {
      msg = append(msg, fmt.Sprintf("%d => %v", i, v))
    }
  }
  //msg = append(msg, "END  <<")
  log.Debugf(strings.Join(msg, "\n"))
}

func IsNoRowsError(err error) bool {
  return "sql: no rows in result set" == err.Error()
}
