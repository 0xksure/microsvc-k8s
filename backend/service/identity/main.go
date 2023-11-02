// / Identity package is meant to extract the identity of a user from the solana blockchain
package identity

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"os"

	bin "github.com/gagliardetto/binary"
	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
)

type Identity struct {
	Discriminator [8]uint8
	Address       solana.PublicKey
	Social        string
	UserId        uint64
	Username      string
	Bump          uint8
	SocialRaw     string
}

// Identity is a struct that contains the identity of a user
func getIdentityPDA(social string, userId uint64) (solana.PublicKey, error) {
	identityProgramID, err := solana.PublicKeyFromBase58("3Nt1tyTJ6VBf4APaPPWixXFJr6DtfGvvTwHY1aGYT4Ws")
	if err != nil {
		return solana.PublicKey{}, errors.Wrapf(err, "failed to get identity program ID")
	}
	userIdb := make([]byte, 8)
	binary.LittleEndian.PutUint64(userIdb, userId)
	seeds := [][]byte{[]byte("identity"), []byte(social), userIdb}

	pubKey, _, err := solana.FindProgramAddress(seeds, identityProgramID)
	return pubKey, err
}

// / getIdentityPDA gets the identity PDA
func GetIdentity(rpcUrl string, social string, userId uint64) (Identity, error) {
	var identity Identity
	cluster := rpc.New(rpcUrl)
	ctx := context.Background()

	// Get the identity program ID:
	identityPDA, err := getIdentityPDA(social, userId)
	if err != nil {
		return identity, errors.Wrapf(err, "failed to get identity PDA for social %s and userId %d", social, userId)
	}

	// Get the identity account:
	identityAccount, err := cluster.GetAccountInfo(ctx, identityPDA)
	if err != nil {
		return identity, errors.Wrapf(err, "failed to get identity account %s for social: %s and userId: %d", identityPDA, social, userId)
	}
	json.NewEncoder(os.Stdout).Encode(identityAccount)

	binaryData := identityAccount.GetBinary()
	dec := bin.NewBorshDecoder(binaryData)
	err = dec.Decode(&identity)
	if err != nil {
		return identity, errors.Wrapf(err, "failed to deserialize identity account %v with binary data %v", identityAccount, binaryData)
	}

	return identity, nil

}
