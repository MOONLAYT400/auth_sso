package auth

import (
	"auth-service/internal/domain/models"
	hashPassword "auth-service/internal/helpers/hash-password"
	"auth-service/internal/lib/jwt"
	"auth-service/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	log *slog.Logger
	userSaver UserSaver
	userProvider UserProvider
	appProvider AppProvider
	tokenTTL time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context ,email string, passHash []byte) (uid int64, err error)
}


type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appId int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppId = errors.New("invalid app id")
	ErrUserExists = errors.New("user already exists")
)

// New returns a new instance of Auth.
// tokenTTL is the time to live for login tokens.
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{log: log, userSaver: userSaver, userProvider: userProvider, appProvider: appProvider, tokenTTL: tokenTTL}
}

// Login returns a login token for the user identified by the given email and
// password, for the given app ID. The token is valid for a duration of
// a.tokenTTL. The returned error is non-nil if the email or password is
// invalid, or if the user is not enabled for the given app ID.
func (a *Auth) Login(ctx context.Context, email,password string, appId int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op))

	log.Info("Logging in user", slog.String("email", email))


// get user by email
	user, err := a.userProvider.User(ctx,email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("User not found",slog.String("email",email))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("Failed to get user",slog.String("error",err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}


	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Error("Invalid credentials",slog.String("error",err.Error()))
	}

	log.Info("User logged in",slog.Int("user_id",user.Id))

	app, err := a.appProvider.App(ctx,appId)
	if err != nil {
			if errors.Is(err, storage.ErrAppNotFound) {
				return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		  }
		log.Error("Failed to get app",slog.String("error",err.Error()))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	token,err := jwt.NewToken(user,app,a.tokenTTL)
	if err != nil {
		a.log.Error("Failed to create token",slog.String("error",err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

// RegisterUser creates a new user with the given email and password.
// It returns the user ID of the newly created user or an error if the registration fails.
func (a *Auth) RegisterUser(ctx context.Context,email,password string) (userId int, err error) {
	const op = "auth.RegisterUser"

	log := a.log.With(slog.String("op", op))

	log.Info("Registering user", slog.String("email", email))

	hashPass,err := hashPassword.HashPassword(log,password)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	uid, err := a.userSaver.SaveUser(ctx,email, hashPass)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}


		log.Error("Failed to save user",slog.String("error",err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User registered",slog.Int("user_id",int(uid)))

	return int(uid), nil
}


// IsAdmin returns true if the user with the given user ID is an admin.
// The returned error is non-nil if the user ID is invalid.
func (a *Auth) IsAdmin(ctx context.Context, userId int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(slog.String("op", op))

	log.Info("Checking if user is admin",slog.Int64("user_id",userId))

	isAdmin, err := a.userProvider.IsAdmin(ctx, userId)
	if err != nil {
		log.Error("Failed to check if user is admin",slog.String("error",err.Error()))

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}