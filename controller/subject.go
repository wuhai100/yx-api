package controller

import (
  "github.com/yb7/echoswg"
  "yx-api/model"
  "yx-api/service"
  "yx-api/util"
)

func getSubjectTitle() *util.ResponseData {
  return util.ResultData(service.GetSubjectTitle())
}

func getSubjectPagination(req *struct{ model.SubjectDto }) *util.ResponseData {
  return util.ResultPageData(service.GetSubjectPagination(req.SubjectDto))
}

func addSubject(user *service.AdminUser, req *struct{ Body service.SubjectDto }) *util.ResponseData {
  req.Body.Uid = service.ExtractAdminUid(user)
  return util.ResultData(nil, service.AddSubject(req.Body))
}

func editSubject(user *service.AdminUser, req *struct{ Body service.SubjectDto }) *util.ResponseData {
  req.Body.Uid = service.ExtractAdminUid(user)
  return util.ResultData(nil, service.EditSubject(req.Body))
}

func editSubjectStatus(req *struct{ Body service.SubjectDto }) *util.ResponseData {
  return util.ResultData(nil, service.EditSubject(req.Body))
}

func getColumnList() *util.ResponseData {
  return util.ResultData(service.GetColumnList())
}

func addColumn(req *struct{ Body service.Column }) *util.ResponseData {
  return util.ResultData(nil, service.AddColumn(req.Body))
}

func editColumn(req *struct{ Body service.Column }) *util.ResponseData {
  return util.ResultData(nil, service.EditColumn(req.Body))
}

func editColumnStatus(req *struct{ Body service.Column }) *util.ResponseData {
  return util.ResultData(nil, service.EditColumn(req.Body))
}

func editColumnSort(req *struct { Body service.ColumnSortDto }) *util.ResponseData {
  return util.ResultData(nil, service.EditColumnSort(req.Body))
}

func init() {
  f := echoswg.NewApiGroup(util.EchoInst, "话题相关接口", "/admin")
  f.SetDescription("Admin版本")
  f.GET("/subject/title", verifyAdminUser, getSubjectTitle, "话题列表")
  f.GET("/subject/list", verifyAdminUser, getSubjectPagination, "话题列表")
  f.POST("/subject", preventFrequent2s, verifyAdminUser, addSubject, "新增话题")
  f.PUT("/subject", preventFrequent2s, verifyAdminUser, editSubject, "编辑话题")
  f.PUT("/subject/status", preventFrequent2s, verifyAdminUser, editSubjectStatus, "编辑话题状态")

  c := echoswg.NewApiGroup(util.EchoInst, "栏目相关接口", "/admin")
  c.SetDescription("Admin版本")
  c.GET("/column/list", verifyAdminUser, getColumnList, "栏目列表")
  c.POST("/column", preventFrequent2s, verifyAdminUser, addColumn, "新增栏目")
  c.PUT("/column", preventFrequent2s, verifyAdminUser, editColumn, "编辑栏目")
  c.PUT("/column/status", preventFrequent2s, verifyAdminUser, editColumnStatus, "编辑栏目状态")
  f.PUT("/column/sort", preventFrequent2s, verifyAdminUser, editColumnSort, "调整栏目顺序")

  p := echoswg.NewApiGroup(util.EchoInst, "APP端话题信息", "/v1")
  p.SetDescription("APP版本")
}
