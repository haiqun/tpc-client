package routers

import (
	"github.com/gin-gonic/gin"
	"tcp-tls-project/http/controllers"
	"tcp-tls-project/http/middles"
)

func NewApis(api *gin.RouterGroup) {
	//user := controllers.User{}
	//  ===  开始路由注册 start====

	userGroup := api.Group("user", middles.CheckIp())
	{
		t := controllers.Test{}
		userGroup.GET("get", t.TestA)
	}
	//
	//proGroup := api.Group("pro", middles.CheckIp())
	//{
	//	newPro(proGroup)
	//}
	//
	//vipGroup := api.Group("vip", middles.CheckIp())
	//{
	//	newVip(vipGroup)
	//}
	//
	//commonGroup := api.Group("common", middles.CheckIp())
	//{
	//	newCommon(commonGroup)
	//}
}

//  === 产品中心 start ====
//func newPro(proGroup *gin.RouterGroup) {
//
//	product := controllers.Product{}
//	audit := controllers.Audit{}
//	adminService := controllers.AdminService{}
//	serviceSystem := controllers.ServiceSystem{}
//
//	// 产品
//	productGroup := proGroup.Group("product")
//	{
//		// 同步产品到呼叫系统
//		productGroup.POST("sync", product.Sync)
//	}
//
//	// 审核
//	auditGroup := proGroup.Group("audit/:area", middles.GetArea())
//	{
//		// 商务审核
//		auditGroup.POST("commerce", audit.Commerce)
//		// 合规审核
//		auditGroup.POST("compliance", audit.Compliance)
//	}
//
//	// 客服
//	adminServiceGroup := proGroup.Group("admin-service")
//	{
//		// 同步客服到客服系统
//		adminServiceGroup.POST("sync", adminService.Sync)
//		// 通过区域、CRM线下用户id、产品中心产品id获取对应客服系统对应 客服id（其实就是绑定产品中心的客服id）
//		adminServiceGroup.GET("customer-service-id", adminService.GetCustomerServiceId)
//	}
//
//	// 客服系统接口
//	serviceSystemGroup := proGroup.Group("service-system")
//	{
//		// 拉取客服系统的产品上级和投顾直属上级，以及不同所属系统值
//		serviceSystemGroup.POST("info", serviceSystem.SystemInfo)
//		// 同步用户客服关系到客服系统
//		serviceSystemGroup.POST("customer-service-user", serviceSystem.CustomerServiceUser)
//		// 触发客服系统分配客服
//		serviceSystemGroup.POST("assign-customer-service", serviceSystem.AssignCustomerService)
//	}
//}
//
////  === Vip操盘 start ====
//func newVip(vipGroup *gin.RouterGroup) {
//
//	product := controllers.VipProduct{}
//	audit := controllers.VipAudit{}
//
//	// 产品
//	productGroup := vipGroup.Group("product", middles.CheckIp())
//	{
//		// 同步产品到呼叫系统
//		productGroup.POST("sync", product.Sync)
//	}
//
//	// 审核
//	auditGroup := vipGroup.Group("audit/:area", middles.GetArea())
//	{
//		// 商务审核
//		auditGroup.POST("commerce", audit.Commerce)
//		// 合规审核
//		auditGroup.POST("compliance", audit.Compliance)
//	}
//}
//
////  === 产品中心/Vip操盘 共用接口 start ====
//func newCommon(commonGroup *gin.RouterGroup) {
//
//	order := controllers.Order{}
//	common := controllers.Common{}
//	customer := controllers.Customer{}
//
//	//订单
//	orderGroup := commonGroup.Group("order/:area", middles.GetArea())
//	{
//		// 推送订单信息到客服系统
//		orderGroup.POST("push", order.ServiceOrderSync)
//		// 同步订单信息的部分数据
//		orderGroup.POST("sync-info", order.SyncOrderInfo)
//		// 通过用户ID，获取订单列表
//		orderGroup.GET("order-list-by-user", order.OrderListByUser)
//		// 呼叫系统订单列表
//		orderGroup.GET("order-list", order.OrderList)
//		// 呼叫系统订单文件列表
//		orderGroup.GET("order-file", order.OrderFile)
//		// 更新系统订单文件存储状态
//		orderGroup.POST("update-order-file-storage-status", order.UpdateOrderFileStorageStatus)
//		// 呼叫系统订单数据
//		orderGroup.GET("order-by-condition", order.OrderByCondition)
//		// 从呼叫系统获取补充数据
//		// account_time 到账时间 20200106 task-516
//		orderGroup.GET("supplement-fields", order.GetSupplementFields)
//	}
//
//	// 客戶
//	customerGroup := commonGroup.Group("customer/:area", middles.GetArea())
//	{
//		// 获取客户信息
//		customerGroup.GET("get", customer.CustomerInfo)
//		// 根据客户ID 置空crm用户编号
//		customerGroup.GET("unset-user-number", customer.CustomerUnsetUserNumber)
//	}
//
//	// 公共模块
//	baseGroup := commonGroup.Group("base/:area", middles.GetArea())
//	{
//		baseGroup.GET("phone-call", common.PhoneCall) // 呼出电话
//	}
//}
