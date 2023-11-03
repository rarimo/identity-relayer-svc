package core

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	merkle "github.com/rarimo/go-merkle"
	"github.com/rarimo/rarimo-core/x/rarimocore/crypto/pkg"
	rarimocore "github.com/rarimo/rarimo-core/x/rarimocore/types"
)

func TestName(t *testing.T) {
	t1 := &rarimocore.IdentityDefaultTransfer{
		Contract:                "0x134B1BE34911E39A8397ec6289782989729807a4",
		Chain:                   "Mumbai",
		GISTHash:                "0xb659950110ef59e462f787f875dc31a187274658a3cef911af0be6e800e2ab",
		Id:                      "0x0e4104e4e3b1894ae9b3b0acdd827c255527866111390a78f5b4fbf2211202",
		StateHash:               "0x20c78cda046b738a5bd795f041d903841f626dff91369fa41905ae63bfe9b5bf",
		StateCreatedAtTimestamp: "1689800813",
		StateCreatedAtBlock:     "38121346",
		StateReplacedBy:         "0x",
		GISTReplacedBy:          "0x",
		GISTCreatedAtTimestamp:  "1689800813",
		GISTCreatedAtBlock:      "38121346",
		ReplacedStateHash:       "0x26c64fa4ded37f23f46d037ec29465ce74a675bf99819fd46a2e97c7663c1eb8",
		ReplacedGISTHash:        "0x0e3a140c031cf5d1b0380fa2b0bac7e4b928ab3abaf299172f847adb812c3d7c",
	}

	t2 := &rarimocore.IdentityDefaultTransfer{
		Contract:                "0x134B1BE34911E39A8397ec6289782989729807a4",
		Chain:                   "Mumbai",
		GISTHash:                "0xb659950110ef59e462f787f875dc31a187274658a3cef911af0be6e800e2ab",
		Id:                      "0x0e4104e4e3b1894ae9b3b0acdd827c255527866111390a78f5b4fbf2211202",
		StateHash:               "0x26c64fa4ded37f23f46d037ec29465ce74a675bf99819fd46a2e97c7663c1eb8",
		StateCreatedAtTimestamp: "1689800763",
		StateCreatedAtBlock:     "38121325",
		StateReplacedBy:         "0x20c78cda046b738a5bd795f041d903841f626dff91369fa41905ae63bfe9b5bf",
		GISTReplacedBy:          "0x",
		GISTCreatedAtTimestamp:  "1689800813",
		GISTCreatedAtBlock:      "38121346",
		ReplacedStateHash:       "0x09d15d3014f00122dd671f24398832dd98178d5123a55ce2c497b34fa8888267",
		ReplacedGISTHash:        "0x0e3a140c031cf5d1b0380fa2b0bac7e4b928ab3abaf299172f847adb812c3d7c",
	}

	content1, _ := pkg.GetIdentityDefaultTransferContent(t1)
	fmt.Println(hexutil.Encode(content1.CalculateHash()))

	content2, _ := pkg.GetIdentityDefaultTransferContent(t2)
	fmt.Println(hexutil.Encode(content2.CalculateHash()))

	c1 := *content1
	c2 := *content2

	tree := merkle.NewTree(crypto.Keccak256, c2, c1)
	fmt.Println(hexutil.Encode(tree.Root()))

	path, _ := tree.Path(c2)
	fmt.Println(len(path))
	fmt.Println(hexutil.Encode(path[0]))

	pathHashes := make([]common.Hash, 0, len(path))
	for _, p := range path {
		pathHashes = append(pathHashes, common.BytesToHash(p))
	}

	signature := hexutil.MustDecode("0x636d3d9e1352ec42e8fd604ce95d9ac015c6c8839896d79bcf4aec771bc2d4d451455d4fd1faa73f68ffa0f475aebcd3cb496f363c958ee83b435c6bdc32f92901")
	signature[64] += 27

	proof, err := proofArgs.Pack(pathHashes, signature)
	if err != nil {
		panic(err)
	}

	fmt.Println(hexutil.Encode(proof))
}
