package controller

import (
  "github.com/yb7/echoswg"
  "yx-api/model"
  "yx-api/service"
  "yx-api/util"
)

func getCommentPagination(req *struct{ model.CommentDto }) *util.ResponseData {
  return util.ResultPageData(service.GetCommentPagination(req.CommentDto))
}

func getCommentList(user *service.User, req *struct{ model.CommentFilter }) *util.ResponseData {
  req.Uid = service.ExtractUid(user)
  return util.ResultData(service.GetCommentList(req.CommentFilter))
}

func addComment(req *struct{ Body service.CommentDto }) *util.ResponseData {
  return util.ResultData(nil, service.AddComment(req.Body))
}

func removeComment(req *struct{ Body struct{ ID int } }) *util.ResponseData {
  return util.ResultData(nil, service.EditCommentStatus(req.Body.ID, 2))
}

func init() {
  f := echoswg.NewApiGroup(util.EchoInst, "评论相关接口", "/admin")
  f.SetDescription("Admin版本")
  f.GET("/comment/list", verifyAdminUser, getCommentPagination, "评论列表")

  p := echoswg.NewApiGroup(util.EchoInst, "APP端评论信息", "/v1")
  p.SetDescription("APP版本")
  p.GET("/comment/list", verifyUserNotLogin, getCommentList, "评论列表")
  p.POST("/comment", preventFrequent2s, verifyUserToken, addComment, "新增评论")
  p.DELETE("/comment", preventFrequent2s, verifyUserToken, removeComment, "删除评论")

}
