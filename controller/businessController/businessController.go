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
	ids := append(req.Ids, strconv.Parselnt(c.Param("id"), 10, 64))
	res := businessService.DeleteGoodsType(ids)
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

// @Summary 查询商家商品种类
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
