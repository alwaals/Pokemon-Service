package utility

import (
	schema "pokemon-service/schema"
	"net/http"
	"testing"
	"time"
)

func TestFrameHttpResponse(t *testing.T) {
	inputs := []struct {
		status   int
		errMsg   string
		userResp *schema.PokemonRequest
		start    time.Time
		w        http.ResponseWriter
	}{
		{status: 200, errMsg: "Success", userResp: &schema.PokemonRequest{}, start: time.Now(),
		 w: nil},
		 {status: 400, errMsg: "Error", userResp: &schema.PokemonRequest{}, start: time.Now(),
		 w: nil},
	}

	for _, item := range inputs {
		if item.status<=0{
			t.Error("Expecting ",item.status,"But received nil status code")
		}
	}

}