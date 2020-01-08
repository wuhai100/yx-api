package model

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "github.com/Masterminds/structable"
  "time"
  "yx-api/util"
)

type (
  Subject struct {
    ID     int    `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT" desc:"ID"`
    Title  string `stbl:"title" desc:"标题"`
    Brief  string `stbl:"brief" desc:"简介"`
    Img    string `stbl:"img" desc:"配图"`
    CTime  int64  `stbl:"ctime" desc:"创建时间"`
    UTime  int64  `stbl:"utime" desc:"最后修改时间"`
    Status int    `stbl:"status" desc:"状态"`
    Uid    int    `stbl:"uid" desc:"用户ID"`
    Pid    int    `stbl:"pid" desc:"父级话题"`
    rec    structable.Recorder
  }

  ArticleSubject struct {
    ArticleID int    `json:"-" desc:"文章ID"`
    SubjectID int    `json:"id" desc:"话题ID"`
    Title     string `json:"title" desc:"标签名称"`
  }

  SubjectArticleDto struct {
    ArticleID  int
    SubjectIds []int
  }

  SubjectDto struct {
    ID    int
    Pid   int
    Title string
    Brief string
    util.Pagination
  }

  Column struct {
    ID        int    `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT" desc:"栏目ID"`
    SubjectId int    `stbl:"subject_id" desc:"话题ID"`
    Title     string `stbl:"title" desc:"栏目名称"`
    Sort      int64  `stbl:"sort" desc:"排序"`
    Status    int    `stbl:"status" desc:"状态"`
    rec       structable.Recorder
  }
)

func NewSubject(db squirrel.DBProxyBeginner) *Subject {
  vo := new(Subject)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableSubject, vo)
  return vo
}

func (vo *Subject) Load() error {
  return vo.rec.Load()
}

func (vo *Subject) Insert() error {
  return vo.rec.Insert()
}

func (vo *Subject) Update() error {
  return vo.rec.Update()
}

func NewColumn(db squirrel.DBProxyBeginner) *Column {
  vo := new(Column)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableColumn, vo)
  return vo
}

func (vo *Column) Load() error {
  return vo.rec.Load()
}

func (vo *Column) Insert() error {
  return vo.rec.Insert()
}

func (vo *Column) Update() error {
  return vo.rec.Update()
}

func SaveSubjectArticle(db squirrel.DBProxyBeginner, dto SubjectArticleDto) error {
  log := subjectLog.With("func", "SaveSubjectArticle")
  if dto.ArticleID == 0 {
    return nil
  }

  remove := squirrel.Delete(tableSubjectArticle).Where("article_id = ?", dto.ArticleID).Where(squirrel.NotEq{"subject_id": dto.SubjectIds})
  util.PrintDelete(log, remove)
  _, err := remove.RunWith(db).Exec()
  if err != nil {
    return err
  }

  if len(dto.SubjectIds) == 0 {
    return nil
  }

  nowUnix := time.Now().Unix()
  insert := squirrel.Insert(tableSubjectArticle).Columns("article_id, subject_id, sort")
  for _, v := range dto.SubjectIds {
    insert = insert.Values(dto.ArticleID, v, nowUnix)
    nowUnix++
  }
  insert = insert.Suffix("ON DUPLICATE KEY UPDATE article_id=VALUES(article_id), subject_id=VALUES(subject_id), sort=VALUES(sort)")
  util.PrintInsert(log, insert)
  _, err = insert.RunWith(db).Exec()
  return err
}

func GetSubjectArticleMap(db squirrel.DBProxyBeginner, articleIds []int) (map[int][]ArticleSubject, error) {
  log := subjectLog.With("func", "GetSubjectArticleMap")
  result := map[int][]ArticleSubject{}
  if len(articleIds) == 0 {
    return result, nil
  }
  query := squirrel.Select("article_id, id, title").From(tableSubjectArticle).Join(tableSubject + " ON subject_id = id")
  query = query.Where(squirrel.Eq{"article_id":articleIds})
  query = query.OrderBy("sort")
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    vo := ArticleSubject{}
    err = rows.Scan(&vo.ArticleID, &vo.SubjectID, &vo.Title)
    if err != nil {
      return result, err
    }
    result[vo.ArticleID] = append(result[vo.ArticleID], vo)
  }
  return result, nil
}

func GetSubjectIdsMap(db squirrel.DBProxyBeginner, articleIds []int) (map[int][]int, error) {
  log := subjectLog.With("func", "GetSubjectIdsMap")
  result := map[int][]int{}
  if len(articleIds) == 0 {
    return result, nil
  }
  query := squirrel.Select("article_id, subject_id").From(tableSubjectArticle).Where(squirrel.Eq{"article_id":articleIds})
  query = query.Where("sort")
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    articleId, subjectId := 0, 0
    err = rows.Scan(&articleId, &subjectId)
    if err != nil {
      return result, err
    }
    result[articleId] = append(result[articleId], subjectId)
  }
  return result, nil
}

func GetSubjectTitle(db squirrel.DBProxyBeginner) (map[int]string, error) {
  log := subjectLog.With("func", "GetSubjectTitle")
  result := map[int]string{}
  query := squirrel.Select("id, title").From(tableSubject).Where("status = 1")
  query = query.OrderBy("ctime desc")
  util.PrintQuery(log, query)

  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    id, title := 0, ""
    err = rows.Scan(&id, &title)
    if err != nil {
      return result, err
    }
    result[id] = title
  }
  return result, err
}

func GetSubjectPagination(db squirrel.DBProxyBeginner, filter SubjectDto) ([]Subject, util.Pagination, error) {
  log := subjectLog.With("func", "GetSubjectPagination")
  result := make([]Subject, 0)
  count := uint64(0)
  column := NewSubject(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableSubject)
  query = query.OrderBy("ctime desc")
  query = paramToSubject(filter, query)
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

  query = squirrel.Select("count(1)").From(tableSubject)
  query = paramToSubject(filter, query)
  err = query.RunWith(db).QueryRow().Scan(&count)
  filter.BuildBy(count)
  return result, filter.Pagination, err
}

func paramToSubject(filter SubjectDto, query squirrel.SelectBuilder) squirrel.SelectBuilder {
  if filter.ID > 0 {
    query = query.Where("id = ?", filter.ID)
  }
  if filter.Pid > 0 {
    query = query.Where("pid = ?", filter.Pid)
  }
  if filter.Title != "" {
    query = query.Where("title like ?", fmt.Sprint("%", filter.Title, "%"))
  }
  if filter.Brief != "" {
    query = query.Where("brief like ?", fmt.Sprint("%", filter.Brief, "%"))
  }
  return query
}

func GetColumnList(db squirrel.DBProxyBeginner) ([]Column, error) {
  log := subjectLog.With("func", "GetColumnList")
  result := make([]Column, 0)
  column := NewColumn(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableColumn)
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

func EditColumnSort(db squirrel.DBProxyBeginner, id int, sort int64) error {
  log := subjectLog.With("func", "EditColumnSort")
  update := squirrel.Update(tableColumn).Set("sort", sort).Where("id = ?", id)
  util.PrintUpdate(log, update)
  _, err := update.RunWith(db).Exec()
  return err
}
