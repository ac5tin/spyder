package api

import (
	"spyder/crawler"

	"github.com/gofiber/fiber/v2"
)

func scrapeFull(c *fiber.Ctx) error {
	type inputReq struct {
		URL string `json:"url"`
	}

	input := new(inputReq)

	if err := c.BodyParser(input); err != nil {
		c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"result": "error",
			"error":  err.Error(),
		})
		return nil
	}

	r := new(crawler.Results)
	crawler.C.Full(input.URL, r)

	// return result
	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"result": "success",
		"data":   r,
	})
	return nil
}
