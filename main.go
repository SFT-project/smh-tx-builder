package main

import (
	"context"
	"encoding/hex"
	"fmt"
	api "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"github.com/spacemeshos/go-spacemesh/common/types"
	"github.com/spacemeshos/go-spacemesh/common/util"
	"github.com/spacemeshos/go-spacemesh/genvm/sdk"
	walletSdk "github.com/spacemeshos/go-spacemesh/genvm/sdk/wallet"
	"github.com/spacemeshos/go-spacemesh/signing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var privatekey = util.FromHex("")
	var address = ""
	var amounts []uint64
	amounts = append(amounts)
	amounts = append(amounts)
	var addressList []string
	addressList = append(addressList, "")
	addressList = append(addressList, "")
	sendMoney(string(privatekey), address, amounts, addressList)
}

func sendMoney(privatekey string, address string, amounts []uint64, addressList []string) {
	ctx := context.Background()
	if len(amounts) != len(addressList) {
		panic("发放的收益和地址要保持一样")
	}
	cc, _ := grpc.Dial("101.36.104.253:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer cc.Close()

	meshClient := api.NewMeshServiceClient(cc)
	meshResp, _ := meshClient.GenesisID(ctx, &api.GenesisIDRequest{})
	var genesisId types.Hash20
	copy(genesisId[:], meshResp.GenesisId)

	client := api.NewNodeServiceClient(cc)
	statusResp, _ := client.Status(context.Background(), &api.StatusRequest{})
	fmt.Sprintf("%+v", statusResp)
	gstate := api.NewGlobalStateServiceClient(cc)
	resp, _ := gstate.Account(ctx, &api.AccountRequest{AccountId: &api.AccountId{Address: address}})
	nonce := resp.AccountWrapper.StateProjected.Counter
	balance := resp.AccountWrapper.StateProjected.Balance
	fmt.Printf("Sender nonce %d balance %d\n", nonce, balance.Value)
	for i := 0; i < len(amounts); i++ {
		recipientAddress, _ := types.StringToAddress(addressList[i])
		// generate the tx
		tx := walletSdk.Spend(signing.PrivateKey(privatekey), recipientAddress, amounts[i], nonce+1+uint64(i),
			sdk.WithGenesisID(genesisId),
		)
		fmt.Printf("Generated signed tx: %s\n", hex.EncodeToString(tx))

		// parse it
		txService := api.NewTransactionServiceClient(cc)
		txResp, _ := txService.ParseTransaction(ctx, &api.ParseTransactionRequest{Transaction: tx})
		fmt.Printf("parsed tx: principal: %s, gasprice: %d, maxgas: %d, nonce: %d\n",
			txResp.Tx.Principal.Address, txResp.Tx.GasPrice, txResp.Tx.MaxGas, txResp.Tx.Nonce.Counter)
		// broadcast it
		sendResp, _ := txService.SubmitTransaction(ctx, &api.SubmitTransactionRequest{Transaction: tx})
		// return the txid
		fmt.Printf("status code: %d, txid: %s, tx state: %s\n",
			sendResp.Status.Code, hex.EncodeToString(sendResp.Txstate.Id.Id), sendResp.Txstate.State.String())
	}

}

//func origin()  {
//	recipientAddress, _ := types.StringToAddress(recipientAddressStr)
//	cc, _ := grpc.Dial(nodeUri, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	defer cc.Close()
//
//	meshClient := api.NewMeshServiceClient(cc)
//	meshResp, _ := meshClient.GenesisID(ctx, &api.GenesisIDRequest{})
//	var genesisId types.Hash20
//	copy(genesisId[:], meshResp.GenesisId)
//
//	client := api.NewNodeServiceClient(cc)
//	statusResp, _ := client.Status(ctx, &api.StatusRequest{})
//
//	gstate := api.NewGlobalStateServiceClient(cc)
//	resp, _ := gstate.Account(ctx, &api.AccountRequest{AccountId: &api.AccountId{Address: principal.String()}})
//	nonce := resp.AccountWrapper.StateProjected.Counter
//	balance := resp.AccountWrapper.StateProjected.Balance
//	fmt.Printf("Sender nonce %d balance %d\n", nonce, balance.Value)
//
//	// generate the tx
//	tx := walletSdk.Spend(signing.PrivateKey(privkey), recipientAddress, amount, nonce+1,
//		sdk.WithGenesisID(genesisId),
//	)
//	fmt.Printf("Generated signed tx: %s\n", hex.EncodeToString(tx))
//
//	// parse it
//	txService := api.NewTransactionServiceClient(cc)
//	txResp, _ := txService.ParseTransaction(ctx, &api.ParseTransactionRequest{Transaction: tx})
//	fmt.Printf("parsed tx: principal: %s, gasprice: %d, maxgas: %d, nonce: %d\n",
//		txResp.Tx.Principal.Address, txResp.Tx.GasPrice, txResp.Tx.MaxGas, txResp.Tx.Nonce.Counter)
//
//	// broadcast it
//	sendResp, _ := txService.SubmitTransaction(ctx, &api.SubmitTransactionRequest{Transaction: tx})
//
//	// return the txid
//	fmt.Printf("status code: %d, txid: %s, tx state: %s\n",
//		sendResp.Status.Code, hex.EncodeToString(sendResp.Txstate.Id.Id), sendResp.Txstate.State.String())
//}
