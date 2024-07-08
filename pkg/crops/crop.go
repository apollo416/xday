package crops

import (
	"time"

	"github.com/google/uuid"
)

type Crop struct {
	Id             uuid.UUID `json:"id"`
	Status         string    `json:"status"`
	Cultivar       uuid.UUID `json:"cultivar"`
	CultivarStart  time.Time `json:"cultivar_start"`
	CultivarEnd    time.Time `json:"cultivar_end"`
	MaturationTime int       `json:"maturation_time"`
	Created        time.Time `json:"created"`
	Generation     int       `json:"generation"`
}
