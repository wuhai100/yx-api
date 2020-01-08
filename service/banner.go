package service

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "time"
  "yx-api/model"
  "yx-api/util"
)

type (
  Banner struct {
    ID        int    `json:"id" desc:"ID"`
    Title     string `json:"title" desc:"标题"`
    Img       string `json:"img" desc:"配图"`
    StartTime int64  `json:"startTime" desc:"开始时间"`
    EndTime   int64  `json:"endTime" desc:"结束时间"`
    CTime     int64  `json:"ctime" desc:"创建时间"`
    Type      int    `json:"type" desc:"类型"`
    Status    int    `json:"status" desc:"状态"`
    Sort      int64  `json:"sort" desc:"排序"`
    URL       string `json:"url" desc:"关联URL"`
    UName     string `json:"uname" desc:"用户名称"`
  }

  BannerDto struct {
    ID        int    `json:"id" desc:"ID"`
    Title     string `json:"title" desc:"标题"`
    Img       string `json:"img" desc:"配图"`
    StartTime int64  `json:"startTime" desc:"开始时间"`
    EndTime   int64  `json:"endTime" desc:"结束时间"`
    Type      int    `json:"type" desc:"类型"`
    Status    int    `json:"status" desc:"状态"`
    URL       string `json:"url" desc:"关联URL"`
    Sort      int    `json:"sort" desc:"排序"`
    Uid       int    `json:"-" desc:"用户ID"`
  }

  BannerSortDto struct {
    ID       int   `json:"id" desc:"当前id"`
    Sort     int64 `json:"sort" desc:"当前排序号"`
    NextID   int   `json:"nextId" desc:"交换排序的ID"`
    NextSort int64 `json:"nextSort" desc:"交换排序的排序号"`
  }

  EffectiveBanner struct {
    ID        int    `json:"id" desc:"ID"`
    Title     string `json:"title" desc:"标题"`
    Img       string `json:"img" desc:"配图"`
    URL       string `json:"url" desc:"关联URL"`
  }
)

// Admin后台文章列表
func GetBannerPagination(dto model.BannerDto) ([]Banner, util.Pagination, error) {
  result := make([]Banner, 0)
  data, page, err := model.GetBannerPagination(dbCache, dto)
  if err != nil {
    return result, page, fmt.Errorf("model.GetBannerPagination error %v ", err)
  }

  ids := make([]int, 0)
  userIds := make([]int, 0)
  for _, v := range data {
    ids = append(ids, v.ID)
    userIds = append(userIds, v.Uid)
  }
  userIds = util.DeDuplicationInt(userIds)
  userMap, err := model.GetAdminUserNameMap(dbCache, userIds)
  if err != nil {
    return result, page, fmt.Errorf("model.GetAdminUserNameMap error %v ", err)
  }

  for _, v := range data {
    vo := Banner{}
    vo.ID = v.ID
    vo.Title = v.Title
    vo.Img = obsPrefix(v.Img)
    vo.StartTime = v.StartTime
    vo.EndTime = v.EndTime
    vo.CTime = v.CTime
    vo.Type = v.Type
    vo.Status = v.Status
    vo.Sort = v.Sort
    vo.URL = v.URL
    vo.UName = userMap[v.Uid]
    result = append(result, vo)
  }
  return result, page, err
}

func AddBanner(dto BannerDto) error {
  vo := model.NewBanner(dbCache)
  vo.Status = 1
  vo.CTime = time.Now().Unix()
  vo.UTime = vo.CTime
  vo.Uid = dto.Uid
  vo.Title = dto.Title
  vo.Img = obsPrefix(dto.Img)
  vo.StartTime = dto.StartTime
  vo.EndTime = dto.EndTime
  vo.Sort = vo.CTime
  vo.Type = dto.Type
  vo.URL = dto.URL
  return vo.Insert()
}

func EditBanner(dto BannerDto) error {
  vo := model.NewBanner(dbCache)
  vo.ID = dto.ID
  err := vo.Load()
  if err != nil {
    return err
  }
  vo.UTime = time.Now().Unix()
  vo.Uid = dto.Uid
  vo.Title = dto.Title
  vo.Img = obsPrefix(dto.Img)
  vo.StartTime = dto.StartTime
  vo.EndTime = dto.EndTime
  vo.Type = dto.Type
  vo.URL = dto.URL
  return vo.Update()
}

func EditBannerStatus(id, status int) error {
  vo := model.NewBanner(dbCache)
  vo.ID = id
  err := vo.Load()
  if err != nil {
    return err
  }
  vo.Status = status
  vo.UTime = time.Now().Unix()
  return vo.Update()
}

// Admin后台修改排序
func EditBannerSort(dto BannerSortDto) error {
  return inTx(func(beginner squirrel.DBProxyBeginner) error {
    err := model.EditBannerSort(dbCache, dto.ID, dto.NextSort)
    if err != nil {
      return err
    }
    return model.EditBannerSort(dbCache, dto.NextID, dto.Sort)
  })
}

func getEffectiveBanner(typeId int) ([]EffectiveBanner, error) {
  result := make([]EffectiveBanner, 0)
  list, err := model.GetEffectiveBanner(dbCache, typeId)
  if err != nil {
    return result, err
  }
  for _, v := range list {
    vo := EffectiveBanner{}
    vo.ID = v.ID
    vo.Title = v.Title
    vo.Img = v.Img
    vo.URL = v.URL
    result = append(result, vo)
  }
  return result, nil
}
