// / Identity package is meant to extract the identity of a user from the solana blockchain
package identity

import (
	"context"
	"encoding/binary"

	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/near/borsh-go"
)

type Identity struct {
	Address   solana.PublicKey
	Social    string
	UserId    uint64
	Username  string
	Bump      uint8
	SocialRaw string
}

// Identity is a struct that contains the identity of a user
func getIdentityPDA(social string, userId uint64) (solana.PublicKey, error) {
	identityProgramID, err := solana.PublicKeyFromBase58("identity")
	if err != nil {
		return solana.PublicKey{}, err
	}
	userIdb := make([]byte, 8)
	seeds := [][]byte{[]byte("identity"), []byte(social), userIdb}
	binary.LittleEndian.PutUint64(userIdb, userId)
	return solana.CreateProgramAddress(seeds, identityProgramID)
}

// / getIdentityPDA gets the identity PDA
func GetIdentity(rpcUrl string, social string, userId uint64) (Identity, error) {
	var identity Identity
	cluster := rpc.New(rpcUrl)
	ctx := context.Background()

	// Get the identity program ID:
	identityPDA, err := getIdentityPDA(social, userId)
	if err != nil {
		return identity, err
	}

	// Get the identity account:
	identityAccount, err := cluster.GetAccountInfo(ctx, identityPDA)

	err = borsh.Deserialize(identity, identityAccount.Bytes())
	if err != nil {
		return identity, err
	}
	return identity, nil

}
