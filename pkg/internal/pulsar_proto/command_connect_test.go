package pulsar_proto

import (
	"encoding/hex"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeCommandConnectSuccess(t *testing.T) {
	decodeString, _ := hex.DecodeString("00000019080212150a05322e352e321a00200f2a046e6f6e6552020801")
	cmd := &BaseCommand{}
	err := proto.Unmarshal(decodeString[4:], cmd)
	assert.Nil(t, err)
}
