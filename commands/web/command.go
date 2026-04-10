package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebUi struct {
}

// Description implements [main.subcommand].
func (w *WebUi) Description() (desc string) {
	panic("unimplemented")
}

// Name implements [main.subcommand].
func (w *WebUi) Name() (name string) { return "web" }

// Run implements [main.subcommand].
func (w *WebUi) Run(args []string) (result any, err error) {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if globeError != nil {
			fmt.Fprintf(w, "Error:%s", globeError)
		} else {
			fmt.Fprintf(w, "we received %d", count)
		}
	})

	go SyncPGS()

	http.ListenAndServe(":3000", r)

	return nil, nil
}

// Usage implements [main.subcommand].
func (w *WebUi) Usage() (usage string) {
	panic("unimplemented")
}
