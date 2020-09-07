package rest_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/arskode/go-notes/api/config"
	"github.com/arskode/go-notes/api/model"
	"github.com/arskode/go-notes/api/rest"
	"github.com/arskode/go-notes/api/store"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNote(t *testing.T) {
	ts, _, teardown := startup(t)
	defer teardown("notes")

	resp, err := reqDo(t, ts.URL+"/notes", "POST", `{"title": "title1", "description": "description1"}`)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var m map[string]interface{}
	err = json.Unmarshal(body, &m)
	assert.NoError(t, err)

	require.Equal(t, http.StatusCreated, resp.StatusCode, string(body))
	require.Equal(t, m["title"], "title1", string(body))
	require.Equal(t, m["description"], "description1", string(body))
}

func TestGetNote(t *testing.T) {
	ts, srv, teardown := startup(t)
	defer teardown("notes")

	_, err := seedOneNote(srv)
	assert.NoError(t, err)
	resp, err := reqDo(t, ts.URL+"/notes/1", "GET", "")
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var note model.Note
	err = json.Unmarshal(body, &note)
	assert.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, note.ID, uint64(1), string(body))
	require.Equal(t, note.Title, "title1", string(body))
	require.Equal(t, note.Description, "description1", string(body))

}

func TestListNote(t *testing.T) {
	ts, srv, teardown := startup(t)
	defer teardown("notes")

	_, err := seedNotes(srv)
	assert.NoError(t, err)

	resp, err := reqDo(t, ts.URL+"/notes", "GET", "")
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var notes []model.Note
	err = json.Unmarshal(body, &notes)
	assert.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, len(notes), 2, string(body))

}

func TestUpdateNote(t *testing.T) {
	ts, srv, teardown := startup(t)
	defer teardown("notes")

	_, err := seedOneNote(srv)
	assert.NoError(t, err)
	resp, err := reqDo(t, ts.URL+"/notes/1", "PUT", `{"title": "title3", "description": "description3"}`)
	assert.NoError(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	var note model.Note
	err = json.Unmarshal(body, &note)
	assert.NoError(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	require.Equal(t, note.ID, uint64(1), string(body))
	require.Equal(t, note.Title, "title3", string(body))
	require.Equal(t, note.Description, "description3", string(body))
}

func TestDeleteNote(t *testing.T) {
	ts, srv, teardown := startup(t)
	defer teardown("notes")

	_, err := seedOneNote(srv)
	assert.NoError(t, err)
	resp, err := reqDo(t, ts.URL+"/notes/1", "DELETE", "")
	assert.NoError(t, err)

	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func reqDo(t *testing.T, url, method string, body string) (*http.Response, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	assert.NoError(t, err)
	return client.Do(req)
}

func startup(t *testing.T) (*httptest.Server, *rest.Server, func(...string)) {

	conf := config.Config{DbURL: "postgresql://postgres:postgres@pg/postgres?sslmode=disable"}

	db, err := sqlx.Connect("postgres", conf.DbURL)
	assert.Equal(t, nil, err)

	store := store.NewStore(db)

	srv := &rest.Server{Config: &conf, Store: store}
	ts := httptest.NewServer(srv.Routes())

	teardown := func(tables ...string) {
		if len(tables) > 0 {
			_, _ = db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
			_, _ = db.Exec(`ALTER SEQUENCE notes_id_seq RESTART WITH 1`)
		}
		db.Close()
	}
	return ts, srv, teardown
}

func seedOneNote(srv *rest.Server) (model.Note, error) {

	note := model.Note{
		Title:       "title1",
		Description: "description1",
	}

	_, err := srv.Store.Note.Create(&note)
	if err != nil {
		return model.Note{}, err
	}
	return note, nil
}

func seedNotes(srv *rest.Server) ([]model.Note, error) {
	notes := []model.Note{
		model.Note{
			Title:       "title1",
			Description: "description1",
		},
		model.Note{
			Title:       "title2",
			Description: "description2",
		},
	}
	for _, note := range notes {
		_, err := srv.Store.Note.Create(&note)
		if err != nil {
			return []model.Note{}, err
		}
	}
	return notes, nil
}
