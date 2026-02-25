package main

import (
	"quickwit-go-demo/logger"
	"quickwit-go-demo/middleware" // เรียกใช้ middleware ที่เราแก้กันตะกี้
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

var jwtSecret = []byte("supersecretkey")

func main() {
	// 1. Initialize Zap logger
	log, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	app := fiber.New()

	// 2. ใช้ Middleware สำหรับเก็บ Log ท่องจำไว้ว่าต้องวางไว้บนสุดเสมอ
	// เพื่อให้มันจับ Log ของทุก Request รวมถึง /login และ /profile
	app.Use(middleware.RequestLogger(log))

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

		// ดึง Request ID จาก Context (ที่ถูกสร้างใน Middleware) มาใช้ใน Log นี้
		// เพื่อให้ Search ใน Quickwit แล้วเจอข้อมูลที่เชื่อมโยงกัน
		requestId := c.GetRespHeader("X-Request-ID")

		// Hardcoded user (no DB)
		if input.Username != "admin" || input.Password != "1234" {
			log.Warn("login_failed",
				zap.String("request_id", requestId), // เชื่อมโยงกับ access log
				zap.String("username", input.Username),
				zap.String("reason", "invalid credentials"),
				zap.String("ip", c.IP()),
			)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		// Create JWT
		claims := jwt.MapClaims{
			"username": input.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString(jwtSecret)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		log.Info("login_success",
			zap.String("request_id", requestId),
			zap.String("username", input.Username),
		)

		return c.JSON(fiber.Map{
			"token": t,
		})
	})

	// Protected routes (Grouping ช่วยให้อ่านง่ายขึ้น)
	api := app.Group("/profile")
	api.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtSecret,
	}))

	api.Get("/", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		return c.JSON(fiber.Map{
			"username": claims["username"],
			"message":  "This is a protected profile",
		})
	})

	log.Info("server_started", zap.String("port", "3000"))
	app.Listen(":8080")
}
