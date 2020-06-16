package businessController

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Biubiubiuuuu/orderingSystem/entity"
	"github.com/Biubiubiuuuu/orderingSystem/helper/configHelper"
	"github.com/Biubiubiuuuu/orderingSystem/helper/fileHelper"
	"github.com/Biubiubiuuuu/orderingSystem/service/businessService"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary 商家注册
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.BusinessLoginOrRegisterRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/register [POST]
func Register(c *gin.Context) {
	req := entity.BusinessLoginOrRegisterRequest{}
	res := entity.ResponseData{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.Register(req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 商家手机验证码登录
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.BusinessLoginOrRegisterRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/codelogin [POST]
func CodeLogin(c *gin.Context) {
	req := entity.BusinessLoginOrRegisterRequest{}
	res := entity.ResponseData{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.CodeLogin(req, c.ClientIP())
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 商家账号密码登录
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.BusinessPassLoginRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/passlogin [POST]
func PassLogin(c *gin.Context) {
	req := entity.BusinessPassLoginRequest{}
	res := entity.ResponseData{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.PassLogin(req, c.ClientIP())
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 查询商家门店信息
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/store [GET]
// @Security ApiKeyAuth
func QueryBusinessStoreInfo(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	res = businessService.QueryBusinessStoreInfo(token)
	c.JSON(http.StatusOK, res)
}

// @Summary 更新商家门店信息
// @tags 商家
// @Accept  multipart/form-data
// @Produce  json
// @Param store_name formData string true "门店名称"
// @Param store_address formData string true "门店详细地址"
// @Param store_logo formData file false "门店logo"
// @Param store_contact_name formData string false "门店联系人姓名"
// @Param store_contact_tel formData string false "门店联系人电话"
// @Param store_start_banking_hours formData string true "门店开始营业时间"
// @Param store_end_banking_hours formData string true "门店结束营业时间"
// @Param store_face_photo formData file false "门脸照"
// @Param in_store_photos formData file false "店内照"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/store [PUT]
// @Security ApiKeyAuth
func UpdateBusinessStoreInfo(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	// 获取主机头
	r := c.Request
	host := r.Host
	if strings.HasPrefix(host, "http://") == false {
		host = "http://" + host
	}
	//店铺logo
	var store_logo string
	if file, err := c.FormFile("store_logo"); err == nil {
		// 文件名 避免重复取uuid
		var filename string
		uuid, _ := uuid.NewUUID()
		arr := strings.Split(file.Filename, ".")
		if strings.EqualFold(arr[1], "png") {
			filename = uuid.String() + ".png"
		} else {
			filename = uuid.String() + ".jpg"
		}
		pathFile := configHelper.ImageDir
		if !fileHelper.IsExist(pathFile) {
			fileHelper.CreateDir(pathFile)
		}
		pathFile = pathFile + filename
		if err := c.SaveUploadedFile(file, pathFile); err == nil {
			store_logo = host + "/" + pathFile
		}
	}
	//门脸照
	var store_face_photo string
	if file, err := c.FormFile("store_face_photo"); err == nil {
		// 文件名 避免重复取uuid
		var filename string
		uuid, _ := uuid.NewUUID()
		arr := strings.Split(file.Filename, ".")
		if strings.EqualFold(arr[1], "png") {
			filename = uuid.String() + ".png"
		} else {
			filename = uuid.String() + ".jpg"
		}
		pathFile := configHelper.ImageDir
		if !fileHelper.IsExist(pathFile) {
			fileHelper.CreateDir(pathFile)
		}
		pathFile = pathFile + filename
		if err := c.SaveUploadedFile(file, pathFile); err == nil {
			store_face_photo = host + "/" + pathFile
		}
	}
	// 店内照
	var inStorePhotos []entity.BusinessStoreRequestInStorePhoto
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["in_store_photos"]
		for _, file := range files {
			// 文件名 避免重复取uuid
			var filename string
			uuid, _ := uuid.NewUUID()
			arr := strings.Split(file.Filename, ".")
			if strings.EqualFold(arr[1], "png") {
				filename = uuid.String() + ".png"
			} else {
				filename = uuid.String() + ".jpg"
			}
			pathFile := configHelper.ImageDir
			if !fileHelper.IsExist(pathFile) {
				fileHelper.CreateDir(pathFile)
			}
			pathFile = pathFile + filename

			if err := c.SaveUploadedFile(file, pathFile); err == nil {
				photo := entity.BusinessStoreRequestInStorePhoto{
					Url: host + "/" + pathFile,
				}
				inStorePhotos = append(inStorePhotos, photo)
			}
		}
	}
	req := entity.BusinessStoreRequest{
		StoreName:              c.PostForm("store_name"),
		StoreAddress:           c.PostForm("store_address"),
		StoreLogo:              store_logo,
		StoreContactName:       c.PostForm("store_contact_name"),
		StoreContactTel:        c.PostForm("store_contact_tel"),
		StoreStartBankingHours: c.PostForm("store_start_banking_hours"),
		StoreEndBankingHours:   c.PostForm("store_end_banking_hours"),
		StoreFacePhoto:         store_face_photo,
		InStorePhotos:          inStorePhotos,
	}
	res = businessService.UpdateBusinessStoreInfo(token, req)
	c.JSON(http.StatusOK, res)
}

// @Summary 商品种类添加
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.GoodsTypeRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goodstype [POST]
// @Security ApiKeyAuth
func AddGoodsType(c *gin.Context) {
	req := entity.GoodsTypeRequest{}
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.AddGoodsType(token, req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 商品种类修改
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "商品种类ID"
// @Param body body entity.GoodsTypeRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goodstype/{id} [PUT]
// @Security ApiKeyAuth
func UpdateGoodsType(c *gin.Context) {
	req := entity.GoodsTypeRequest{}
	res := entity.ResponseData{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.UpdateGoodsType(token, id, req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 商品种类删除
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "商品种类ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goodstype/{id} [DELETE]
// @Security ApiKeyAuth
func DeleteGoodsType(c *gin.Context) {
	req := entity.DeleteIds{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ids := append(req.Ids, id)
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	res = businessService.DeleteGoodsType(token, ids)
	c.JSON(http.StatusOK, res)
}

// @Summary 查询商品种类By ID
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "商品种类ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goodstype/{id} [GET]
// @Security ApiKeyAuth
func QueryGoodsTypeByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res := businessService.QueryGoodsTypeByID(id)
	c.JSON(http.StatusOK, res)
}

// @Summary 分页查询商家商品种类
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param pageSize query string false "页大小"
// @Param page query string false "页数"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goodstype [GET]
// @Security ApiKeyAuth
func QueryGoodsType(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	res = businessService.QueryGoodsType(token, pageSize, page)
	c.JSON(http.StatusOK, res)
}

// @Summary 查询商家商品种类ID和名称
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goodstypes [GET]
// @Security ApiKeyAuth
func QueryGoodsTypeIDAndName(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	res = businessService.QueryGoodsTypeIDAndNameByAdminID(token)
	c.JSON(http.StatusOK, res)
}

// @Summary 添加商品
// @tags 商家
// @Accept  multipart/form-data
// @Produce  json
// @Param goods_name formData string true "商品名称"
// @Param goods_photo formData file string "商品图片"
// @Param goods_description formData string false "商品描述"
// @Param goods_listing formData bool false "是否上架"
// @Param goods_price formData float64 false "商品价格"
// @Param goods_unit formData string true "商品单位 份、杯、瓶"
// @Param goods_sort formData string true "排序"
// @Param goods_type_id formData int false "商品种类ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goods [POST]
// @Security ApiKeyAuth
func AddGoods(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	// 获取主机头
	r := c.Request
	host := r.Host
	if strings.HasPrefix(host, "http://") == false {
		host = "http://" + host
	}
	// 商品图片
	var goods_photo string
	if file, err := c.FormFile("goods_photo"); err == nil {
		// 文件名 避免重复取uuid
		var filename string
		uuid, _ := uuid.NewUUID()
		arr := strings.Split(file.Filename, ".")
		if strings.EqualFold(arr[1], "png") {
			filename = uuid.String() + ".png"
		} else {
			filename = uuid.String() + ".jpg"
		}
		pathFile := configHelper.ImageDir
		if !fileHelper.IsExist(pathFile) {
			fileHelper.CreateDir(pathFile)
		}
		pathFile = pathFile + filename
		if err := c.SaveUploadedFile(file, pathFile); err == nil {
			goods_photo = host + "/" + pathFile
		}
	}
	goods_listing, _ := strconv.ParseBool(c.DefaultPostForm("goods_listing", "false"))
	goods_price, _ := strconv.ParseFloat(c.PostForm("goods_price"), 64)
	goods_sort, _ := strconv.ParseInt(c.PostForm("goods_sort"), 10, 64)
	goods_type_id, _ := strconv.ParseInt(c.PostForm("goods_type_id"), 10, 64)
	req := entity.GoodsRequest{
		GoodsName:        c.PostForm("goods_name"),
		GoodsPhoto:       goods_photo,
		GoodsDescription: c.PostForm("goods_description"),
		GoodsListing:     goods_listing,
		GoodsPrice:       goods_price,
		GoodsUnit:        c.PostForm("goods_unit"),
		GoodsSort:        goods_sort,
		GoodsTypeID:      goods_type_id,
	}
	res = businessService.AddGoods(token, req)
	c.JSON(http.StatusOK, res)
}

// @Summary 修改商品
// @tags 商家
// @Accept  multipart/form-data
// @Produce  json
// @Param id path int true "商品ID"
// @Param goods_name formData string true "商品名称"
// @Param goods_photo formData file string "商品图片"
// @Param goods_description formData string false "商品描述"
// @Param goods_listing formData bool false "是否上架"
// @Param goods_price formData float64 false "商品价格"
// @Param goods_unit formData string true "商品单位 份、杯、瓶"
// @Param goods_sort formData string true "排序"
// @Param goods_type_id formData int false "商品种类ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goods/{id} [PUT]
// @Security ApiKeyAuth
func UpdateGoods(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	// 获取主机头
	r := c.Request
	host := r.Host
	if strings.HasPrefix(host, "http://") == false {
		host = "http://" + host
	}
	// 商品图片
	var goods_photo string
	if file, err := c.FormFile("goods_photo"); err == nil {
		// 文件名 避免重复取uuid
		var filename string
		uuid, _ := uuid.NewUUID()
		arr := strings.Split(file.Filename, ".")
		if strings.EqualFold(arr[1], "png") {
			filename = uuid.String() + ".png"
		} else {
			filename = uuid.String() + ".jpg"
		}
		pathFile := configHelper.ImageDir
		if !fileHelper.IsExist(pathFile) {
			fileHelper.CreateDir(pathFile)
		}
		pathFile = pathFile + filename
		if err := c.SaveUploadedFile(file, pathFile); err == nil {
			goods_photo = host + "/" + pathFile
		}
	}
	goods_listing, _ := strconv.ParseBool(c.DefaultPostForm("goods_listing", "false"))
	goods_price, _ := strconv.ParseFloat(c.PostForm("goods_price"), 64)
	goods_sort, _ := strconv.ParseInt(c.PostForm("goods_sort"), 10, 64)
	goods_type_id, _ := strconv.ParseInt(c.PostForm("goods_type_id"), 10, 64)
	req := entity.GoodsRequest{
		GoodsName:        c.PostForm("goods_name"),
		GoodsPhoto:       goods_photo,
		GoodsDescription: c.PostForm("goods_description"),
		GoodsListing:     goods_listing,
		GoodsPrice:       goods_price,
		GoodsUnit:        c.PostForm("goods_unit"),
		GoodsSort:        goods_sort,
		GoodsTypeID:      goods_type_id,
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res = businessService.UpdateGoods(token, id, req)
	c.JSON(http.StatusOK, res)
}

// @Summary 商品删除
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "商品ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goods/{id} [DELETE]
// @Security ApiKeyAuth
func DeleteGoods(c *gin.Context) {
	req := entity.DeleteIds{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ids := append(req.Ids, id)
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	res = businessService.DeleteGoods(token, ids)
	c.JSON(http.StatusOK, res)
}

// @Summary 查询商品By ID
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "商品ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goods/{id} [GET]
// @Security ApiKeyAuth
func QueryGoodsByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res := businessService.QueryGoodsByID(id)
	c.JSON(http.StatusOK, res)
}

// @Summary 分页查询商家商品
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param pageSize query string false "页大小"
// @Param page query string false "页数"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goods [GET]
// @Security ApiKeyAuth
func QueryGoods(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	res = businessService.QueryGoods(token, pageSize, page)
	c.JSON(http.StatusOK, res)
}

// @Summary 上架/下架商品
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "商品ID"
// @Param downorup path bool false "是否上架"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/goods/{id}/{downorup} [PUT]
// @Security ApiKeyAuth
func DownOrUpGoods(c *gin.Context) {
	res := entity.ResponseData{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	downOrup, _ := strconv.ParseBool(c.Param("downorup"))
	res = businessService.DownOrUpGoods(id, downOrup)
	c.JSON(http.StatusOK, res)
}

// @Summary 添加餐桌种类
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.TableTypeRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/tabletype [POST]
// @Security ApiKeyAuth
func AddTableType(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	req := entity.TableTypeRequest{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.AddTableType(token, req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 修改餐桌种类
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "餐桌种类ID"
// @Param body body entity.TableTypeRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/tabletype/{id} [PUT]
// @Security ApiKeyAuth
func UpdateTableType(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	req := entity.TableTypeRequest{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.UpdateTableType(token, id, req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 删除餐桌种类
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "餐桌种类ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/tabletype/{id} [DELETE]
// @Security ApiKeyAuth
func DeleteTableType(c *gin.Context) {
	req := entity.DeleteIds{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ids := append(req.Ids, id)
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	res = businessService.DeleteTableType(token, ids)
	c.JSON(http.StatusOK, res)
}

// @Summary 查询餐桌种类ID和名称
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/tabletypes [GET]
// @Security ApiKeyAuth
func QueryTableTypeIDAndName(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	res = businessService.QueryTableTypeIDAndNameByAdminID(token)
	c.JSON(http.StatusOK, res)
}

// @Summary 分页查询餐桌种类
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param pageSize query string false "页大小"
// @Param page query string false "页数"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/tabletype [GET]
// @Security ApiKeyAuth
func QueryTableType(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	res = businessService.QueryTableType(token, pageSize, page)
	c.JSON(http.StatusOK, res)
}

// @Summary 查询餐桌种类By ID
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "餐桌种类ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/tabletype/{id} [GET]
// @Security ApiKeyAuth
func QueryTableTypeByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res := businessService.QueryTableTypeByID(id)
	c.JSON(http.StatusOK, res)
}

// @Summary 添加餐桌
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param body body entity.TableRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/table [POST]
// @Security ApiKeyAuth
func AddTable(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	req := entity.TableRequest{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.AddTable(token, req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 修改餐桌
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "餐桌ID"
// @Param body body entity.TableRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/table/{id} [PUT]
// @Security ApiKeyAuth
func UpdateTable(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	req := entity.TableRequest{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = businessService.UpdateTable(token, id, req)
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 删除餐桌
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "餐桌ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/table/{id} [DELETE]
// @Security ApiKeyAuth
func DeleteTable(c *gin.Context) {
	req := entity.DeleteIds{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ids := append(req.Ids, id)
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	res = businessService.DeleteTable(token, ids)
	c.JSON(http.StatusOK, res)
}

// @Summary 查询餐桌By ID
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "餐桌ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/table/{id} [GET]
// @Security ApiKeyAuth
func QueryTableByID(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res := businessService.QueryTableByID(id)
	c.JSON(http.StatusOK, res)
}

// @Summary 分页查询餐桌
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param pageSize query string false "页大小"
// @Param page query string false "页数"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/table [GET]
// @Security ApiKeyAuth
func QueryTable(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	res = businessService.QueryTable(token, pageSize, page)
	c.JSON(http.StatusOK, res)
}

// @Summary 生成餐桌二维码
// @tags 商家
// @Accept  application/json
// @Produce  json
// @Param id path int true "餐桌ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/business/table/{id}/qrcode [GET]
// @Security ApiKeyAuth
func GetTableqrcode(c *gin.Context) {
	res := entity.ResponseData{}
	token := c.Query("token")
	if token == "" {
		authToken := c.GetHeader("Authorization")
		if authToken == "" {
			res.Message = "Query not 'token' param OR header Authorization has not Bearer token"
			c.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}
		token = strings.TrimSpace(authToken)
	}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	// 获取主机头
	r := c.Request
	host := r.Host
	if strings.HasPrefix(host, "http://") == false {
		host = "http://" + host
	}
	res = businessService.SettingTableqrcode(token, host, id)
	c.JSON(http.StatusOK, res)
}
