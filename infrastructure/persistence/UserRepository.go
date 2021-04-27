package persistence

import (
	"github.com/smarest/smarest-account/domain/entity"
	"github.com/smarest/smarest-account/domain/repository"
	"github.com/smarest/smarest-common/infrastructure/persistence"
	"gopkg.in/gorp.v3"
)

type UserRepositoryImpl struct {
	*persistence.DAOImpl
}

func NewUserRepository(dbMap *gorp.DbMap) repository.UserRepository {
	return &UserRepositoryImpl{persistence.NewDAOImpl("`user`", dbMap)}
}

func (r *UserRepositoryImpl) FindByUserName(userName string) (*entity.User, error) {
	var user entity.User
	err := r.DbMap.SelectOne(&user, "SELECT * FROM "+r.Table+" WHERE user_name=? AND available=1", userName)

	if err == nil {
		return &user, nil
	}
	return nil, err
}

func (r *UserRepositoryImpl) FindByUserNameAndPassword(userName string, password string) (*entity.User, error) {
	var user entity.User
	err := r.DbMap.SelectOne(&user, "SELECT * FROM "+r.Table+" WHERE user_name=? AND password=? AND available=1", userName, password)

	if err == nil {
		return &user, nil
	}
	return nil, err
}
