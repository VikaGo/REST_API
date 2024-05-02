package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"

	model "github.com/VikaGo/REST_API/model"
	"github.com/VikaGo/REST_API/pkg/types"
	"github.com/VikaGo/REST_API/store"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

const (
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 24 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
}

// UserWebService ...
type UserWebService struct {
	ctx   context.Context
	store *store.Store
}
type CustomError struct {
	Code    int
	Message string
}

func (e *CustomError) Error() string {
	return e.Message
}

// NewUserWebService creates a new user web service
func NewUserWebService(ctx context.Context, store *store.Store) *UserWebService {
	return &UserWebService{
		ctx:   ctx,
		store: store,
	}
}

// GetUser ...
func (svc *UserWebService) GetUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	userDB, err := svc.store.User.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.GetUser")
	}
	if userDB == nil {
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%s' not found", userID.String()))
	}

	return userDB.ToWeb(), nil
}

// CreateUser ...
func (svc *UserWebService) CreateUser(ctx context.Context, reqUser *model.User) (*model.User, error) {

	if reqUser == nil {
		return nil, errors.New("reqUser is nil")
	}

	reqUser.ID = uuid.New()

	if svc.store.User == nil {
		return nil, errors.New("svc.store.User is nil")
	}

	dbUser := reqUser.ToDB()
	if dbUser == nil {
		return nil, errors.New("conversion from reqUser to DBUser resulted in nil pointer")
	}

	_, err := svc.store.User.CreateUser(ctx, reqUser.ToDB())
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.CreateUser error")
	}

	// get created user by ID
	createdDBUser, err := svc.store.User.GetUser(ctx, reqUser.ID)
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.GetUser error")
	}

	if createdDBUser == nil {
		return nil, errors.New("createdDBUser is nil")
	}

	webUser := createdDBUser.ToWeb()
	if webUser == nil {
		return nil, errors.New("conversion from DBUser to User (web) resulted in nil pointer")
	}

	return webUser, nil
}

// UpdateUser ...
func (svc *UserWebService) UpdateUser(ctx context.Context, reqUser *model.User) (*model.User, error) {
	// Perform the update in the store
	updatedUserDB, err := svc.store.User.UpdateUser(ctx, reqUser.ToDB())
	if err != nil {
		return nil, errors.Wrap(err, "svc.user.UpdateUser error")
	}

	return updatedUserDB.ToWeb(), nil
}

// DeleteUser ...
func (svc *UserWebService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Check if user exists
	userDB, err := svc.store.User.GetUser(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "svc.user.GetUser error")
	}
	if userDB == nil {
		return errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%s' not found", userID.String()))
	}

	err = svc.store.User.DeleteUser(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "svc.user.DeleteUser error")
	}

	return nil
}

// Get Password
func (svc *UserWebService) GetPassword(ctx context.Context, userID uuid.UUID) (string, error) {
	password, err := svc.store.User.GetPassword(ctx, userID)
	if err != nil {
		return "", err
	}
	return password, nil
}

// ChangePassword updates the password for a user with the specified userID.
func (svc *UserWebService) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// First, retrieve the user from the database
	userDB, err := svc.store.User.GetUser(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "Error fetching user")
	}
	if userDB == nil {
		return errors.Wrap(types.ErrNotFound, fmt.Sprintf("User '%s' not found", userID.String()))
	}

	// Update the user's password
	userDB.Password = newPassword

	// Use the UpdateUser function from the data store to update the user in the database
	_, err = svc.store.User.UpdateUser(ctx, userDB)
	if err != nil {
		return errors.Wrap(err, "Error updating user's password")
	}

	return nil
}

// Checking LogIn
func (svc *UserWebService) GenerateToken(ctx context.Context, nickname, password string) (string, error) {

	user, err := svc.store.User.GetUserByNickname(ctx, nickname)
	if err != nil {
		return "", errors.Wrap(err, "error getting user by nickname")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.Wrap(err, "incorrect password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	return token.SignedString([]byte(signingKey))
}

func (svc *UserWebService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

func (svc *UserWebService) GetUserByNickname(ctx context.Context, nickname string) (*model.User, error) {
	// Log the start of the function for debugging purposes.
	log.Printf("GetUserByNickname: Retrieving user with nickname '%s'", nickname)

	userDB, err := svc.store.User.GetUserByNickname(ctx, nickname)
	if err != nil {
		// Log the error for debugging purposes.
		log.Printf("GetUserByNickname: Error while fetching user: %v", err)
		return nil, errors.Wrap(err, "svc.user.GetUserByNickname")
	}
	if userDB == nil {
		// Log that the user was not found.
		log.Printf("GetUserByNickname: User with nickname '%s' not found", nickname)
		return nil, errors.Wrap(types.ErrNotFound, fmt.Sprintf("User with nickname '%s' not found", nickname))
	}

	// Log the successful retrieval of the user for debugging purposes.
	log.Printf("GetUserByNickname: User with nickname '%s' retrieved successfully", nickname)

	return userDB.ToWeb(), nil
}
