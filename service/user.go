package service

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "github.com/labstack/echo/v4"
  "strconv"
  "time"
  "yx-api/model"
  "yx-api/util"
)

type (
  AdminUser struct {
    ID     int    `json:"id"`
    Name   string `json:"name" desc:"用户名"`
    Avatar string `json:"avatar" desc:"头像"`
    Token  string `json:"token" desc:"登录令牌"`
    ApiEnv string `json:"apiEnv" desc:"系统环境变量test|prod"`
  }

  AdminUserList struct {
    ID     int    `json:"id"`
    Name   string `json:"name" desc:"用户名"`
    Avatar string `json:"avatar" desc:"头像"`
    Status int    `json:"status" desc:"状态"`
    CTime  int64  `json:"ctime" desc:"新增时间"`
  }

  User struct {
    ID         int    `json:"id"`
    Account    string `json:"account" desc:"账号"`
    OpenID     string `json:"openId" desc:"微信OpenID"`
    NickName   string `json:"nickName" desc:"用户昵称"`
    Avatar     string `json:"avatar" desc:"头像"`
    Gender     int    `json:"gender" desc:"性别"`
    CountryID  int    `json:"CountryId" desc:"国家"`
    Token      string `json:"token" desc:"用户登录凭据"`
    SessionKey string `json:"sessionKey" desc:"解密秘钥"`
  }

  Country struct {
    ID   int    `json:"id"`
    Code string `json:"code" desc:"国家编码"`
    Name string `json:"name" desc:"国家名称"`
  }
)

// Admin后台管理员登录
func AdminLogin(dto model.AdminUserDto) (AdminUser, error) {
  result := AdminUser{}
  user, err := model.GetAdminUser(dbCache, dto)
  if err != nil || user.ID == 0 {
    fmt.Println("err = ", err)
    return result, util.CodeBizError("601", "用户名密码错误或不存在")
  }
  if user.Status != 1 {
    return result, util.CodeBizError("602", "用户已经失效")
  }
  result.ID = user.ID
  result.Name = user.Name
  result.Avatar = user.Avatar
  result.Token = "admin" + util.GetRandomString(30, 3)
  result.ApiEnv = conf.ApiEnv
  userBytes, err := json.Marshal(result)
  redisCache.Set(result.Token, userBytes, time.Hour * 12)
  return result, err
}

// Admin后台管理员登出
func AdminLogout(token string) error {
  return redisCache.Del(token).Err()
}

// Admin后台管理员列表
func GetAdminUserList(dto model.AdminUserDto) ([]AdminUserList, error) {
  result := make([]AdminUserList, 0)
  list, err := model.GetAdminUserList(dbCache, dto)
  if err != nil {
    return result, fmt.Errorf("model.GetAdminUserList error %v ", err)
  }
  for _, v := range list {
    vo := AdminUserList{}
    vo.ID = v.ID
    vo.Name = v.Name
    vo.Avatar = v.Avatar
    vo.Status = v.Status
    vo.CTime = v.CTime
    result = append(result, vo)
  }
  return result, err
}

// 新增Admin后台管理员
func AddAdminUser(dto model.AdminUserDto) error {
  if model.ExistsAdminUser(dbCache, dto) {
    return fmt.Errorf("账号【%s】已经存在", dto.Account)
  }
  vo := model.NewAdminUser(dbCache)
  vo.ID = dto.ID
  vo.Account = dto.Account
  vo.Name = dto.Name
  vo.Avatar = dto.Avatar
  vo.Pwd = dto.Pwd
  vo.Status = 1
  vo.CTime = time.Now().Unix()
  return vo.Insert()
}

// 编辑Admin后台管理员
func EditAdminUser(dto model.AdminUserDto) error {
  vo := model.NewAdminUser(dbCache)
  vo.ID = dto.ID
  vo.Name = dto.Name
  vo.Avatar = dto.Avatar
  vo.Status = dto.Status
  if dto.Pwd != "" && vo.Pwd == dto.Pwd2 {
    vo.Pwd = dto.Pwd
  }
  return vo.Update()
}

