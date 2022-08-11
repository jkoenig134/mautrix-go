package mautrix

// RespPublicRooms is the JSON response for http://matrix.org/speculator/spec/HEAD/client_server/unstable.html#get-matrix-client-unstable-publicrooms
type RespPublicRooms struct {
	TotalRoomCountEstimate int          `json:"total_room_count_estimate"`
	PrevBatch              string       `json:"prev_batch"`
	NextBatch              string       `json:"next_batch"`
	Chunk                  []PublicRoom `json:"chunk"`
}
