package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/javierdelapuente/test-go-12-factor/config"
	"github.com/javierdelapuente/test-go-12-factor/internal/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type mainHandler struct {
	counter     prometheus.Counter
	charmConfig config.CharmConfig
	service     service.Service
}

func (h mainHandler) serveHelloWorld(w http.ResponseWriter, r *http.Request) {
	h.counter.Inc()
	log.Println("root handler")
	log.Printf("Service %#v\n", h.service)
	log.Printf("Counter %#v\n", h.counter)

	fmt.Fprintf(w, "Hello world! Path: %s\n", r.URL.Path)
}

func (h mainHandler) serveEnvs(w http.ResponseWriter, r *http.Request) {
	h.counter.Inc()

	type EnvVar struct {
		Name, Value string
	}

	envVars := []EnvVar{}

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		log.Printf("pair 1: %s pair2 %s\n", pair[0], pair[1])
		envVars = append(envVars, EnvVar{
			Name:  pair[0],
			Value: pair[1]})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(envVars)
}

func (h mainHandler) serveSleep(w http.ResponseWriter, r *http.Request) {
	val := r.URL.Query().Get("duration")
	duration, err := strconv.Atoi(val)
	if err != nil {
		http.Error(w, "Bad Request. Wrong duration", http.StatusBadRequest)
		return
	}
	time.Sleep(time.Duration(duration) * time.Second)
}

func (h mainHandler) serveConfig(w http.ResponseWriter, r *http.Request) {
	configToSearch := r.PathValue("config")
	for _, config := range h.charmConfig.Configs {
		if config.Name == configToSearch {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(config.Value)
			return
		}
	}
	http.Error(w, fmt.Sprintf("Bad Request. Config %s not found.", configToSearch), http.StatusBadRequest)
	return
}

func (h mainHandler) serveMysql(w http.ResponseWriter, r *http.Request) {
	err := h.service.CheckMysqlStatus()
	if err != nil {
		log.Printf(err.Error())
		io.WriteString(w, "FAILURE")
		return
	} else {
		io.WriteString(w, "SUCCESS")
	}
}

func (h mainHandler) servePostgresql(w http.ResponseWriter, r *http.Request) {
	err := h.service.CheckPostgresqlStatus()
	if err != nil {
		log.Printf(err.Error())
		io.WriteString(w, "FAILURE")
		return
	} else {
		io.WriteString(w, "SUCCESS")
	}
}

func (h mainHandler) notImplemented(w http.ResponseWriter, r *http.Request) {
	h.counter.Inc()
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Not implemented yet!"))
}

func main() {
	log.Println("standard logger")

	requestCounter := prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "request_count",
			Help: "No of request handled",
		})
	prometheus.MustRegister(requestCounter)
	prometheusMux := http.NewServeMux()
	prometheusMux.Handle("/metrics", promhttp.Handler())
	prometheusServer := &http.Server{
		Addr:    ":8081",
		Handler: prometheusMux,
	}
	go func() {
		if err := prometheusServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Prometheus HTTP server error: %v", err)
		}
		log.Println("Prometheus HTTP Stopped serving new connections.")
	}()

	mux := http.NewServeMux()

	// build manually for now.
	// val1 := "val1"
	// charmConfig := config.CharmConfig{
	// 	// Configs: []{config.Config{
	// 	// 	Name: "a",
	// 	// 	Value: nil,}},
	// 	Configs: []config.Config{
	// 		{Name: "name1", Value: &val1},
	// 		{Name: "name2"},
	// 	},
	// 	Integrations: config.Integrations{},
	// }
	charmConfig := config.BuildCharmConfig(os.Environ())
	mainHandler := mainHandler{
		counter:     requestCounter,
		charmConfig: charmConfig,
		service:     service.Service{CharmConfig: charmConfig},
	}

	mux.HandleFunc("/", mainHandler.serveHelloWorld)
	mux.HandleFunc("/env", mainHandler.serveEnvs)
	mux.HandleFunc("/sleep", mainHandler.serveSleep)
	// the config variable is in the path.
	mux.HandleFunc("/config/{config}", mainHandler.serveConfig)
	mux.HandleFunc("/mysql/status", mainHandler.serveMysql)
	mux.HandleFunc("/postgresql/status", mainHandler.servePostgresql)
	mux.HandleFunc("/s3/status", mainHandler.notImplemented)
	mux.HandleFunc("/mongodb/status", mainHandler.notImplemented)
	mux.HandleFunc("/redis/status", mainHandler.notImplemented)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
		log.Println("Stopped serving new connections.")
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")

}
