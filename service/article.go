package service

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "time"
  "yx-api/model"
  "yx-api/util"
  "github.com/yb7/asd"
)

type (
  Article struct {
    ID         int    `json:"id" desc:"ID"`
    Title      string `json:"title" desc:"标题"`
    Brief      string `json:"brief" desc:"简介"`
    Content    string `json:"content" desc:"内容"`
    Img        string `json:"img" desc:"配图"`
    CTime      int64  `json:"ctime" desc:"创建时间"`
    Status     int    `json:"status" desc:"状态"`
    UName      string `json:"uname" desc:"用户名称"`
    SubjectIds []int  `json:"subjectIds" desc:"关联话题"`
    TagIds     []int  `json:"tagIds" desc:"关联标签"`
  }

  ArticleDto struct {
    ID         int    `json:"id" desc:"ID"`
    Title      string `json:"title" desc:"标题"`
    Brief      string `json:"brief" desc:"简介"`
    Content    string `json:"content" desc:"内容"`
    Img        string `json:"img" desc:"配图"`
    Status     int    `json:"status" desc:"状态"`
    SubjectIds []int  `json:"subjectIds" desc:"话题Ids"`
    TagIds     []int  `json:"tagIds" desc:"标签Ids"`
    Uid        int    `json:"-" desc:"用户ID"`
  }

  // 首页数据
  HomeAssembledView struct {
    BannerList  []EffectiveBanner `json:"bannerList" desc:"轮播图"`
    TopArticle  []ArticleVo       `json:"topArticle" desc:"置顶文章列表"`
    ArticleList []ArticleVo       `json:"articleList" desc:"文章流列表"`
  }

  ArticleVo struct {
    ID          int                    `json:"id" desc:"ID"`
    Title       string                 `json:"title" desc:"标题"`
    Brief       string                 `json:"brief" desc:"简介"`
    Img         string                 `json:"img" desc:"配图"`
    CTime       int64                  `json:"ctime" desc:"创建时间"`
    SubjectList []model.ArticleSubject `json:"subjectList" desc:"关联话题"`
    TagList     []model.ArticleTag     `json:"tagList" desc:"关联标签"`
  }

  AssembledDto struct {
    SubjectId int   `json:"subjectId" desc:"话题ID"`
    LastTime  int64 `json:"lastTime" desc:"上拉分页时间"`
  }
)

