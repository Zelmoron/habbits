package routers

import (
	"Trecker/internal/controllers"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

const (
	contextKeyUser = "user"
)

var jwtSecretKey = []byte("very-secret-key")

func Routers(app *fiber.App) {

	public := app.Group("")   // routers without jwt
	public.Use(recover.New()) // middleware for recover panic

	public.Post("/registration", controllers.Registration)
	public.Post("/auth", controllers.Auth)

	///////////////////////////////////////////////////////////

	private := app.Group("")                //routers with jwt
	private.Use(recover.New())              // middleware for recover panic
	private.Use(jwtware.New(jwtware.Config{ //middleware for jwt
		SigningKey: jwtware.SigningKey{
			Key: jwtSecretKey, // sercret key
		},
		ContextKey: contextKeyUser, //context key
	}))

	private.Get("/app", controllers.MainPage)
	private.Get("/profile", controllers.Profile)
	private.Get("/habits", controllers.GetHabits)
	private.Post("/habits", controllers.AddHabits)
	private.Post("/updateday", controllers.UpdateHabits)

	// authorizedGroup.Get("/profile", middleware.JWTMiddleware, controllers.Profile) // старый middleware

}
