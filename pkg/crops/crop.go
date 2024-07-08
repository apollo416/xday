package crops

import (
	"time"

	"github.com/google/uuid"
)

type Crop struct {
	ID             uuid.UUID `json:"id" dynamodbav:"id"`
	Status         string    `json:"status" dynamodbav:"status"`
	Cultivar       uuid.UUID `json:"cultivar" dynamodbav:"cultivar"`
	CultivarStart  time.Time `json:"cultivar_start" dynamodbav:"cultivar_start"`
	CultivarEnd    time.Time `json:"cultivar_end" dynamodbav:"cultivar_end"`
	MaturationTime int       `json:"maturation_time" dynamodbav:"maturation_time"`
	Created        time.Time `json:"created" dynamodbav:"created"`
	Generation     int       `json:"generation" dynamodbav:"generation"`
}
