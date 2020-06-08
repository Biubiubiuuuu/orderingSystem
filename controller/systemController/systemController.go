package systemController

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Biubiubiuuuu/orderingSystem/entity"
	"github.com/Biubiubiuuuu/orderingSystem/helper/configHelper"
	"github.com/Biubiubiuuuu/orderingSystem/helper/fileHelper"
	"github.com/Biubiubiuuuu/orderingSystem/service/systemService"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// @Summary 管理员登录
// @tags 系统管理员
// @Accept  application/json
// @Produce  json
// @Param body body entity.SystemAdminLoginRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/login [POST]
func Login(c *gin.Context) {
	req := entity.SystemAdminLoginRequest{}
	res := entity.ResponseData{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = systemService.Login(req, c.ClientIP())
	}
	c.JSON(http.StatusOK, res)
}

// @Summary 添加管理员
// @tags 系统管理员
// @Accept  multipart/form-data
// @Produce  json
// @Param username formData string true "用户名"
// @Param nikename formData string false "昵称"
// @Param password formData string true "密码"
// @Param manager formData string false "操作权限 Y | N"
// @Param avatar formData file false "用户头像"
// @Param is_enable formData bool false "是否启用"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admin [POST]
// @Security ApiKeyAuth
func AddAdmin(c *gin.Context) {
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
	var avatar string
	if file, err := c.FormFile("avatar"); err == nil {
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
		// 获取主机头
		r := c.Request
		host := r.Host
		if err := c.SaveUploadedFile(file, pathFile); err == nil {
			if strings.HasPrefix(host, "http://") == false {
				host = "http://" + host
			}
			avatar = host + "/" + pathFile
		}
	}
	is_enable, _ := strconv.ParseBool(c.DefaultPostForm("is_enable", "false"))
	req := entity.SystemAdminAddRequest{
		Username: c.PostForm("username"),
		Nikename: c.PostForm("nikename"),
		Password: c.DefaultPostForm("password", "123456"),
		Manager:  c.DefaultPostForm("manager", "N"),
		Avatar:   avatar,
		IsEnable: is_enable,
	}
	res = systemService.Add(token, req)
	c.JSON(http.StatusOK, res)
}

// @Summary 修改管理员密码
// @tags 系统管理员
// @Accept  application/json
// @Produce  json
// @Param body body entity.SystemAdminUpdatePassRequest true "body"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admin/password [PUT]
// @Security ApiKeyAuth
func UpdatePass(c *gin.Context) {
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
	req := entity.SystemAdminUpdatePassRequest{}
	if c.ShouldBindJSON(&req) != nil {
		res.Message = "请求参数json错误"
	} else {
		res = systemService.UpdatePass(token, req)
	}
	c.JSON(http.StatusCreated, res)
}

// @Summary 分页查询管理员(默认前100条) 并返回总记录数
// @tags 系统管理员
// @Accept application/x-www-form-urlencoded
// @Produce  json
// @Param username query string false "用户名"
// @Param created_at_start query string false "创建开始时间"
// @Param created_at_end query string false "创建结束时间"
// @Param manager query string false "操作权限"
// @Param created_by query string false "创建人"
// @Param is_enable query bool false "是否启用"
// @Param pageSize query string false "页大小"
// @Param page query string false "页数"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admins [GET]
// @Security ApiKeyAuth
func QueryAdmins(c *gin.Context) {
	username := c.Query("username")
	created_at_start := c.Query("created_at_start")
	created_at_end := c.Query("created_at_end")
	manager := c.Query("manager")
	created_by := c.Query("created_by")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "100"))
	args := map[string]interface{}{
		"username":         username,
		"created_at_start": created_at_start,
		"created_at_end":   created_at_end,
		"manager":          manager,
		"created_by":       created_by,
	}
	is_enableStr := c.Query("is_enable")
	if is_enableStr != "" {
		is_enable, _ := strconv.ParseBool(is_enableStr)
		args["is_enable"] = is_enable
	}
	res := systemService.QueryByLimitOffset(args, pageSize, page)
	c.JSON(http.StatusCreated, res)
}

