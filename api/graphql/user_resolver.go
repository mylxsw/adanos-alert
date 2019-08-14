package graphql

import (
	"context"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/coll"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) AddUser(ctx context.Context, user NewUser) (*User, error) {
	id, err := r.userRepo.Add(UserToRepo(repository.User{}, user))
	if err != nil {
		return nil, err
	}

	newUser, err := r.userRepo.Get(id)
	if err != nil {
		return nil, err
	}

	return RepoToUser(newUser), err
}

func (r *mutationResolver) UpdateUser(ctx context.Context, id string, user NewUser) (*User, error) {
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	original, err := r.userRepo.Get(mid)
	if err != nil {
		return nil, err
	}

	if err = r.userRepo.Update(mid, UserToRepo(original, user)); err != nil {
		return nil, err
	}

	newUser, err := r.userRepo.Get(mid)
	if err != nil {
		return nil, err
	}

	return RepoToUser(newUser), err
}

func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*User, error) {
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user, err := r.userRepo.Get(mid)
	if err != nil {
		return nil, err
	}

	if err = r.userRepo.DeleteID(mid); err != nil {
		return nil, err
	}

	return RepoToUser(user), nil
}

func (r *queryResolver) Users(ctx context.Context, first *int, limit *int) (*UserPagination, error) {
	users, next, err := r.userRepo.Paginate(bson.M{}, int64(*first), int64(*limit))
	if err != nil {
		return nil, err
	}

	var userPagination UserPagination
	userPagination.Cursor = int(next)
	userPagination.HasNext = next > 0

	if err := coll.Map(users, &userPagination.Users, func(u repository.User) *User {
		return RepoToUser(u)
	}); err != nil {
		return nil, err
	}

	return &userPagination, nil
}

func (r *queryResolver) User(ctx context.Context, id string) (*User, error) {
	mid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	user, err := r.userRepo.Get(mid)
	if err != nil {
		return nil, err
	}

	return RepoToUser(user), nil
}
