package persistence

import (
	"github.com/smarest/smarest-account/domain/entity"
	"github.com/smarest/smarest-account/domain/repository"
	"github.com/smarest/smarest-common/infrastructure/persistence"
	"gopkg.in/gorp.v3"
)

type RestaurantRepositoryImpl struct {
	*persistence.DAOImpl
}

func NewRestaurantRepository(dbMap *gorp.DbMap) repository.RestaurantRepository {
	return &RestaurantRepositoryImpl{persistence.NewDAOImpl("`restaurant`", dbMap)}
}

func (r *RestaurantRepositoryImpl) FindByIDAndAccessKey(id int64, accessKey string) (*entity.Restaurant, error) {
	var restaurant entity.Restaurant
	err := r.DbMap.SelectOne(&restaurant, "SELECT * FROM "+r.Table+" WHERE id=? AND access_key=? AND available=1", id, accessKey)

	if err == nil {
		return &restaurant, nil
	}
	return nil, err
}
