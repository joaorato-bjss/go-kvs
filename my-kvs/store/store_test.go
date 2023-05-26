package store

import (
	"strconv"
	"testing"
	"time"
)

func TestDoStorePut(t *testing.T) {
	InitStore(0)
	words := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	go listen()

	t.Run("InitialPut", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			resp := DoStorePut(strconv.Itoa(i), "rato", words[i])
			if resp.Error != nil {
				t.Errorf("expected no error got '%s'", resp.Error.Error())
			}
		}
	})

	t.Run("NotOwnerPut", func(t *testing.T) {
		resp := DoStorePut("0", "joao", words[3])
		if resp.Error == nil {
			t.Errorf("expected not owner error got none")
		}
	})

	t.Run("OwnerPut", func(t *testing.T) {
		resp := DoStorePut("0", "rato", words[3])
		if resp.Error != nil {
			t.Errorf("expected no error got '%s'", resp.Error.Error())
		}
	})

	t.Run("AdminPut", func(t *testing.T) {
		resp := DoStorePut("0", "admin", words[5])
		if resp.Error != nil {
			t.Errorf("expected no error got '%s'", resp.Error.Error())
		}
	})

}

func TestDoStoreGet(t *testing.T) {
	InitStore(0)
	words := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	go listen()

	for i := 0; i < 10; i++ {
		resp := DoStorePut(strconv.Itoa(i), "rato", words[i])
		if resp.Error != nil {
			t.Errorf("expected no error got '%s'", resp.Error.Error())
		}
	}

	t.Run("NotFoundGet", func(t *testing.T) {
		resp := DoStoreGet("10")
		if resp.Error == nil {
			t.Errorf("expected not found error got none")
		}
	})

	t.Run("ValidGets", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			resp := DoStoreGet(strconv.Itoa(i))
			if resp.Error != nil {
				t.Errorf("expected no error got '%s'", resp.Error.Error())
			}
			if resp.Data != words[i] {
				t.Errorf("expected %s, got %s", words[i], resp.Data)
			}
		}
	})

}

func TestDoStoreDelete(t *testing.T) {
	InitStore(0)
	words := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	go listen()

	for i := 0; i < 10; i++ {
		resp := DoStorePut(strconv.Itoa(i), "rato", words[i])
		if resp.Error != nil {
			t.Errorf("expected no error got '%s'", resp.Error.Error())
		}
	}

	t.Run("NotFoundDelete", func(t *testing.T) {
		resp := DoStoreDelete("10", "rato")
		if resp.Error == nil {
			t.Errorf("expected not found error got none")
		}
	})

	t.Run("NotOwnerDelete", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			resp := DoStoreDelete(strconv.Itoa(i), "joao")
			if resp.Error == nil {
				t.Errorf("expected not owner error got none")
			}
		}
	})

	t.Run("OwnerDelete", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			resp := DoStoreDelete(strconv.Itoa(i), "rato")
			if resp.Error != nil {
				t.Errorf("expected no error got %s", resp.Error.Error())
			}

			respGet := DoStoreGet(strconv.Itoa(i))
			if respGet.Error == nil {
				t.Errorf("expected not found, got none")
			}
		}
	})

}

func TestDoListGet(t *testing.T) {
	InitStore(0)
	words := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	go listen()

	for i := 0; i < 10; i++ {
		resp := DoStorePut(strconv.Itoa(i), "rato", words[i])
		if resp.Error != nil {
			t.Errorf("expected no error got '%s'", resp.Error.Error())
		}
	}

	t.Run("NotFoundList", func(t *testing.T) {
		resp := DoListGet("10")
		if resp.Error == nil {
			t.Errorf("expected not found error got none")
		}
	})

	t.Run("ValidList", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			resp := DoListGet(strconv.Itoa(i))
			if resp.Error != nil {
				t.Errorf("expected no error got '%s'", resp.Error.Error())
			}

			if resp.Data.Owner != "rato" {
				t.Errorf("expected %s, got %s", "rato", resp.Data.Owner)
			}
		}
	})
}

func TestDoListGetAll(t *testing.T) {
	InitStore(0)
	words := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	go listen()

	for i := 0; i < 10; i++ {
		resp := DoStorePut(strconv.Itoa(i), "rato", words[i])
		if resp.Error != nil {
			t.Errorf("expected no error got '%s'", resp.Error.Error())
		}
	}

	t.Run("ValidListAll", func(t *testing.T) {
		resp := DoListGetAll()

		if len(resp.Data) != 10 {
			t.Errorf("expected 10 elements, got %d", len(resp.Data))
		}
	})
}

func TestDoStorePut2(t *testing.T) {
	InitStore(10)
	words := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight"}
	go listen()

	for i := 0; i < 9; i++ {
		resp := DoStorePut(strconv.Itoa(i), "rato", words[i])
		if resp.Error != nil {
			t.Errorf("expected no error got '%s'", resp.Error.Error())
		}
		time.Sleep(time.Millisecond)
	}

	t.Run("NewPutUnderDepth", func(t *testing.T) {
		resp := DoStorePut("9", "joao", "nine")
		if resp.Error != nil {
			t.Errorf("expected no error, got %s", resp.Error.Error())
		}
		if len(storage.Data) != 10 {
			t.Errorf("expected 10 elements, got %d", len(storage.Data))
		}
	})

	t.Run("NewPutOverDepth", func(t *testing.T) {
		resp := DoStorePut("10", "joao", "ten")
		if resp.Error != nil {
			t.Errorf("expected no error, got %s", resp.Error.Error())
		}
		resp = DoStorePut("11", "pedro", "eleven")
		if resp.Error != nil {
			t.Errorf("expected no error, got %s", resp.Error.Error())
		}
		if len(storage.Data) != 10 {
			t.Errorf("expected 10 elements, got %d", len(storage.Data))
		}
	})

}
