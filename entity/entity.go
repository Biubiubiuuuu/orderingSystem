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
	Ids []int64 `json:"ids"` // ids
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

// 商家门店信息更新结构体
type BusinessStoreRequest struct {
	StoreName              string                             `json:"store_name"`                // 门店名称
	StoreAddress           string                             `json:"store_address"`             // 门店详细地址
	StoreLogo              string                             `json:"store_logo"`                // 门店logo
	StoreContactName       string                             `json:"store_contact_name"`        // 门店联系人姓名
	StoreContactTel        string                             `json:"store_contact_tel"`         // 门店联系人电话
	StoreStartBankingHours string                             `json:"store_start_banking_hours"` // 门店开始营业时间
	StoreEndBankingHours   string                             `json:"store_end_banking_hours"`   // 门店结束营业时间
	StoreFacePhoto         string                             `json:"store_face_photo"`          // 门脸照
	InStorePhotos          []BusinessStoreRequestInStorePhoto `json:"in_store_photos"`           // 店内照
}

type BusinessStoreRequestInStorePhoto struct {
	Url     string `json:"Url"`      // 图片保存地址
	StoreID int64  `json:"store_id"` // 门店ID
}

// 商品种类请求结构体
type GoodsTypeRequest struct {
	Name         string `json:"name"`           // 分类名称
	TypeSort     int64  `json:"type_sort"`      // 分类排序
	DisplayOrNot bool   `json:"display_or_not"` // 是否显示
}

// 商品请求结构体
type GoodsRequest struct {
	GoodsName        string  `json:"goods_name"`        // 商品名称
	GoodsPhoto       string  `json:"goods_photo"`       // 商品图片
	GoodsDescription string  `json:"goods_description"` // 商品描述
	GoodsListing     bool    `json:"goods_listing"`     // 是否上架
	GoodsPrice       float64 `json:"goods_price"`       // 商品价格
	GoodsUnit        string  `json:"goods_unit"`        // 商品单位 份、杯
	GoodsSort        int64   `json:"goods_sort"`        // 商品排序
	GoodsTypeID      int64   `json:"goods_type_id"`     // 商品种类ID
}

// 商品种类结构体
type GoodsTypeResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"` // 分类名称
}
