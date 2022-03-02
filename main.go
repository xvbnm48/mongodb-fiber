package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	gotenv "github.com/joho/godotenv"
	"github.com/xvbnm48/mongodb-fiber/config"
	"github.com/xvbnm48/mongodb-fiber/routes"
	"log"
	"os"
)

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "you are at the root end point",
		})
	})
	api := app.Group("/api")
	routes.CatchprasesRoute(api.Group("/catchprases"))
	app.Listen(":3000")
}

func main() {
	if os.Getenv("APP_ENV") == "production" {
		err := gotenv.Load()
		if err != nil {
			log.Fatalln("error loading env file")
		}
	}

	app := fiber.New()
	app.Use(cors.New())
	app.Use(logger.New())

	config.ConnectDB()

	setupRoutes(app)

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalln("error starting server")
		panic(err)
	}
}
