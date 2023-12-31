package repository

import (
	"context"
	"log"

	"github.com/brestmatias/iot-libs/model"
	"github.com/brestmatias/iot-libs/wrappers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SensorHistoryRepository interface {
	InsertOne(in model.SensorHistory) (*model.SensorHistory, error)
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

func (s *sensorHistoryRepository) getStationCollection() *mongo.Collection {
	return s.MongoDB.Collection("sensor_history")
}
