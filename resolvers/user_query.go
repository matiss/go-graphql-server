package resolvers

import (
	"context"
	"fmt"

	"github.com/matiss/go-graphql-server/utils"
)

func (r *Resolver) User(ctx context.Context, args struct{ Email string }) (*userResolver, error) {
	// Auth
	if isAuthorized, err := utils.IsAuthorized(ctx, utils.AuthUser); !isAuthorized {
		return nil, err
	}

	user, err := r.UserService.FindByEmail(args.Email)
	if err != nil {
		return nil, fmt.Errorf("Could not find user with email %s", args.Email)
	}

	return &userResolver{user}, nil
}

func (r *Resolver) Users(ctx context.Context, args struct {
	Limit  *int32
	Offset *int32
	Role   *int32
}) (*usersConnectionResolver, error) {
	// Auth
	if isAuthorized, err := utils.IsAuthorized(ctx, utils.AuthAdmin); !isAuthorized {
		return nil, err
	}

	limit := 0
	if args.Limit != nil {
		limit = int(*args.Limit)
	}

	offset := 0
	if args.Offset != nil {
		offset = int(*args.Offset)
	}

	userRole := 0
	if args.Role != nil {
		userRole = int(*args.Role)
	}

	users, err := r.UserService.List(limit, offset, userRole)
	if err != nil {
		return nil, err
	}

	count, err := r.UserService.Count(userRole)
	if err != nil {
		return nil, err
	}

	return &usersConnectionResolver{users: users, totalCount: count}, nil
}
