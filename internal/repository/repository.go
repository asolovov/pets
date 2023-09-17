package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"pets/internal/config"
	"pets/internal/model"
	"pets/pkg/logger"
)

// IRepository is a repository layer interface
type IRepository interface {
	// GetPets is used to get pet from DB. Pagination can be used by setting limit and offset values. Order should be
	// "asc" or "desc" in any register, all other values will be ignored. 0 limit will be ignored.
	GetPets(limit int, offset int, order string) (pets []*model.Pet, err error)
	// GetPet is used to get pet from DB by given ID
	GetPet(id int) (pet *model.Pet, err error)
	// AddPet is used to add new pet to the DB. Only "name" field will be used. Fields id and created_at will be set automatically
	AddPet(pet *model.Pet) error
	// UpdatePet is used to update existing pet to the DB by given id filed. Only "name" field will be used. Fields id and
	// updated_at will be set automatically
	UpdatePet(pet *model.Pet) error
	// DeletePet is used to delete pet from the DB by given id
	DeletePet(pet *model.Pet) error
	// Stop is used to stop repository work
	Stop()
}

// Repository is a repository struct, implements IRepository interface
type Repository struct {
	db *sqlx.DB
}

// NewRepository is used to get new Repository instance
func NewRepository(conf *config.DB) IRepository {
	if conf == nil {
		logger.Log().WithField("layer", "Repository-Init").Fatalf("nil config err")
	}

	db, err := sqlx.Open(conf.Driver, conf.Addr)
	if err != nil {
		logger.Log().WithField("layer", "Repository-Init").Fatalf("err open db: %v", err.Error())
	}

	logger.Log().WithField("layer", "Repository-Init").Infof("connected to %v", conf.Driver)

	if err = db.Ping(); err != nil {
		logger.Log().WithField("layer", "Repository-Init").Fatalf("err ping db: %v", err.Error())
	}

	logger.Log().WithField("layer", "Repository-Init").Infof("ping db ok")

	return &Repository{
		db: db,
	}
}

// Stop is implementing IRepository.Stop function. It will close the DB connection and log error if occurred
func (r *Repository) Stop() {
	if r.db != nil {
		if err := r.db.Close(); err != nil {
			logger.Log().WithField("layer", "Repository-Stop").Warningf("err closing db: %v", err.Error())
		} else {
			logger.Log().WithField("layer", "Repository-Stop").Infof("db closed")
		}
	}
}
