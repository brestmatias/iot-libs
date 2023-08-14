package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/brestmatias/iot-libs/model"
	"github.com/brestmatias/iot-libs/wrappers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SensorLastStatusRepository interface {
	UpsertReportedStatus(stationId string, interfaceId string, value string) int64
	FindByField(field string, value interface{}) *[]model.SensorLastStatus
}

type sensorLastStatusRepository struct {
	MongoDB    *mongo.Database
	Collection *mongo.Collection
}

func NewSensorLastStatusRepository(mongodb *wrappers.MongoClientWrapper) SensorLastStatusRepository {
	return &sensorLastStatusRepository{
		MongoDB:    mongodb.GetDatabase(),
		Collection: mongodb.GetDatabase().Collection("sensor_last_status"),
	}
}

func (h *sensorLastStatusRepository) UpsertReportedStatus(stationId string, sensorID string, value string) int64 {
	method := "UpsertDispatcherStatus"
	filter := bson.M{"sensor_id": sensorID, "station_id": stationId}

	update := bson.D{{"$set", bson.D{
		{"sensor_id", sensorID},
		{"station_id", stationId},
		{"reported_value", value},
		{"last_report", primitive.NewDateTimeFromTime(time.Now())},
	}}}
	opts := options.Update().SetUpsert(true)
	result, err := h.Collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		log.Println(fmt.Errorf("[method:%s]Error!!", method), err)
	}
	return result.UpsertedCount
}

func (h *sensorLastStatusRepository) FindByField(field string, value interface{}) *[]model.SensorLastStatus {
	var result []model.SensorLastStatus
	filter := bson.M{field: value}
	findResult, err := h.Collection.Find(context.Background(), filter, nil)
	if err != nil {
		log.Println(err)
	}
	if findResult.Err() != nil && errors.Is(findResult.Err(), mongo.ErrNoDocuments) {
		return nil
	}

	defer findResult.Close(context.TODO())
	err = findResult.All(context.TODO(), &result)

	if err != nil {
		log.Println(err)
	}
	return &result
}
