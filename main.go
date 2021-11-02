package spyder

import (
	"flag"
	"fmt"
	"log"
	"spyder/crawler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// init crawler
	crawler.C = new(crawler.Crawler)
	crawler.C.Init()

	// start REST API Server
	port := flag.Uint("p", 3131, "Port number")
	prefork := flag.Bool("prefork", false, "Prefork")

	app := fiber.New(fiber.Config{
		Prefork: *prefork,
	})
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(cors.New())

	// ==== API Routes ====
	app.Get("/ping", func(c *fiber.Ctx) error { c.Status(200).Send([]byte("pong")); return nil })

	log.Println(fmt.Sprintf("Listening on PORT %d", *port))
	if err := app.Listen(fmt.Sprintf(":%d", *port)); err != nil {
		log.Fatal(err.Error())
	}
}
