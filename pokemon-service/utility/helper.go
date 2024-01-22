package utility

import (
	schema "pokemon-service/schema"
	"encoding/json"
	"net/http"
	"time"
)

func FrameHttpResponse(status int, errMsg string, userResp *schema.PokemonResponse, start time.Time, w http.ResponseWriter) {
	w.WriteHeader(status)
	userResp.RequestTs = start.Format(time.RFC3339)
	userResp.RespMessage = errMsg
	userResp.RespCode = status
	userResp.Latency = time.Since(start).String()
	json.NewEncoder(w).Encode(userResp)
}