package smartcontract

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	logger "github.com/ElrondNetwork/elrond-go-logger"
	"github.com/ElrondNetwork/elrond-go/data/state"
	"github.com/ElrondNetwork/elrond-go/data/transaction"
	"github.com/ElrondNetwork/elrond-go/integrationTests"
	"github.com/ElrondNetwork/elrond-go/integrationTests/vm/arwen"
	"github.com/stretchr/testify/require"
)

func TestDNS_Register(t *testing.T) {
	expectedDNSAddress := []byte{0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 180, 108, 178, 102, 195, 67, 184, 127, 204, 159, 104, 123, 190, 33, 224, 91, 255, 244, 118, 95, 24, 217}

	var empty struct{}
	arwen.DNSAddresses[string(expectedDNSAddress)] = empty
	arwen.GasSchedulePath = "../../vm/arwen/gasSchedule.toml"

	context := arwen.SetupTestContext(t)
	defer context.Close()

	context.GasLimit = 40000000
	err := context.DeploySC("dns.wasm", "0064")
	require.Nil(t, err)
	require.True(t, bytes.Equal(expectedDNSAddress, context.ScAddress))

	name := "thisisalice398"
	testname := hex.EncodeToString([]byte(name))
	context.GasLimit = 40000000
	err = context.ExecuteSCWithValue(&context.Alice, "register@"+testname, big.NewInt(100))
	require.Nil(t, err)

	context.GasLimit = 8000000
	err = context.ExecuteSCWithValue(&context.Alice, "resolve@"+testname, big.NewInt(0))
	require.Nil(t, err)

	for _, scr := range context.LastSCResults {
		if bytes.Equal(scr.GetOriginalTxHash(), context.LastTxHash) {
			data := scr.GetData()
			if len(data) > 0 {
				// The first 6 characters of data are '@6f6b@', where 6f6b means 'ok';
				// the resolved address comes after.
				resolvedAddress, err := hex.DecodeString(string(data[6:]))
				require.Nil(t, err)
				require.True(t, bytes.Equal(context.Alice.Address, resolvedAddress))
			}
		}
	}
}

func TestDNS_IOTimeout(t *testing.T) {
	logger.SetLogLevel("*:TRACE")
	integrationTests.MaxGasLimitPerBlock = uint64(30000000)

	user, err := hex.DecodeString("e21ac250ef528573b860c439d8bc1711f3f378e1a5e51d07696e90742d282655")
	require.Nil(t, err)

	relayer, err := hex.DecodeString("edd147c39f9b16542e6f0a090aa9ce1269172d34bf60d110d00c8590911db067")
	require.Nil(t, err)

	var empty struct{}
	dns, err := hex.DecodeString("0000000000000000050065071a7588f64b526dd11289738d9b719a9ac71f0015")
	require.Nil(t, err)
	integrationTests.ExtraDNSAddresses[string(dns)] = empty

	network := integrationTests.NewOneNodeNetwork()
	defer network.Stop()

	network.Mint(user, big.NewInt(1000000000000000000))
	network.Mint(relayer, big.NewInt(1000000000000000000))

	network.GoToRoundOne()

	scPath := "/var/work/Elrond/elrond-go/integrationTests/singleShard/smartContract/dns.wasm"
	scCode, err := hex.DecodeString(arwen.GetSCCode(scPath))
	require.Nil(t, err)

	scDNSAccount, err := network.Node.AccntState.LoadAccount(dns)
	require.Nil(t, err)
	scDNS := scDNSAccount.(state.UserAccountHandler)
	scDNS.SetCode(scCode)
	scDNS.SetOwnerAddress(relayer)
	network.Node.AccntState.SaveAccount(scDNS)

	relayerAccount, err := network.Node.AccntState.LoadAccount(relayer)
	relayerUserAccount := relayerAccount.(state.UserAccountHandler)
	relayerUserAccount.IncreaseNonce(3)
	network.Node.AccntState.SaveAccount(relayerUserAccount)

	network.Continue(t, 1)

	txData := "relayedTx@7b226e6f6e6365223a332c2276616c7565223a302c227265636569766572223a2241414141414141414141414641475548476e5749396b74536264455369584f4e6d3347616d7363664142553d222c2273656e646572223a223764464877352b62466c517562776f4a43716e4f456d6b584c54532f594e4551304179466b4a45647347633d222c226761735072696365223a313030303030303030302c226761734c696d6974223a32353030303030302c2264617461223a22636d566e61584e305a584a414e6a49324e545a6c4e6a6b324d545a6b4e6a6b325a54597a4e6d59334d7a5a6b4e6a453d222c22636861696e4944223a226447567a6447356c6443316c62484a76626d5174595778734c576c754c5739755a513d3d222c2276657273696f6e223a312c227369676e6174757265223a2274593277466b5874785064734271422f4a6e31375a515431393469384941493932325436424c53483035426a423638546357726d6d764f4e2b52764977303577734e47484b6435386c75677a4e5739343861674a44413d3d227d"
	network.AddTxToPool(&transaction.Transaction{
		Nonce:    0,
		Value:    big.NewInt(0),
		RcvAddr:  relayer,
		SndAddr:  user,
		GasPrice: 1000000000,
		GasLimit: 25001808,
		Data:     []byte(txData),
	})

	network.Continue(t, 1)
}
