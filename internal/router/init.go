package router

import (
	"net/http"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/dto"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/handler"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/apperror"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/jwt"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/router/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Init(opts *handler.HandlerOpts, config *config.Config) http.Handler {
	r := gin.New()
	r.ContextWithFallback = true

	middlewares := []gin.HandlerFunc{
		middleware.Monitoring(),
		middleware.Logger(),
		gin.Recovery(),
		middleware.Error(),
	}
	r.Use(cors.New(config.Cors))
	r.Use(middlewares...)

	r.NoRoute(func(ctx *gin.Context) {
		ctx.Error(apperror.ErrNotFound)
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.Response{Message: "pong"})
	})

	r.POST("/auth/register", opts.UserHandler.Register)
	r.POST("/auth/verify", opts.UserHandler.Verify)
	r.POST("/auth/google/register", opts.UserHandler.GoogleRegister)
	r.GET("/auth/google/register", opts.UserHandler.RedirectGoogleRegister)

	r.POST("/auth/login", opts.UserHandler.Login)
	r.POST("/auth/google/login", opts.UserHandler.GoogleLogin)
	r.GET("/auth/google/login", opts.UserHandler.RedirectGoogleLogin)
	r.GET("/auth/google/callback", opts.UserHandler.GoogleLoginCallback)

	r.POST("/auth/reset-password", opts.UserHandler.ResetPassword)
	r.POST("/auth/confirm-reset-password", opts.UserHandler.ConfirmReset)

	r.GET("/logout", opts.UserHandler.Logout)

	r.GET("/location/provinces", opts.LocationHandler.ListProvinces)
	r.GET("/location/cities", opts.LocationHandler.ListCities)
	r.GET("/location/districts", opts.LocationHandler.ListDistricts)
	r.GET("/location/sub-districts", opts.LocationHandler.ListSubDistricts)

	r.GET("/products/popular", opts.UserHandler.ListPopularProduct)
	r.GET("/products", opts.UserHandler.ListProduct)
	r.GET("/pharmacies/:id/products/:product_id", opts.UserHandler.GetProductDetails)
	r.GET("/pharmacies/bestseller/products/:product_id", opts.UserHandler.AvailablePharmacy)

	r.GET("/logistics", opts.LogisticHandler.ListLogistics)

	jwt := jwt.NewJwtProvider(config.Jwt)
	authMiddleware := middleware.NewAuthMiddleware(jwt)
	cloudinaryMiddleware := middleware.NewCloudinaryMiddleware(&config.Cloudinary)
	workerMiddleware := middleware.NewWorkerMiddleware(&config.Worker)

	j := r.Group("", authMiddleware.RequireToken())
	j.GET("/pharmacies/:id/logistics", opts.PharmacyHandler.ListLogistics)

	c := j.Group("", middleware.Authorize("customer"))

	cv := c.Group("", middleware.VerifiedOnly())
	cv.GET("/carts", opts.UserHandler.GetCartItems)
	cv.POST("/carts", opts.UserHandler.AddCartItem)
	cv.PUT("/carts/pharmacies/:pharmacy_id/products/:product_id", opts.UserHandler.UpdateCartItem)
	cv.DELETE("/carts/pharmacies/:pharmacy_id/products/:product_id", opts.UserHandler.RemoveCartItem)
	cv.POST("/orders", opts.UserHandler.CreateOrder)
	cv.GET("/orders", opts.UserHandler.ListOrder)
	cv.GET("/unpaid-orders", opts.UserHandler.ListUnpaidOrder)
	cv.PATCH("/payments/:payment_id", cloudinaryMiddleware.Upload("payment"), opts.UserHandler.UploadPaymentProof)
	cv.POST("/profile/addresses", opts.UserHandler.AddAddress)
	cv.GET("/profile/addresses", opts.UserHandler.GetAddresses)
	cv.PATCH("/user/orders/:order_id", opts.UserHandler.ConfirmOrder)

	a := j.Group("/admin", middleware.Authorize("admin"))
	a.GET("/users", opts.AdminHandler.ListUsers)
	a.POST("/partners", cloudinaryMiddleware.Upload("partner"), opts.AdminHandler.AddPartner)
	a.GET("/partners", opts.AdminHandler.ListPartners)
	a.GET("/partners/:id", opts.AdminHandler.GetPartnerByID)
	a.PUT("/partners/:id", opts.AdminHandler.EditPartner)
	a.POST("/pharmacists", opts.AdminHandler.AddPharmacist)
	a.GET("/pharmacists", opts.AdminHandler.ListPharmacists)
	a.GET("/pharmacists/:id", opts.AdminHandler.GetPharmacist)
	a.PUT("/pharmacists/:id", opts.AdminHandler.EditPharmacist)
	a.DELETE("/pharmacists/:id", opts.AdminHandler.RemovePharmacist)
	a.POST("/pharmacies", opts.AdminHandler.AddPharmacy)
	a.GET("/pharmacies", opts.AdminHandler.ListPharmacies)
	a.GET("/categories", opts.AdminHandler.ListProductCategories)
	a.POST("/categories", opts.AdminHandler.AddProductCategory)
	a.PUT("/categories/:id", opts.AdminHandler.UpdateProductCategory)
	a.DELETE("/categories/:id", opts.RemoveProductCategory)
	a.GET("/products", opts.AdminHandler.ListProduct)
	a.POST("/products", cloudinaryMiddleware.Upload("product"), opts.AdminHandler.AddProduct)
	a.GET("/manufacturers", opts.AdminHandler.ListManufacturer)
	a.GET("/product-classifications", opts.AdminHandler.ListProductClassification)
	a.GET("/product-forms", opts.AdminHandler.ListProductForm)

	p := j.Group("/pharmacist", middleware.Authorize("pharmacist"))
	p.GET("/master-products", opts.PharmacistHandler.ListMasterProducts)
	p.GET("/pharmacy-products", opts.PharmacistHandler.ListPharmacyProducts)
	p.GET("/pharmacy-products/:product_id", opts.PharmacistHandler.GetPharmacyProduct)
	p.POST("/pharmacy-products", opts.PharmacistHandler.AddPharmacyProduct)
	p.PUT("/pharmacy-products/:product_id", opts.PharmacistHandler.UpdatePharmacyProduct)
	p.DELETE("/pharmacy-products/:product_id", opts.PharmacistHandler.DeletePharmacyProduct)
	p.GET("/orders", opts.PharmacistHandler.ListOrder)
	p.PATCH("/orders/:order_id", opts.PharmacistHandler.SendOrder)

	w := r.Group("", workerMiddleware.Bypass())
	w.PATCH("/orders/payments/:payment_id", opts.AdminHandler.ProcessOrder)
	w.PATCH("/orders/:order_id", workerMiddleware.Bypass(), opts.AdminHandler.ConfirmOrder)
	w.PATCH("/partners/:id", opts.AdminHandler.EditPartnerDaysAndHours)

	ch := j.Group("/cache", middleware.Authorize("admin"))
	ch.DELETE("/delete", opts.CacheHandler.DeleteAll)
	ch.DELETE("/flush", opts.CacheHandler.FlushAll)

	pprof.Register(r)
	return r
}
