package hashPassword

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(log *slog.Logger ,password string) ([]byte,error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password",slog.String("error",err.Error()))
	}
	return hashed, nil
}