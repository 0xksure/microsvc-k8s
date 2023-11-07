package bounty_program

import (
	"context"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/err/generated/bounty"
	solana "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/pkg/errors"
)

type BountyData struct{}

type BountyAccounts struct {
}

func GetBountyProgramId() solana.PublicKey {
	return bounty.ProgramID
}

func GetSignerKeysFromEnv() (solana.PrivateKey, error) {
	return solana.PrivateKeyFromBase58(os.Getenv("WALLET_SECRET_KEY"))
}

func GetProtocolPDA() (solana.PublicKey, uint8, error) {
	bountyProgramId := GetBountyProgramId()

	// Get the protocol PDA:
	return solana.FindProgramAddress([][]byte{[]byte("BOUNTY_SANDBLIZZARD")}, bountyProgramId)
}

func GetfeeCollectorPDA(mint solana.PublicKey) (solana.PublicKey, uint8, error) {
	bountyProgramId := GetBountyProgramId()
	seeds := [][]byte{[]byte("BOUNTY_SANDBLIZZARD"), []byte("FEE_COLLECTOR"), mint[:]}

	return solana.FindProgramAddress(seeds, bountyProgramId)
}

func GetBountyDenominationPDA(mint solana.PublicKey) (solana.PublicKey, uint8, error) {
	bountyProgramId := GetBountyProgramId()
	seeds := [][]byte{[]byte("BOUNTY_SANDBLIZZARD"), []byte("DENOMINATION"), mint[:]}
	// Get the bounty denomination PDA:
	return solana.FindProgramAddress(seeds, bountyProgramId)
}

func GetBountyPDA(bountyId uint64) (solana.PublicKey, uint8, error) {
	bountyProgramId := GetBountyProgramId()

	bountyIdb := make([]byte, 8)
	binary.LittleEndian.PutUint64(bountyIdb, bountyId)
	seeds := [][]byte{[]byte("BOUNTY_SANDBLIZZARD"), bountyIdb}
	// Get the bounty PDA:
	return solana.FindProgramAddress(seeds, bountyProgramId)
}

func GetEscrowPDA(bountyPk solana.PublicKey) (solana.PublicKey, uint8, error) {
	bountyProgramId := GetBountyProgramId()
	// Get the escrow PDA:
	return solana.FindProgramAddress([][]byte{[]byte("BOUNTY_SANDBLIZZARD"), bountyPk[:]}, bountyProgramId)
}

func GetSolverTokenAccounts(mint solana.PublicKey, solvers []solana.PublicKey) ([]solana.PublicKey, error) {

	/// get token accounts for solvers
	solverAtas := make([]solana.PublicKey, len(solvers))
	for _, solver := range solvers {
		ata, _, err := solana.FindAssociatedTokenAddress(solver, mint)
		if err != nil {
			break
		}
		solverAtas = append(solverAtas, ata)
	}

	// Get the escrow PDA:
	return solverAtas, nil
}

// CollectErrors collects errors into a single error
func CollectErrors(errors []error) error {
	if len(errors) > 0 {
		// combine errors
		var errStr string
		for _, err := range errors {
			errStr += err.Error()
		}
		return fmt.Errorf(errStr)
	}
	return nil
}

// GetAndCheckSolverTokenAccounts gets the solver token accounts and checks that they exist
func GetAndCheckSolverTokenAccounts(ctx context.Context, mint solana.PublicKey, solvers []solana.PublicKey, rpcClient *rpc.Client) ([]solana.PublicKey, error) {
	/// get token accounts for solvers
	solverAtas := make([]solana.PublicKey, len(solvers))
	var errs []error
	for _, solver := range solvers {
		ata, _, err := solana.FindAssociatedTokenAddress(solver, mint)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "Failed to find ata for solver %s and mint %s", solver, mint))
			break
		}

		if _, err := rpcClient.GetAccountInfo(ctx, ata); err != nil {
			errs = append(errs, errors.Wrapf(err, "Failed to get account info for solver %s and mint %s and ata %s", solver, mint, ata.String()))
			break
		}
		solverAtas = append(solverAtas, ata)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("Failed to get account info for solvers. Cause: %s", CollectErrors(errs))
	}
	// Get the escrow PDA:
	return solverAtas, nil
}

// CompleteBountyAsRelayer completes the bounty as the relayer
func CompleteBountyAsRelayer(rpcUrl string, bountyId uint64, solverPks []solana.PublicKey, mint solana.PublicKey) error {
	cluster := rpc.New(rpcUrl)
	ctx := context.Background()

	signer, err := GetSignerKeysFromEnv()
	if err != nil {
		return err
	}

	protocol, _, err := GetProtocolPDA()
	if err != nil {
		return err
	}
	bountyProgramId := GetBountyProgramId()

	feeCollector, _, err := GetfeeCollectorPDA(mint)
	if err != nil {
		return err
	}

	bountyDenomination, _, err := GetBountyDenominationPDA(mint)
	if err != nil {
		return err
	}

	bountyPk, _, err := GetBountyPDA(bountyId)
	if err != nil {
		return err
	}

	escrow, _, err := GetEscrowPDA(bountyPk)
	if err != nil {
		return err
	}

	solvers, err := GetAndCheckSolverTokenAccounts(ctx, mint, solverPks, cluster)
	if err != nil {
		return err
	}

	var accountMetaSlice solana.AccountMetaSlice
	accountMetaSlice.Append(solana.NewAccountMeta(signer.PublicKey(), true, true))
	accountMetaSlice.Append(solana.NewAccountMeta(protocol, false, false))
	accountMetaSlice.Append(solana.NewAccountMeta(bountyProgramId, false, false))
	accountMetaSlice.Append(solana.NewAccountMeta(feeCollector, false, true))
	accountMetaSlice.Append(solana.NewAccountMeta(bountyDenomination, false, false))
	accountMetaSlice.Append(solana.NewAccountMeta(bountyPk, false, true))
	accountMetaSlice.Append(solana.NewAccountMeta(escrow, false, true))
	for _, solver := range solvers {
		accountMetaSlice.Append(solana.NewAccountMeta(solver, false, true))
	}
	accountMetaSlice.Append(solana.NewAccountMeta(solana.SystemProgramID, false, false))
	accountMetaSlice.Append(solana.NewAccountMeta(solana.TokenProgramID, false, false))

	builder := bounty.NewAddBountyDenominationInstructionBuilder()
	err = builder.SetAccounts(accountMetaSlice)
	if err != nil {
		return err
	}
	ix := builder.Build()
	txBuilder := solana.NewTransactionBuilder()
	txBuilder = txBuilder.AddInstruction(ix)
	tx, err := txBuilder.Build()
	if err != nil {
		return err
	}

	sig, err := cluster.SendTransaction(ctx, tx)
	if err != nil {
		return err
	}
	fmt.Println("Signature: ", sig)
	out, err := cluster.GetConfirmedTransaction(ctx, sig)
	if err != nil {
		return err
	}
	fmt.Println("Transaction: ", out)
	return nil

}
