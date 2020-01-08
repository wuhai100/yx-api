package model

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "github.com/Masterminds/structable"
  "time"
  "yx-api/util"
)

type (
  Tag struct {
    ID        int    `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT" desc:"ID"`
    Title     string `stbl:"title" desc:"标题"`
    CTime     int64  `stbl:"ctime" desc:"创建时间"`
    Status    int    `stbl:"status" desc:"状态"`
    rec       structable.Recorder
  }

  TagDto struct {
    ID      int
    Title   string
    Status  int
    util.Pagination
  }

  ArticleTagDto struct {
    ArticleID  int
    TagIds []int
  }

  ArticleTag struct {
    ArticleID int    `json:"-" desc:"文章ID"`
    TagID     int    `json:"id" desc:"标签ID"`
    Title     string `json:"title" desc:"标签名称"`
  }
)

func NewTag(db squirrel.DBProxyBeginner) *Tag {
  vo := new(Tag)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableTag, vo)
  return vo
}

func (vo *Tag) Load() error {
  return vo.rec.Load()
}

func (vo *Tag) Insert() error {
  return vo.rec.Insert()
}

func (vo *Tag) Update() error {
  return vo.rec.Update()
}

func SaveArticleTag(db squirrel.DBProxyBeginner, dto ArticleTagDto) error {
  log := tagLog.With("func", "SaveArticleTag")
  if dto.ArticleID == 0 {
    return nil
  }

  remove := squirrel.Delete(tableArticleTag).Where("article_id = ?", dto.ArticleID).Where(squirrel.NotEq{"tag_id": dto.TagIds})
  util.PrintDelete(log, remove)
  _, err := remove.RunWith(db).Exec()
  if err != nil {
    return err
  }

  if len(dto.TagIds) == 0 {
    return nil
  }

  insert := squirrel.Insert(tableArticleTag).Columns("article_id, tag_id, sort")
  sort := time.Now().Unix()
  for i, v := range dto.TagIds {
    // 保障sort是有序存储,如果==1为"置顶"标签，特定置顶为最上层
    if v == 1 {
      sort += 1000
    } else {
      sort += int64(i)
    }
    insert = insert.Values(dto.ArticleID, v, sort)
  }
  insert = insert.Suffix("ON DUPLICATE KEY UPDATE article_id=VALUES(article_id), tag_id=VALUES(tag_id), sort=VALUES(sort)")
  util.PrintInsert(log, insert)
  _, err = insert.RunWith(db).Exec()
  return err
}

func EditTagSort(db squirrel.DBProxyBeginner, articleId, tagId int, sort int64) error {
  log := tagLog.With("func", "EditTagSort")
  update := squirrel.Update(tableArticleTag).Set("sort", sort).Where("article_id = ? and tag_id = ?", articleId, tagId)
  util.PrintUpdate(log, update)
  _, err := update.RunWith(db).Exec()
  return err
}

func GetTagPagination(db squirrel.DBProxyBeginner, filter TagDto) ([]Tag, util.Pagination, error) {
  log := tagLog.With("func", "GetTagPagination")
  result := make([]Tag, 0)
  count := uint64(0)
  column := NewTag(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableTag)
  query = query.OrderBy("ctime desc")
  query = paramToTag(filter, query)
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

  query = squirrel.Select("count(1)").From(tableTag)
  query = paramToTag(filter, query)
  err = query.RunWith(db).QueryRow().Scan(&count)
  filter.BuildBy(count)
  return result, filter.Pagination, err
}

func paramToTag(filter TagDto, query squirrel.SelectBuilder) squirrel.SelectBuilder {
  if filter.ID > 0 {
    query = query.Where("id = ?", filter.ID)
  }
  if filter.Title != "" {
    query = query.Where("title like ?", fmt.Sprint("%", filter.Title, "%"))
  }
  if filter.Status > 0 {
    query = query.Where("status = ?", filter.Status)
  }
  return query
}

func GetEffectiveTag(db squirrel.DBProxyBeginner) (map[int]string, error) {
  log := tagLog.With("func", "GetEffectiveTag")
  result := map[int]string{}
  column := NewTag(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableTag).Where("status = 1").OrderBy("ctime desc")
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
    result[column.ID] = column.Title
  }
  return result, err
}

func GetArticleTagMap(db squirrel.DBProxyBeginner, articleIds []int) (map[int][]ArticleTag, error) {
  log := tagLog.With("func", "GetArticleTagMap")
  result := map[int][]ArticleTag{}
  query := squirrel.Select("article_id, id, title").From(tableArticleTag).Join(tableTag + " ON id = tag_id")
  query = query.Where("status = 1")
  query = query.Where(squirrel.Eq{"article_id": articleIds})
  query = query.OrderBy("sort desc")
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    vo := ArticleTag{}
    err = rows.Scan(&vo.ArticleID, &vo.TagID, &vo.Title)
    if err != nil {
      return result, err
    }
    result[vo.ArticleID] = append(result[vo.ArticleID], vo)
  }
  return result, err
}

func GetTagIdsMap(db squirrel.DBProxyBeginner, articleIds []int) (map[int][]int, error) {
  log := tagLog.With("func", "GetTagIdsMap")
  result := map[int][]int{}
  if len(articleIds) == 0 {
    return result, nil
  }
  query := squirrel.Select("article_id, tag_id").From(tableArticleTag).Where(squirrel.Eq{"article_id":articleIds}).OrderBy("sort desc")
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    articleId, tagId := 0, 0
    err = rows.Scan(&articleId, &tagId)
    if err != nil {
      return result, err
    }
    result[articleId] = append(result[articleId], tagId)
  }
  return result, nil
}
