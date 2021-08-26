package hugo

import (
  "context"
  "github.com/gorilla/mux"
  "log"
  "net/http"
)

// Webserver provides a webserver if required for other tasks.
// It serves the public directory
type Webserver struct {
  enabled bool         // true then run webserver
  server  *http.Server // enabled
}

// Enable enables the webserver if v is true. If false nothing happens as this is
// a one way action. e.g. Enable(false) will not disable if it's already enabled.
func (w *Webserver) Enable(v bool) {
  // Only set, never disable
  if v {
    w.enabled = true
  }
}

func (w *Webserver) Name() string {
  return "webserver"
}

func (w *Webserver) Start() error {
  if w.enabled {
    log.Println("Starting webserver")

    router := mux.NewRouter()
    router.PathPrefix("/").
      Handler(http.FileServer(http.Dir("public")))

    w.server = &http.Server{
      Addr:    "0.0.0.0:8080",
      Handler: router,
    }

    go func() {
      _ = w.server.ListenAndServe()
    }()

  }

  return nil
}

func (w *Webserver) Stop() {
  if w.server != nil {
    log.Println("Stopping webserver")
    err := w.server.Shutdown(context.Background())
    log.Println(err)
  }
}
