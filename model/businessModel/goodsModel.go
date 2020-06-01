package businessModel

import "github.com/Biubiubiuuuu/orderingSystem/model"

type Goods struct {
	model.Model
	GoodsName        string `json:"goods_name"`        // 商品名称
	GoodsPhoto       string `json:"goods_photo"`       // 商品图片
	GoodsDescription string `json:"doods_description"` // 商品描述

}
