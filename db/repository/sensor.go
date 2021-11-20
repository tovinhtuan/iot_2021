package repository

import (
	"iot-project/db/model"
	"iot-project/db/storage"
	"time"
)

type SensorRepository interface {
	InsertSensor(sensor *model.Sensor) error
	GetSensorById(Id int64) (*model.Sensor, error)
	GetSensorByTime(timeRequest time.Time)(*model.Sensor, error)
}
type sensorRepository struct {
	db *storage.PSQLManager
}
func NewSensorRepository(db *storage.PSQLManager) SensorRepository {
	return &sensorRepository{
		db : db,
	}
}
func(s *sensorRepository) InsertSensor(sensor *model.Sensor) error {
	if err := s.db.Create(sensor).Error; err != nil {
		return err
	}
	return nil
}
func(s *sensorRepository) GetSensorById(id int64) (*model.Sensor, error) {
	sensor := &model.Sensor{}
	if err := s.db.Where(&model.Sensor{Id: id}).First(sensor).Error; err != nil {
		return nil, err
	}
	return sensor, nil
}
func (s *sensorRepository) GetSensorByTime(timeRequest time.Time) (*model.Sensor, error) {
	sensor := &model.Sensor{}
	if err := s.db.Where(&model.Sensor{UpdatedAt: timeRequest}).First(sensor).Error; err != nil {
		return nil, err
	}
	return sensor, nil
}