package registry

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func TestParseVehicleMinted(t *testing.T) {
	abix, err := AbiMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}

	l1 := ceLog{
		Topics: []string{
			"0x0c2616265c4fd089569533525abc7b19b9f82b423d7cdb61801490b8f9e0ce59",
			"0xae25b7c67f95ee4f67c09e3ce3c623354a29ca628b2ec02e873f3fc18dbb1ee3",
			"0x00000000000000000000000000000000000000000000000000000000000000d3",
		},
		Data: "",
	}

	nodeMintedEvent := abix.Events["NodeMinted"]

	l2 := convertLog(&l1)
	out := new(AbiNodeMinted)
	if len(l2.Data) > 0 {
		if err := abix.UnpackIntoInterface(out, nodeMintedEvent.Name, l2.Data); err != nil {
			t.Fatal(err)
		}
	}
	var indexed abi.Arguments
	for _, arg := range nodeMintedEvent.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	if err := abi.ParseTopics(out, indexed, l2.Topics[1:]); err != nil {
		t.Fatal(err)
	}

	if out.NodeId != big.NewInt(211) {
		t.Errorf("Node id should have been %d but was %s", 211, out.NodeId)
	}
}
