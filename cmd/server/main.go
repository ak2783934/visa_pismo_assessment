package main

import (
	"log"

	"github.com/ak2783934/visa-pismo-assessment/internal/account"
	"github.com/ak2783934/visa-pismo-assessment/internal/database"
	"github.com/ak2783934/visa-pismo-assessment/internal/transaction"
	_ "github.com/ak2783934/visa-pismo-assessment/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			Visa/Pismo Assessment API
// @version		1.0
// @description	REST API for managing accounts and financial transactions.
// @host			localhost:8080
// @BasePath		/v1
func main() {
	db, err := database.Open()
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close()

	accountRepo := account.NewRepository(db)
	transactionRepo := transaction.NewRepository(db)

	accountService := account.NewService(accountRepo)
	transactionService := transaction.NewService(accountRepo, transactionRepo)

	accountHandler := account.NewHandler(accountService)
	transactionHandler := transaction.NewHandler(transactionService)

	router := gin.Default()
	v1 := router.Group("/v1")
	accountHandler.RegisterRoutes(v1)
	transactionHandler.RegisterRoutes(v1)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("server: %v", err)
	}
}
