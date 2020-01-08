
### 后台结构说明
#### 系统接口层级
* controller API接口RESTController
* service 服务逻辑处理层
* model 数据库操作处理层
* util 通用工具类，包括一些外部SDK包
* deploy 用户k8s+docker部署的文件
* swagger 接口调试工具，访问URL:http://localhost:1111/swagger/index.html

#### 系统工具和配置
* 统一错误处理 /controller/error.go
* 自定义错误处理 /util/bizerror.go
* 系统初始配置 /util/setting.go
* 数据库工具包 /util/db_ext.go
* 系统自定义工具 /util/comm.go
* 华为obs存储SDK /util/obs/*

#### Token类型有2种，针对2种类型用户
* 管理后台用户token (verifyAdminUser)
* App用户token (verifyUserToken)

### ResultData公共字段说明
| 字段|类型 |选项|说明 |
|---------|---------|---------|---------|
| errno | int | 必需| 错误码。正常返回0 异常返回相应错误代码 |
| data | interface{} | 必需 | 返回数据内容。 如果有返回数据，可以是字符串或者数组JSON等等 |
| page | Pagination | 可选 | 返回有data集合的总分页数量。后台分页显示才会有此结果 |
#### 注：如果接口本身没有数据需要返回，则对象为空，默认errno,data是必须会在JSON对象中存在的。如果返回其他格式(XML，File)需要自定义处理。

### 接口
#### 管理后端用户登录和信息Admin版本
* 用户登录 POST /admin/login
* 用户登出 POST /admin/logout

#### App用户登录和信息v1版本
* 用户登录 POST /v1/login
