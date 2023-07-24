package types

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

// go test ./types -v to test the password creation and other data.
func TestNewAccount(t *testing.T){
	acc, err := NewAccount("test", "test", "abc")
	assert.Nil(t, err)

	log.Println(acc)
}
