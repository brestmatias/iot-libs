package repository

import (
	"context"
	"log"
	"time"

	"github.com/brestmatias/iot-libs/model"
	"github.com/brestmatias/iot-libs/wrappers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SensorHistoryRepository interface {
	InsertOne(in model.SensorHistory) (*model.SensorHistory, error)
	RemoveOlderThanDays(days int) (int64, error)
}

type sensorHistoryRepository struct {
	MongoDB    *mongo.Database
	Collection *mongo.Collection
}

func NewSensorHistoryRepository(mongodb *wrappers.MongoClientWrapper) SensorHistoryRepository {
	return &sensorHistoryRepository{
		MongoDB:    mongodb.GetDatabase(),
		Collection: mongodb.GetDatabase().Collection("sensor_history"),
	}
}

func (h *sensorHistoryRepository) InsertOne(in model.SensorHistory) (*model.SensorHistory, error) {
	res, err := h.getStationCollection().InsertOne(context.Background(), in)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	in.DocId = res.InsertedID.(primitive.ObjectID)

	return &in, nil
}

func (h *sensorHistoryRepository) RemoveOlderThanDays(days int) (int64, error) {
	now := time.Now()
	start := now.AddDate(0, 0, days*-1)
	minutes := now.Sub(start).Minutes()

	filter := bson.D{{"date", bson.D{{"$lt", minutes}}}}
	res, err := h.getStationCollection().DeleteMany(context.Background(), filter)

	if err != nil {
		log.Println(err)
		return 0, err
	}
	return res.DeletedCount, nil
}

func (s *sensorHistoryRepository) getStationCollection() *mongo.Collection {
	return s.MongoDB.Collection("sensor_history")
}
