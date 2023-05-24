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
var requestChannel chan interface{}

func InitStore() {
	storage = Store{
		Data: make(map[string]Entry),
	}

	requestChannel = make(chan interface{})

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
		case ListGetRequest:
		case ListGetAllRequest:
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

func authorised(user, owner string) bool {
	if user == "admin" {
		return true
	}

	return user == owner
}
