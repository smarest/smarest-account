package repository

import (
	"github.com/smarest/smarest-account/domain/entity"
)

type UserRepository interface {
	FindByUserName(userName string) (*entity.User, error)
	FindByUserNameAndPassword(userName string, password string) (*entity.User, error)
}
