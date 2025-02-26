package presentation

import (
	"github.com/gofiber/fiber/v2"
	"github.com/pahmiudahgede/senggoldong/config"
	"github.com/pahmiudahgede/senggoldong/internal/handler"
	"github.com/pahmiudahgede/senggoldong/internal/repositories"
	"github.com/pahmiudahgede/senggoldong/internal/services"
	"github.com/pahmiudahgede/senggoldong/middleware"
	"github.com/pahmiudahgede/senggoldong/utils"
)

func ProductRouter(api fiber.Router) {
	productRepo := repositories.NewProductRepository(config.DB)
	storeRepo := repositories.NewStoreRepository(config.DB)
	productService := services.NewProductService(productRepo, storeRepo)
	productHandler := handler.NewProductHandler(productService)

	productAPI := api.Group("/productinstore")

	productAPI.Post("/add-product", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), productHandler.CreateProduct)

	productAPI.Get("/getproductbyuser", middleware.AuthMiddleware, productHandler.GetAllProductsByStoreID)
	productAPI.Get("getproduct/:product_id", middleware.AuthMiddleware, productHandler.GetProductByID)

	productAPI.Put("updateproduct/:product_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), productHandler.UpdateProduct)

	productAPI.Delete("/delete/:product_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), productHandler.DeleteProduct)

	productAPI.Delete("/delete-products", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), productHandler.DeleteProducts)

	productAPI.Delete("/delete-image/:image_id", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), productHandler.DeleteProductImage)

	productAPI.Delete("/delete-images", middleware.AuthMiddleware, middleware.RoleMiddleware(utils.RoleAdministrator, utils.RolePengelola, utils.RolePengepul), productHandler.DeleteProductImages)
}
