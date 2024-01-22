package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	_ "log"
	"net/http"
	schema "pokemon-service/schema"
	utility "pokemon-service/utility"
	"runtime/debug"
	"time"

	"bytes"
	"github.com/allegro/bigcache"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const (
	contentType = "Content-Type"
	application = "Application/json"
)

// Received cache and logger from main file
type Service struct {
	Cache  *bigcache.BigCache
	Logger *schema.Logger
}

// Retrieves existing pokemon record from cache
func (service *Service) GetByID(w http.ResponseWriter, req *http.Request) {
	_, cancelFunc := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancelFunc()

	w.Header().Set(contentType, application)
	var pokemonResp schema.PokemonResponse
	start := time.Now()

	//In case of any panic errors, gracefully recovers and prints stack to response
	defer func() {
		if err := recover(); err != nil {
			utility.FrameHttpResponse(500, string(debug.Stack()), &pokemonResp, start, w)
			return
		}
	}()

	//Retrieving params data from URL
	id := mux.Vars(req)["Id"]

	if len(id) <= 0 {
		utility.FrameHttpResponse(422, "Id is expected in endpoint", &pokemonResp, start, w)
		return
	}
	//Setting new Request ID for every request using uuid library when reqId is not sent by user
	xRequestID := uuid.New().String()
	pokemonResp.RequestId = xRequestID

	//Getting data from cache
	pokemondById, _ := service.Cache.Get(id)

	//Unmarshalling cache data to display in response
	buf := bytes.NewBuffer(pokemondById)
	err := json.NewDecoder(buf).Decode(&pokemonResp)
	if err != nil {
		utility.FrameHttpResponse(404, fmt.Sprintf("Unable to get data from cache for Id:%v", id), &pokemonResp, start, w)
		return
	}

	utility.FrameHttpResponse(200, "Success", &pokemonResp, start, w)
}

// Retrieves existing pokemon record from cache
func (service *Service) GetByName(w http.ResponseWriter, req *http.Request) {
	_, cancelFunc := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancelFunc()
	w.Header().Set(contentType, application)
	var pokemonResp schema.PokemonResponse
	start := time.Now()

	//In case of any panic errors, gracefully recovers and prints stack to response
	defer func() {
		if err := recover(); err != nil {
			utility.FrameHttpResponse(500, string(debug.Stack()), &pokemonResp, start, w)
			return
		}
	}()

	//Retrieving params data from URL
	name := mux.Vars(req)["Name"]
	if len(name) <= 0 {
		utility.FrameHttpResponse(422, "Name is expected in request param", &pokemonResp, start, w)
		return
	}

	//Setting new Request ID for every request using uuid library when reqId is not sent by user
	xRequestID := uuid.New().String()
	pokemonResp.RequestId = xRequestID

	//Getting data from cache
	pokemonByName, _ := service.Cache.Get(name)

	//Unmarshalling cache data to display in response
	buf := bytes.NewBuffer(pokemonByName)
	err := json.NewDecoder(buf).Decode(&pokemonResp)
	if err != nil {
		utility.FrameHttpResponse(400, fmt.Sprintf("Unable to get data from cache for Name:%v", name), &pokemonResp, start, w)
		return
	}

	utility.FrameHttpResponse(200, "Success", &pokemonResp, start, w)
}

// Deletes existing pokemon record from cache
func (service *Service) DeleteByID(w http.ResponseWriter, req *http.Request) {
	_, cancelFunc := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancelFunc()
	w.Header().Set(contentType, application)
	var pokemonResp schema.PokemonResponse
	start := time.Now()

	//In case of any panic errors, gracefully recovers and prints stack to response
	defer func() {
		if err := recover(); err != nil {
			utility.FrameHttpResponse(500, string(debug.Stack()), &pokemonResp, start, w)
			return
		}
	}()

	//Retrieving params data from URL
	id := mux.Vars(req)["Id"]
	if len(id) <= 0 {
		utility.FrameHttpResponse(422, "Id is expected in endpoint", &pokemonResp, start, w)
		return
	}

	//Setting new Request ID for every request using uuid library when reqId is not sent by user
	xRequestID := uuid.New().String()
	pokemonResp.RequestId = xRequestID

	//Deletes record only when its present, else not found error
	pokemondById, _ := service.Cache.Get(id)
	buf := bytes.NewBuffer(pokemondById)
	err := json.NewDecoder(buf).Decode(&pokemonResp)
	if err != nil {
		utility.FrameHttpResponse(400, fmt.Sprintf("Unable to get data from cache for Id to delete:%v", id), &pokemonResp, start, w)
		return
	}

	//If record found, deletes it from cache
	err = service.Cache.Delete(id)
	if err != nil {
		utility.FrameHttpResponse(400, fmt.Sprintf("Unable to get data from cache for Id:%v", id), &pokemonResp, start, w)
		return
	}

	utility.FrameHttpResponse(200, "Success", &pokemonResp, start, w)
}

// Adding new pokemon data into cache
func (service *Service) AddPokemon(w http.ResponseWriter, req *http.Request) {
	_, cancelFunc := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancelFunc()

	w.Header().Set(contentType, application)
	var pokemonResp schema.PokemonResponse
	start := time.Now()

	//In case of any panic errors, gracefully recovers and prints stack to response
	defer func() {
		if err := recover(); err != nil {
			utility.FrameHttpResponse(500, string(debug.Stack()), &pokemonResp, start, w)
			return
		}
	}()

	var pokemonReq schema.PokemonRequest
	err := json.NewDecoder(req.Body).Decode(&pokemonReq)
	if err != nil {
		utility.FrameHttpResponse(500, "Invalid Json request", &pokemonResp, start, w)
		return
	}

	if len(pokemonReq.RequestId) <= 0 {
		// Setting new Request ID for every request using uuid library when reqId is not sent by user
		xRequestID := uuid.New().String()
		pokemonResp.RequestId = xRequestID
	}

	resp, _ := json.Marshal(pokemonReq)

	//Adds this new pokemon record into cache
	service.Cache.Set(pokemonReq.Name, resp)
	service.Cache.Set(pokemonReq.Id, resp)

	pokemonResp.Id = pokemonReq.Id
	pokemonResp.Type = pokemonReq.Type
	pokemonResp.Name = pokemonReq.Name
	pokemonResp.Height = pokemonReq.Height
	pokemonResp.Weight = pokemonReq.Weight
	pokemonResp.Abilities = pokemonReq.Abilities

	utility.FrameHttpResponse(200, "Success", &pokemonResp, start, w)
}

// Health check function
func (service *Service) HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Health Check success"))
}
