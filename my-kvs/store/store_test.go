package store

import (
	"strconv"
	"testing"
)

func TestDoStorePut(t *testing.T) {
	InitStore()
	words := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	go listen()

	for i := 0; i < 10; i++ {
		resp := DoStorePut(strconv.Itoa(i), "rato", words[i])
		if resp.Error != nil {
			t.Errorf("expected no error got '%s'", resp.Error.Error())
		}
	}
}
