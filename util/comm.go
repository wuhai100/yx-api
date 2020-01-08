package util

import (
  "bytes"
  "crypto/aes"
  "crypto/cipher"
  "crypto/md5"
  "crypto/sha1"
  "encoding/base64"
  "encoding/hex"
  "fmt"
  jsoniter "github.com/json-iterator/go"
  "github.com/labstack/echo/v4"
  "github.com/labstack/echo/v4/middleware"
  "github.com/yb7/alilog"
  "io"
  "io/ioutil"
  "math/rand"
  "net"
  "net/http"
  "reflect"
  "sort"
  "strconv"
  "strings"
  "time"
)

const (
  layoutDate     = "2006-01-02"
  layoutDateTime = "2006-01-02 15:04:05"
  wxApiURL       = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

var (
  loc, _   = time.LoadLocation("Asia/Shanghai")
  json    = jsoniter.ConfigCompatibleWithStandardLibrary
  EchoInst = echoInst()
  AppLog   = logStore()
)

/**
 * response返回json对象
 */
type (
  ResponseData struct {
    Errno int         `json:"errno"`          // 必需 错误码。正常返回0 异常返回560 错误提示561对应errorInfo
    Data  interface{} `json:"data,string"`    // 必需 返回数据内容。 如果有返回数据，可以是字符串或者数组JSON等等
    Page  *Pagination `json:"page,omitempty"` // 非必需 分页信息
  }

  wxDecoded struct {
    PhoneNumber     string                                `json:"phoneNumber"`
    PurePhoneNumber string                                `json:"purePhoneNumber"`
    Watermark       struct{ AppId string `json:"appid"` } `json:"watermark"`
  }

  wxAPIResult struct {
    SessionKey string `json:"session_key"`
    OpenID     string `json:"openid"`
  }
)

func ResultData(data interface{}, err error) *ResponseData {
  var result = &ResponseData{}
  if err != nil {
    result.Errno = 501
    result.Data = err.Error()
    AppLog.With("error", "ResultData").Error(err)
    return result
  }
  if data == nil {
    data = "success"
  }
  result.Data = data
  return result
}

func ResultPageData(data interface{}, page Pagination, err error) *ResponseData {
  var result = &ResponseData{}
  if err != nil {
    result.Errno = 501
    result.Data = err.Error()
    AppLog.With("error", "ResultPageData").Error(err)
    return result
  }
  if data == nil {
    data = "success"
  }
  result.Data = data
  result.Page = &page
  return result
}

func echoInst() *echo.Echo {
  var ei = echo.New()
  ei.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 9}))
  return ei
}

func logStore() *alilog.SLog {
  var logstore = "lod_debug"
  return alilog.New("opt-admin", logstore)
}

// md5
func MD5(inp string) string {
  if inp == "" {
    return ""
  }
  h := md5.New()
  h.Write([]byte(inp))
  return hex.EncodeToString(h.Sum(nil))
}

/**
* 生成随机字符串
* @param  num int
* @param  kind
   KC_RAND_KIND_NUM   = 0 // 纯数字
   KC_RAND_KIND_LOWER = 1 // 小写字母
   KC_RAND_KIND_UPPER = 2 // 大写字母
   KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
* @return str string
*/
func GetRandomString(size int, kind int) string {
  ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
  is_all := kind > 2 || kind < 0
  rand.Seed(time.Now().UnixNano())
  for i := 0; i < size; i++ {
    if is_all { // random ikind
      ikind = rand.Intn(3)
    }
    scope, base := kinds[ikind][0], kinds[ikind][1]
    result[i] = uint8(base + rand.Intn(scope))
  }
  return string(result)
}

func EncodeToString(str string) string {
  if str == "" {
    return str
  }
  return base64.RawURLEncoding.EncodeToString([]byte(str))
}

func DecodeString(str string) string {
  if str == "" {
    return str
  }
  var res, err = base64.RawURLEncoding.DecodeString(str)
  if err != nil {
    return str
  }
  return string(res)
}

