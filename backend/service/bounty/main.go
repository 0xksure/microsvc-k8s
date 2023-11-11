package bounty_program

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

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

func GetRelayerPDA(owner solana.PublicKey) (solana.PublicKey, uint8, error) {
	bountyProgramId := GetBountyProgramId()
	// Get the relayer PDA:
	return solana.FindProgramAddress([][]byte{[]byte("BOUNTY_SANDBLIZZARD"), owner[:]}, bountyProgramId)
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

	// convert bountyId to string
	seeds := [][]byte{[]byte("BOUNTY_SANDBLIZZARD"), []byte(strconv.Itoa(int(bountyId)))}
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
	var solverAtas []solana.PublicKey
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
func CompleteBountyAsRelayer(rpcUrl string, bountyId uint64, solverPks []solana.PublicKey, mint solana.PublicKey) (solana.Signature, error) {
	cluster := rpc.New(rpcUrl)
	ctx := context.Background()
	var signature solana.Signature
	signer, err := GetSignerKeysFromEnv()
	if err != nil {
		return signature, err
	}
	fmt.Println("Signer: ", signer.PublicKey().String())

	protocol, _, err := GetProtocolPDA()
	if err != nil {
		return signature, err
	}
	//bountyProgramId := GetBountyProgramId()

	feeCollector, _, err := GetfeeCollectorPDA(mint)
	if err != nil {
		return signature, err
	}

	bountyDenomination, _, err := GetBountyDenominationPDA(mint)
	if err != nil {
		return signature, err
	}

	bountyPk, _, err := GetBountyPDA(bountyId)
	if err != nil {
		return signature, err
	}

	escrow, _, err := GetEscrowPDA(bountyPk)
	if err != nil {
		return signature, err
	}

	relayer, _, err := GetRelayerPDA(signer.PublicKey())
	if err != nil {
		return signature, err
	}

	solvers, err := GetAndCheckSolverTokenAccounts(ctx, mint, solverPks, cluster)
	if err != nil {
		return signature, err
	}

	if len(solvers) < 1 {
		return signature, errors.Errorf("Expected at least one solver")
	}
	fmt.Printf("Solver ata: %v  \n Solver wallets %v \n ", solvers, solverPks)

	completeBounty := bounty.NewCompleteBountyAsRelayerInstruction(
		signer.PublicKey(),
		protocol,
		feeCollector,
		bountyDenomination,
		bountyPk,
		escrow,
		solvers[0],
		solvers[0],
		solvers[0],
		solvers[0],
		solana.SystemProgramID,
		solana.TokenProgramID,
		relayer,
	)

	ix, err := completeBounty.ValidateAndBuild()
	if err != nil {
		return signature, err
	}

	fmt.Printf("Accounts %v", ix.Accounts())
	data, err := ix.Data()
	if err != nil {
		return signature, errors.Wrapf(err, "failed to get data")
	}
	fmt.Printf("Data: %v", data)
	fmt.Printf("Data string %s", string(data))

	recentBlockhash, err := cluster.GetRecentBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		return signature, errors.Wrapf(err, "failed to get recent blockhash")
	}
	blockhash := recentBlockhash.Value.Blockhash
	if blockhash.IsZero() {
		return signature, errors.Errorf("blockhash is zero")
	}
	tx, err := solana.NewTransactionBuilder().AddInstruction(ix).SetFeePayer(signer.PublicKey()).SetRecentBlockHash(blockhash).Build()
	if err != nil {
		return signature, errors.Wrapf(err, "failed to create transaction")
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		return &signer
	})
	if err != nil {
		return signature, errors.Wrapf(err, "failed to sign transaction")
	}
	println("Blockhash ", tx.Message.RecentBlockhash.String())
	signature, err = cluster.SendTransaction(ctx, tx)
	if err != nil {
		return signature, errors.Wrapf(err, "failed to send transaction for bountyId %d", bountyId)
	}
	fmt.Println("Signature: ", signature)

	out, err := retryGetSignatureConfirmation(ctx, signature, cluster)
	if err != nil {
		return signature, errors.Wrapf(err, "failed to get transaction")
	}
	fmt.Println("Transaction: ", out)
	return signature, nil

}

func retryGetSignatureConfirmation(ctx context.Context, signature solana.Signature, cluster *rpc.Client) (*rpc.TransactionWithMeta, error) {
	maxRetries := 3
	retriesLeft := maxRetries
	var err error
	for retriesLeft > 0 {
		out, err := cluster.GetConfirmedTransactionWithOpts(ctx, signature,
			&rpc.GetTransactionOpts{
				Commitment: rpc.CommitmentConfirmed,
			})
		if err != nil {
			return out, nil
		}
		retriesLeft--
		time.Sleep(200 * time.Microsecond)
	}

	return nil, errors.Wrapf(err, "failed to get transaction after %d retries", maxRetries)
}
