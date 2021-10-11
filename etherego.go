package etherego

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	WeiEth  = big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil)
	GweiEth = big.NewInt(0).Exp(big.NewInt(10), big.NewInt(9), nil)
)

// Accounts it's a mapping of account address to account private key
type Accounts map[string]string

// ETHClient wraps *ethclient.Client
type ETHClient struct {
	// *ethclient.Client composition
	*ethclient.Client
	Accounts *Accounts
}

// New ETHClient
func New(addr string, accs *Accounts) (*ETHClient, error) {
	client, err := ethclient.Dial(addr)
	if err != nil {
		return nil, err
	}
	return &ETHClient{client, accs}, nil
}

// TransactionFrom returns transaction sender address
func (e *ETHClient) TransactionFrom(ctx context.Context, tx *types.Transaction) (string, error) {
	chainID, err := e.ChainID(ctx)
	if err != nil {
		return "", nil
	}

	msg, err := tx.AsMessage(types.NewEIP155Signer(chainID), tx.GasFeeCap())
	if err != nil {
		return "", nil
	}

	return msg.From().Hex(), nil

}

// BlocksRange returns blocks by range
func (e *ETHClient) BlocksRange(ctx context.Context, beg, end int) ([]*types.Block, error) {
	blocks := make([]*types.Block, 0)
	for i := beg; i < end; i++ {
		b, err := e.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			// maybe return blocks, err if exception occured on non zero block?
			return nil, err
		}
		blocks = append(blocks, b)
	}
	return blocks, nil
}

// TransferTokens transfer tokens from one account to another, and returns transaction hash
func (e *ETHClient) TransferTokens(ctx context.Context,
	fromAddr, toAddr string, ethVal *big.Int) (txHash common.Hash, err error) {

	pk := (*e.Accounts)[fromAddr][2:]
	privateKey, err := crypto.HexToECDSA(pk)
	if err != nil {
		log.Fatal(err)
	}

	pubKey := privateKey.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := e.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return
	}

	gasLimit := uint64(6721975)
	gasPrice, err := e.SuggestGasPrice(ctx)
	if err != nil {
		return
	}

	dstAddr := common.HexToAddress(toAddr)
	tx := types.NewTransaction(nonce, dstAddr, ethVal, gasLimit, gasPrice, []byte("donation"))

	chainID, err := e.ChainID(ctx)
	if err != nil {
		return
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return
	}

	err = e.SendTransaction(ctx, signedTx)
	if err != nil {
		return
	}
	return signedTx.Hash(), nil
}

// SubscribeNewBlocks subscribes to new block creation events.
// Works only with wss:// connection scheme
func (e *ETHClient) SubscribeNewBlocks(ctx context.Context) (
	<-chan *types.Header, <-chan error, error) {

	var (
		headsChan = make(chan *types.Header)
		outChan   = make(chan *types.Header)
		errChan   = make(chan error)
		sub, err  = e.SubscribeNewHead(ctx, headsChan)
	)
	if err != nil {
		return nil, nil, err
	}

	// event listener goroutine
	go func(sub ethereum.Subscription, ch chan *types.Header) {
		for {
			select {
			case err = <-sub.Err():
				errChan <- err
			case header := <-headsChan:
				outChan <- header
			}
		}
	}(sub, headsChan)

	return outChan, errChan, nil
}
