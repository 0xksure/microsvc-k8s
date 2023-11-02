package identity

import (
	"testing"

	bin "github.com/gagliardetto/binary"
	solana "github.com/gagliardetto/solana-go"
)

func TestMain(t *testing.T) {

	t.Run("TestMain", func(t *testing.T) {
		t.Log("Testing main")
		expectedPK := "2hjRDP8CFChV9kC54wpN6J2mKqJqVHPg8MH6Jwe5FLca"
		social := "github"
		userId := uint64(47750504)
		identity, err := getIdentityPDA(social, userId)
		if err != nil {
			t.Errorf("failed to get identity PDA for social %s and userId %d", social, userId)
		}
		if identity == (solana.PublicKey{}) {
			t.Errorf("identity is empty")
		}
		if identity.String() != expectedPK {
			t.Errorf("identity is not equal to expectedPK: expected: %s, got: %s", expectedPK, identity.String())
		}
	})

	t.Run("Test deserialization", func(t *testing.T) {
		binaryData := []byte{58, 132, 5, 12, 176, 164, 85, 112, 241, 195, 98, 152, 244, 207, 69, 134, 124, 157, 209, 178, 222, 193, 129, 53, 67, 124, 60, 135, 82, 173, 111, 181, 43, 68, 115, 130, 234, 135, 238, 222, 32, 0, 0, 0, 103, 105, 116, 104, 117, 98, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104, 157, 216, 2, 0, 0, 0, 0, 32, 0, 0, 0, 48, 120, 107, 115, 117, 114, 101, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 255, 6, 0, 0, 0, 103, 105, 116, 104, 117, 98, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		var identity Identity
		borshDec := bin.NewBorshDecoder(binaryData)

		err := borshDec.Decode(&identity)
		if err != nil {
			t.Errorf("failed to deserialize binary data %v %s", binaryData, err)
		}
		if identity.UserId != 47750504 {
			t.Errorf("userId is not equal to 47750504, but %d", identity.UserId)
		}

	})

	t.Run("Integration test ", func(t *testing.T) {
		rpcUrl := "https://api.devnet.solana.com"
		social := "github"
		userId := uint64(47750504)
		identity, err := GetIdentity(rpcUrl, social, userId)
		if err != nil {
			t.Errorf("failed to get identity for social %s and userId %d and %v", social, userId, err)
		}

		if identity == (Identity{}) {
			t.Errorf("identity is empty")
		}

		if identity.UserId != userId {
			t.Errorf("userId is not equal to %d, but %d", userId, identity.UserId)
		}

	})

}
