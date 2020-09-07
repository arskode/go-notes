package postgres

import (
	"database/sql"

	"github.com/arskode/go-notes/api/model"
	"github.com/jmoiron/sqlx"
)

type NoteStore struct {
	DB *sqlx.DB
}

func (ns *NoteStore) Create(note *model.Note) (*model.Note, error) {
	query := `INSERT INTO "notes" (title, description) VALUES ($1, $2) RETURNING id`
	err := ns.DB.QueryRowx(query, note.Title, note.Description).Scan(&note.ID)
	if err != nil {
		return nil, err
	}
	return note, nil
}

func (ns *NoteStore) Get(noteID uint64) (*model.Note, error) {
	var note model.Note
	query := `SELECT id, title, description FROM "notes" WHERE id = $1`
	err := ns.DB.Get(&note, query, noteID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err

	}
	return &note, nil
}

func (ns *NoteStore) List() ([]model.Note, error) {
	notes := []model.Note{}

	query := `SELECT id, title, description FROM "notes" LIMIT 100`
	err := ns.DB.Select(&notes, query)
	if err != nil {
		return nil, err
	}

	return notes, nil
}

func (ns *NoteStore) Update(noteID uint64, note *model.Note) (*model.Note, error) {
	query := `UPDATE "notes" SET title=$1, description=$2 WHERE id = $3 RETURNING id`
	err := ns.DB.QueryRowx(query, note.Title, note.Description, noteID).Scan(&note.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return note, nil
}

func (ns *NoteStore) Delete(noteID uint64) (bool, error) {
	var deleted bool
	query := `DELETE FROM "notes" WHERE id = $1 RETURNING true`
	err := ns.DB.QueryRowx(query, noteID).Scan(&deleted)
	if err != nil {
		if err == sql.ErrNoRows {
			return deleted, nil
		}
		return deleted, err
	}
	return deleted, nil
}
