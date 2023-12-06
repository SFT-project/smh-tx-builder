package main

// grpcurl -plaintext localhost:9092 spacemesh.v1.MeshService.GenesisID
// grpcurl -plaintext localhost:9092 spacemesh.v1.DebugService.NetworkInfo
// grpcurl --plaintext -d "{}" localhost:9092 spacemesh.v1.NodeService.Status
// grpcurl --plaintext -d "{}" localhost:9092 spacemesh.v1.ActivationService.Highest

// https://hakedev.substack.com/p/common-spacemesh-grpcurl-commands
// https://hakedev.substack.com/p/reduce-spacemesh-node-traffic
// https://configs.spacemesh.network/config.mainnet.json
// https://discover.spacemesh.io/networks.json

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	address := os.Getenv("ADDRESS")                      // wallet smh-rpc
	privatekey := util.FromHex(os.Getenv("PRIVATE_KEY")) //38rk
	// target := "localhost:9092"
	target := os.Getenv("RPC_URL")
	fmt.Print("send : ", address, "\ntarget : ", target, "\n")

	var addressList []string
	var amounts []uint64

	recipientsArrayString := "sm1qqqqqqpp3rwldgtjp3aq7vg6r3g0nzj8528rjng6vxfgr,sm1qqqqqqpydwlzautxe9uet6a0f72x3e9fsuvmwjqjsvk29"
	hashArrayString := "1,2"

	addressArray := strings.Split(recipientsArrayString, ",")
	hashArray := strings.Split(hashArrayString, ",")
	// 1 Smesh equals 10^9 Smidge.
	amount := uint64(20012711720) // 20.01271172 SMH
	// amount := uint64(00001271172) // 00.01271172 SMH
	for index, recipient := range addressArray {
		addressList = append(addressList, recipient)
		hash, _ := strconv.ParseUint(hashArray[index], 0, 64)
		smhAmount := amount * hash
		amounts = append(amounts, smhAmount)
		fmt.Print(index, " : ", recipient, "-", smhAmount, "\n")
	}

	sendMoney(target, string(privatekey), address, amounts, addressList)
}

func sendMoney(target string, privatekey string, address string, amounts []uint64, addressList []string) {
	ctx := context.Background()
	if len(amounts) != len(addressList) {
		panic("the size of address should be equal with the length of amount array.")
	}
	cc, _ := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer cc.Close()

	meshClient := api.NewMeshServiceClient(cc)
	meshResp, _ := meshClient.GenesisID(ctx, &api.GenesisIDRequest{})
	var genesisId types.Hash20
	copy(genesisId[:], meshResp.GenesisId)

	client := api.NewNodeServiceClient(cc)
	statusResp, _ := client.Status(context.Background(), &api.StatusRequest{})
	fmt.Printf("statusResp: %+v\n", statusResp)

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
		sendResp, err := txService.SubmitTransaction(ctx, &api.SubmitTransactionRequest{Transaction: tx})
		if err != nil {
			log.Fatalf("broadcast tx:%d, Err: %s", txResp.Tx.Nonce.Counter, err)
			panic("broadcast tx failed.")
		}
		// return the txid
		fmt.Printf("status code: %d, txid: %s, tx state: %s\n",
			sendResp.Status.Code, hex.EncodeToString(sendResp.Txstate.Id.Id), sendResp.Txstate.State.String())
	}

}
