package message

type StoreQuery struct {
	Key string `json:"key"`
}

type StoreWrite struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
