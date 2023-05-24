package store

// Store

// PUT
type StorePutRequest struct {
	Key         string
	User        string
	Data        any
	RespChannel chan StorePutResponse
}

type StorePutResponse struct {
	Error error `json:"error"`
}

// GET
type StoreGetRequest struct {
	Key         string
	RespChannel chan StoreGetResponse
}

type StoreGetResponse struct {
	Data  any   `json:"data"`
	Error error `json:"error"`
}

// DELETE
type StoreDeleteRequest struct {
	Key         string
	User        string
	RespChannel chan StoreDeleteResponse
}

type StoreDeleteResponse struct {
	Error error `json:"error"`
}

// List

// GET key
type ListGetRequest struct {
	Key         string
	RespChannel chan ListGetResponse
}

type ListGetResponse struct {
	Data struct {
		Key   string `json:"key"`
		Owner string `json:"owner"`
	}
	Error error
}

// GET all
type ListGetAllRequest struct {
	RespChannel chan ListGetAllResponse
}

type ListGetAllResponse struct {
	Data []struct {
		Key   string `json:"key"`
		Owner string `json:"owner"`
	}
	Error error
}
