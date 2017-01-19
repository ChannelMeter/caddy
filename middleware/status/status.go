package status

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/mholt/caddy/middleware"
)

var launchTime = time.Now()

type Status struct {
	Next    middleware.Handler
	Configs []Config
}

type Config struct {
	PathScope string
	Code      int
	Body      string
}

// ServeHTTP implements the middleware.Handler interface.
func (s Status) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	for _, sc := range s.Configs {
		if !middleware.Path(r.URL.Path).Matches(sc.PathScope) {
			continue
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("X-Start-Time", launchTime.String())
		w.Header().Set("X-Uptime", time.Now().Sub(launchTime).String())
		w.Header().Set("X-Mesos-Host", os.Getenv("HOST"))
		w.Header().Set("X-Mesos-Id", os.Getenv("MESOS_TASK_ID"))
		if len(sc.Body) > 0 {
			fmt.Fprint(w, sc.Body)
		}

		return sc.Code, nil
	}

	// Didn't qualify; pass-thru
	return s.Next.ServeHTTP(w, r)
}
