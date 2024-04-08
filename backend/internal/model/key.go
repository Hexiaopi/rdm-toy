package model

type Key struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
	TTL   int64       `json:"ttl"`
}

type ZSetValue struct {
	Score  float64 `json:"score"`
	Member string  `json:"member"`
}
