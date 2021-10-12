package etherego_test

import (
	"context"
	"log"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/i0tool5/etherego"
)

/*

Available Accounts
(0) 0xe840A468E935C38892f7ddcefE5184b943cA56f0 (100 ETH)
(1) 0x687ea5fD9b32591b178291d047655700C46C33eF (100 ETH)
(2) 0x35A621A4c17f0D96a994Bf9b189736f04C50400b (100 ETH)
(3) 0x444b34473f7C2aDf6e5575b30aEddaE467a2Bc48 (100 ETH)
(4) 0x18b3605D8cD5C885d961d79E0367c020C6D9cf0E (100 ETH)
(5) 0x42c6F42155d7f1462a295Ff475C3888DeaB73118 (100 ETH)
(6) 0x1c7502e769416ea80251d64Ec3D119E8f96cf020 (100 ETH)
(7) 0x6410AFBA3966D758A34EaFedF842f1E2B8777fA4 (100 ETH)
(8) 0x52437495C2ab8863A960B2E15e1F2A854a33Dd58 (100 ETH)
(9) 0x8bceb930a0dC49348f16b8532e25791698F7d0c8 (100 ETH)

Private Keys
==================
(0) 0xc4d1862347ae81d6508a4d91568a707955e403ad4b8b4ee2355f3b6163fd9b8c
(1) 0x78d6f12bafdd6f4632e6e52d2845baae1026567d17410addf53f99ba6370216d
(2) 0x98bcd20a64877c4a37fa0de91e7e4a0559ee96660c7495b521b551c31da55abb
(3) 0x13a58f6b6066f9c2cca94a4daf0af4a20b639523b00fd626cd6765779788c604
(4) 0xffc2bbf21ff6f00fcce0b9b6c6687b2bd789aac4ee7279cef650ec51c3bf1510
(5) 0x8dad58bd48cab77b40ae1569207de7572b06d83f9ce96f3e9e133e99975f9f78
(6) 0x66292af8ceb2b29c244d7cecacb74a84a5cc9bb1a9c02aaa8b8c4e178e6769d2
(7) 0x07f64f384471db753868231f56492bc7de9761c788a01c8fd213f7118acedc5e
(8) 0x993c9403d2b05107fd0005c81a01efc87ac7ca023d421ff9c8c6b6c7e0d323c1
(9) 0x24bc22fef5b22062be43d305c90ea37baea80ef622409952883582c30b453f9f

*/

var (
	// Accounts (got from Ganache)
	Accounts = map[string]string{
		"0xe840A468E935C38892f7ddcefE5184b943cA56f0": "0xc4d1862347ae81d6508a4d91568a707955e403ad4b8b4ee2355f3b6163fd9b8c",
		"0x687ea5fD9b32591b178291d047655700C46C33eF": "0x78d6f12bafdd6f4632e6e52d2845baae1026567d17410addf53f99ba6370216d",
		"0x35A621A4c17f0D96a994Bf9b189736f04C50400b": "0x98bcd20a64877c4a37fa0de91e7e4a0559ee96660c7495b521b551c31da55abb",
		"0x444b34473f7C2aDf6e5575b30aEddaE467a2Bc48": "0x13a58f6b6066f9c2cca94a4daf0af4a20b639523b00fd626cd6765779788c604",
		"0x18b3605D8cD5C885d961d79E0367c020C6D9cf0E": "0xffc2bbf21ff6f00fcce0b9b6c6687b2bd789aac4ee7279cef650ec51c3bf1510",
		"0x42c6F42155d7f1462a295Ff475C3888DeaB73118": "0x8dad58bd48cab77b40ae1569207de7572b06d83f9ce96f3e9e133e99975f9f78",
		"0x1c7502e769416ea80251d64Ec3D119E8f96cf020": "0x66292af8ceb2b29c244d7cecacb74a84a5cc9bb1a9c02aaa8b8c4e178e6769d2",
		"0x6410AFBA3966D758A34EaFedF842f1E2B8777fA4": "0x07f64f384471db753868231f56492bc7de9761c788a01c8fd213f7118acedc5e",
		"0x52437495C2ab8863A960B2E15e1F2A854a33Dd58": "0x993c9403d2b05107fd0005c81a01efc87ac7ca023d421ff9c8c6b6c7e0d323c1",
		"0x8bceb930a0dC49348f16b8532e25791698F7d0c8": "0x24bc22fef5b22062be43d305c90ea37baea80ef622409952883582c30b453f9f",
	}
)