// Admin后台文章列表
func GetArticlePagination(dto model.ArticleDto) ([]Article, util.Pagination, error) {
  result := make([]Article, 0)
  data, page, err := model.GetArticlePagination(dbCache, dto)
  if err != nil {
    return result, page, fmt.Errorf("model.GetArticlePagination error %v ", err)
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
  subjectMap, err := model.GetSubjectIdsMap(dbCache, ids)
  if err != nil {
    return result, page, fmt.Errorf("model.GetSubjectArticleMap error %v ", err)
  }
  tagMap, err := model.GetTagIdsMap(dbCache, ids)
  if err != nil {
    return result, page, fmt.Errorf("model.GetTagIdsMap error %v ", err)
  }

  for _, v := range data {
    vo := Article{}
    vo.ID = v.ID
    vo.CTime = v.CTime
    vo.Title = v.Title
    vo.Brief = v.Brief
    vo.Content = v.Content
    vo.Img = obsPrefix(v.Img)
    vo.Status = v.Status
    vo.UName = userMap[v.Uid]
    vo.SubjectIds = subjectMap[v.ID]
    vo.TagIds = tagMap[v.ID]
    if vo.SubjectIds == nil {
      vo.SubjectIds = []int{}
    }
    if vo.TagIds == nil {
      vo.TagIds = []int{}
    }
    result = append(result, vo)
  }
  return result, page, err
}

func AddArticle(dto ArticleDto) error {
  return inTx(func(dbProxyBeginner squirrel.DBProxyBeginner) error {
    vo := model.NewArticle(dbProxyBeginner)
    vo.Status = 1
    vo.CTime = time.Now().Unix()
    vo.UTime = vo.CTime
    vo.Uid = dto.Uid
    vo.Title = dto.Title
    vo.Brief = dto.Brief
    vo.Content = dto.Content
    vo.Img = obsPrefix(dto.Img)
    err := vo.Insert()
    if err != nil {
      return err
    }
    err = model.SaveSubjectArticle(dbProxyBeginner, model.SubjectArticleDto{ArticleID:vo.ID, SubjectIds:dto.SubjectIds})
    if err != nil {
      return err
    }
    return model.SaveArticleTag(dbProxyBeginner, model.ArticleTagDto{ ArticleID:vo.ID, TagIds:dto.TagIds })
  })
}

func EditArticle(dto ArticleDto) error {
  return inTx(func(dbProxyBeginner squirrel.DBProxyBeginner) error {
    vo := model.NewArticle(dbCache)
    vo.ID = dto.ID
    err := vo.Load()
    if err != nil {
      return err
    }
    vo.UTime = time.Now().Unix()
    vo.Title = dto.Title
    vo.Brief = dto.Brief
    vo.Content = dto.Content
    vo.Img = obsPrefix(dto.Img)
    err = vo.Update()
    if err != nil {
      return err
    }
    err = model.SaveSubjectArticle(dbProxyBeginner, model.SubjectArticleDto{ArticleID:vo.ID, SubjectIds:dto.SubjectIds})
    if err != nil {
      return err
    }
    return model.SaveArticleTag(dbProxyBeginner, model.ArticleTagDto{ ArticleID:vo.ID, TagIds:dto.TagIds })
  })
}

func EditArticleStatus(id, status int) error {
  vo := model.NewArticle(dbCache)
  vo.ID = id
  err := vo.Load()
  if err != nil {
    return err
  }
  vo.Status = status
  vo.UTime = time.Now().Unix()
  return vo.Update()
}

func HomeAssembled(dto AssembledDto) (*HomeAssembledView, error) {
  var err error
  result := &HomeAssembledView{}
  key := util.UniqueKey("article", "HomeAssembled", dto.SubjectId, dto.LastTime)
  err = asd.OnceInRedis(key, 30 * time.Second, func() (interface{}, error) {
    dataSource := &HomeAssembledView{}
    dataSource.BannerList, err = getEffectiveBanner(dto.SubjectId)
    if err != nil {
      return dataSource, err
    }

    // 置顶文章
    topList, err := model.GetSubjectArticleTop(dbCache)
    if err != nil {
      return dataSource, err
    }

    // 普通文章流
    articleList, err := model.GetSubjectArticleList(dbCache, model.ArticleFilter{SubjectId:dto.SubjectId})
    if err != nil {
      return dataSource, err
    }

    data := append(topList, articleList...)

    ids := make([]int, len(data))
    for _, v := range data {
      ids = append(ids, v.ID)
    }
    ids = util.DeDuplicationInt(ids)

    // 文章相关话题
    subjectMap, err := model.GetSubjectArticleMap(dbCache, ids)
    if err != nil {
      return dataSource, err
    }

    // 文章相关标签
    tagMap, err := model.GetArticleTagMap(dbCache, ids)
    if err != nil {
      return dataSource, err
    }

    for _, v := range topList {
      vo := ArticleVo{}
      vo.ID = v.ID
      vo.Title = v.Title
      vo.Brief = v.Brief
      vo.Img = obsPrefix(v.Img)
      vo.CTime = v.CTime
      vo.SubjectList = subjectMap[v.ID]
      vo.TagList = tagMap[v.ID]
      dataSource.TopArticle = append(dataSource.TopArticle, vo)
    }

    for _, v := range articleList {
      vo := ArticleVo{}
      vo.ID = v.ID
      vo.Title = v.Title
      vo.Brief = v.Brief
      vo.Img = obsPrefix(v.Img)
      vo.CTime = v.CTime
      vo.SubjectList = subjectMap[v.ID]
      vo.TagList = tagMap[v.ID]
      dataSource.ArticleList = append(dataSource.ArticleList, vo)
    }
    return dataSource, err
  }, &result)
  return result, err
}
