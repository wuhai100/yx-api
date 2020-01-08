package service

import (
  "fmt"
  "time"
  "yx-api/model"
  "yx-api/util"
)

type (
  Comment struct {
    ID       int             `json:"id" desc:"ID"`
    CTime    int64           `json:"ctime" desc:"评论时间"`
    Content  string          `json:"content" desc:"评论内容"`
    Status   int             `json:"status" desc:"状态"`
    RNum     int             `json:"rnum" desc:"被回复数量"`
    PId      int             `json:"pid" desc:"回复评论ID"`
    ZNum     int             `json:"znum" desc:"点赞数量"`
    FUser    CommentUserInfo `json:"fuser" desc:"评论用户信息"`
    RUser    interface{}     `json:"ruser" desc:"回复用户信息"`
    RComment []Comment       `json:"rcomment" desc:"回复评论列表"`
  }

  CommentUserInfo struct {
    ID       int    `json:"id" desc:"用户ID"`
    NickName string `json:"nickName" desc:"用户昵称"`
    Avatar   string `json:"avatar" desc:"头像"`
  }

  CommentDto struct {
    Type    int    `json:"type" desc:"资源类型"`
    Sid     int    `json:"sid" desc:"资源ID"`
    Rid     int    `json:"rid" desc:"回复用户ID"`
    CTime   int64  `json:"ctime" desc:"评论时间"`
    Content string `json:"content" desc:"评论内容"`
    Status  int    `json:"status" desc:"状态"`
    Pid     int    `json:"pid" desc:"回复评论ID"`
    Uid     int    `json:"-" desc:"评论用户ID"`
  }
)

// Admin后台评论列表
func GetCommentPagination(dto model.CommentDto) (*[]Comment, util.Pagination, error) {
  result := make([]Comment, 0)
  data, page, err := model.GetCommentPagination(dbCache, dto)
  if err != nil {
    return &result, page, fmt.Errorf("model.GetCommentPagination error %v ", err)
  }

  ids := make([]int, 0)
  userIds := make([]int, 0)
  for _, v := range *data {
    ids = append(ids, v.ID)
    userIds = append(userIds, v.Uid, v.Rid)
  }
  userIds = util.DeDuplicationInt(userIds)
  userMap, err := getCommentUserInfo(userIds)
  if err != nil {
    return &result, page, fmt.Errorf("getCommentUserInfo error %v ", err)
  }
  rcommentMap, err := getReplyCommentMap(ids)
  if err != nil {
    return &result, page, fmt.Errorf("getReplyCommentMap error %v ", err)
  }

  for _, v := range *data {
    vo := Comment{}
    vo.ID = v.ID
    vo.Content = v.Content
    vo.CTime = v.CTime
    vo.Status = v.Status
    vo.FUser = userMap[v.Uid]
    vo.RUser = struct{}{}
    vo.RComment = []Comment{}
    if v, ok := userMap[v.Rid]; ok {
      vo.RUser = v
    }
    if v, ok := rcommentMap[v.ID]; ok {
      vo.RComment = v
    }
    result = append(result, vo)
  }
  return &result, page, err
}

// App评论
func GetCommentList(dto model.CommentFilter) (*[]Comment, error) {
  result := make([]Comment, 0)
  data, err := model.GetCommentList(dbCache, dto)
  if err != nil {
    return &result, fmt.Errorf("model.GetCommentList error %v ", err)
  }

  ids := make([]int, 0)
  userIds := make([]int, 0)
  for _, v := range *data {
    ids = append(ids, v.ID)
    userIds = append(userIds, v.Uid, v.Rid)
  }
  userIds = util.DeDuplicationInt(userIds)
  userMap, err := getCommentUserInfo(userIds)
  if err != nil {
    return &result, fmt.Errorf("getCommentUserInfo error %v ", err)
  }

  for _, v := range *data {
    vo := Comment{}
    vo.ID = v.ID
    vo.Content = v.Content
    vo.CTime = v.CTime
    vo.Status = v.Status
    vo.FUser = userMap[v.Uid]
    vo.RUser = struct{}{}
    if v, ok := userMap[v.Rid]; ok {
      vo.RUser = v
    }
    result = append(result, vo)
  }
  return &result, err
}

// 获取评论用户和回复评论用户的基本信息
func getCommentUserInfo(userIds []int) (map[int]CommentUserInfo, error) {
  result := map[int]CommentUserInfo{}
  userMap, err := model.GetUserMap(dbCache, userIds)
  if err != nil {
    return result, fmt.Errorf("model.GetUserMap error %v ", err)
  }
  for k, v := range userMap {
    result[k] = CommentUserInfo{ID:v.ID, NickName:v.NickName, Avatar:v.Avatar}
  }
  return result, err
}

// 获取回复评论列表
func getReplyCommentMap(ids []int) (map[int][]Comment, error) {
  result := map[int][]Comment{}
  commentMap, err := model.GetCommentIdsMap(dbCache, ids)
  if err != nil {
    return result, fmt.Errorf("model.GetCommentIdsMap error %v ", err)
  }
  userIds := make([]int, 0)
  for _, v := range *commentMap {
    for _, c := range v {
      userIds = append(userIds, c.Uid, c.Rid)
    }
  }
  userIds = util.DeDuplicationInt(userIds)
  userMap, err := getCommentUserInfo(userIds)
  if err != nil {
    return result, fmt.Errorf("getCommentUserInfo error %v ", err)
  }

  for k, temp := range *commentMap {
    for _, c := range temp {
      vo := Comment{}
      vo.ID = c.ID
      vo.Content = c.Content
      vo.CTime = c.CTime
      vo.Status = c.Status
      vo.FUser = userMap[c.Uid]
      vo.RUser = struct{}{}
      vo.RComment = []Comment{}
      if v, ok := userMap[c.Rid]; ok {
        vo.RUser = v
      }
      result[k] = append(result[k], vo)
    }
  }
  return result, err
}


func AddComment(dto CommentDto) error {
  vo := model.NewComment(dbCache)
  vo.Status = 1
  vo.CTime = time.Now().Unix()
  vo.Type = dto.Type
  vo.Sid = dto.Sid
  vo.Uid = dto.Uid
  vo.Rid = dto.Rid
  vo.Pid = dto.Pid
  vo.Content = dto.Content
  return vo.Insert()
}

func EditCommentStatus(id, status int) error {
  vo := model.NewComment(dbCache)
  vo.ID = id
  err := vo.Load()
  if err != nil {
    return err
  }
  vo.Status = status
  return vo.Update()
}
