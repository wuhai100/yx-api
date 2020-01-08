package model

import (
  "fmt"
  "github.com/Masterminds/squirrel"
  "github.com/Masterminds/structable"
  "yx-api/util"
)

type (
  AdminUser struct {
    ID      int    `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT"`
    Account string `stbl:"account" desc:"账号"`
    Pwd     string `stbl:"pwd" desc:"密码"`
    Name    string `stbl:"name" desc:"用户名"`
    Avatar  string `stbl:"avatar" desc:"头像"`
    CTime   int64  `stbl:"ctime" desc:"新增时间"`
    Status  int    `stbl:"status" desc:"状态"`
    rec     structable.Recorder
  }

  AdminUserDto struct {
    ID      int    `json:"id" desc:"ID"`
    Account string `json:"account" desc:"账号"`
    Name    string `json:"name" desc:"用户名称"`
    Pwd     string `json:"pwd" desc:"密码"`
    Status  int    `json:"status" desc:"状态"`
    Avatar  string `json:"avatar" desc:"头像"`
    Pwd2    string `json:"pwd2" desc:"确认密码"`
  }

  User struct {
    ID        int    `stbl:"id,PRIMARY_KEY,AUTO_INCREMENT"`
    Account   string `stbl:"account" desc:"账号"`
    Pwd       string `stbl:"pwd" desc:"密码"`
    OpenID    string `stbl:"open_id" desc:"微信OpenID"`
    NickName  string `stbl:"nick_name" desc:"用户昵称"`
    Avatar    string `stbl:"avatar" desc:"头像"`
    Gender    int    `stbl:"gender" desc:"性别"`
    CountryID int    `stbl:"country_id" desc:"国家"`
    rec       structable.Recorder
  }

  UserDto struct {
    ID        int    `json:"id"`
    Account   string `json:"account" desc:"账号"`
    Pwd       string `json:"pwd" desc:"密码"`
    JSCode    string `json:"jsCode" desc:"微信登录凭据"`
    OpenID    string `json:"openId" desc:"微信OpenID"`
    NickName  string `json:"nickName" desc:"用户昵称"`
    Avatar    string `json:"avatar" desc:"头像"`
    Gender    int    `json:"gender" desc:"性别"`
    CountryID int    `json:"countryId" desc:"国家"`
  }

  Country struct {
    ID   int    `stbl:"id"`
    Code string `stbl:"code" desc:"国家编码"`
    Name string `stbl:"name" desc:"国家名称"`
    rec  structable.Recorder
  }

  UserFollowDto struct {
    Option bool `json:"option" desc:"true:添加关注, false:取消关注"`
    Uid    int  `json:"uid" desc:"用户ID"`
    Fid    int  `json:"fid" desc:"关注用户ID"`
  }
)

func NewAdminUser(db squirrel.DBProxyBeginner) *AdminUser {
  vo := new(AdminUser)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableAdminUser, vo)
  return vo
}

func NewUser(db squirrel.DBProxyBeginner) *User {
  vo := new(User)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableUser, vo)
  return vo
}

func NewCountry(db squirrel.DBProxyBeginner) *Country {
  vo := new(Country)
  vo.rec = structable.New(db, conf.DriverName).Bind(tableCountry, vo)
  return vo
}

func (vo *AdminUser) Insert() error {
  return vo.rec.Insert()
}

func (vo *User) Insert() error {
  return vo.rec.Insert()
}

func (vo *User) LoadByAccount() error {
  return vo.rec.LoadWhere("account = ?", vo.Account)
}

func (vo *User) ExistsByAccount() bool {
  b, _ := vo.rec.ExistsWhere("account = ?", vo.Account)
  return b
}

func (vo *AdminUser) Update() error {
  if vo.ID == 0 {
    return fmt.Errorf("ID 不能为：0")
  }
  update := squirrel.Update(tableAdminUser).Where("id = ?", vo.ID)
  if vo.Name != "" {
    update = update.Set("name", vo.Name)
  }
  if vo.Pwd != "" {
    update = update.Set("pwd", vo.Pwd)
  }
  if vo.Avatar != "" {
    update = update.Set("avatar", vo.Avatar)
  }
  if vo.Status > 0 {
    update = update.Set("status", vo.Status)
  }
  fmt.Println(update.ToSql())
  _, err := update.RunWith(vo.rec.DB()).Exec()
  return err
}

func ExistsAdminUser(db squirrel.DBProxyBeginner, filter AdminUserDto) bool {
  log := userLog.With("func", "ExistsAdminUser")
  query := squirrel.Select("count(1)").From(tableAdminUser).Where("account=?", filter.Account)
  util.PrintQuery(log, query)
  count := 0
  err := query.RunWith(db).QueryRow().Scan(&count)
  if err != nil {
    log.Errorf("ExistsAdminUser error %v ", err)
  }
  return count > 0
}

func GetAdminUser(db squirrel.DBProxyBeginner, filter AdminUserDto) (*AdminUser, error) {
  log := userLog.With("func", "GetAdminUser")
  column := NewAdminUser(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableAdminUser).Where("account=? and pwd=?", filter.Account, filter.Pwd)
  util.PrintQuery(log, query)
  err := query.RunWith(db).QueryRow().Scan(column.rec.FieldReferences(true)...)
  return column, err
}