// @Summary 删除管理员
// @tags 系统管理员
// @Accept  application/json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admin/{id} [DELETE]
// @Security ApiKeyAuth
func DeleteAdmin(c *gin.Context) {
	req := entity.DeleteIds{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ids := append(req.Ids, id)
	res := systemService.DeleteAdmin(ids)
	c.JSON(http.StatusOK, res)
}

// @Summary 批量删除管理员
// @tags 系统管理员
// @Accept  application/json
// @Produce  json
// @Param ids path string true "Account ID 多个用,分开"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/admin/admins/{ids} [POST]
// @Security ApiKeyAuth
func DeleteAdmins(c *gin.Context) {
	res := entity.ResponseData{}
	ids := c.Param("ids")
	idArr := strings.Split(",", ids)
	var arr []int64
	for _, v := range idArr {
		item, _ := strconv.ParseInt(v, 10, 64)
		arr = append(arr, item)
	}
	res = systemService.DeleteAdmin(arr)
	c.JSON(http.StatusOK, res)
}

// @Summary 启用/禁用管理员
// @tags 系统管理员
// @Accept  application/json
// @Produce  json
// @Param id path int true "Account ID"
// @Param is_enable query bool true "启用/禁用管理员"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admin/enable/{id} [PUT]
// @Security ApiKeyAuth
func IsEnableAdmin(c *gin.Context) {
	res := entity.ResponseData{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	is_enable, _ := strconv.ParseBool(c.Query("is_enable"))
	res = systemService.IsEnableAdmin(id, is_enable)
	c.JSON(http.StatusOK, res)
}

// @Summary 授权/禁止管理员权限
// @tags 系统管理员
// @Accept  application/json
// @Produce  json
// @Param id path int true "Account ID"
// @Param is_manager query string true "Y：授权/ N：禁止"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admin/manager/{id} [PUT]
// @Security ApiKeyAuth
func IsManagerAdmin(c *gin.Context) {
	res := entity.ResponseData{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	is_manager := c.Query("is_manager")
	res = systemService.IsManagerAdmin(id, is_manager)
	c.JSON(http.StatusOK, res)
}

// @Summary 修改管理员
// @tags 系统管理员
// @Accept  multipart/form-data
// @Produce  json
// @Param nikename formData string false "昵称"
// @Param manager formData string false "操作权限 Y | N"
// @Param avatar formData file false "用户头像"
// @Param is_enable formData bool true "是否启用"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admin [PUT]
// @Security ApiKeyAuth
func UpdateAdmin(c *gin.Context) {
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
	var avatar string
	if file, err := c.FormFile("avatar"); err == nil {
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
		// 获取主机头
		r := c.Request
		host := r.Host
		if err := c.SaveUploadedFile(file, pathFile); err == nil {
			if strings.HasPrefix(host, "http://") == false {
				host = "http://" + host
			}
			avatar = host + "/" + pathFile
		}
	}
	is_enable, _ := strconv.ParseBool(c.DefaultPostForm("is_enable", "false"))
	args := map[string]interface{}{
		"nikename":  c.PostForm("nikename"),
		"manager":   c.DefaultPostForm("manager", "N"),
		"avatar":    avatar,
		"is_enable": is_enable,
	}
	res = systemService.UpdateAdmin(token, args)
	c.JSON(http.StatusCreated, res)
}

// @Summary 查询管理员by token
// @tags 系统管理员
// @Accept  application/json
// @Produce  json
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admin [GET]
// @Security ApiKeyAuth
func QueryAdmin(c *gin.Context) {
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
	res = systemService.QueryAdminByToken(token)
	c.JSON(http.StatusCreated, res)
}

// @Summary 查询管理员by id
// @tags 系统管理员
// @Accept  application/json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} entity.ResponseData "desc"
// @Router /api/v1/system/admin/{id} [GET]
// @Security ApiKeyAuth
func QueryAdminByID(c *gin.Context) {
	res := entity.ResponseData{}
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res = systemService.QueryAdminByID(id)
	c.JSON(http.StatusCreated, res)
}
