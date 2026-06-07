package dto

type HealthResponse struct {
	ServerActive   bool `json:"server_active"`
	RedisActive    bool `json:"redis_active"`
	DatabaseActive bool `json:"database_active"`
}

func NewHealthResponse(ServerActive bool, RedisActive bool, DatabaseActive bool) *HealthResponse {
	return &HealthResponse{
		ServerActive:   ServerActive,
		RedisActive:    RedisActive,
		DatabaseActive: DatabaseActive,
	}
}
