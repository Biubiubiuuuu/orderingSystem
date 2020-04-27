package responseCommon

import "github.com/Biubiubiuuuu/orderingSystem/entity"

// {"status":true,"data":{},"message":""}
// {"status":true,"data":nil,"message":""}
// {"status":false,"data":nil,"message":""}
func Response(status bool, data map[string]interface{}, message string) (responseData entity.ResponseData) {
	responseData.Status = status
	responseData.Data = data
	responseData.Message = message
	return responseData
}
