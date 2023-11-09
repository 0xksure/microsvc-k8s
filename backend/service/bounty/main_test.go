package bounty_program

import (
	"os"
	"testing"

	"github.com/gagliardetto/solana-go"
)

func TestCompleteBounty(t *testing.T) {
	t.Log("Testing complete bounty")

	t.Run("Test Complete bounty as relayer against an rpc", func(t *testing.T) {
		t.Log("Testing complete bounty as relayer against an rpc")
		os.Setenv("WALLET_SECRET_KEY", "nk6qnXoyWT1FjHFqNVUeBB2Kzn4RUDA6kcfp79ymMzZrpzvji2PPhxkncvSdUSj9n26FTMtKT4KNw29xd9mRGXv")
		rpcUrl := "https://api.devnet.solana.com"
		bountyId := uint64(1981506938)
		solverPks := []solana.PublicKey{
			solana.MustPublicKeyFromBase58("CNY467c6XURCPjiXiKRLCvxdRf3bpunagYTJpr685gPv"),
		}
		mint := solana.MustPublicKeyFromBase58("sandphoQsRiNd85VgRrdSXdhS56d58Xa9iDKwdnKfWR")
		err := CompleteBountyAsRelayer(rpcUrl, bountyId, solverPks, mint)
		if err != nil {
			t.Errorf("failed to complete bounty as relayer %v", err)
		}
	})
}
