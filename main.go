package main

import (
	"quickwit-go-demo/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var jwtSecret = []byte("supersecretkey")

func main() {

	// Zap logger
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	app := fiber.New()

	// Login route
	app.Post("/login", func(c *fiber.Ctx) error {
		type LoginInput struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var input LoginInput
		if err := c.BodyParser(&input); err != nil {
			return fiber.ErrBadRequest
		}

		// Hardcoded user (no DB)
		if input.Username != "admin" || input.Password != "1234" {
			log.Warn("login_failed",
				zap.String("username", input.Username),
				zap.String("reason", "invalid credentials"),
			)
			return fiber.ErrUnauthorized
		}

		// Create JWT
		claims := jwt.MapClaims{
			"username": input.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, _ := token.SignedString(jwtSecret)

		log.Info("login_success",
			zap.String("username", input.Username),
		)

		return c.JSON(fiber.Map{
			"token": t,
		})
	})

	// Protected route
	app.Use("/profile", jwtware.New(jwtware.Config{
		SigningKey: jwtSecret,
	}))

	app.Get("/profile", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		return c.JSON(fiber.Map{
			"username": claims["username"],
		})
	})

	app.Listen(":3000")
}
