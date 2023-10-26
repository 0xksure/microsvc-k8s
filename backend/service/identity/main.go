// / Identity package is meant to extract the identity of a user from the solana blockchain
package identity

import (
	"context"
	"encoding/binary"

	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// Identity is a struct that contains the identity of a user
func getIdentityPDA(social, username string, userId uint32) (solana.PublicKey, error) {
	identityProgramID, err := solana.PublicKeyFromBase58("identity")
	if err != nil {
		return solana.PublicKey{}, err
	}
	userIdb := make([]byte, 4)
	seeds := [][]byte{[]byte("identity"), []byte(social), userIdb}
	binary.LittleEndian.PutUint32(userIdb, userId)
	return solana.CreateProgramAddress(seeds, identityProgramID)
}

func getIdentity() {
	cluster := rpc.New("https://api.devnet.solana.com")
	ctx := context.Background()

}
