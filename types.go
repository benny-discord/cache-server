package main

type cacheValue struct {
	data    string
	expires int64
}

type wsRequestPayload struct {
	Key     string `json:"key,omitempty"`
	Value   string `json:"value,omitempty"`
	Expires int64  `json:"expires,omitempty"`
	Op      string `json:"op"`
}
