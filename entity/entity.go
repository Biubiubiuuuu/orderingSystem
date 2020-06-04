package entity

// http请求响应数据data
type ResponseData struct {
	Status  bool                   `json:"status"`  // 成功失败标志；true：成功 、false：失败
	Data    map[string]interface{} `json:"data"`    // 返回数据
	Message string                 `json:"message"` // 提示信息
}

// 系统管理员登录请求结构体
type SystemAdminLoginRequest struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

// 系统管理员添加请求结构体
type SystemAdminAddRequest struct {
	Username string `json:"username"`  // 用户名
	Nikename string `json:"nikename"`  // 昵称
	Password string `json:"password"`  // 密码
	Manager  string `json:"manager"`   // 操作权限 Y | N
	Avatar   string `json:"avatar"`    // 头像
	IsEnable bool   `json:"is_enable"` // 是否启用
}

// 系统管理员修改密码请求结构体
type SystemAdminUpdatePassRequest struct {
	OldPassword string `json:"old_password"` // 旧密码
	NewPassword string `json:"new_password"` // 新密码
}

// 删除ids
type DeleteIds struct {
	Ids []string `json:"ids"` // ids
}

// 商家注册请求、手机验证码登录结构体
type BusinessLoginOrRegisterRequest struct {
	Tel  string `json:"tel"`  // 手机号码
	Code string `json:"code"` // 验证码
}

// 商家账号密码登录请求结构体
type BusinessPassLoginRequest struct {
	Tel      string `json:"tel"`      // 手机号码
	Password string `json:"password"` // 密码
}
