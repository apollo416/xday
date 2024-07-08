package products

import "time"

type Product struct {
	Id             string    `json:"id"`
	Slug           string    `json:"slug"`
	Name           string    `json:"name"`
	Source         string    `json:"source"`
	Yield          string    `json:"yield"`
	SellPrice      string    `json:"sell_price"`
	MaturationTime time.Time `json:"maturation_time"`
}
