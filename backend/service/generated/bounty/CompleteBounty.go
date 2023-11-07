// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package bounty

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// complete_bounty
//
// Try to complete bounty
type CompleteBounty struct {

	// [0] = [WRITE, SIGNER] payer
	// ··········· only owners or relayers can complete bounties
	//
	// [1] = [] protocol
	//
	// [2] = [WRITE] feeCollector
	//
	// [3] = [] bountyDenomination
	// ··········· bounty denomination is the allowed denomination of a bounty
	// ··········· it needs to be checked against the fee collector and the mint
	//
	// [4] = [WRITE] bounty
	// ··········· bounty to be completed
	// ··········· FIXME
	//
	// [5] = [WRITE] escrow
	//
	// [6] = [WRITE] solver1
	// ··········· up to 4 receivers
	//
	// [7] = [WRITE] solver2
	//
	// [8] = [WRITE] solver3
	//
	// [9] = [WRITE] solver4
	//
	// [10] = [] systemProgram
	//
	// [11] = [] tokenProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewCompleteBountyInstructionBuilder creates a new `CompleteBounty` instruction builder.
func NewCompleteBountyInstructionBuilder() *CompleteBounty {
	nd := &CompleteBounty{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 12),
	}
	return nd
}

// SetPayerAccount sets the "payer" account.
// only owners or relayers can complete bounties
func (inst *CompleteBounty) SetPayerAccount(payer ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(payer).WRITE().SIGNER()
	return inst
}

// GetPayerAccount gets the "payer" account.
// only owners or relayers can complete bounties
func (inst *CompleteBounty) GetPayerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetProtocolAccount sets the "protocol" account.
func (inst *CompleteBounty) SetProtocolAccount(protocol ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(protocol)
	return inst
}

// GetProtocolAccount gets the "protocol" account.
func (inst *CompleteBounty) GetProtocolAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetFeeCollectorAccount sets the "feeCollector" account.
func (inst *CompleteBounty) SetFeeCollectorAccount(feeCollector ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(feeCollector).WRITE()
	return inst
}

// GetFeeCollectorAccount gets the "feeCollector" account.
func (inst *CompleteBounty) GetFeeCollectorAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetBountyDenominationAccount sets the "bountyDenomination" account.
// bounty denomination is the allowed denomination of a bounty
// it needs to be checked against the fee collector and the mint
func (inst *CompleteBounty) SetBountyDenominationAccount(bountyDenomination ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(bountyDenomination)
	return inst
}

// GetBountyDenominationAccount gets the "bountyDenomination" account.
// bounty denomination is the allowed denomination of a bounty
// it needs to be checked against the fee collector and the mint
func (inst *CompleteBounty) GetBountyDenominationAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetBountyAccount sets the "bounty" account.
// bounty to be completed
// FIXME
func (inst *CompleteBounty) SetBountyAccount(bounty ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(bounty).WRITE()
	return inst
}

// GetBountyAccount gets the "bounty" account.
// bounty to be completed
// FIXME
func (inst *CompleteBounty) GetBountyAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetEscrowAccount sets the "escrow" account.
func (inst *CompleteBounty) SetEscrowAccount(escrow ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(escrow).WRITE()
	return inst
}

// GetEscrowAccount gets the "escrow" account.
func (inst *CompleteBounty) GetEscrowAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetSolver1Account sets the "solver1" account.
// up to 4 receivers
func (inst *CompleteBounty) SetSolver1Account(solver1 ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(solver1).WRITE()
	return inst
}

// GetSolver1Account gets the "solver1" account.
// up to 4 receivers
func (inst *CompleteBounty) GetSolver1Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetSolver2Account sets the "solver2" account.
func (inst *CompleteBounty) SetSolver2Account(solver2 ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(solver2).WRITE()
	return inst
}

// GetSolver2Account gets the "solver2" account.
func (inst *CompleteBounty) GetSolver2Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetSolver3Account sets the "solver3" account.
func (inst *CompleteBounty) SetSolver3Account(solver3 ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(solver3).WRITE()
	return inst
}

// GetSolver3Account gets the "solver3" account.
func (inst *CompleteBounty) GetSolver3Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetSolver4Account sets the "solver4" account.
func (inst *CompleteBounty) SetSolver4Account(solver4 ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(solver4).WRITE()
	return inst
}

// GetSolver4Account gets the "solver4" account.
func (inst *CompleteBounty) GetSolver4Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *CompleteBounty) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *CompleteBounty) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *CompleteBounty) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *CompleteBounty {
	inst.AccountMetaSlice[11] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *CompleteBounty) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(11)
}

func (inst CompleteBounty) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_CompleteBounty,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst CompleteBounty) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *CompleteBounty) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Payer is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Protocol is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.FeeCollector is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.BountyDenomination is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.Bounty is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.Escrow is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.Solver1 is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.Solver2 is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.Solver3 is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.Solver4 is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
		if inst.AccountMetaSlice[11] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
	}
	return nil
}

func (inst *CompleteBounty) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("CompleteBounty")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=12]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("             payer", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("          protocol", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("      feeCollector", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("bountyDenomination", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("            bounty", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("            escrow", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("           solver1", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("           solver2", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("           solver3", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("           solver4", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("     systemProgram", inst.AccountMetaSlice.Get(10)))
						accountsBranch.Child(ag_format.Meta("      tokenProgram", inst.AccountMetaSlice.Get(11)))
					})
				})
		})
}

func (obj CompleteBounty) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *CompleteBounty) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewCompleteBountyInstruction declares a new CompleteBounty instruction with the provided parameters and accounts.
func NewCompleteBountyInstruction(
	// Accounts:
	payer ag_solanago.PublicKey,
	protocol ag_solanago.PublicKey,
	feeCollector ag_solanago.PublicKey,
	bountyDenomination ag_solanago.PublicKey,
	bounty ag_solanago.PublicKey,
	escrow ag_solanago.PublicKey,
	solver1 ag_solanago.PublicKey,
	solver2 ag_solanago.PublicKey,
	solver3 ag_solanago.PublicKey,
	solver4 ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey) *CompleteBounty {
	return NewCompleteBountyInstructionBuilder().
		SetPayerAccount(payer).
		SetProtocolAccount(protocol).
		SetFeeCollectorAccount(feeCollector).
		SetBountyDenominationAccount(bountyDenomination).
		SetBountyAccount(bounty).
		SetEscrowAccount(escrow).
		SetSolver1Account(solver1).
		SetSolver2Account(solver2).
		SetSolver3Account(solver3).
		SetSolver4Account(solver4).
		SetSystemProgramAccount(systemProgram).
		SetTokenProgramAccount(tokenProgram)
}
