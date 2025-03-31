package model

type Client struct {
	Id    int    `json:"id"`
	Addr  string `json:"addr"`
	DB    int    `json:"db"`
	Age   int64  `json:"age"`
	Idle  int64  `json:"idle"`
	Flags string `json:"flags"`
	Cmd   string `json:"cmd"`
	Sub   int    `json:"sub"`
	Psub  int    `json:"psub"`
}
