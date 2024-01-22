package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	handlers "pokemon-service/handlers"
	middlewares "pokemon-service/middlewares"
	s "pokemon-service/schema"
	"syscall"
	"time"

	"github.com/allegro/bigcache"
	"github.com/gorilla/mux"
)

var (
	loggerFileName = "logger.text"
	logger         = s.Logger{}
)

// Logging every transaction details in logger file for observing ongoing traffic
func init() {
	logFile, err := os.OpenFile(loggerFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Unable to create Logger file:", err.Error())
		return
	}
	log.SetOutput(logFile)

	logger.InfoLogger = log.New(logFile, "Info:", log.Ldate|log.Ltime|log.Lshortfile)
	logger.WarnLogger = log.New(logFile, "Warn:", log.Ldate|log.Ltime|log.Lshortfile)
	logger.DebugLogger = log.New(logFile, "Debug:", log.Ldate|log.Ltime|log.Lshortfile)
	logger.ErrorLogger = log.New(logFile, "Error:", log.Ldate|log.Ltime|log.Lshortfile)
	logger.FatalLogger = log.New(logFile, "Fatal:", log.Ldate|log.Ltime|log.Lshortfile)
}
func main() {
	r := mux.NewRouter()
	//cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(24 * time.Hour))
	cache, err := customerConfigBigCache()
	if err != nil {
		log.Fatal("Unable to load cache data:", err.Error())
	}
	loadingInMemCache(cache)
	service := &handlers.Service{Cache: cache, Logger: &logger}

	commonMiddleware := []middlewares.Middleware{
		middlewares.LoggingRequest,
		middlewares.LoggingResponse,
	}
	r.HandleFunc("/health-check", middlewares.Chain(service.HealthCheckHandler, logger, commonMiddleware...)).Methods("GET")
	r.HandleFunc("/pokemon-service/getByID/{Id}", middlewares.Chain(service.GetByID, logger, commonMiddleware...)).Methods("GET")
	r.HandleFunc("/pokemon-service/getByName/{Name}", middlewares.Chain(service.GetByName, logger, commonMiddleware...)).Methods("GET")
	r.HandleFunc("/pokemon-service/{Id}", middlewares.Chain(service.DeleteByID, logger, commonMiddleware...)).Methods("DELETE")
	r.HandleFunc("/pokemon-service/Add", middlewares.Chain(service.AddPokemon, logger, commonMiddleware...)).Methods("POST")

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Start the server in a separate Goroutine.
	go func() {
		service.Logger.InfoLogger.Println("Starting the server on :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	// Implement graceful shutdown for the server to handle fault tolerance
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	service.Logger.InfoLogger.Println("Shutting down the server...")

	// Set a timeout for shutdown (for example, 5 seconds).
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		service.Logger.InfoLogger.Fatalf("Server shutdown error: %v", err)
	}
	service.Logger.InfoLogger.Println("Server gracefully stopped")
}
func loadingInMemCache(cache *bigcache.BigCache) {
	pokemons := loadSamplePokemonData(cache)
	for _, val := range pokemons {
		resp, _ := json.Marshal(val)
		cache.Set(val.Name, resp)
		cache.Set(val.Id, resp)
	}
}
func loadSamplePokemonData(cache *bigcache.BigCache) []s.Pokemon {
	//generateUniqueIds := fmt.Sprintf("PK%v", rand.Intn(100000))
	return []s.Pokemon{
		{Id: fmt.Sprintf("PK%v", 10001), Name: "Chespin", Type: "TT", Height: "20.9", Weight: "30.9", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10002), Name: "Fennekin", Type: "PP", Height: "10.9", Weight: "31.1", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10003), Name: "Froakie", Type: "JJ", Height: "30.9", Weight: "32.0", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10004), Name: "Sylveon", Type: "KK", Height: "60.9", Weight: "34.8", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10005), Name: "Xerneas", Type: "UU", Height: "80.9", Weight: "31.5", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10006), Name: "Yveltal", Type: "LL", Height: "50.9", Weight: "37.3", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10007), Name: "Zygarde", Type: "WW", Height: "10.9", Weight: "33.9", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10008), Name: "PokemonX", Type: "QQ", Height: "90.9", Weight: "31.1", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10009), Name: "PokemonY", Type: "WY", Height: "20.9", Weight: "33.2", Abilities: "Eat&Sleep"},
	}
}
func customerConfigBigCache() (*bigcache.BigCache, error) {
	config := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: 1024,

		// time after which entry can be evicted
		LifeWindow: 3 * time.Minute,

		// Interval between removing expired entries (clean up).
		// If set to <= 0 then no action is performed.
		// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
		CleanWindow: 5 * time.Second,

		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: 12,

		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: 10,

		// prints information about additional memory allocation
		Verbose: true,

		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: 10,

		// callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,

		// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
		// for the new entry, or because delete was called. A constant representing the reason will be passed through.
		// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
		// Ignored if OnRemove is specified.
		OnRemoveWithReason: nil,
	}
	return bigcache.NewBigCache(config)
}