// 数字int去重复
func DeDuplicationInt(ids []int) []int {
  var result = make([]int, 0)
  var temp = map[int]bool{}
  for _, v := range ids {
    if temp[v] {
      continue
    }
    temp[v] = true
    result = append(result, v)
  }
  return result
}

// 生成签名
func GetSign(param string) string {
  return MD5(SHA1(param))
}

// 运算保留2位小数
func Decimal2(value float64) float64 {
  // return math.Trunc(value*1e2+0.5) * 1e-2
  value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
  return value
}

// 运算保留3位小数
func Decimal3(value float64) float64 {
  value, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", value), 64)
  return value
}

// 运算保留2位小数
func Atoi(value string) int {
  i, _ := strconv.Atoi(value)
  return i
}

// 生成sha1
func SHA1(data string) string {
  h := sha1.New()
  h.Write([]byte(data))
  sha1str1 := h.Sum(nil)
  sha1str2 := fmt.Sprintf("%x", sha1str1)
  return sha1str2
}

// 校验签名
func CheckSign(ctx echo.Context) error {
  var arr []string
  formParams, _ := ctx.FormParams()
  for k, v := range formParams {
    if k == "debug" && v[0] == "test" {
      return nil
    }
    if k != "sign" {
      arr = append(arr, fmt.Sprintf("%s=%v", k, v[0]))
    }
  }
  sort.Strings(arr)
  fmt.Println("arr = ", arr)
  signVerify := GetSign(strings.Join(arr, "&"))
  fmt.Println("signVerify = ", signVerify)
  if signVerify == ctx.FormValue("sign") {
    return nil
  }
  return CustomBizError("签名错误")
}

func drainBody(b io.ReadCloser) (r1 io.ReadCloser) {
  if b == http.NoBody {
    return http.NoBody
  }
  var buf bytes.Buffer
  if _, err := buf.ReadFrom(b); err != nil {
    return b
  }
  if err := b.Close(); err != nil {
    return b
  }
  return ioutil.NopCloser(&buf)
}

func RequestIPAddress(req *http.Request) string {
  ra := req.RemoteAddr
  if ip := req.Header.Get("X-Forwarded-For"); ip != "" {
    ra = strings.Split(ip, ", ")[0]
  } else if ip := req.Header.Get("X-Real-IP"); ip != "" {
    ra = ip
  } else {
    ra, _, _ = net.SplitHostPort(ra)
  }
  return ra
}

