package store

import (
	"errors"
	"fmt"
)

type Store struct {
	Data map[string]Entry
}

type Entry struct {
	Owner string `json:"owner"`
	Value any    `json:"value"`
}

var ErrNotFound = errors.New("not found")
var ErrNotOwner = errors.New("not owner")

var storage Store
var requestChannel chan any

func InitStore() {
	storage = Store{
		Data: make(map[string]Entry),
	}

	requestChannel = make(chan any)

	go listen()
}

func listen() {
	for request := range requestChannel {
		switch event := request.(type) {
		case StorePutRequest:
			err := put(event.Key, event.User, event.Data)
			storePutResponse := StorePutResponse{
				Error: err,
			}

			event.RespChannel <- storePutResponse
			close(event.RespChannel)

		case StoreGetRequest:
			data, err := get(event.Key)
			storeGetResponse := StoreGetResponse{
				Data:  data,
				Error: err,
			}

			event.RespChannel <- storeGetResponse
			close(event.RespChannel)

		case StoreDeleteRequest:
			err := del(event.Key, event.User)
			storeDeleteResponse := StoreDeleteResponse{
				Error: err,
			}

			event.RespChannel <- storeDeleteResponse
			close(event.RespChannel)

		case ListGetRequest:
			owner, err := list(event.Key)
			listGetResponse := ListGetResponse{
				Data: struct {
					Key   string `json:"key"`
					Owner string `json:"owner"`
				}{event.Key, owner},
				Error: err,
			}

			event.RespChannel <- listGetResponse
			close(event.RespChannel)

		case ListGetAllRequest:
			data := listAll()
			listGetAllResponse := ListGetAllResponse{
				Data: data,
			}

			event.RespChannel <- listGetAllResponse
			close(event.RespChannel)
		}
	}
}

func put(key string, user string, value any) error {
	var entry Entry

	element, ok := storage.Data[key]

	if ok {
		// check if user is same as owner and update value if so
		if !authorised(user, element.Owner) {
			return fmt.Errorf("put: %q %w of %q", user, ErrNotOwner, key)
		}
		element.Value = value
	} else {
		// create value anew
		entry = Entry{Owner: user, Value: value}
		storage.Data[key] = entry
	}

	return nil
}

func get(key string) (any, error) {
	entry, ok := storage.Data[key]

	if !ok {
		return "", fmt.Errorf("get: key: %q: %w", key, ErrNotFound)
	}

	return entry.Value, nil
}

func del(key string, user string) error {
	entry, ok := storage.Data[key]

	if !ok {
		return fmt.Errorf("delete: key %q: %w", key, ErrNotFound)
	}

	if !authorised(user, entry.Owner) {
		return fmt.Errorf("delete: %q %w of %q", user, ErrNotOwner, key)
	}

	delete(storage.Data, key)

	return nil
}

func list(key string) (string, error) {

	entry, ok := storage.Data[key]

	if !ok {
		return "", fmt.Errorf("list: key %q: %w", key, ErrNotFound)
	}

	return entry.Owner, nil
}

func listAll() []struct {
	Key   string `json:"key"`
	Owner string `json:"owner"`
} {

	data := storage.Data
	entries := make([]struct {
		Key   string `json:"key"`
		Owner string `json:"owner"`
	}, 0, len(data))

	for key, entry := range data {
		entries = append(entries, struct {
			Key   string `json:"key"`
			Owner string `json:"owner"`
		}{key, entry.Owner})
	}

	return entries
}

func authorised(user, owner string) bool {
	if user == "admin" {
		return true
	}

	return user == owner
}
