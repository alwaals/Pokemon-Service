package schema

type Pokemon struct {
	Id        string `json:"ID"`
	Name      string `json:"Name"`
	Type      string `json:"Type"`
	Height    string `json:"Height"`
	Weight    string `json:"Weight"`
	Abilities string `json:"Abilities"`
}
type PokemonRequest struct {
	Pokemon
	RequestId string `json:"RequestID,omitempty"`
	RequestTs string `json:"RequestTS,omitempty"`
}
type PokemonResponse struct {
	PokemonRequest
	RespMessage string `json:"RespMessage"`
	RespCode    int    `json:"RespCode"`
	Latency     string `json:"Latency"`
}
