package repository

import (
	"github.com/smarest/smarest-account/domain/entity"
)

type RestaurantRepository interface {
	FindByIDAndAccessKey(id int64, accessKey string) (*entity.Restaurant, error)
}
