package model

import (
	"errors"
	"strings"
)

type Note struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (n *Note) Prepare() {
	n.ID = 0
	n.Title = strings.TrimSpace(n.Title)
	n.Description = strings.TrimSpace(n.Description)
}

func (n *Note) Validate() error {
	if n.Title == "" {
		return errors.New("Required Title")
	}
	if n.Description == "" {
		return errors.New("Required Description")
	}
	return nil
}
