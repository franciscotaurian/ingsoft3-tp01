package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"realico-comidas-backend/config"
	"realico-comidas-backend/database"
	"realico-comidas-backend/handler"
	"realico-comidas-backend/middleware"
	"realico-comidas-backend/repository"
	"realico-comidas-backend/service"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Cargar Configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error al cargar configuración: %v", err)
	}

	// 2. Conectar e Inicializar Base de Datos
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Printf("Advertencia: No se pudo conectar a la base de datos PostgreSQL: %v", err)
		log.Printf("Asegúrate de que PostgreSQL esté corriendo y la base de datos '%s' exista.\n", cfg.DBName)
	}

	// 3. Poblar la tabla 'admins' con el usuario del .env si no existe
	database.SeedAdmin(db, cfg)

	// 4. Inyección de Dependencias
	// Repositories
	adminRepo := repository.NewAdminRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	productRepo := repository.NewProductRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	// Services
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)
	orderService := service.NewOrderService(orderRepo, productRepo)

	// Handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)
	orderHandler := handler.NewOrderHandler(orderService)
	metricsHandler := handler.NewMetricsHandler(orderService)

	// Inicializar Sentry (si existe la variable de entorno)
	sentryDsn := os.Getenv("SENTRY_DSN")
	if sentryDsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              sentryDsn,
			EnableTracing:    true,
			TracesSampleRate: 1.0, // medir 100% de las transacciones en demo
			Environment:      "production",
		}); err != nil {
			log.Printf("Advertencia: No se pudo inicializar Sentry: %v", err)
		} else {
			log.Println("Sentry inicializado correctamente")
		}
		defer sentry.Flush(2 * time.Second)
	}

	// 5. Configurar Router Gin
	// 5. Configurar Router Gin
	r := gin.Default()

	// Middleware de Sentry para captura de errores y performance
	if sentryDsn != "" {
		r.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	// Ruta de salud (soporta GET y HEAD para UptimeRobot)
	r.Match([]string{"GET", "HEAD"}, "/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"app":     "Realico Comidas API",
			"version": cfg.Version,
		})
	})

	// Endpoint de prueba temporal para Sentry (§3.5 de la guía - sacar al terminar el TP)
	r.GET("/api/falla", func(c *gin.Context) {
		panic("Falla de prueba del TP9 para verificar Sentry en producción")
	})

	// 6. Definir Rutas
	api := r.Group("/api")
	{
		// --- Rutas Públicas ---
		api.GET("/products", productHandler.GetAvailable) // Solo productos con stock > 0
		api.GET("/categories", categoryHandler.GetAll)
		api.POST("/orders", orderHandler.Create) // Crear pedido (valida stock y teléfono)

		// --- Rutas de Admin (Protegidas: valida credenciales contra tabla 'admins' en DB) ---
		admin := api.Group("/admin")
		admin.Use(middleware.BasicAuthMiddleware(adminRepo))
		{
			// Productos Admin (devuelve TODOS los productos, incluyendo stock <= 0)
			admin.GET("/products", productHandler.GetAll)
			admin.POST("/products", productHandler.Create)
			admin.PUT("/products/:id", productHandler.Update)
			admin.DELETE("/products/:id", productHandler.Delete)

			// Categorías Admin
			admin.POST("/categories", categoryHandler.Create)
			admin.PUT("/categories/:id", categoryHandler.Update)
			admin.DELETE("/categories/:id", categoryHandler.Delete)

			// Pedidos Admin (historial protegido: requiere autenticación de administrador)
			admin.GET("/orders", orderHandler.GetAll)
			admin.PATCH("/orders/:id/status", orderHandler.UpdateStatus)
			admin.POST("/orders/:id/cancel", orderHandler.CancelOrder)

			// Métricas Admin
			admin.GET("/metrics", metricsHandler.GetMetrics)
		}
	}

	// 7. Iniciar Servidor
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Servidor corriendo en http://localhost%s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Error al arrancar el servidor: %v", err)
	}
}
