package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/internal/controllers"
	"github.com/pahmiudahgede/senggoldong/internal/middleware"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func AppRouter(app *fiber.App) {
	// # init
	pointRepo := repositories.NewPointRepository()
	pointService := services.NewPointService(pointRepo)
	pointController := controllers.NewPointController(pointService)

	bannerRepo := repositories.NewBannerRepository()
	bannerService := services.NewBannerService(bannerRepo)
	bannerController := controllers.NewBannerController(bannerService)

	// # api group domain endpoint #
	api := app.Group("/apirijikid")

	// # API Secure #
	api.Use(middleware.APIKeyMiddleware)
	api.Use(middleware.RateLimitMiddleware)

	// # user initial coint #
	api.Get("/user/initial-coint", pointController.GetAllPoints)
	api.Get("/user/initial-coint/:id", pointController.GetPointByID)
	api.Post("/user/initial-coint", middleware.RoleRequired(utils.RoleAdministrator), pointController.CreatePoint)
	api.Put("/user/initial-coint/:id", middleware.RoleRequired(utils.RoleAdministrator), pointController.UpdatePoint)
	api.Delete("/user/initial-coint/:id", middleware.RoleRequired(utils.RoleAdministrator), pointController.DeletePoint)

	//# coverage area #
	api.Get("/coverage-areas", controllers.GetCoverageAreas)
	api.Get("/coverage-areas-district/:id", controllers.GetCoverageAreaByIDProvince)
	api.Get("/coverage-areas-subdistrict/:id", controllers.GetCoverageAreaByIDDistrict)
	api.Post("/coverage-areas", middleware.RoleRequired(utils.RoleAdministrator), controllers.CreateCoverageArea)
	api.Post("/coverage-areas-district", middleware.RoleRequired(utils.RoleAdministrator), controllers.CreateCoverageDistrict)
	api.Post("/coverage-areas-subdistrict", middleware.RoleRequired(utils.RoleAdministrator), controllers.CreateCoverageSubdistrict)
	api.Put("/coverage-areas/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.UpdateCoverageArea)
	api.Put("/coverage-areas-district/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.UpdateCoverageDistrict)
	api.Put("/coverage-areas-subdistrict/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.UpdateCoverageSubdistrict)
	api.Delete("/coverage-areas/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.DeleteCoverageArea)
	api.Delete("/coverage-areas-district/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.DeleteCoverageDistrict)
	api.Delete("/coverage-areas-subdistrict/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.DeleteCoverageSubdistrict)

	// # role #
	api.Get("/roles", middleware.RoleRequired(utils.RoleAdministrator), controllers.GetAllUserRoles)
	api.Get("/role/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.GetUserRoleByID)

	// # authentication #
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Post("/logout", controllers.Logout)

	// # userinfo #
	api.Get("/user", middleware.AuthMiddleware, controllers.GetUserInfo)
	api.Post("/user/update-password", middleware.AuthMiddleware, controllers.UpdatePassword)
	api.Put("/user/update-user", middleware.AuthMiddleware, controllers.UpdateUser)

	// # view all user (admin)
	api.Get("/user/listallusers", middleware.RoleRequired(utils.RoleAdministrator), controllers.GetListUsers)
	api.Get("/user/listalluser/:roleid", middleware.RoleRequired(utils.RoleAdministrator), controllers.GetUsersByRole)
	api.Get("/user/listuser/:userid", middleware.RoleRequired(utils.RoleAdministrator), controllers.GetUserByUserID)

	// # user set pin #
	api.Get("/user/verif-pin", middleware.AuthMiddleware, controllers.GetPin)
	api.Get("/user/cek-pin-status", middleware.AuthMiddleware, controllers.GetPinStatus)
	api.Post("/user/set-pin", middleware.AuthMiddleware, controllers.CreatePin)
	api.Put("/user/update-pin", middleware.AuthMiddleware, controllers.UpdatePin)
	api.Put("/user/update-pin", middleware.AuthMiddleware, controllers.UpdatePin)

	// # address routing #
	api.Get("/addresses", middleware.AuthMiddleware, controllers.GetListAddress)
	api.Get("/address/:id", middleware.AuthMiddleware, controllers.GetAddressByID)
	api.Post("/address/create-address", middleware.AuthMiddleware, controllers.CreateAddress)
	api.Put("/address/update-address/:id", middleware.AuthMiddleware, controllers.UpdateAddress)
	api.Delete("/address/delete-address/:id", middleware.AuthMiddleware, controllers.DeleteAddress)

	// # article #
	api.Get("/articles", controllers.GetArticles)
	api.Get("/article/:id", controllers.GetArticleByID)
	api.Post("/article/create-article", middleware.RoleRequired(utils.RoleAdministrator), controllers.CreateArticle)
	api.Put("/article/update-article/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.UpdateArticle)
	api.Delete("/article/delete-article/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.DeleteArticle)

	// # trash type #
	api.Get("/trash-categorys", controllers.GetTrashCategories)
	api.Get("/trash-category/:id", controllers.GetTrashCategoryDetail)
	api.Post("/trash-category/create-trash-category", middleware.RoleRequired(utils.RoleAdministrator), controllers.CreateTrashCategory)
	api.Post("/trash-category/create-trash-categorydetail", middleware.RoleRequired(utils.RoleAdministrator), controllers.CreateTrashDetail)
	api.Put("/trash-category/update-trash-category/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.UpdateTrashCategory)
	api.Put("/trash-category/update-trash-detail/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.UpdateTrashDetail)
	api.Delete("/trash-category/delete-trash-category/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.DeleteTrashCategory)
	api.Delete("/trash-category/delete-trash-detail/:id", middleware.RoleRequired(utils.RoleAdministrator), controllers.DeleteTrashDetail)

	// # banner #
	api.Get("/banners", bannerController.GetAllBanners)
	api.Get("/banner/:id", bannerController.GetBannerByID)
	api.Post("/banner/create-banner", middleware.RoleRequired(utils.RoleAdministrator), bannerController.CreateBanner)
	api.Put("/banner/update-banner/:id", middleware.RoleRequired(utils.RoleAdministrator), bannerController.UpdateBanner)
	api.Delete("/banner/delete-banner/:id", middleware.RoleRequired(utils.RoleAdministrator), bannerController.DeleteBanner)

	// # wilayah di indonesia #
	api.Get("/wilayah-indonesia/provinces", controllers.GetProvinces)
	api.Get("/wilayah-indonesia/regencies", controllers.GetRegencies)
	api.Get("/wilayah-indonesia/subdistricts", controllers.GetDistricts)
	api.Get("/wilayah-indonesia/villages", controllers.GetVillages)
	api.Get("/wilayah-indonesia/provinces/:id", controllers.GetProvinceByID)
	api.Get("/wilayah-indonesia/regencies/:id", controllers.GetRegencyByID)
	api.Get("/wilayah-indonesia/subdistricts/:id", controllers.GetDistrictByID)
	api.Get("/wilayah-indonesia/villages/:id", controllers.GetVillageByID)

	// # request pickup by user (masyarakat) #
	api.Get("/requestpickup", middleware.RoleRequired(utils.RoleMasyarakat, utils.RolePengepul), controllers.GetRequestPickupsByUser)
	api.Post("/addrequestpickup", middleware.RoleRequired(utils.RoleMasyarakat), controllers.CreateRequestPickup)
	api.Delete("/deleterequestpickup/:id", middleware.RoleRequired(utils.RoleMasyarakat), controllers.DeleteRequestPickup)

	// # product post by pengepul
	api.Get("/post/products", middleware.RoleRequired(utils.RolePengepul), controllers.GetAllProducts)
	api.Get("/post/product/:productid", middleware.RoleRequired(utils.RolePengepul), controllers.GetProductByID)
	api.Get("/view/product/:storeid", middleware.RoleRequired(utils.RolePengepul), controllers.GetProductsByStore)
	api.Post("/post/addproduct", middleware.RoleRequired(utils.RolePengepul), controllers.CreateProduct)
	api.Put("/post/product/:productid", middleware.RoleRequired(utils.RolePengepul), controllers.UpdateProduct)
	api.Delete("/delete/product/:productid", middleware.RoleRequired(utils.RolePengepul), controllers.DeleteProduct)

	// # marketplace
	api.Get("/store/marketplace", middleware.RoleRequired(utils.RolePengepul), controllers.GetStoresByUserID)
	api.Get("/store/marketplace/:storeid", middleware.RoleRequired(utils.RolePengepul), controllers.GetStoreByID)
}
