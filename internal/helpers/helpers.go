package helpers

import (
	"fmt"
	"net/http"
	"server/everydaymuslimappserver/internal/config"

	"golang.org/x/crypto/bcrypt"
)

var app *config.AppConfig

//NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "userId")
	return exists
}

func HashPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	fmt.Println(string(hashedPassword))
}
