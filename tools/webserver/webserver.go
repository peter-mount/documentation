package webserver

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/peter-mount/go-kernel/v2/log"
	"net/http"
	"strconv"
)

// Webserver provides a webserver if required for other tasks.
// It serves the public directory
type Webserver struct {
	Foreground bool             // If true run in foreground, used by webserver tool
	config     *WebserverConfig `kernel:"config,webserver"`
	server     *http.Server
}

type WebserverConfig struct {
	Address string `yaml:"address"` // Webserver bind address, default "127.0.0.1"
	Port    int    `yaml:"port"`    // Port defaults to 8080
}

func (w *Webserver) Start() error {
	log.Println("Starting webserver")

	router := mux.NewRouter()
	router.PathPrefix("/").
		Handler(http.FileServer(http.Dir("public")))

	w.server = &http.Server{
		Addr:    w.config.Address + ":" + strconv.Itoa(w.config.Port),
		Handler: router,
	}

	if w.Foreground {
		return w.server.ListenAndServe()
	} else {
		go func() {
			_ = w.server.ListenAndServe()
		}()

		return nil
	}
}

func (w *Webserver) Stop() {
	if w.server != nil {
		log.Println("Stopping webserver")
		_ = w.server.Shutdown(context.Background())
	}
}

func (c *Webserver) WebPath(f string, a ...interface{}) string {
	return fmt.Sprintf("http://%s:%d/"+f, append([]interface{}{c.config.Address, c.config.Port}, a...)...)
}
