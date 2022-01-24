package web

import (
  "context"
  "github.com/gorilla/mux"
  "github.com/peter-mount/documentation/tools/hugo"
  "log"
  "net/http"
  "strconv"
)

// Webserver provides a webserver if required for other tasks.
// It serves the public directory
type Webserver struct {
  config *hugo.Config `kernel:"inject"` // Config
  server *http.Server // enabled
}

func (w *Webserver) Name() string {
  return "webserver"
}

func (w *Webserver) Start() error {
  log.Println("Starting webserver")

  router := mux.NewRouter()
  router.PathPrefix("/").
    Handler(http.FileServer(http.Dir("public")))

  w.server = &http.Server{
    Addr:    w.config.Webserver.Address + ":" + strconv.Itoa(w.config.Webserver.Port),
    Handler: router,
  }

  go func() {
    _ = w.server.ListenAndServe()
  }()

  return nil
}

func (w *Webserver) Stop() {
  if w.server != nil {
    log.Println("Stopping webserver")
    _ = w.server.Shutdown(context.Background())
  }
}
