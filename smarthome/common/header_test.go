package common_test

import (
	"github.com/betom84/go-alexa/smarthome/common"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeader(t *testing.T) {
	h := common.NewHeader("Name", "Namespace")

	assert.NotEmpty(t, h.MessageID)
	assert.Equal(t, h.PayloadVersion, "3")
	assert.Equal(t, h.Name, "Name")
	assert.Equal(t, h.Namespace, "Namespace")
}
