package api

import "github.com/gofiber/fiber/v2"

func Routes(router *fiber.Router) {
	(*router).Post("/full", scrapeFull)
}
