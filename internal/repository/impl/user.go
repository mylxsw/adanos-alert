package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"github.com/mylxsw/coll"
	"github.com/mylxsw/go-utils/str"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepo struct {
	col *mongo.Collection
}

func NewUserRepo(db *mongo.Database) repository.UserRepo {
	return &UserRepo{col: db.Collection("user")}
}

func (u UserRepo) Add(user repository.User) (id primitive.ObjectID, err error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	if user.Status == "" {
		user.Status = repository.UserStatusEnabled
	}

	rs, err := u.col.InsertOne(context.TODO(), user)
	if err != nil {
		return
	}

	return rs.InsertedID.(primitive.ObjectID), nil
}

func (u UserRepo) Get(id primitive.ObjectID) (user repository.User, err error) {
	err = u.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	return
}

func (u UserRepo) GetByEmail(email string) (user repository.User, err error) {
	err = u.col.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	return
}

func (u UserRepo) Find(filter bson.M) (users []repository.User, err error) {
	users = make([]repository.User, 0)
	cur, err := u.col.Find(context.TODO(), filter)
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user repository.User
		if err = cur.Decode(&user); err != nil {
			return
		}

		users = append(users, user)
	}

	return
}

func (u UserRepo) Paginate(filter bson.M, offset, limit int64) (users []repository.User, next int64, err error) {
	users = make([]repository.User, 0)
	cur, err := u.col.Find(context.TODO(), filter, options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"created_at": -1}))
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var user repository.User
		if err = cur.Decode(&user); err != nil {
			return
		}

		users = append(users, user)
	}

	if int64(len(users)) == limit {
		next = offset + limit
	}

	return
}

func (u UserRepo) DeleteID(id primitive.ObjectID) error {
	return u.Delete(bson.M{"_id": id})
}

func (u UserRepo) Delete(filter bson.M) error {
	_, err := u.col.DeleteMany(context.TODO(), filter)
	return err
}

func (u UserRepo) Update(id primitive.ObjectID, user repository.User) error {
	old, err := u.Get(id)
	if err != nil {
		return err
	}
	user.CreatedAt = old.CreatedAt
	user.UpdatedAt = time.Now()
	if user.Password == "" {
		user.Password = old.Password
	}

	_, err = u.col.ReplaceOne(context.TODO(), bson.M{"_id": id}, user)
	return err
}

func (u UserRepo) Count(filter bson.M) (int64, error) {
	return u.col.CountDocuments(context.TODO(), filter)
}

func (u UserRepo) getUserIDWithMetas(queryK string, queryV interface{}, field string) ([]repository.UserIDWithMeta, error) {
	filter := bson.M{}
	if queryK == "id" || queryK == "_id" {
		id, err := primitive.ObjectIDFromHex(queryV.(string))
		if err != nil {
			return nil, err
		}

		filter["_id"] = id
	} else if str.In(queryK, []string{"name", "phone", "email", "role", "status"}) {
		filter[queryK] = queryV
	} else {
		filter["metas.key"] = queryK
		filter["metas.value"] = queryV
	}

	users, err := u.Find(filter)
	if err != nil {
		return nil, err
	}

	var res []repository.UserIDWithMeta
	_ = coll.MustNew(users).Map(func(u repository.User) repository.UserIDWithMeta {
		res := repository.UserIDWithMeta{UserID: u.ID.Hex(), Meta: make([]string, 0)}
		switch field {
		case "id", "_id":
			res.Meta = append(res.Meta, u.ID.Hex())
		case "name":
			res.Meta = append(res.Meta, u.Name)
		case "phone":
			res.Meta = append(res.Meta, u.Phone)
		case "email":
			res.Meta = append(res.Meta, u.Email)
		case "role":
			res.Meta = append(res.Meta, u.Role)
		case "status":
			res.Meta = append(res.Meta, string(u.Status))
		default:
			for _, m := range u.Metas {
				if m.Key == field {
					res.Meta = append(res.Meta, m.Value)
				}
			}
		}

		return res
	}).Filter(func(v repository.UserIDWithMeta) bool {
		for _, m := range v.Meta {
			if m != "" {
				return true
			}
		}

		return false
	}).All(&res)
	return res, nil
}

func (u UserRepo) GetUserIDWithMetasRegex(queryK, queryVRegex, field string) ([]repository.UserIDWithMeta, error) {
	if queryK == "id" || queryK == "_id" {
		return u.getUserIDWithMetas(queryK, queryVRegex, field)
	}

	return u.getUserIDWithMetas(queryK, primitive.Regex{Pattern: queryVRegex, Options: "im"}, field)
}

func (u UserRepo) GetUserMetasRegex(queryK, queryVRegex, field string) ([]string, error) {
	userWithMetas, err := u.GetUserIDWithMetasRegex(queryK, queryVRegex, field)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, u := range userWithMetas {
		res = append(res, u.Meta...)
	}

	return str.Distinct(res), nil
}

func (u UserRepo) GetUserIDWithMetas(queryK, queryV, field string) ([]repository.UserIDWithMeta, error) {
	return u.getUserIDWithMetas(queryK, queryV, field)
}

func (u UserRepo) GetUserMetas(queryK, queryV, field string) ([]string, error) {
	userWithMetas, err := u.GetUserIDWithMetas(queryK, queryV, field)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0)
	for _, u := range userWithMetas {
		res = append(res, u.Meta...)
	}

	return str.Distinct(res), nil
}
