package model

type PushMessage struct {
	Op       int32
	Servrice string
	Keys     []string
	Message  []byte
}

type BraodcastMsg struct {
	Op      int32
	Speed   int32
	Room    string
	Message []byte
}

type Mapping struct {
	Mid    int64
	Key    string
	Server string
	Online *Online
}

type Online struct {
	Server    string           `json:"server"`
	RoomCount map[string]int32 `json:"room_count"`
	Updated   int64            `json:"updated"`
}

type SendPluginOption struct {
	MaxTry, Mid                      int
	RedisKey1, RedisKey2, Key, Field string
	Expire                           int32
}