// Admin后台检测用户Token是否有效
func CheckAdminLogin(token string) (*AdminUser, error) {
  userBytes, err := redisCache.Get(token).Bytes()
  if err != nil {
    return nil, util.CodeBizError("603", "用户未登录")
  }
  user := &AdminUser{}
  err = json.Unmarshal(userBytes, &user)
  return user, err
}

// APP登录
func Login(dto model.UserDto) (*User, error) {
  result := &User{}
  if dto.Account == "" || dto.Pwd == "" {
    return result, util.CodeBizError("606", "用户密码不能为空")
  }
  user, err := model.GetUserByAccountAndPassword(dbCache, dto.Account, dto.Pwd)
  if err != nil || user.ID == 0 {
    fmt.Println("get user error ", err)
    vo := model.NewUser(dbCache)
    vo.NickName = dto.NickName
    vo.Avatar = dto.Avatar
    vo.Gender = dto.Gender
    vo.CountryID = dto.CountryID
    err = vo.Insert()
    if err != nil {
      return result, err
    }
    user.ID = vo.ID
  }

  result.ID = user.ID
  result.Account = user.Account
  result.OpenID = user.OpenID
  result.NickName = user.NickName
  result.Avatar = user.Avatar
  result.Gender = user.Gender
  result.CountryID = user.CountryID
  result.Token = "user" + util.GetRandomString(30, 3)

  // 写入redis缓存
  redisCache.Set(result.Token, result.ID, time.Hour * 12)
  userBytes, err := json.Marshal(result)
  if err == nil {
    redisCache.HSet(keyUserInfo, strconv.Itoa(result.ID), userBytes)
  }
  return result, err
}

// APP注册
func Register(dto model.UserDto) error {
  vo := model.NewUser(dbCache)
  vo.Account = dto.Account
  if vo.ExistsByAccount() {
    return util.CodeBizError("605", "账号已经存在")
  }

  vo.Pwd = dto.Pwd
  vo.NickName = dto.NickName
  vo.Avatar = dto.Avatar
  vo.Gender = dto.Gender
  vo.CountryID = dto.CountryID
  err := vo.Insert()
  return err
}

// Token通过校验 返回用户详细信息
func VerifyUserToken(ctx echo.Context) (*User, error) {
  if err := util.CheckSign(ctx); err != nil {
    return nil, err
  }
  token := ctx.FormValue("token")
  var uid = redisCache.Get(token).Val()
  if uid == "" {
    return nil, util.CodeBizError("603", "用户未登录")
  }
  bytes, err := redisCache.HGet(keyUserInfo, uid).Bytes()
  if err != nil {
    return nil, util.CodeBizError("603", "用户未登录")
  }
  user := &User{}
  err = json.Unmarshal(bytes, &user)
  return user, err
}

// App端用户取uid
func ExtractUid(user *User) int {
  var uid = 0
  if user != nil {
    uid = user.ID
  }
  return uid
}

// Admin端用户取uid
func ExtractAdminUid(user *AdminUser) int {
  var uid = 0
  if user != nil {
    uid = user.ID
  }
  return uid
}

func GetCountry() ([]Country, error) {
  result := make([]Country, 0)
  data, err := model.GetCountry(dbCache)
  if err != nil {
    return result, err
  }
  for _, v:= range data {
    vo := Country{ID:v.ID, Code:v.Code, Name:v.Name}
    result = append(result, vo)
  }
  return result, nil
}

// 用户关注和取消关注
func SetUserFollow(dto model.UserFollowDto) error {
  return inTx(func(dbProxyBeginner squirrel.DBProxyBeginner) error {
    return model.SetUserFollow(dbProxyBeginner, dto)
  })
}
