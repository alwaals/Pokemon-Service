package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"pokemon-service/schema"
	"testing"
	"time"
)

func TestHealthCheckHandler(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	service := &Service{}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(service.HealthCheckHandler)
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Check the response body is what we expect.
	expected := `Health Check success`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAddPokemon(t *testing.T) {

	service := loadBigCache()
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	inputs := []struct {
		testName string
		status   int
		respMesg string
		req      schema.PokemonRequest
	}{
		{testName: "TestAddPokemonSuccess1", status: 200, respMesg: "Success", req: schema.PokemonRequest{Pokemon: schema.Pokemon{Id: "111", Name: "100111"}}},
		{testName: "TestAddPokemonSuccess2", status: 200, respMesg: "Success", req: schema.PokemonRequest{Pokemon: schema.Pokemon{Id: "222", Name: "2222"}}},
		{testName: "TestAddPokemonFailure", status: 200, respMesg: "Success", req: schema.PokemonRequest{}},
	}

	for _, item := range inputs {
		body, _ := json.Marshal(item.req)
		req, err := http.NewRequest("POST", "/pokemon-service/add", bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		handler := http.HandlerFunc(service.AddPokemon)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the response body is what we expect.
		if rr.Code != item.status {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), item.respMesg)
		}
	}

}
func TestGetByID(t *testing.T) {

	service := loadBigCache()
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	inputs := []struct {
		status   int
		respMesg string
		id       string
	}{
		{status: 200, respMesg: "Success", id: "PK10002"},
		{status: 200, respMesg: "Success", id: "PK10001"},
		{status: 200, respMesg: "Unable to get data from cache for Id:PK1000908", id: "PK1000908"},
	}

	for _, item := range inputs {
		vars := map[string]string{
			"Id": item.id,
		}

		req, err := http.NewRequest("GET", "/pokemon-service/{Id}", nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(service.GetByID)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if rr.Code != item.status {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, item.status)
		}

	}

}
func TestGetByName(t *testing.T) {

	service := loadBigCache()
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	inputs := []struct {
		status   int
		respMesg string
		name     string
	}{
		{status: 200, respMesg: "Success", name: "Picachoo1"},
		{status: 200, respMesg: "Success", name: "Picachoo1"},
		{status: 200, respMesg: "Unable to get data from cache for Name:PK1000908", name: "PK1000908"},
	}

	for _, item := range inputs {
		vars := map[string]string{
			"Name": item.name,
		}

		req, err := http.NewRequest("GET", "/pokemon-service/{Name}", nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(service.GetByName)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if rr.Code != item.status {
			t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, item.status)
		}
	}

}
func TestDeleteByID(t *testing.T) {

	service := loadBigCache()
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	inputs := []struct {
		status   int
		respMesg string
		id       string
	}{
		{status: 200, respMesg: "Success", id: "PK10002"},
		{status: 200, respMesg: "Success", id: "PK10001"},
		{status: 200, respMesg: "Unable to get data from cache for Name:PK1000908", id: "PK1000908"},
	}

	for _, item := range inputs {
		vars := map[string]string{
			"Id": item.id,
		}

		req, err := http.NewRequest("DELETE", "/pokemon-service/delete/", nil)
		if err != nil {
			t.Fatal(err)
		}

		req = mux.SetURLVars(req, vars)

		handler := http.HandlerFunc(service.DeleteByID)
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		handler.ServeHTTP(rr, req)

		// Check the status code is what we expect.
		if status := rr.Code; status != item.status {
			t.Errorf("handler returned wrong status code: got %v want %v", status, item.status)
		}
	}

}
func loadBigCache() *Service {
	cache, _ := bigcache.NewBigCache(bigcache.DefaultConfig(24 * time.Hour))
	ps := []schema.Pokemon{
		{Id: fmt.Sprintf("PK%v", 10001), Name: "Picachoo1", Type: "TT", Height: "20.9", Weight: "30.9", Abilities: "Eat&Sleep"},
		{Id: fmt.Sprintf("PK%v", 10002), Name: "Picachoo2", Type: "PP", Height: "10.9", Weight: "31.1", Abilities: "Eat&Sleep"},
	}
	for _, val := range ps {
		resp, _ := json.Marshal(val)
		cache.Set(val.Name, resp)
		cache.Set(val.Id, resp)
	}
	return &Service{Cache: cache}
}
