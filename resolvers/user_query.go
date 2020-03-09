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
	First *int32
	After *string
}) (*[]*userResolver, error) {
	// Auth
	if isAuthorized, err := utils.IsAuthorized(ctx, utils.AuthUser); !isAuthorized {
		return nil, err
	}

	users, err := r.UserService.List(args.First)
	if err != nil {
		return nil, err
	}

	var resolvers = make([]*userResolver, 0, len(users))

	for _, user := range users {
		resolver := &userResolver{user}
		resolvers = append(resolvers, resolver)
	}

	return &resolvers, nil
}