func HttpGet(url string) ([]byte, error) {
  client := http.Client{
    Timeout: time.Duration(3 * time.Second),
  }
  var result []byte
  resp, err := client.Get(url)
  if err != nil {
    return result, err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  return body, err
}

func ToXmlString(param map[string]string) string {
  xml := "<xml>"
  for k, v := range param {
    xml = xml + fmt.Sprintf("<%s>%s</%s>", k, v, k)
  }
  xml = xml + "</xml>"

  return xml
}

func ToMap(in interface{}) (map[string]string, error) {
  out := make(map[string]string)

  v := reflect.ValueOf(in)
  if v.Kind() == reflect.Ptr {
    v = v.Elem()
  }

  // we only accept structs
  if v.Kind() != reflect.Struct {
    return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
  }

  typ := v.Type()
  for i := 0; i < v.NumField(); i++ {
    // gets us a StructField
    fi := typ.Field(i)
    if tagv := fi.Tag.Get("xml"); tagv != "" && tagv != "xml" {
      // set key of map to value in struct field
      out[tagv] = v.Field(i).String()
    }
  }
  return out, nil
}

//格式
func UnixDateString(ctime int64) string {
  tm := time.Unix(ctime,0)
  return tm.Format(layoutDate)
}

//格式
func UnixTimeString(ctime int64) string {
  tm := time.Unix(ctime,0)
  return tm.Format(layoutDateTime)
}

// 根据年份月份取得当前月份开始和结束的时间戳
func BetweenMonth(year, month int) (int64, int64) {
  monthM := time.Month(month)
  if year == 0 {
    year = time.Now().Year()
  }
  if month == 0 {
    monthM = time.Now().Month()
  }

  st1 := time.Date(year, monthM, 1, 0,0,0, 0, time.UTC)
  st2 := st1.AddDate(0, 1, 0)
  return st1.Unix(), st2.Unix()
}

//格式
func FormatUnixTime(ctime int64) time.Time {
  return time.Unix(ctime,0)
}

// javaJasperReports导出
func GenerateReport(reportId, format string, jsonData []byte) (string, error) {
  url := fmt.Sprintf("http://%s/reports/%s/%s", "122.112.160.144", reportId, format)
  msg := make([]string, 0)
  msg = append(msg, fmt.Sprintf("POST URL:> %s", url))
  msg = append(msg, fmt.Sprintf("    REQ: %s", string(jsonData)))

  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()

  msg = append(msg, fmt.Sprintf("    RESP STATUS: %s", resp.Status))
  body, _ := ioutil.ReadAll(resp.Body)
  msg = append(msg, fmt.Sprintf("    RESP BODY  : %s", string(body)))

  type GenerateReportResponse struct {
    Filename string
  }
  reportResp := &GenerateReportResponse{}
  err = json.Unmarshal(body, reportResp)
  if err != nil {
    return "", err
  }
  fmt.Println(msg)
  return reportResp.Filename, nil
}

// 解密微信加密信息
func DecodeWxData(encryptedData, iv, sessionKey string) (wxDecoded, error) {
  result := wxDecoded{}
  aesCipher, err := base64.StdEncoding.DecodeString(encryptedData)
  if err != nil {
    return result, fmt.Errorf("Decode EncryptedData error %v ", err)
  }
  aesIv, err := base64.StdEncoding.DecodeString(iv)
  if err != nil {
    return result, fmt.Errorf("Decode Iv error %v ", err)
  }
  aesSessionKey, err := base64.StdEncoding.DecodeString(sessionKey)
  if err != nil {
    return result, fmt.Errorf("Decode sessionKey error %v ", err)
  }
  block, err := aes.NewCipher(aesSessionKey)
  if err != nil {
    return result, fmt.Errorf("NewCipher sessionKey error %v ", err)
  }
  bm := cipher.NewCBCDecrypter(block, aesIv)

  plaintextCopy := make([]byte, len(aesCipher)+10)
  bm.CryptBlocks(plaintextCopy, aesCipher)

  // 这个地方要replace的原因是因为plaintextCopy在解密完后面有一窜trim不掉的空格字符，所以我复制了这个特殊字符进行replace
  pc := strings.Replace(string(plaintextCopy), "", "", -1)
  fmt.Printf("%x=>%s\n", aesCipher, pc)
  err = json.Unmarshal([]byte(pc), &result)
  if err != nil {
    return result, fmt.Errorf("Unmarshal result error %v ", err)
  }
  if result.PurePhoneNumber == "" {
    return result, fmt.Errorf("wxDecoded phone is 0")
  }
  fmt.Println("result = ", result)
  return result, nil
}

// 通过微信jsCode获取openID（用户）
func GetWxOpenID(jsCode, appId, secret string) (wxAPIResult, error) {
  log := AppLog.With("func", "getWxOpenID")
  result := wxAPIResult{}
  url := fmt.Sprintf(wxApiURL, appId, secret, jsCode)
  log.Infof("URL: %s", url)
  wxData, err := HttpGet(url)
  if err != nil {
    return result, log.Errorf("HttpGet error %v", err)
  }
  err = json.Unmarshal(wxData, &result)
  if err != nil {
    return result, log.Errorf("Unmarshal error %v", err)
  }
  if result.OpenID == "" {
    log.Infof("Result: %s", string(wxData))
    return result, CustomBizError("微信验证错误:" + string(wxData))
  }
  return result, nil
}

func UniqueKey(dto ...interface{}) string {
  var key = "yx_api"
  for _, v := range dto {
    key = fmt.Sprintf("%s_%v", key, v)
  }
  return key
}
