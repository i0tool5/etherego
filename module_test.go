package etherego_test

import (
	"context"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/i0tool5/etherego"
)

var (
	// Accounts (got from Ganache)
	Accounts = map[string]string{
		"0xc91a7a9a34574607c58311B9492Efad0a8ca580D": "0xabf7e203e101698b51334bf0e031286c5372a8b2f1fcdd779cd126209a539d97",
		"0xa1A0f120d399C83b02a7f1DD64D088A49a17d0Ef": "0x6f0297f7e236d296333a8543206bfedb779d2850da3da852d7e0a081f47bbcb0",
		"0xe1EC783275Ca755aBCF06C08dAB3e7611B889948": "0x74bf019fe8b01ee5e7b95bec4aaf589e011bfbf5ff2115d2c3bfd29d5c3a61d2",
		"0x3c851aE8c79e792AD202A21a5640aA10Dbe0d4c8": "0xd1cd2a01167736aa810527caf41fe297f0a93bd9754f941f3e244e8ede1344c4",
		"0x408EBae0F603FDe53399962225c82CD57366ecd7": "0xd1658559f680c756be858025f956bdec2a4bf98b828559ba7f45f4fbd92f6929",
		"0xe4A102Ed256973276FA89b768C3199d6383b41C9": "0x72f577d4ebae7b35086136869ba925b8f29a83d9b7e053ff9361f7dd40d6838f",
		"0x3284B2cF004cE4De13D64f814798Bbfe18b90088": "0x6e82755a7a81189c1462d58277be44652e5f7b2f97251c505ddcc790f34e4141",
		"0x0799018C14F63006A43DB399dcA06fD0DFc861d8": "0xf313299bca1a8ea56bfa7e970002e902313a927e82330ad7ddd3c5924268c3e6",
		"0x97c71e285C015fE617527911b0f150865Eb4b476": "0xfad90c047bfe0e1f6439e524188d49a26f6cf4e5a51fd7ae6a1e660ca11ecd52",
		"0xecba04Fed95F856fD557C085D340AC9Eaa79Be62": "0x2504fa6088073d44e75d791bc6e0267ed9928b3aa43f244d99110d90d1ed99e1",
	}
)

func TestAll(t *testing.T) {
	var (
		ctx   = context.Background()
		accs  = etherego.Accounts(Accounts)
		addrs = make([]string, 0)
	)

	for k := range Accounts {
		addrs = append(addrs, k)
	}

	client, err := etherego.New("http://localhost:8545", &accs)

	if err != nil {
		log.Fatal(err)
	}

	// block, err := client.BlockByNumber(ctx, big.NewInt(5))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// amount := fromEthToWei(big.NewFloat(0.1))
	// txHash := client.transferTokens(ctx, accounts[1], accounts[0], amount)

	// trx, pending, err := client.TransactionByHash(ctx, txHash)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if pending {
	// 	fmt.Println("Transaction is pending:", trx.Hash().Hex())
	// }

	// b, err := client.PendingBalanceAt(ctx, common.HexToAddress(accounts[0]))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Pending balance:", b)

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
