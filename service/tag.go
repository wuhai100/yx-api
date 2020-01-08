package service

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "time"
  "yx-api/model"
  "yx-api/util"
)

type (
  Tag struct {
    ID        int    `json:"id" desc:"ID"`
    Title     string `json:"title" desc:"标题"`
    CTime     int64  `json:"ctime" desc:"创建时间"`
    Status    int    `json:"status" desc:"状态"`
  }

  TagDto struct {
    ID        int    `json:"id" desc:"ID"`
    Title     string `json:"title" desc:"标题"`
    Status    int    `json:"status" desc:"状态"`
  }

  TagSortDto struct {
    ArticleID     int   `json:"articleId" desc:"当前articleId"`
    TagID         int   `json:"tagId" desc:"当前tagId"`
    Sort          int64 `json:"sort" desc:"当前排序号"`
    NextArticleID int   `json:"nextArticleId" desc:"交换排序的ID"`
    NextTagID     int   `json:"nextTagId" desc:"交换排序的ID"`
    NextSort      int64 `json:"nextSort" desc:"交换排序的排序号"`
  }

  EffectiveTag struct {
    ID        int    `json:"id" desc:"ID"`
    Title     string `json:"title" desc:"标题"`
  }
)

// Admin后台文章列表
func GetTagPagination(dto model.TagDto) ([]Tag, util.Pagination, error) {
  result := make([]Tag, 0)
  data, page, err := model.GetTagPagination(dbCache, dto)
  if err != nil {
    return result, page, fmt.Errorf("model.GetTagPagination error %v ", err)
  }

  for _, v := range data {
    vo := Tag{}
    vo.ID = v.ID
    vo.Title = v.Title
    vo.CTime = v.CTime
    vo.Status = v.Status
    result = append(result, vo)
  }
  return result, page, err
}

func AddTag(dto TagDto) error {
  vo := model.NewTag(dbCache)
  vo.Status = 1
  vo.CTime = time.Now().Unix()
  vo.Title = dto.Title
  return vo.Insert()
}

func EditTag(dto TagDto) error {
  vo := model.NewTag(dbCache)
  vo.ID = dto.ID
  err := vo.Load()
  if err != nil {
    return err
  }
  vo.Title = dto.Title
  return vo.Update()
}

func EditTagStatus(id, status int) error {
  vo := model.NewTag(dbCache)
  vo.ID = id
  err := vo.Load()
  if err != nil {
    return err
  }
  vo.Status = status
  return vo.Update()
}

// Admin后台修改排序
func EditTagSort(dto TagSortDto) error {
  return inTx(func(beginner squirrel.DBProxyBeginner) error {
    err := model.EditTagSort(dbCache, dto.ArticleID, dto.TagID, dto.NextSort)
    if err != nil {
      return err
    }
    return model.EditTagSort(dbCache, dto.NextArticleID, dto.NextTagID, dto.Sort)
  })
}

func GetEffectiveTag() (map[int]string, error) {
  return  model.GetEffectiveTag(dbCache)
}
