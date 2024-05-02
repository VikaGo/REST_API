package controller

import (
	"context"
	"github.com/VikaGo/REST_API/logger"
	"github.com/VikaGo/REST_API/model"
	"github.com/VikaGo/REST_API/pkg/types"
	"github.com/VikaGo/REST_API/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
)

// UserController ...
type UserController struct {
	ctx      context.Context
	services *service.Manager
	logger   *logger.Logger
}

// NewUsers creates a new user controller.
func NewUsers(ctx context.Context, services *service.Manager, logger *logger.Logger) *UserController {
	return &UserController{
		ctx:      ctx,
		services: services,
		logger:   logger,
	}
}

type LogInInput struct {
	Nickname string `json:"nickname" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Create new user
func (ctr *UserController) Create(ctx echo.Context) error {
	var user model.User
	err := ctx.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode user data"))
	}

	// Validate user input, including checking for a strong password
	if err := validatePassword(&user); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	// Generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not generate hashed password"))
	}

	// Set the hashed password in the user data
	user.Password = string(hashedPassword)

	createdUser, err := ctr.services.User.CreateUser(ctx.Request().Context(), &user)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not create user"))
		}
	}

	ctr.logger.Debug().Msgf("Created user '%s'", createdUser.ID.String())

	return ctx.JSON(http.StatusCreated, createdUser)
}

// Get returns user by ID
func (ctr *UserController) Get(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse user UUID"))
	}
	user, err := ctr.services.User.GetUser(ctx.Request().Context(), userID)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err)
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not get user"))
		}
	}
	return ctx.JSON(http.StatusOK, user)
}

// Update user by ID
func (ctr *UserController) Update(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse user UUID"))
	}

	basicAuthMiddleware("your_nickname", "your_password")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	var updatedUser model.User
	err = ctx.Bind(&updatedUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode updated user data"))
	}
	err = ctx.Validate(&updatedUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	password := updatedUser.Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not hash password"))
	}
	updatedUser.Password = string(hashedPassword)

	updatedUser.ID = userID
	u, err := ctr.services.User.UpdateUser(ctx.Request().Context(), &updatedUser)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not update user"))
		}
	}

	ctr.logger.Debug().Msgf("Updated user '%s'", u.ID.String())

	return ctx.JSON(http.StatusOK, updatedUser)
}

// Delete deletes user by ID
func (ctr *UserController) Delete(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not parse user UUID"))
	}
	err = ctr.services.User.DeleteUser(ctx.Request().Context(), userID)
	if err != nil {
		switch {
		case errors.Cause(err) == types.ErrNotFound:
			return echo.NewHTTPError(http.StatusNotFound, err)
		case errors.Cause(err) == types.ErrBadRequest:
			return echo.NewHTTPError(http.StatusBadRequest, err)
		default:
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not delete user"))
		}
	}

	ctr.logger.Debug().Msgf("Deleted user '%s'", userID.String())

	return ctx.JSON(http.StatusOK, "OK")
}

// Login
func (ctr *UserController) LogIn(ctx echo.Context) error {
	var input LogInInput

	if err := ctx.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	//user, err := ctr.services.User.GetUserByNickname(ctx.Request().Context(), input.Nickname)
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusInternalServerError, err)
	//}
	//
	//err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	//if err != nil {
	//	return echo.NewHTTPError(http.StatusUnauthorized, "Пароль неправильний")
	//}

	token, err := ctr.services.User.GenerateToken(ctx.Request().Context(), input.Nickname, input.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
	return nil
}

func basicAuthMiddleware(username, password string) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(u, p string, c echo.Context) (bool, error) {
		if u == username && p == password {
			return true, nil
		}
		return false, nil
	})
}

func validatePassword(user *model.User) error {

	// at least eight characters
	if len(user.Password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	// at least one figure
	hasNumber := regexp.MustCompile(`\d`).MatchString(user.Password)
	if !hasNumber {
		return errors.New("Password must include at least one figure")
	}

	// at least one special character
	hasSpecialChar := regexp.MustCompile(`[^a-zA-Z0-9\s]`).MatchString(user.Password)
	if !hasSpecialChar {
		return errors.New("Password must include at least one special character")
	}
	return nil
}

// Change password endpoint
func (ctr *UserController) ChangePassword(ctx echo.Context) error {
	userID, err := uuid.Parse(ctx.Param("id"))

	// Check Basic Auth
	basicAuthMiddleware("your_nickname", "your_password")
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}

	// Get the new and existing passwords from the request
	var newPassword struct {
		NewPassword      string `json:"new_password"`
		ExistingPassword string `json:"existing_password"`
	}

	if err := ctx.Bind(&newPassword); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Wrap(err, "could not decode new password"))
	}

	// Validate the new password
	if err := validatePassword(&model.User{Password: newPassword.NewPassword}); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	// Generate a hashed password for the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not generate hashed password"))
	}

	// Retrieve the existing password from the database
	existingPasswordHash, err := ctr.services.User.GetPassword(ctx.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "could not retrieve existing password"))
	}

	// Compare the existing password hash with the new hashed password (both as strings)
	if existingPasswordHash != string(hashedNewPassword) {
		// The passwords do not match, updating the password in the database
		if err := ctr.services.User.UpdatePassword(ctx.Request().Context(), userID, string(hashedNewPassword)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "Could not update password hash"))
		}
		return ctx.JSON(http.StatusOK, "Password changed successfully")
	}

	// Passwords match, please choose a new password
	return echo.NewHTTPError(http.StatusUnauthorized, "Passwords match, please choose a new password")
}
