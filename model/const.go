package model

import "yx-api/util"

const (
  tablePrefix         = "yx_"
  tableAdminUser      = tablePrefix + "admin_user"
  tableUser           = tablePrefix + "user"
  tableCountry        = tablePrefix + "country"
  tableSubject        = tablePrefix + "subject"
  tableArticle        = tablePrefix + "article"
  tableSubjectArticle = tablePrefix + "subject_article"
  tableBanner         = tablePrefix + "banner"
  tableColumn         = tablePrefix + "column"
  tableTag            = tablePrefix + "tag"
  tableArticleTag     = tablePrefix + "article_tag"
  tableComment        = tablePrefix + "comment"
  tableUserFollow     = tablePrefix + "user_follow"
)

var (
  conf       = util.Get()
  userLog    = util.AppLog.With("file", "model.user.go")
  subjectLog = util.AppLog.With("file", "model.subject.go")
  articleLog = util.AppLog.With("file", "model.article.go")
  bannerLog  = util.AppLog.With("file", "model.banner.go")
  tagLog     = util.AppLog.With("file", "model.tag.go")
  commentLog = util.AppLog.With("file", "model.comment.go")
)
