package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidJWT(t *testing.T) {
	id := uuid.New()
	duration, err := time.ParseDuration("5s")
	if err != nil {
		t.Errorf("%v", err)
	}

	token, err := MakeJWT(id, "secret", duration)
	if err != nil {
		t.Errorf("%v", err)
	}

	resultId, err := ValidateJWT(token, "secret")
	if err != nil {
		t.Errorf("%v", err)
	}

	if resultId != id {
		t.Errorf("Got: %v Want: %v\n", resultId, id)
	}
}

func TestExpiredJWT(t *testing.T) {
	id := uuid.New()
	duration, err := time.ParseDuration("2ms")
	if err != nil {
		t.Errorf("%v", err)
	}

	token, err := MakeJWT(id, "secret", duration)
	if err != nil {
		t.Errorf("%v", err)
	}

	wait, err := time.ParseDuration("10ms")
	if err != nil {
		t.Errorf("%v", err)
	}

	time.Sleep(wait)

	_, err = ValidateJWT(token, "secret")
	if err == nil {
		t.Errorf("did not return error for expired token")
	}

}

func TestWrongKeyJWT(t *testing.T) {
	id := uuid.New()
	duration, err := time.ParseDuration("5s")
	if err != nil {
		t.Errorf("%v", err)
	}

	token, err := MakeJWT(id, "secret", duration)
	if err != nil {
		t.Errorf("%v", err)
	}

	_, err = ValidateJWT(token, "diffrentSecret")
	if err == nil {
		t.Errorf("did not return error for wrong secret")
	}
}

func TestGetBearer(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "Bearer Token")

	result, err := GetBearerToken(header)
	if err != nil {
		t.Errorf("%v", err)
	}

	if result != "Token" {
		t.Errorf("want: Token got: %v", result)
	}
}
