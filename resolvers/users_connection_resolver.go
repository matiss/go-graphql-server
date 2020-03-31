package resolvers

import (
	"github.com/matiss/go-graphql-server/models"
)

type usersConnectionResolver struct {
	users      []*models.User
	totalCount int
}

func (r *usersConnectionResolver) TotalCount() int32 {
	return int32(r.totalCount)
}

func (r *usersConnectionResolver) Edges() *[]*userResolver {
	l := make([]*userResolver, len(r.users))
	for i := range l {
		l[i] = &userResolver{
			u: r.users[i],
		}
	}
	return &l
}
