package repository

import "plairsty/backend/domain/entity"

type UserRepository interface {
	SaveUser(server *entity.Server)
}
