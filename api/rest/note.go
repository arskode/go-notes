package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/arskode/go-notes/api/model"
	"github.com/arskode/go-notes/api/responses"
	"github.com/go-chi/chi"
)

func (s *Server) createNote(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	note := model.Note{}
	err = json.Unmarshal(body, &note)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	note.Prepare()
	err = note.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	var noteCreated *model.Note
	noteCreated, err = s.Store.Note.Create(&note)
	if err != nil {
		fmt.Println(err)

		responses.ERROR(w, http.StatusInternalServerError, errors.New("Something went wrong"))
		return
	}

	responses.JSON(w, http.StatusCreated, noteCreated)
}

func (s *Server) getNote(w http.ResponseWriter, r *http.Request) {

	noteID, err := strconv.ParseUint(chi.URLParam(r, "noteID"), 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	noteGotten, err := s.Store.Note.Get(noteID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Something went wrong"))
		return
	}
	if noteGotten == nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Note Not Found"))
		return
	}
	responses.JSON(w, http.StatusOK, noteGotten)
}

func (s *Server) listNotes(w http.ResponseWriter, r *http.Request) {

	notes, err := s.Store.Note.List()
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Something went wrong"))
		return
	}

	responses.JSON(w, http.StatusOK, notes)
}

func (s *Server) updateNote(w http.ResponseWriter, r *http.Request) {

	noteID, err := strconv.ParseUint(chi.URLParam(r, "noteID"), 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	note := model.Note{}
	err = json.Unmarshal(body, &note)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	note.Prepare()
	err = note.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	updatedNote, err := s.Store.Note.Update(noteID, &note)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Something went wrong"))
		return
	}
	if updatedNote == nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Note Not Found"))
		return
	}

	responses.JSON(w, http.StatusOK, updatedNote)
}

func (s *Server) deleteNote(w http.ResponseWriter, r *http.Request) {

	noteID, err := strconv.ParseUint(chi.URLParam(r, "noteID"), 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	deleted, err := s.Store.Note.Delete(noteID)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, errors.New("Something went wrong"))
		return
	}
	if !deleted {
		responses.ERROR(w, http.StatusNotFound, errors.New("Note Not Found"))
		return
	}

	responses.JSON(w, http.StatusNoContent, "")
}
