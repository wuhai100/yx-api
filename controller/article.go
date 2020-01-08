package controller

import (
  "github.com/yb7/echoswg"
  "yx-api/model"
  "yx-api/service"
  "yx-api/util"
)


func getArticlePagination(req *struct{ model.ArticleDto }) *util.ResponseData {
  return util.ResultPageData(service.GetArticlePagination(req.ArticleDto))
}

func addArticle(user *service.AdminUser, req *struct{ Body service.ArticleDto }) *util.ResponseData {
  req.Body.Uid = service.ExtractAdminUid(user)
  return util.ResultData(nil, service.AddArticle(req.Body))
}

func editArticle(user *service.AdminUser, req *struct{ Body service.ArticleDto }) *util.ResponseData {
  req.Body.Uid = service.ExtractAdminUid(user)
  return util.ResultData(nil, service.EditArticle(req.Body))
}

func editArticleStatus(req *struct{ Body struct{ ID, Status int } }) *util.ResponseData {
  return util.ResultData(nil, service.EditArticleStatus(req.Body.ID, req.Body.Status))
}

// 首页
func homeAssembled(req *struct{ service.AssembledDto }) *util.ResponseData {
  return util.ResultData(service.HomeAssembled(req.AssembledDto))
}

func init() {
  f := echoswg.NewApiGroup(util.EchoInst, "文章相关接口", "/admin")
  f.SetDescription("Admin版本")
  f.GET("/article/list", verifyAdminUser, getArticlePagination, "文章列表")
  f.POST("/article", preventFrequent2s, verifyAdminUser, addArticle, "新增文章")
  f.PUT("/article", preventFrequent2s, verifyAdminUser, editArticle, "编辑文章")
  f.PUT("/article/status", preventFrequent2s, verifyAdminUser, editArticleStatus, "编辑文章状态")

  p := echoswg.NewApiGroup(util.EchoInst, "APP端文章信息", "/v1")
  p.SetDescription("APP版本")
  p.GET("/home/assembled/:ID", verifyUserNotLogin, homeAssembled, "首页")
}
