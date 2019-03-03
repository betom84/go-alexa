package helpers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/betom84/go-alexa/smarthome/common"
)

// AssertEqualsGolden compares content of golden file with marshaled response
func AssertEqualsGolden(t *testing.T, goldenFile string, response *common.Response) {
	t.Helper()

	resp := normalizeResponse(response)

	current, err := json.MarshalIndent(resp, " ", "  ")
	if err != nil {
		t.Fatalf("failed to assert equality, unable to marshal response; %v", err)
	}

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatalf("failed to assert equality, unable to read golden; %v", err)
	}

	current = removePrettyPrint(current)
	expected = removePrettyPrint(expected)

	if c := bytes.Compare(current, expected); c != 0 {
		t.Errorf("response doesnt match golden; %v", c)
		t.Errorf("response:\n%s\ngolden: \n%s\n doesnt match golden", string(current), string(expected))
	}
}

func removePrettyPrint(target []byte) []byte {
	result := bytes.Replace(target, []byte("\x0a"), []byte(""), -1)
	result = bytes.Replace(result, []byte("\x0d"), []byte(""), -1)
	result = bytes.Replace(result, []byte("\x20"), []byte(""), -1)
	//result = bytes.Replace(result, []byte("\x09"), []byte(""), -1)

	return result
}

// UpdateGolden saves marshaled response to golden file
func UpdateGolden(t *testing.T, goldenFile string, response *common.Response) {
	t.Helper()

	resp := normalizeResponse(response)

	bytes, err := json.MarshalIndent(resp, " ", "  ")
	if err != nil {
		t.Fatalf("failed to update golden file, unable to marshal response; %v", err)
	}

	err = ioutil.WriteFile(goldenFile, bytes, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to update golden file; %v", err)
	}
}

func normalizeResponse(response *common.Response) common.Response {
	resp := *response

	// need to unset the message id to make response comparable
	resp.Event.Header.MessageID = ""

	return resp
}