func TestDefaultConn(t *testing.T) {
	var (
		ctx   = context.Background()
		accs  = etherego.Accounts(Accounts)
		addrs = make([]string, 0)
	)

	// prepare wallets
	for k := range Accounts {
		addrs = append(addrs, k)
	}

	client, err := etherego.New("http://localhost:8545", &accs)

	if err != nil {
		t.Fatalf("error creating client %s", err)
	}

	rand.Seed(time.Now().Unix())

	for k := 1; k < len(addrs); k++ {
		amount, err := etherego.EthToWei(big.NewFloat(float64(rand.Intn(100)) / 100))
		if err != nil {
			t.Fatal(err)
		}
		txHash, err := client.TransferTokens(ctx, addrs[k], addrs[k-1], nil, amount)
		if err != nil {
			t.Fatal(err)
		}
		if txHash.Hex() == "" {
			t.Fatal("no transaction created")
		}
	}

	for _, a := range addrs {
		bal, err := client.BalanceAt(ctx, common.HexToAddress(a), nil)
		if err != nil {
			t.Fatal(err)
		}
		b, err := etherego.WeiToEth(bal)
		if err != nil {
			log.Fatal(err)
		}
		t.Logf("Balance of %s -> %v\n", a, b)
	}

	latest, err := client.BlockNumber(ctx)
	if err != nil {
		log.Fatal(err)
	}

	blocks, err := client.BlocksRange(ctx, 6, int(latest))
	if err != nil {
		t.Fatalf("getting blocks: %v", err)
	}

	for _, block := range blocks {
		for _, tr := range block.Transactions() {
			from, err := client.TransactionFrom(ctx, tr)
			if err != nil {
				t.Fatalf("getting transactin from: %v", err)
			}
			val, err := etherego.WeiToEth(tr.Value())
			if err != nil {
				t.Fatalf("wei to eth: %v", err)
			}
			t.Logf(
				"BLOCK %s TRANS HASH: %s FROM %s TO %s AMOUNT -> %v\n",
				block.Number(), tr.Hash().Hex(), from, tr.To(), val,
			)
		}
	}
}

func TestWssConn(t *testing.T) {
	var (
		ctx           = context.Background()
		accs          = etherego.Accounts(Accounts)
		breaker       = make(chan struct{})
		headers       = make([]*types.Header, 0)
		infuraAddress = "wss://mainnet.infura.io/ws/v3/<project_id>" // replace <project_id> to walid project id
	)

	client, err := etherego.New(infuraAddress, &accs)

	if err != nil {
		t.Fatalf("error creating client %s", err)
	}

	hds, errs, err := client.SubscribeNewBlocks(ctx)
	if err != nil {
		t.Fatalf("subscribing new blocks %s", err)
	}

	go func() {
		for {
			select {
			case err := <-errs:
				t.Logf("listener has error %s", err)
			case <-breaker:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case h := <-hds:
				t.Logf("got new block header")
				headers = append(headers, h)
			case <-breaker:
				return
			}
		}
	}()

	// waiting 30 seconds for new block (15s eth)
	<-time.After(30 * time.Second)
	close(breaker)

	t.Logf("headers: %v", headers)
	if len(headers) == 0 {
		t.Fatal("no new blocks received")
	}
}
