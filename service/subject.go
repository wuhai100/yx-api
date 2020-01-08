package service

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "time"
  "yx-api/model"
  "yx-api/util"
)

type (
  Subject struct {
    ID     int    `json:"id" desc:"ID"`
    Title  string `json:"title" desc:"标题"`
    Brief  string `json:"brief" desc:"简介"`
    Img    string `json:"img" desc:"配图"`
    CTime  int64  `json:"ctime" desc:"创建时间"`
    Status int    `json:"status" desc:"状态"`
    UName  string `json:"uname" desc:"用户"`
    Pid    int    `json:"pid" desc:"父级话题"`
  }

  SubjectTitle struct {
    ID     int    `json:"id" desc:"ID"`
    Title  string `json:"title" desc:"标题"`
  }

  SubjectDto struct {
    ID     int    `json:"id" desc:"ID"`
    Title  string `json:"title" desc:"标题"`
    Brief  string `json:"brief" desc:"简介"`
    Img    string `json:"img" desc:"配图"`
    CTime  int64  `json:"ctime" desc:"创建时间"`
    Status int    `json:"status" desc:"状态"`
    Uid    int    `json:"uid" desc:"用户ID"`
  }

  Column struct {
    ID        int    `json:"id" desc:"栏目ID"`
    SubjectId int    `json:"subjectId" desc:"话题ID"`
    Title     string `json:"title" desc:"标题"`
    Sort      int64  `json:"sort" desc:"排序"`
    Status    int    `json:"status" desc:"状态"`
  }

  ColumnSortDto struct {
    ID       int   `json:"id" desc:"当前id"`
    Sort     int64 `json:"sort" desc:"当前排序号"`
    NextID   int   `json:"nextId" desc:"交换排序的ID"`
    NextSort int64 `json:"nextSort" desc:"交换排序的排序号"`
  }
)

// Admin后台文章列表
func GetSubjectPagination(dto model.SubjectDto) ([]Subject, util.Pagination, error) {
  result := make([]Subject, 0)
  data, page, err := model.GetSubjectPagination(dbCache, dto)
  if err != nil {
    return result, page, fmt.Errorf("model.GetSubjectPagination error %v ", err)
  }

  userIds := make([]int, 0)
  for _, v := range data {
    userIds = append(userIds, v.Uid)
  }
  userIds = util.DeDuplicationInt(userIds)
  userMap, err := model.GetAdminUserNameMap(dbCache, userIds)
  if err != nil {
    return result, page, fmt.Errorf("model.GetAdminUserNameMap error %v ", err)
  }

  for _, v := range data {
    vo := Subject{}
    vo.ID = v.ID
    vo.Pid = v.Pid
    vo.CTime = v.CTime
    vo.Title = v.Title
    vo.Brief = v.Brief
    vo.Img = obsPrefix(v.Img)
    vo.Status = v.Status
    vo.UName = userMap[v.Uid]
    result = append(result, vo)
  }
  return result, page, err
}

// Admin后台文章列表
func GetSubjectTitle() (map[int]string, error) {
  return model.GetSubjectTitle(dbCache)
}

func AddSubject(dto SubjectDto) error {
  vo := model.NewSubject(dbCache)
  vo.Status = 1
  vo.CTime = time.Now().Unix()
  vo.UTime = vo.CTime
  vo.Uid = dto.Uid
  vo.Title = dto.Title
  vo.Brief = dto.Brief
  vo.Img = obsPrefix(dto.Img)
  return vo.Insert()
}

func EditSubject(dto SubjectDto) error {
  vo := model.NewSubject(dbCache)
  vo.ID = dto.ID
  err := vo.Load()
  if err != nil {
    return err
  }

  if dto.Status > 0 {
    vo.Status = dto.Status
  }
  if dto.Title != "" {
    vo.Title = dto.Title
  }
  if dto.Brief != "" {
    vo.Brief = dto.Brief
  }
  if dto.Img != "" {
    vo.Img = obsPrefix(dto.Img)
  }
  vo.UTime = time.Now().Unix()
  return vo.Update()
}


// Admin后台文章列表
func GetColumnList() ([]Column, error) {
  result := make([]Column, 0)
  data, err := model.GetColumnList(dbCache)
  if err != nil {
    return result, err
  }

  for _, v := range data {
    vo := Column{}
    vo.ID = v.ID
    vo.SubjectId = v.SubjectId
    vo.Title = v.Title
    vo.Status = v.Status
    vo.Sort = v.Sort
    result = append(result, vo)
  }
  return result, nil
}

func AddColumn(dto Column) error {
  vo := model.NewColumn(dbCache)
  vo.Status = 1
  vo.SubjectId = dto.SubjectId
  vo.Title = dto.Title
  vo.Sort = time.Now().Unix()
  return vo.Insert()
}

func EditColumn(dto Column) error {
  vo := model.NewColumn(dbCache)
  vo.ID = dto.ID
  err := vo.Load()
  if err != nil {
    return err
  }

  if dto.SubjectId > 0 {
    vo.SubjectId = dto.SubjectId
  }
  if dto.Status > 0 {
    vo.Status = dto.Status
  }
  if dto.Title != "" {
    vo.Title = dto.Title
  }
  return vo.Update()
}

// Admin后台修改排序
func EditColumnSort(dto ColumnSortDto) error {
  return inTx(func(beginner squirrel.DBProxyBeginner) error {
    err := model.EditColumnSort(dbCache, dto.ID, dto.NextSort)
    if err != nil {
      return err
    }
    return model.EditColumnSort(dbCache, dto.NextID, dto.Sort)
  })
}