func GetAdminUserList(db squirrel.DBProxyBeginner, filter AdminUserDto) ([]AdminUser, error) {
  log := userLog.With("func", "GetAdminUserList")
  result := make([]AdminUser, 0)
  column := NewAdminUser(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableAdminUser)
  if filter.ID > 0 {
    query = query.Where("id=?", filter.ID)
  }
  if filter.Account != "" {
    query = query.Where("account=?", filter.Account)
  }
  if filter.Name != "" {
    query = query.Where("name=?", filter.Name)
  }
  if filter.Status > 0 {
    query = query.Where("status=?", filter.Status)
  }
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    err := rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return result, err
    }
    result = append(result, *column)
  }
  return result, err
}

func GetUserByAccountAndPassword(db squirrel.DBProxyBeginner, account, pwd string) (*User, error) {
  log := userLog.With("func", "GetUser")
  column := NewUser(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableUser).Where("account = ? and pwd = ?", account, pwd)
  util.PrintQuery(log, query)
  err := query.RunWith(db).QueryRow().Scan(column.rec.FieldReferences(true)...)
  return column, err
}

func GetUser(db squirrel.DBProxyBeginner, filter UserDto) (*User, error) {
  log := userLog.With("func", "GetUser")
  column := NewUser(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableUser)
  if filter.ID > 0 {
    query = query.Where("id = ?", filter.ID)
  }
  if filter.OpenID != "" {
    query = query.Where("open_id = ?", filter.OpenID)
  }
  if filter.NickName != "" {
    query = query.Where("nick_name = ?", filter.NickName)
  }
  util.PrintQuery(log, query)
  err := query.RunWith(db).QueryRow().Scan(column.rec.FieldReferences(true)...)
  return column, err
}

func GetUserMap(db squirrel.DBProxyBeginner, ids []int) (map[int]User, error) {
  log := userLog.With("func", "GetUserMap")
  result := map[int]User{}
  if len(ids) == 0 {
    return result, log.Errorf("ids length is 0")
  }
  column := NewUser(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableUser).Where(squirrel.Eq{"id": ids})
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, log.Errorf("query error %v", err)
  }
  defer rows.Close()

  for rows.Next() {
    err := rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return result, log.Errorf("scan error %v", err)
    }
    result[column.ID] = *column
  }
  return result, err
}

func GetAdminUserNameMap(db squirrel.DBProxyBeginner, ids []int) (map[int]string, error) {
  log := userLog.With("func", "GetAdminUserNameMap")
  result := map[int]string{}
  if len(ids) == 0 {
    return result, nil
  }
  query := squirrel.Select("id, name").From(tableAdminUser).Where(squirrel.Eq{"id": ids})
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, log.Errorf("query error %v", err)
  }
  defer rows.Close()

  for rows.Next() {
    id, name := 0, ""
    err := rows.Scan(&id, &name)
    if err != nil {
      return result, log.Errorf("scan error %v", err)
    }
    result[id] = name
  }
  return result, err
}

func GetCountry(db squirrel.DBProxyBeginner) ([]Country, error) {
  log := userLog.With("func", "GetCountry")
  result := make([]Country, 0)
  column := NewCountry(db)
  query := squirrel.Select(column.rec.Columns(true)...).From(tableCountry).Where("status=1")
  util.PrintQuery(log, query)
  rows, err := query.RunWith(db).Query()
  if err != nil {
    return result, err
  }
  defer rows.Close()

  for rows.Next() {
    err := rows.Scan(column.rec.FieldReferences(true)...)
    if err != nil {
      return result, err
    }
    result = append(result, *column)
  }
  return result, err
}

// 用户关注和取消关注
func SetUserFollow(db squirrel.DBProxyBeginner, dto UserFollowDto) error {
  log := userLog.With("func", "SetUserFollow")
  if dto.Option {
    insert := squirrel.Insert(tableUserFollow).Columns("uid, fid").Values(dto.Uid, dto.Fid)
    util.PrintInsert(log, insert)
    _, err := insert.RunWith(db).Exec()
    if err != nil {
      return err
    }
    exec := squirrel.Update(tableUser).Set("followers", squirrel.Expr("followers + 1")).Where("id = ?", dto.Uid)
    util.PrintUpdate(log, exec)
    _, err = exec.RunWith(db).Exec()
    if err != nil {
      return err
    }
    exec = squirrel.Update(tableUser).Set("following", squirrel.Expr("following + 1")).Where("id = ?", dto.Fid)
    _, err = exec.RunWith(db).Exec()
    util.PrintUpdate(log, exec)
    return err
  }
  del := squirrel.Delete(tableUserFollow).Where("uid = ? AND fid = ?", dto.Uid, dto.Fid)
  util.PrintDelete(log, del)
  _, err := del.RunWith(db).Exec()
  if err != nil {
    return err
  }
  exec := squirrel.Update(tableUser).Set("followers", squirrel.Expr("followers - 1")).Where("id = ?", dto.Uid)
  util.PrintUpdate(log, exec)
  _, err = exec.RunWith(db).Exec()
  if err != nil {
    return err
  }
  exec = squirrel.Update(tableUser).Set("following", squirrel.Expr("following - 1")).Where("id = ?", dto.Fid)
  util.PrintUpdate(log, exec)
  _, err = exec.RunWith(db).Exec()
  return err
}
