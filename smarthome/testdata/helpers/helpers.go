package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/betom84/go-alexa/smarthome/common"
)

// CreateDirective initializes an directive with the given payload
func CreateDirective(t *testing.T, payload string) *common.Directive {
	t.Helper()

	dir, err := common.NewDirective([]byte(payload))
	if err != nil {
		t.Fatalf("could not create directive; %v", err)
	}

	if dir.Header == nil {
		t.Fatal("could not parse payload; directive doesn't contain a header")
	}

	return dir
}

// LoadRequest unmarshals an directive from the given file
func LoadRequest(t *testing.T, file string) *common.Directive {
	t.Helper()

	payload, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("could read response file; %v", err)
	}

	var raw map[string]*json.RawMessage
	err = json.Unmarshal(payload, &raw)
	if err != nil {
		t.Fatalf("could not unmarshal response file; %v", err)
	}

	var dir *common.Directive
	if rawDir, ok := raw["directive"]; ok {
		dir, err = common.NewDirective(*rawDir)
	} else {
		err = fmt.Errorf("missing directive in request file")
	}

	if err != nil {
		t.Fatalf("could not create directive; %v", err)
	}

	if dir.Header == nil {
		t.Fatal("failed to parse request file; directive doesn't contain a header")
	}

	return dir
}

// FailOnPanic recovers a panic and let the test fail with the recovered error
func FailOnPanic(t *testing.T) {
	if r := recover(); r != nil {
		t.Fatal(r)
	}
}
