package selectors_test

import (
	"testing"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/constants"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/selectors"
	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/protocol/tlv"
	"github.com/stretchr/testify/assert"
)

// -------------------------------------- Positive tests --------------------------------------

func TestNewSingleNodeSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewSingleNodeSelector(0x42)
	assert.NotNil(selector)
	assert.Equal(constants.SINGLE_NODE, selector.Type)
	assert.Equal([]uint8{0x42}, selector.NodeIDs)

	// test encoding
	encoded, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)
	assert.Equal(uint8(constants.SINGLE_NODE), encoded.Type())
	assert.Equal(uint16(1), encoded.Length())
	assert.Equal([]byte{0x42}, encoded.Value())

	// test decoding
	decoded, err := selectors.DecodeNodeSelector(encoded)
	assert.NoError(err)
	assert.Equal(constants.SINGLE_NODE, decoded.Type)
	assert.Equal([]uint8{0x42}, decoded.NodeIDs)
}

func TestNewNodeListSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewNodeListSelector([]uint8{0x01, 0x02, 0x03})
	assert.NotNil(selector)
	assert.Equal(constants.NODE_LIST, selector.Type)
	assert.Equal([]uint8{0x01, 0x02, 0x03}, selector.NodeIDs)

	// test encoding
	encoded, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)
	assert.Equal(uint8(constants.NODE_LIST), encoded.Type())
	assert.Equal(uint16(3), encoded.Length())
	assert.Equal([]uint8{0x01, 0x02, 0x03}, encoded.Value())

	// test decoding
	decoded, err := selectors.DecodeNodeSelector(encoded)
	assert.NoError(err)
	assert.Equal(constants.NODE_LIST, decoded.Type)
	assert.Equal([]uint8{0x01, 0x02, 0x03}, decoded.NodeIDs)
}

func TestNewAllNodesSelector(t *testing.T) {
	assert := assert.New(t)

	selector := selectors.NewAllNodesSelector()
	assert.NotNil(selector)
	assert.Equal(constants.ALL_NODES, selector.Type)
	assert.Nil(selector.NodeIDs)

	// test encoding
	encoded, err := selector.Encode()
	assert.NoError(err)
	assert.NotNil(encoded)
	assert.Equal(uint8(constants.ALL_NODES), encoded.Type())
	assert.Equal(uint16(0), encoded.Length())
	assert.Equal([]byte{}, encoded.Value())

	// test decoding
	decoded, err := selectors.DecodeNodeSelector(encoded)
	assert.NoError(err)
	assert.Equal(constants.ALL_NODES, decoded.Type)
	assert.Nil(decoded.NodeIDs)
}

// -------------------------------------- Negative tests --------------------------------------

func TestSingleNodeSelector_InvalidCases(t *testing.T) {
	assert := assert.New(t)

	// Empty NodeIDs
	selector := &selectors.NodeSelector{Type: constants.SINGLE_NODE, NodeIDs: []uint8{}}
	_, err := selector.Encode()
	assert.Error(err)

	// Multiple NodeIDs
	selector = &selectors.NodeSelector{Type: constants.SINGLE_NODE, NodeIDs: []uint8{0x01, 0x02}}
	_, err = selector.Encode()
	assert.Error(err)

	// Decode with length=0
	tlvData, _ := tlv.DecodeTLV([]byte{uint8(constants.SINGLE_NODE), 0x00, 0x00})
	_, err = selectors.DecodeNodeSelector(tlvData)
	assert.Error(err)

	// Decode with length=2
	tlvData, _ = tlv.DecodeTLV([]byte{uint8(constants.SINGLE_NODE), 0x00, 0x02, 0x01, 0x02})
	_, err = selectors.DecodeNodeSelector(tlvData)
	assert.Error(err)
}

func TestNodeListSelector_InvalidCases(t *testing.T) {
	assert := assert.New(t)

	// Empty list
	selector := &selectors.NodeSelector{Type: constants.NODE_LIST, NodeIDs: []uint8{}}
	_, err := selector.Encode()
	assert.Error(err)

	// Decode with length=0
	tlvData, _ := tlv.DecodeTLV([]byte{uint8(constants.NODE_LIST), 0x00, 0x00})
	_, err = selectors.DecodeNodeSelector(tlvData)
	assert.Error(err)
}

func TestAllNodesSelector_InvalidCases(t *testing.T) {
	assert := assert.New(t)

	// Decode with non-empty value
	tlvData, _ := tlv.DecodeTLV([]byte{uint8(constants.ALL_NODES), 0x00, 0x01, 0x42})
	_, err := selectors.DecodeNodeSelector(tlvData)
	assert.Error(err)
}

func TestNodeSelector_UnknownType(t *testing.T) {
	assert := assert.New(t)

	// Encode with unknown type
	selector := &selectors.NodeSelector{Type: constants.NodeSelector(0xFF), NodeIDs: []uint8{0x01}}
	_, err := selector.Encode()
	assert.Error(err)

	// Decode with unknown type
	tlvData, _ := tlv.DecodeTLV([]byte{0xFF, 0x00, 0x01, 0x42})
	_, err = selectors.DecodeNodeSelector(tlvData)
	assert.Error(err)
}
