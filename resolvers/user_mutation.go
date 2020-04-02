package resolvers

import (
	"context"
	"fmt"
	"time"

	"github.com/matiss/go-graphql-server/models"
	"github.com/matiss/go-graphql-server/utils"
)

func (r *Resolver) CreateUser(ctx context.Context, args *struct {
	Email    string
	Password string
	Name     string
}) (*userResolver, error) {
	user := &models.User{
		Email:    args.Email,
		Password: args.Password,
		Name:     args.Name,
	}

	user, err := r.UserService.Create(user)
	if err != nil {
		return nil, err
	}

	return &userResolver{user}, nil
}

func (r *Resolver) LoginUser(ctx context.Context, args *struct {
	Email    string
	Password string
}) (*userLoginResolver, error) {
	user, err := r.UserService.ComparePassword(args.Email, args.Password)
	if err != nil {
		return nil, err
	}

	ipAddress := ctx.Value("IP").(string)

	err = r.UserService.UpdateLogin(user, ipAddress)
	if err != nil {
		return nil, err
	}

	// Create JWT token
	exp := time.Now().Add(time.Second * time.Duration(r.TokenTTL)).Unix()
	token, err := utils.GenerateJWT(exp, r.JWTSecret, user.ID, utils.AuthLevel(user.Role))
	if err != nil {
		return nil, err
	}

	return &userLoginResolver{token}, nil
}

func (r *Resolver) RenewToken(ctx context.Context) (*userLoginResolver, error) {
	// Auth
	if isAuthorized, err := utils.IsAuthorized(ctx, utils.AuthUser); !isAuthorized {
		return nil, err
	}

	userID, ok := ctx.Value("userID").(int32)
	if !ok || userID == 0 {
		return nil, fmt.Errorf("Invalid userID")
	}

	// Find user
	user, err := r.UserService.Find(userID)
	if err != nil || !user.Active() {
		return nil, fmt.Errorf("Could not find user with id %d", userID)
	}

	// Create JWT token
	exp := time.Now().Add(time.Second * time.Duration(r.TokenTTL)).Unix()
	tokenStr, err := utils.GenerateJWT(exp, r.JWTSecret, user.ID, utils.AuthLevel(user.Role))
	if err != nil {
		return nil, err
	}

	return &userLoginResolver{tokenStr}, nil
}
