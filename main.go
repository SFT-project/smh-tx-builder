package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"os"
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
	fmt.Print("send : ", address, "-", privatekey, "\n")

	var addressList []string
	var amounts []uint64

	recipientsArrayString := "sm1qqqqqqq0mk66e0k6rg84p5spc9vhet6s8nxnk7cge46nj,sm1qqqqqqpydwlzautxe9uet6a0f72x3e9fsuvmwjqjsvk29,sm1qqqqqqpydwlzautxe9uet6a0f72x3e9fsuvmwjqjsvk29,sm1qqqqqqputfs7rj6xwvjp0jku9xnef9ukyqz94nszekp5z,sm1qqqqqqyf59u3hqfvwm4jwtxn5ae0zyhct4gtnvqv6xcfq,sm1qqqqqqyrfh3r08jceuy2w2y0pwjwtae2ndjv8tquh5hd0,sm1qqqqqqy8qqqjwq89uxsgaejprnhjdc0enhkd97gg63akz,sm1qqqqqqya6zttr6m6auhpxc7xntlu60yy02u27pcgzg5wg,sm1qqqqqqydmp77g6u8at04jzkhntr0r7kpalpm8mg9f7v96,sm1qqqqqqzlk4nx9z5m948zkkgpdm03hzpt4435pdqzn088a,sm1qqqqqqr56ydar0tfl9gerxrru9thhxxrjl9ym6g4n9cpv,sm1qqqqqqyn36shlpsm2p720zuzarxff8lhxa9uqgc2jte4k,sm1qqqqqqrs08dkp4lwl6hkcpqvzn8c2zzc7kyw6ucasgmvc,sm1qqqqqqppcp5zeug5tj8sw4xgqfpfqp3strk0gnqk583eg,sm1qqqqqqya6zttr6m6auhpxc7xntlu60yy02u27pcgzg5wg,sm1qqqqqqya6zttr6m6auhpxc7xntlu60yy02u27pcgzg5wg,sm1qqqqqqya6zttr6m6auhpxc7xntlu60yy02u27pcgzg5wg,sm1qqqqqqya6zttr6m6auhpxc7xntlu60yy02u27pcgzg5wg,sm1qqqqqqywq359dhpyrqm4s09xh5ndcrmtsv4x05szcu4a5,sm1qqqqqqp2wgueqdp7ywdnryglecfczzx684rnafc8ffddn,sm1qqqqqqz0mtd7tgjr00s04vdl67pw0h37sc9mh4shrremy,sm1qqqqqqxaa8ecfj3z6slkgg76m2prg6qzmcwx2wqra44fn,sm1qqqqqqx3cm0fns4zckxyds9crthz4cec25aealgaxzaqc,sm1qqqqqqqt60yhdcuhz7q8z53vwh0fz5h7h3xh8kqwc5vm4,sm1qqqqqq95peqadnp5trur9qv9n4nv6elpqw9z0vqdg3nqk,sm1qqqqqqqes8h5tsckxyeryrewp78yahxey2c2eechzvrsd,sm1qqqqqq8hw6gsd7qpwqx9kesuw7k53tse4vgcz4g4gcg5x,sm1qqqqqqr0tf9m3394cxwz6jycwte4x6vxz3r4tes9c2vxd,sm1qqqqqqrlpmfqt8ts6y390w55x2zekye722kdkuqzez6p2,sm1qqqqqqzh76hv0kvww3t6hm0wgrunu4myatm4vgg4zqx07,sm1qqqqqqp72w5qxjc676p370uemt966decntgdjmshedr80,sm1qqqqqq9t9y7utx4v4t2wgf0yg44na4hqmypcj2cysmmmh,sm1qqqqqq9dc47nq3s74gj0e3y82xxnnfry3ngjd6g4tl9dd,sm1qqqqqq9cfhgw0a9kf0d7nrdj8nc3zam34t5hwkshx50vc,sm1qqqqqqytnmp6rpt3vhxdmatx08sax0a4la482hgg480qw,sm1qqqqqqrrx77nf89zsyhn36meh0hw3vjms2j2qucvpdp6x,sm1qqqqqqy74j3lcpdrd5rdzg2ey9463d6aefzs7ugm9enkn,sm1qqqqqq8g76nmepghcaqtdz77806z76t88kawdhq8f2mg3,sm1qqqqqq88uuzyrcx7k05h5xse27sdkrjry45zrkswy67py,sm1qqqqqq83pwce0t5778uvsduehhwvlds64gqghysqh4a5u,sm1qqqqqqyqpv6cs7cngkmgrue4fnlp70j666czj6s8970uw,sm1qqqqqq9pngcfvq89fpgmqsus70c097fp9uy0mvc9gp64c,sm1qqqqqqreyrr9ks3pjh8s50jww6nleets5g59fhq58xxk4,sm1qqqqqqqgcqud7l56f5eq477rc9scz0uzr93f4jq7way4x,sm1qqqqqqzhvkyhfdjf8myuc5q29ngmn7tf6g9kffca3d69p,sm1qqqqqqxnfh2gt43klq49dvg2eguk5vnttc76waqt5t0qx,sm1qqqqqqrljlvxf50fzls3yrsuflfg9vvjpdd29lc8vymht,sm1qqqqqqpe695l5fag5fq69x5kx89amhemrvn7y6gssn5q0,sm1qqqqqqxduzz2g97xzc8z0mxlwhkug86vp6du0tq44n3dl,sm1qqqqqqzm35qg9sdx4fp7h542gg5wvynvng4kdpg4wr9ku,sm1qqqqqq86yvrj5kd0gd5fdkszaeznk06sap4y2jqkvn69m,sm1qqqqqq8wr9uywneta8wpjlhyx56ez390s7qlvfcxcrpca,sm1qqqqqqr945vgsdc6dd5v5cvd35kfd32dwxel6cqfu550x,sm1qqqqqq94wjqjsvywv8uxuch304r5n2ga76w549qvstuuk,sm1qqqqqqykmwfmtmf0c9w9a0ypapkxmd3zfp0kwps2gk92j,sm1qqqqqqreyrr9ks3pjh8s50jww6nleets5g59fhq58xxk4,sm1qqqqqqreyrr9ks3pjh8s50jww6nleets5g59fhq58xxk4,sm1qqqqqqreyrr9ks3pjh8s50jww6nleets5g59fhq58xxk4,sm1qqqqqq9pngcfvq89fpgmqsus70c097fp9uy0mvc9gp64c,sm1qqqqqqqgcqud7l56f5eq477rc9scz0uzr93f4jq7way4x,sm1qqqqqqzhvkyhfdjf8myuc5q29ngmn7tf6g9kffca3d69p,sm1qqqqqqxnfh2gt43klq49dvg2eguk5vnttc76waqt5t0qx,sm1qqqqqqrljlvxf50fzls3yrsuflfg9vvjpdd29lc8vymht,sm1qqqqqqpe695l5fag5fq69x5kx89amhemrvn7y6gssn5q0,sm1qqqqqqxduzz2g97xzc8z0mxlwhkug86vp6du0tq44n3dl,sm1qqqqqqzm35qg9sdx4fp7h542gg5wvynvng4kdpg4wr9ku,sm1qqqqqq86yvrj5kd0gd5fdkszaeznk06sap4y2jqkvn69m,sm1qqqqqq8wr9uywneta8wpjlhyx56ez390s7qlvfcxcrpca,sm1qqqqqqr945vgsdc6dd5v5cvd35kfd32dwxel6cqfu550x,sm1qqqqqqreyrr9ks3pjh8s50jww6nleets5g59fhq58xxk4,sm1qqqqqq9pngcfvq89fpgmqsus70c097fp9uy0mvc9gp64c,sm1qqqqqqqgcqud7l56f5eq477rc9scz0uzr93f4jq7way4x,sm1qqqqqqzhvkyhfdjf8myuc5q29ngmn7tf6g9kffca3d69p,sm1qqqqqqxnfh2gt43klq49dvg2eguk5vnttc76waqt5t0qx,sm1qqqqqqrljlvxf50fzls3yrsuflfg9vvjpdd29lc8vymht,sm1qqqqqqpe695l5fag5fq69x5kx89amhemrvn7y6gssn5q0,sm1qqqqqqxduzz2g97xzc8z0mxlwhkug86vp6du0tq44n3dl,sm1qqqqqqzm35qg9sdx4fp7h542gg5wvynvng4kdpg4wr9ku,sm1qqqqqq86yvrj5kd0gd5fdkszaeznk06sap4y2jqkvn69m,sm1qqqqqq8wr9uywneta8wpjlhyx56ez390s7qlvfcxcrpca,sm1qqqqqqr945vgsdc6dd5v5cvd35kfd32dwxel6cqfu550x,sm1qqqqqqreyrr9ks3pjh8s50jww6nleets5g59fhq58xxk4,sm1qqqqqq9pngcfvq89fpgmqsus70c097fp9uy0mvc9gp64c,sm1qqqqqqqgcqud7l56f5eq477rc9scz0uzr93f4jq7way4x,sm1qqqqqqpws27ww5pzvcxue4979cmlnfwmsen3z3qs2qrxz,sm1qqqqqqzhvkyhfdjf8myuc5q29ngmn7tf6g9kffca3d69p,sm1qqqqqq9ktf24xgmt87036fzyh38m08983cjwm0cyyhdrs"

	addressArray := strings.Split(recipientsArrayString, ",")
	// amountsString := "1000000000" // 1 Smesh equals 10^9 Smidge.
	amount := uint64(28776332030) // 28.77633203 SMH
	for index, recipient := range addressArray {
		addressList = append(addressList, recipient)
		amounts = append(amounts, amount)
		fmt.Print(index, " : ", recipient, "-", amount, "\n")
	}

	// sendMoney(string(privatekey), address, amounts, addressList)
}

func sendMoney(privatekey string, address string, amounts []uint64, addressList []string) {
	ctx := context.Background()
	if len(amounts) != len(addressList) {
		panic("发放的收益和地址要保持一样")
	}
	cc, _ := grpc.Dial("localhost:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
