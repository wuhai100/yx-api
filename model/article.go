package model

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "github.com/Masterminds/structable"
  "yx-api/util"
)

type (
  Article struct {
    ID      int    `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT" desc:"ID"`
    Title   string `stbl:"title" desc:"标题"`
    Brief   string `stbl:"brief" desc:"简介"`
    Content string `stbl:"content" desc:"内容"`
    Img     string `stbl:"img" desc:"配图"`
    CTime   int64  `stbl:"ctime" desc:"创建时间"`
    UTime   int64  `stbl:"utime" desc:"最后修改时间"`
    Uid     int    `stbl:"uid" desc:"用户ID"`
    Status  int    `stbl:"status" desc:"状态"`
    rec     structable.Recorder
  }

  ArticleDto struct {
    ID      int
    Title   string
    Content string
    Status  int
    util.Pagination
  }

  ArticleFilter struct {
    SubjectId int
    LastTime  int64
  }
)

func NewArticle(db squirrel.DBProxyBeginner) *Article {
  vo := new(Article)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableArticle, vo)
  return vo
}

func (vo *Article) Load() error {
  return vo.rec.Load()
}

func (vo *Article) Insert() error {
  return vo.rec.Insert()
}

func (vo *Article) Update() error {
  return vo.rec.Update()
}

func GetArticlePagination(db squirrel.DBProxyBeginner, filter ArticleDto) ([]Article, util.Pagination, error) {
  log := articleLog.With("func", "GetArticlePagination")
  result := make([]Article, 0)
  count := uint64(0)
  column := NewArticle(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableArticle)
  query = query.OrderBy("ctime desc")
  query = paramToArticle(filter, query)
  query = filter.PageLimit(query)
  util.PrintQuery(log, query)

  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, filter.Pagination, err
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return result, filter.Pagination, err
    }
    result = append(result, *column)
  }

  query = squirrel.Select("count(1)").From(tableArticle)
  query = paramToArticle(filter, query)
  err = query.RunWith(db).QueryRow().Scan(&count)
  filter.BuildBy(count)
  return result, filter.Pagination, err
}

func paramToArticle(filter ArticleDto, query squirrel.SelectBuilder) squirrel.SelectBuilder {
  if filter.ID > 0 {
    query = query.Where("id = ?", filter.ID)
  }
  if filter.Title != "" {
    query = query.Where("title like ?", fmt.Sprint("%", filter.Title, "%"))
  }
  if filter.Content != "" {
    query = query.Where("content like ?", fmt.Sprint("%", filter.Content, "%"))
  }
  if filter.Status > 0 {
    query = query.Where("status = ?", filter.Status)
  }
  return query
}

func GetSubjectArticleTop(db squirrel.DBProxyBeginner) ([]Article, error) {
  log := articleLog.With("func", "GetSubjectArticleList")
  result := make([]Article, 0)
  column := NewArticle(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableArticle).Join(tableArticleTag + " ON id = article_id")
  query = query.Where("tag_id = 1")
  query = query.Limit(20)
  query = query.OrderBy("sort desc")
  util.PrintQuery(log, query)

  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return result, err
    }
    result = append(result, *column)
  }
  return result, err
}

func GetSubjectArticleList(db squirrel.DBProxyBeginner, filter ArticleFilter) ([]Article, error) {
  log := articleLog.With("func", "GetSubjectArticleList")
  result := make([]Article, 0)
  column := NewArticle(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableArticle).Join(tableSubjectArticle + " ON id = article_id")
  query = query.GroupBy("id")
  query = query.Limit(20)
  if filter.LastTime > 0 {
    query = query.Where("ctime < ?", filter.LastTime)
  }
  if filter.SubjectId > 0 {
    query = query.Where("subject_id = ?", filter.SubjectId)
  }
  query = query.OrderBy("ctime desc")
  util.PrintQuery(log, query)

  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    err = rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return result, err
    }
    result = append(result, *column)
  }
  return result, err
}
