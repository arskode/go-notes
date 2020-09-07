package store

import (
	"github.com/arskode/go-notes/api/model"
	"github.com/arskode/go-notes/api/store/postgres"
	"github.com/jmoiron/sqlx"
)

type NoteStoreInterface interface {
	Create(note *model.Note) (*model.Note, error)
	Get(noteID uint64) (*model.Note, error)
	List() ([]model.Note, error)
	Update(noteID uint64, note *model.Note) (*model.Note, error)
	Delete(noteID uint64) (bool, error)
}

type Store struct {
	Note NoteStoreInterface
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		Note: &postgres.NoteStore{DB: db},
	}
}
