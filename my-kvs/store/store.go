package store

import (
	"errors"
	"fmt"
	"time"
)

type Store struct {
	Data map[string]*Entry
}

type Entry struct {
	Owner        string    `json:"owner"`
	Value        any       `json:"value"`
	Writes       int       `json:"writes"`
	Reads        int       `json:"reads"`
	LastAccessed time.Time `json:"lastAccessed"`
}

var ErrNotFound = errors.New("not found")
var ErrNotOwner = errors.New("not owner")

var storage Store
var requestChannel chan any
var depth int

func InitStore(d int) {
	storage = Store{
		Data: make(map[string]*Entry),
	}

	requestChannel = make(chan any)

	if d > 0 {
		depth = d
	}

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
			owner, writes, reads, age, err := list(event.Key)
			listGetResponse := ListGetResponse{
				Data: struct {
					Key    string `json:"key"`
					Owner  string `json:"owner"`
					Writes int    `json:"writes"`
					Reads  int    `json:"reads"`
					Age    int64  `json:"age"`
				}{event.Key, owner, writes, reads, age},
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
	var entry *Entry

	element, ok := storage.Data[key]

	if ok {
		// check if user is same as owner and update value if so
		if !authorised(user, element.Owner) {
			return fmt.Errorf("put: %q %w of %q", user, ErrNotOwner, key)
		}
		element.Value = value
		element.Writes += 1
		element.LastAccessed = time.Now()
	} else {
		// create value anew
		entry = &Entry{Owner: user, Value: value, Writes: 1, Reads: 0, LastAccessed: time.Now()}

		// check if size is equal to depth
		storeSize := len(storage.Data)

		if depth > 0 && storeSize >= depth {
			deleteLeastRecent()
		}

		// insert new value into store
		storage.Data[key] = entry

	}

	return nil
}

func get(key string) (any, error) {
	entry, ok := storage.Data[key]

	if !ok {
		return "", fmt.Errorf("get: key: %q: %w", key, ErrNotFound)
	}
	entry.Reads += 1
	entry.LastAccessed = time.Now()
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

func list(key string) (string, int, int, int64, error) {

	entry, ok := storage.Data[key]

	if !ok {
		return "", 0, 0, 0, fmt.Errorf("list: key %q: %w", key, ErrNotFound)
	}

	age := int64(time.Since(entry.LastAccessed) / time.Millisecond)

	return entry.Owner, entry.Writes, entry.Reads, age, nil
}

func listAll() []struct {
	Key    string `json:"key"`
	Owner  string `json:"owner"`
	Writes int    `json:"writes"`
	Reads  int    `json:"reads"`
	Age    int64  `json:"age"`
} {

	data := storage.Data
	entries := make([]struct {
		Key    string `json:"key"`
		Owner  string `json:"owner"`
		Writes int    `json:"writes"`
		Reads  int    `json:"reads"`
		Age    int64  `json:"age"`
	}, 0, len(data))

	for key, entry := range data {
		age := int64(time.Since(entry.LastAccessed) / time.Millisecond)
		entries = append(entries, struct {
			Key    string `json:"key"`
			Owner  string `json:"owner"`
			Writes int    `json:"writes"`
			Reads  int    `json:"reads"`
			Age    int64  `json:"age"`
		}{key, entry.Owner, entry.Writes, entry.Reads, age})
	}

	return entries
}

func authorised(user, owner string) bool {
	if user == "admin" {
		return true
	}

	return user == owner
}

func deleteLeastRecent() {
	lruKey := ""
	lrDate := time.Now()

	for key := range storage.Data {
		lruKey = key
		break
	}

	for key, entry := range storage.Data {
		if entry.LastAccessed.Before(lrDate) {
			lrDate = entry.LastAccessed
			lruKey = key
		}
	}

	delete(storage.Data, lruKey)
}
