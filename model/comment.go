package model

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "github.com/Masterminds/structable"
  "yx-api/util"
)

type (
  Comment struct {
    ID      int    `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT" desc:"ID"`
    Type    int    `stbl:"type" desc:"资源类型"`
    Sid     int    `stbl:"sid" desc:"资源ID"`
    Uid     int    `stbl:"uid" desc:"评论用户ID"`
    Rid     int    `stbl:"rid" desc:"回复用户ID"`
    CTime   int64  `stbl:"ctime" desc:"评论时间"`
    Content string `stbl:"content" desc:"评论内容"`
    Status  int    `stbl:"status" desc:"状态"`
    RNum    int    `stbl:"rnum" desc:"被回复数量"`
    Pid     int    `stbl:"pid" desc:"回复评论ID"`
    ZNum    int    `stbl:"znum" desc:"点赞数量"`
    rec     structable.Recorder
  }

  CommentDto struct {
    ID      int
    Content   string
    Status  int
    util.Pagination
  }

  CommentFilter struct {
    Uid      int   `json:"-" desc:"当前用户ID"`
    Sid      int   `json:"sid" desc:"资源ID"`
    Type     int   `json:"type" desc:"资源类型"`
    LastTime int64 `json:"lastTime" desc:"上拉加载时间"`
  }
)

func NewComment(db squirrel.DBProxyBeginner) *Comment {
  vo := new(Comment)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableComment, vo)
  return vo
}

func (vo *Comment) Load() error {
  return vo.rec.Load()
}

func (vo *Comment) Insert() error {
  return vo.rec.Insert()
}

func (vo *Comment) Update() error {
  return vo.rec.Update()
}

func GetCommentPagination(db squirrel.DBProxyBeginner, filter CommentDto) (*[]Comment, util.Pagination, error) {
  log := commentLog.With("func", "GetCommentPagination")
  result := make([]Comment, 0)
  count := uint64(0)
  column := NewComment(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableComment)
  query = query.Where("status = 1 and pid = 0")
  query = query.OrderBy("ctime desc")
  query = paramToComment(filter, query)
  query = filter.PageLimit(query)
  util.PrintQuery(log, query)

  rows, err := query.RunWith(db).Query()
  if err != nil {
    return &result, filter.Pagination, err
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return &result, filter.Pagination, err
    }
    result = append(result, *column)
  }

  query = squirrel.Select("count(1)").From(tableComment)
  query = paramToComment(filter, query)
  err = query.RunWith(db).QueryRow().Scan(&count)
  filter.BuildBy(count)
  return &result, filter.Pagination, err
}

func paramToComment(filter CommentDto, query squirrel.SelectBuilder) squirrel.SelectBuilder {
  if filter.ID > 0 {
    query = query.Where("id = ?", filter.ID)
  }
  if filter.Content != "" {
    query = query.Where("content like ?", fmt.Sprint("%", filter.Content, "%"))
  }
  if filter.Status > 0 {
    query = query.Where("status = ?", filter.Status)
  }
  return query
}

func GetCommentList(db squirrel.DBProxyBeginner, filter CommentFilter) (*[]Comment, error) {
  log := commentLog.With("func", "GetCommentList")
  result := make([]Comment, 0)
  column := NewComment(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableComment)
  query = query.Where("status = 1 and pid = 0")
  query = query.Where("type = ? and sid = ?", filter.Type, filter.Sid)
  if filter.LastTime > 0 {
    query = query.Where("ctime < ?", filter.LastTime)
  }
  query = query.OrderBy("ctime desc")
  util.PrintQuery(log, query)

  rows, err := query.RunWith(db).Query()
  if err != nil {
    return &result, err
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return &result, err
    }
    result = append(result, *column)
  }
  return &result, err
}

func GetCommentIdsMap(db squirrel.DBProxyBeginner, ids []int) (*map[int][]Comment, error) {
  log := commentLog.With("func", "GetCommentIdsMap")
  result := map[int][]Comment{}
  if len(ids) == 0 {
    return &result, nil
  }
  column := NewComment(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableComment)
  query = query.Where(squirrel.Eq{"pid":ids})
  query = query.OrderBy("ctime desc")
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return &result, err
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return &result, err
    }
    result[column.Pid] = append(result[column.Pid], *column)
  }
  return &result, nil
}
