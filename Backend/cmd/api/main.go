package main

import (
	"Backend/docs"
	"Backend/internal/config"
	"Backend/internal/controller"
	"Backend/internal/middleware"
	"Backend/internal/model"
	"Backend/internal/repository"
	"Backend/internal/service"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db := config.ConnectDB(cfg)

	userRepo := repository.NewUserRepository(db)
	clientRepo := repository.NewApiClientRepository(db)
	personelRepo := repository.NewPersonelRepository(db)
	riwayatRepo := repository.NewRiwayatJabatanRepository(db)
	referensiRepo := repository.NewReferensiRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)
	exportRepo := repository.NewExportRepository(db)

	userService := service.NewUserService(userRepo, cfg)
	clientAuthService := service.NewClientAuthService(clientRepo, cfg)
	apiClientService := service.NewApiClientService(clientRepo)
	userAdminService := service.NewUserAdminService(userRepo)
	personelService := service.NewPersonelService(personelRepo)
	riwayatService := service.NewRiwayatJabatanService(riwayatRepo, personelRepo)

	authCtl := controller.NewAuthController(userService)
	oauthCtl := controller.NewOAuthController(clientAuthService)
	apiClientCtl := controller.NewApiClientController(apiClientService)
	userCtl := controller.NewUserController(userAdminService)
	personelCtl := controller.NewPersonelController(personelService)
	riwayatCtl := controller.NewRiwayatJabatanController(riwayatService)
	referensiCtl := controller.NewReferensiController(referensiRepo)
	dashboardCtl := controller.NewDashboardController(dashboardRepo)
	exportCtl := controller.NewExportController(dashboardRepo, exportRepo)

	gin.SetMode(cfg.GinMode)
	r := gin.Default()
	docs.Register(r) // Swagger UI: http://localhost:<APP_PORT>/docs

	api := r.Group("/api")

	api.POST("/login", authCtl.Login)        // user (Admin SSDM / Operator)
	api.POST("/oauth/token", oauthCtl.Token) // aplikasi lain (client credentials)

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret, userRepo, clientRepo))
	{
		protected.GET("/auth/me", authCtl.Me)

		protected.GET("/dashboard",
			middleware.RequirePermission(model.PermPersonelRead),
			dashboardCtl.Ringkasan)
		protected.GET("/export/personel",
			middleware.RequirePermission(model.PermPersonelRead),
			exportCtl.Personel)

		// data master untuk form
		protected.GET("/referensi/pangkat", referensiCtl.Pangkat)
		protected.GET("/referensi/satker", referensiCtl.Satker)
		protected.GET("/referensi/fungsi", referensiCtl.Fungsi)
		protected.GET("/referensi/nivelering", referensiCtl.Nivelering)

		protected.GET("/personel",
			middleware.RequirePermission(model.PermPersonelRead),
			personelCtl.List)
		protected.GET("/personel/:id",
			middleware.RequirePermission(model.PermPersonelRead),
			personelCtl.GetByID)
		protected.GET("/personel/:id/drh",
			middleware.RequirePermission(model.PermPersonelRead),
			personelCtl.DRH)
		protected.POST("/personel",
			middleware.RequirePermission(model.PermPersonelWrite),
			personelCtl.Create)
		protected.PUT("/personel/:id",
			middleware.RequirePermission(model.PermPersonelWrite),
			personelCtl.Update)
		protected.DELETE("/personel/:id",
			middleware.RequirePermission(model.PermPersonelDelete),
			personelCtl.Delete)

		// riwayat jabatan
		protected.GET("/personel/:id/riwayat-jabatan",
			middleware.RequirePermission(model.PermJabatanRead),
			riwayatCtl.List)
		protected.POST("/personel/:id/riwayat-jabatan",
			middleware.RequirePermission(model.PermJabatanWrite),
			riwayatCtl.Create)
		protected.PUT("/riwayat-jabatan/:id",
			middleware.RequirePermission(model.PermJabatanWrite),
			riwayatCtl.Update)
		protected.DELETE("/riwayat-jabatan/:id",
			middleware.RequirePermission(model.PermJabatanDelete),
			riwayatCtl.Delete)

		// kelola aplikasi lain (khusus Admin SSDM)
		admin := protected.Group("/api-clients", middleware.RequirePermission(model.PermUserManage))
		{
			admin.GET("", apiClientCtl.List)
			admin.POST("", apiClientCtl.Create)
			admin.PUT("/:id", apiClientCtl.Update)
			admin.PATCH("/:id/status", apiClientCtl.SetStatus)
			admin.POST("/:id/rotate-secret", apiClientCtl.RotateSecret)
		}

		// kelola pengguna (khusus Admin SSDM)
		users := protected.Group("/users", middleware.RequirePermission(model.PermUserManage))
		{
			users.GET("", userCtl.List)
			users.POST("", userCtl.Create)
			users.PUT("/:id", userCtl.Update)
			users.PATCH("/:id/status", userCtl.SetStatus)
			users.POST("/:id/reset-password", userCtl.ResetPassword)
		}
	}
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}

}
