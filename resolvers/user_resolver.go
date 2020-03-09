package resolvers

import (
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/matiss/go-graphql-server/models"
)

type userResolver struct {
	u *models.User
}

func (r *userResolver) ID() graphql.ID {
	return graphql.ID(r.u.Email)
}

func (r *userResolver) UserID() *int32 {
	return &r.u.ID
}

func (r *userResolver) Email() *string {
	return &r.u.Email
}

func (r *userResolver) Password() *string {
	maskedPassword := "********"
	return &maskedPassword
}

func (r *userResolver) Name() *string {
	return &r.u.Name
}

func (r *userResolver) Status() int32 {
	status := int32(r.u.Status)
	return status
}

func (r *userResolver) Role() *int32 {
	role := int32(r.u.Role)
	return &role
}

func (r *userResolver) CreatedAt() (*graphql.Time, error) {
	return &graphql.Time{Time: r.u.CreatedAt}, nil
}

func (r *userResolver) UpdatedAt() (*graphql.Time, error) {
	return &graphql.Time{Time: r.u.CreatedAt}, nil
}

type userLoginResolver struct {
	token string
}

func (r *userLoginResolver) Token() string {
	return r.token
}
