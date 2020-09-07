package rest

import (
	"net/http"
	"time"

	"github.com/arskode/go-notes/api/config"
	"github.com/arskode/go-notes/api/responses"
	"github.com/arskode/go-notes/api/store"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	Config *config.Config
	Store  *store.Store
}

func Start() error {

	conf := config.NewConfig()

	db, err := sqlx.Connect("postgres", conf.DbURL)
	if err != nil {
		return err
	}
	defer db.Close()

	store := store.NewStore(db)

	server := Server{Config: conf, Store: store}
	err = server.run()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) run() error {
	httpServer := &http.Server{
		Addr:              ":5000",
		Handler:           s.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
	return httpServer.ListenAndServe()

}

func (s *Server) Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Throttle(1000), middleware.RealIP)
	router.Use(middleware.Recoverer, middleware.RedirectSlashes)

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		responses.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	router.Post("/notes", s.createNote)
	router.Get("/notes", s.listNotes)
	router.Get("/notes/{noteID}", s.getNote)
	router.Put("/notes/{noteID}", s.updateNote)
	router.Delete("/notes/{noteID}", s.deleteNote)

	return router
}
