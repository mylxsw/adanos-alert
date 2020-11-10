package impl

import (
	"context"
	"time"

	"github.com/mylxsw/adanos-alert/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DingdingRobotRepo struct {
	col *mongo.Collection
}

func NewDingdingRobotRepo(db *mongo.Database) repository.DingdingRobotRepo {
	return &DingdingRobotRepo{col: db.Collection("dingding_robot")}
}

func (u DingdingRobotRepo) Add(robot repository.DingdingRobot) (id primitive.ObjectID, err error) {
	robot.CreatedAt = time.Now()
	robot.UpdatedAt = robot.CreatedAt

	rs, err := u.col.InsertOne(context.TODO(), robot)
	if err != nil {
		return
	}

	return rs.InsertedID.(primitive.ObjectID), nil
}

func (u DingdingRobotRepo) Get(id primitive.ObjectID) (robot repository.DingdingRobot, err error) {
	err = u.col.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&robot)
	if err == mongo.ErrNoDocuments {
		err = repository.ErrNotFound
	}

	return
}

func (u DingdingRobotRepo) Find(filter bson.M) (robots []repository.DingdingRobot, err error) {
	robots = make([]repository.DingdingRobot, 0)
	cur, err := u.col.Find(context.TODO(), filter, options.Find().SetSort(bson.M{"name": -1}))
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var robot repository.DingdingRobot
		if err = cur.Decode(&robot); err != nil {
			return
		}

		robots = append(robots, robot)
	}

	return
}

func (u DingdingRobotRepo) Paginate(filter bson.M, offset, limit int64) (robots []repository.DingdingRobot, next int64, err error) {
	robots = make([]repository.DingdingRobot, 0)
	cur, err := u.col.Find(context.TODO(), filter, options.Find().SetSkip(offset).SetLimit(limit).SetSort(bson.M{"name": -1}))
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var robot repository.DingdingRobot
		if err = cur.Decode(&robot); err != nil {
			return
		}

		robots = append(robots, robot)
	}

	if int64(len(robots)) == limit {
		next = offset + limit
	}

	return
}

func (u DingdingRobotRepo) DeleteID(id primitive.ObjectID) error {
	return u.Delete(bson.M{"_id": id})
}

func (u DingdingRobotRepo) Delete(filter bson.M) error {
	_, err := u.col.DeleteMany(context.TODO(), filter)
	return err
}

func (u DingdingRobotRepo) Update(id primitive.ObjectID, robot repository.DingdingRobot) error {
	robot.UpdatedAt = time.Now()
	_, err := u.col.ReplaceOne(context.TODO(), bson.M{"_id": id}, robot)
	return err
}

func (u DingdingRobotRepo) Count(filter bson.M) (int64, error) {
	return u.col.CountDocuments(context.TODO(), filter)
}
