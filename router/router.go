package router

import (
	"net/http"
	"yelloment-admin/controllers"
	"yelloment-admin/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine) {

	router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "error/notfound.html", gin.H{
			"title": "Not Found",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.Done()
	})

	auth := router.Group("/auth")
	{
		auth.GET("/login", controllers.ShowLoginPage)
		auth.POST("/login", controllers.ConfirmLogin)

		auth.GET("/logout", controllers.Logout)
	}

	root := router.Group("", middleware.ConfirmAuth())
	{
		root.GET("", func(c *gin.Context) {
			c.Redirect(http.StatusFound, "/dashboard")
		})

		root.GET("/dashboard", controllers.ShowDashboardPage)

		code := root.Group("/code")
		{
			code.GET("/list", controllers.ShowCodeListPage)
			code.GET("/view/:id", controllers.ShowCodeViewPage)
			code.GET("/modal_view/:id", controllers.ShowCodeViewModalPage)
		}

		banner := root.Group("/banner")
		{
			banner.GET("/list", controllers.ShowBannerListPage)
			banner.GET("/view/:id", controllers.ShowBannerViewPage)
		}

		user := root.Group("/user")
		{
			user.GET("/list", controllers.ShowUserListPage)
			user.GET("/leavelist", controllers.ShowLeaveUserListPage)
			user.GET("/view/:id", controllers.ShowUserViewPage)
			user.GET("/view/:id/pop", controllers.ShowUserViewPopPage)
		}

		item := root.Group("/item")
		{
			item.GET("/list", controllers.ShowItemListPage)
			item.GET("/view/:id", controllers.ShowItemViewPage)

			item.GET("/bundle/list", controllers.ShowBundleListPage)
			item.GET("/bundle/view/:id", controllers.ShowBundleViewPage)
		}

		subs := root.Group("/subs")
		{
			subs.GET("/list", controllers.ShowSubsListPage)
			subs.GET("/cancellist", controllers.ShowSubsCancelListPage)
			subs.GET("/view/:id", controllers.ShowSubsViewPage)
		}

		order := root.Group("/order")
		{
			order.GET("/list", controllers.ShowOrderListPage)                              //통합리스트
			order.GET("/view/:id", controllers.ShowOrderViewPage)                          //상세
			order.GET("/itempickup", controllers.ShowOrderItemPickupPage)                  //상품구성
			order.GET("/itempickup_v2", controllers.ShowOrderItemPickupV2Page)             //상품구성 v2
			order.GET("/payment", controllers.ShowOrderPaymentPage)                        //결제
			order.GET("/purchase/item", controllers.ShowOrderPurchaseItemPage)             //원물 발주서
			order.GET("/purchase/itemunpaid", controllers.ShowOrderPurchaseItemUnpaidPage) //원물 발주서(미결제포함)
			order.GET("/purchase/user", controllers.ShowOrderPurchaseUserPage)             //고객 발주서
			order.GET("/purchase/userunpaid", controllers.ShowOrderPurchaseUserUnpaidPage) //고객 발주서(미결제포함)
			order.GET("/delivery/invoice", controllers.ShowOrderInvoicerPage)              //운송장 관리
			order.GET("/delivery/done", controllers.ShowOrderDonePage)                     //배송완료(종결) 관리
		}
	}

	api := router.Group("/api", middleware.ConfirmAuthApi())
	{
		commonApi := api.Group("/common")
		{
			commonApi.POST("/fileupload", controllers.FileUpload)
		}

		codeApi := api.Group("/code")
		{
			codeApi.POST("/set", controllers.SetCode)
		}

		bannerApi := api.Group("/banner")
		{
			bannerApi.POST("/set", controllers.SetBanner)
		}

		userApi := api.Group("/user")
		{
			userApi.GET("/list", controllers.GetUserList)

			userApi.GET("/address", controllers.GetUserAddress)
			userApi.GET("/address/list", controllers.GetUserAddressList)
		}

		itemApi := api.Group("/item")
		{
			itemApi.GET("/list", controllers.GetItemList)

			itemApi.POST("/set", controllers.SetItem)

			itemApi.GET("/bundle/info/:id", controllers.GetBundleInfo)
			itemApi.GET("/bundle/list", controllers.GetBundleList)

			itemApi.POST("/bundle/set", controllers.SetBundle)
			itemApi.POST("/bundle/remove", controllers.RemoveBundle)
		}

		subsApi := api.Group("/subs")
		{
			subsApi.GET("/list", controllers.GetSubsList)

			subsApi.POST("/set", controllers.SetSubs)
			subsApi.POST("/cancel/:id", controllers.CancelSubs)
		}

		orderApi := api.Group("/order")
		{
			orderApi.GET("/list", controllers.GetOrderList)
			orderApi.GET("/prevlist", controllers.GetOrderPrevList)
			orderApi.GET("/itemlist/:id", controllers.GetOrderItemList)

			orderApi.GET("/purchaseitem", controllers.GetOrderPurchaseItem)

			orderApi.POST("/set", controllers.SetOrder)
			orderApi.POST("/cancel/:id", controllers.CancelOrder)

			orderApi.POST("/gensubsorder", controllers.GenSubsOrder)
			orderApi.POST("/item/set", controllers.SetOrderItem)
			orderApi.POST("/payment/:id", controllers.ProcOrderPayment)
			orderApi.POST("/paymentunpaid", controllers.ProcOrderPaymentForUnpaid)
			orderApi.POST("/invoice/bulk", controllers.UpdateOrderInvoiceNoBulk)
			orderApi.POST("/done/:id", controllers.DoneOrder)
			orderApi.POST("/done_bulk", controllers.DoneOrderBulk)
		}

		tagApi := api.Group("/tag")
		{
			tagApi.GET("/info", controllers.GetTagInfo)
		}

		reviewApi := api.Group("/review")
		{
			reviewApi.GET("/list", controllers.GetReviewList)
		}

		counselApi := api.Group("/counsel")
		{
			counselApi.GET("/list", controllers.GetCounselList)

			counselApi.POST("/add", controllers.AddCounsel)
			counselApi.POST("/remove/:id", controllers.RemoveCounsel)
		}
	}
}
