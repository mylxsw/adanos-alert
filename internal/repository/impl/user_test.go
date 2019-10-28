package impl_test

import (
	"fmt"
	"testing"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/adanos-alert/internal/repository/impl"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepoTestSuite struct {
	suite.Suite
	repo repository.UserRepo
}

func (u *UserRepoTestSuite) TearDownTest() {
	u.NoError(u.repo.Delete(bson.M{}))
}

func (u *UserRepoTestSuite) SetupTest() {
	db, err := Database()
	u.NoError(err)

	u.repo = impl.NewUserRepo(db)
}

func (u *UserRepoTestSuite) TestUserRepo() {
	user := repository.User{
		Name: "Friday",
		Metas: []repository.UserMeta{
			{Key: "phone", Value: "1111111111",},
			{Key: "email", Value: "mylxsw@aicode.cc",},
		},
		Status: repository.UserStatusEnabled,
	}

	// add user
	id, err := u.repo.Add(user)
	u.NoError(err)
	u.NotEmpty(id.String())

	// get user
	user2, err := u.repo.Get(id)
	u.NoError(err)
	u.Equal(user.Name, user2.Name)
	u.NotEmpty(user2.CreatedAt)
	u.NotEmpty(user2.UpdatedAt)

	{
		users, err := u.repo.Find(bson.M{"name": "Friday", "metas.value": "1111111111"})
		u.NoError(err)
		u.Len(users, 1)
	}

	// not found user
	_, err = u.repo.Get(primitive.NewObjectID())
	u.Error(err)
	u.Equal(repository.ErrNotFound, err)

	// batch add users
	for i := 0; i < 10; i++ {
		user.Name = fmt.Sprintf("Friday %d", i/2)

		id, err := u.repo.Add(user)
		u.NoError(err)
		u.NotEmpty(id.String())
	}

	// count
	userCount, err := u.repo.Count(bson.M{})
	u.NoError(err)
	u.EqualValues(11, userCount)

	// Find
	users, err := u.repo.Find(bson.M{"name": "Friday 2"})
	u.NoError(err)
	u.EqualValues(2, len(users))

	// Paginate
	users, next, err := u.repo.Paginate(bson.M{}, 0, 5)
	u.NoError(err)
	u.EqualValues(5, len(users))
	u.EqualValues(5, next)

	users, next, err = u.repo.Paginate(bson.M{}, next, 1000)
	u.NoError(err)
	u.EqualValues(0, next)
	u.EqualValues(6, len(users))

	// Update
	user.Name = "Saturday"
	u.NoError(u.repo.Update(id, user))

	user4, err := u.repo.Get(id)
	u.NoError(err)
	u.Equal("Saturday", user4.Name)

	// Delete
	u.NoError(u.repo.DeleteID(id))
	userCount, err = u.repo.Count(bson.M{})
	u.NoError(err)
	u.EqualValues(10, userCount)
}

func TestUserRepo(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}
