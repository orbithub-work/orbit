package httpapi

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	port int
	srv  *http.Server
	ln   net.Listener
}

func Start(ctx context.Context, basePort int, tryCount int, deps Deps) (*Server, error) {
	if tryCount <= 0 {
		tryCount = 5
	}

	ln, port, err := listenInRange(basePort, tryCount)
	if err != nil {
		return nil, err
	}

	h := NewHandler(deps, port)
	h.StartBackgroundTasks(ctx)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// Wrap with global middleware (CORS, Metrics, TraceID)
	// corsMiddleware inside handler.go (moved to middleware.go but attached to Handler)
	// handles: Metrics begin/end, TraceID injection, CORS headers.
	var rootHandler http.Handler = mux
	rootHandler = h.corsMiddleware(rootHandler)

	s := &http.Server{
		Handler:      rootHandler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		_ = s.Serve(ln)
	}()

	return &Server{
		port: port,
		srv:  s,
		ln:   ln,
	}, nil
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) BaseURL() string {
	return "http://localhost:" + strconv.Itoa(s.port)
}

func (s *Server) Close(ctx context.Context) error {
	if s == nil {
		return nil
	}
	return s.srv.Shutdown(ctx)
}
