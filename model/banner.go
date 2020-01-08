package model

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "github.com/Masterminds/structable"
  "time"
  "yx-api/util"
)

type (
  Banner struct {
    ID        int    `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT" desc:"ID"`
    Title     string `stbl:"title" desc:"标题"`
    Img       string `stbl:"img" desc:"配图"`
    StartTime int64  `stbl:"start_time" desc:"开始时间"`
    EndTime   int64  `stbl:"end_time" desc:"结束时间"`
    CTime     int64  `stbl:"ctime" desc:"创建时间"`
    UTime     int64  `stbl:"utime" desc:"最后修改时间"`
    Type      int    `stbl:"type" desc:"类型"`
    Status    int    `stbl:"status" desc:"状态"`
    URL       string `stbl:"url" desc:"关联URL"`
    Uid       int    `stbl:"uid" desc:"用户ID"`
    Sort      int64  `stbl:"sort" desc:"排序"`
    rec       structable.Recorder
  }

  BannerDto struct {
    ID      int
    Type    int
    Title   string
    Status  int
    util.Pagination
  }
)

func NewBanner(db squirrel.DBProxyBeginner) *Banner {
  vo := new(Banner)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableBanner, vo)
  return vo
}

func (vo *Banner) Load() error {
  return vo.rec.Load()
}

func (vo *Banner) Insert() error {
  return vo.rec.Insert()
}

func (vo *Banner) Update() error {
  return vo.rec.Update()
}

func GetBannerPagination(db squirrel.DBProxyBeginner, filter BannerDto) ([]Banner, util.Pagination, error) {
  log := bannerLog.With("func", "GetBannerPagination")
  result := make([]Banner, 0)
  count := uint64(0)
  column := NewBanner(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableBanner)
  query = query.OrderBy("sort desc, ctime desc")
  query = paramToBanner(filter, query)
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

  query = squirrel.Select("count(1)").From(tableBanner)
  query = paramToBanner(filter, query)
  err = query.RunWith(db).QueryRow().Scan(&count)
  filter.BuildBy(count)
  return result, filter.Pagination, err
}

func paramToBanner(filter BannerDto, query squirrel.SelectBuilder) squirrel.SelectBuilder {
  if filter.ID > 0 {
    query = query.Where("id = ?", filter.ID)
  }
  if filter.Type > 0 {
    query = query.Where("type = ?", filter.Type)
  }
  if filter.Title != "" {
    query = query.Where("title like ?", fmt.Sprint("%", filter.Title, "%"))
  }
  if filter.Status > 0 {
    query = query.Where("status = ?", filter.Status)
  }
  return query
}

func EditBannerSort(db squirrel.DBProxyBeginner, id int, sort int64) error {
  update := squirrel.Update(tableBanner).Set("sort", sort).Where("id = ?", id)
  fmt.Println(update.ToSql())
  _, err := update.RunWith(db).Exec()
  return err
}

func GetEffectiveBanner(db squirrel.DBProxyBeginner, typeId int) ([]Banner, error) {
  log := bannerLog.With("func", "GetEffectiveBanner")
  result := make([]Banner, 0)
  column := NewBanner(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableBanner)
  query = query.Where("? between start_time and end_time", time.Now().Unix())
  query = query.Where("status = 1 and type = ?", typeId)
  query = query.OrderBy("sort desc, ctime desc")
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
