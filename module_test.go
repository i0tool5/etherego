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

var (
	// Accounts (got from Ganache)
	Accounts = map[string]string{
		"0x6eA84B613C5EdD129086a5AB30d02CDf635A63b6": "0xa78b21ef4f87ab590a489ae86cd27cdb5ef755a4060fa3622cc28a1f5209e693",
		"0x6De327DAb7462ee0b898548b7ae64ca094292DC2": "0x7710a4c8b74bbe9aedd78a00d9ab11a348343668534f9d2b1b2f863c8ee395c3",
		"0x423B541dEB63601Cc33E8b4B331263c26C3Bc9Ce": "0xd7fdc6701e2233bf9d75fcfeabc89fff31ce64c7dad115aac090d2685b29eaba",
		"0x8F087Bfa43e78e491666c339a1c9129b03fF8446": "0x30017efa07c155f1de25954b43d265d8efd491c1200785e05c388a9d45bdcdc3",
		"0x43a7835854B077D2647334B10550ebA38bd250B3": "0xe20375aee28830258f8f315528f12b5f8097f05977ff4e1f66b1d60adbf57849",
		"0x0c3d6433f679d6e0DEECAaa8516D1C1e2E222a1d": "0x496c12c0d61d62e0271a327c05c8acb32c39579016ba996012b99f942e39693b",
		"0x02D3daA6740b4bfb73b382364aE307B693B7A310": "0x8425f29a14f9b52a9a325afb232e468b4b6f11eb01f82df59d5d0b076f6d39a2",
		"0xD73d75da6a6cAFEaa67125bdd8978FAde362A477": "0x35768fd4670f86564f9dc69c86f7782d9926cf9d05bf6831ec52398de66b26d3",
		"0x9962aBB93e23c16Aa94bACEd0186c006DCb1333f": "0x8dbd0f4f00fe932cace5ad3ebab3fee1d89e9bd5f4950bcb80a02733d6edaf8a",
		"0x19faD65ceFEe08d38E1B96FF2DBE28b1487b99Df": "0xd2229b917ffef431fcd1688e3d93c8acd66b397a08bad326c6ff560f57ad2646",
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
		txHash, err := client.TransferTokens(ctx, addrs[k], addrs[k-1], amount)
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
