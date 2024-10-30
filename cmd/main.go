package main

import (
	"Trecker/internal/db"
	"Trecker/internal/routers"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	db.Connect()
	app := fiber.New(fiber.Config{ // close connection if will  long read and write
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000", // Замените на ваш фронтенд URL
		AllowCredentials: true,
	}))
	routers.Routers(app)

	app.Listen(":8080")

}
