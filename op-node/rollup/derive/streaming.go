package derive

import (
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/predeploys"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	StreamFuncSignature    = "stream()"
	StreamFuncBytes4       = crypto.Keccak256([]byte(StreamFuncSignature))[:4]
	StreamDepositerAddress = common.HexToAddress("0xdeaddeaddeaddeaddeaddeaddeaddeaddead0001")
	StreamAddress          = predeploys.StreamingAddr
)

func MarshalBinary() ([]byte, error) {
	data := make([]byte, 4)
	offset := 0
	copy(data[offset:4], StreamFuncBytes4)
	return data, nil
}

// StreamDeposit creates a Stream deposit transaction.
func StreamDeposit(seqNumber uint64, streamGasLimit uint64, block eth.BlockInfo) (*types.DepositTx, error) {
	data, err := MarshalBinary()
	if err != nil {
		return nil, err
	}
	source := L1InfoDepositSource{
		L1BlockHash: block.Hash(),
		SeqNumber:   seqNumber,
	}
	return &types.DepositTx{
		SourceHash:          source.SourceHash(),
		From:                StreamDepositerAddress,
		To:                  &StreamAddress,
		Mint:                nil,
		Value:               big.NewInt(0),
		Gas:                 streamGasLimit,
		IsSystemTransaction: true,
		Data:                data,
	}, nil
}

// StreamDepositBytes returns a serialized stream transaction.
func StreamDepositBytes(seqNumber uint64, streamGasLimit uint64, Stream eth.BlockInfo) ([]byte, error) {
	dep, err := StreamDeposit(seqNumber, streamGasLimit, Stream)
	if err != nil {
		return nil, fmt.Errorf("failed to create L1 info tx: %w", err)
	}
	l1Tx := types.NewTx(dep)
	opaqueL1Tx, err := l1Tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to encode L1 info tx: %w", err)
	}
	return opaqueL1Tx, nil
}
