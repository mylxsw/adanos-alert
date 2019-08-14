package graphql

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/container"
)

type Resolver struct {
	cc       *container.Container
	ruleRepo repository.RuleRepo
	userRepo repository.UserRepo
}

func NewResolver(cc *container.Container) *Resolver {
	res := Resolver{cc: cc}
	cc.MustResolve(func(ruleRepo repository.RuleRepo, userRepo repository.UserRepo) {
		res.ruleRepo = ruleRepo
		res.userRepo = userRepo
	})

	return &res
}
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }