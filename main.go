package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/grpcwebexample/internal/hellopb"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer l.Close()

	mux := cmux.New(l)

	httpL := mux.Match(cmux.HTTP1Fast())
	grpcL := mux.MatchWithWriters(
		cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"),
	)

	grpcSrv := grpc.NewServer()
	healthSrv := health.NewServer()
	healthSrv.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(grpcSrv, healthSrv)
	reflection.Register(grpcSrv)
	greeter := newGreeterService()
	hellopb.RegisterGreeterServer(grpcSrv, greeter)
	restMux := restHandler(greeter)
	grpcWebSrv := grpcweb.WrapServer(
		grpcSrv,
		grpcweb.WithOriginFunc(func(string) bool { return true }),
	)
	restSrv := &http.Server{
		Handler: grpcWebAwareHandler(grpcWebSrv, restMux),
	}

	errCh := make(chan error, 3)

	go func() {
		if serveErr := restSrv.Serve(httpL); !errors.Is(serveErr, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http serve: %w", serveErr)
		}
	}()

	go func() {
		if serveErr := grpcSrv.Serve(grpcL); serveErr != nil {
			errCh <- fmt.Errorf("grpc serve: %w", serveErr)
		}
	}()

	go func() {
		log.Print("server listening on :8080")
		if serveErr := mux.Serve(); serveErr != nil {
			errCh <- fmt.Errorf("cmux serve: %w", serveErr)
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("server shutting down...")
	case serveErr := <-errCh:
		log.Printf("server error: %v", serveErr)
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := restSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("http shutdown error: %v", err)
	}
	grpcSrv.GracefulStop()
}

func corsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowOrigin := r.Header.Get("Origin")
		if allowOrigin == "" {
			allowOrigin = "*"
		}
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func restHandler(greeter *greeterService) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		type response struct {
			Message string `json:"message"`
		}
		writeJSON(w, http.StatusOK, response{
			Message: greeter.greeting(name),
		})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	return corsHandler(mux)
}

func grpcWebAwareHandler(grpcWebSrv *grpcweb.WrappedGrpcServer, rest http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if grpcWebSrv.IsGrpcWebRequest(r) || grpcWebSrv.IsGrpcWebSocketRequest(r) {
			grpcWebSrv.ServeHTTP(w, r)
			return
		}
		if grpcWebSrv.IsAcceptableGrpcCorsRequest(r) {
			allowOrigin := r.Header.Get("Origin")
			if allowOrigin == "" {
				allowOrigin = "*"
			}
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "content-type,x-grpc-web,x-user-agent")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		rest.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("json encode error: %v", err)
	}
}
