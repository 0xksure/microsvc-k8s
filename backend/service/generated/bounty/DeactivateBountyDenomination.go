// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package bounty

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// deactivate bounty denomination
type DeactivateBountyDenomination struct {

	// [0] = [WRITE, SIGNER] creator
	//
	// [1] = [] mint
	// ··········· mint to be used for denomination
	//
	// [2] = [WRITE] denomination
	// ··········· bounty denoination to be created
	//
	// [3] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewDeactivateBountyDenominationInstructionBuilder creates a new `DeactivateBountyDenomination` instruction builder.
func NewDeactivateBountyDenominationInstructionBuilder() *DeactivateBountyDenomination {
	nd := &DeactivateBountyDenomination{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetCreatorAccount sets the "creator" account.
func (inst *DeactivateBountyDenomination) SetCreatorAccount(creator ag_solanago.PublicKey) *DeactivateBountyDenomination {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(creator).WRITE().SIGNER()
	return inst
}

// GetCreatorAccount gets the "creator" account.
func (inst *DeactivateBountyDenomination) GetCreatorAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetMintAccount sets the "mint" account.
// mint to be used for denomination
func (inst *DeactivateBountyDenomination) SetMintAccount(mint ag_solanago.PublicKey) *DeactivateBountyDenomination {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(mint)
	return inst
}

// GetMintAccount gets the "mint" account.
// mint to be used for denomination
func (inst *DeactivateBountyDenomination) GetMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetDenominationAccount sets the "denomination" account.
// bounty denoination to be created
func (inst *DeactivateBountyDenomination) SetDenominationAccount(denomination ag_solanago.PublicKey) *DeactivateBountyDenomination {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(denomination).WRITE()
	return inst
}

// GetDenominationAccount gets the "denomination" account.
// bounty denoination to be created
func (inst *DeactivateBountyDenomination) GetDenominationAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *DeactivateBountyDenomination) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *DeactivateBountyDenomination {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *DeactivateBountyDenomination) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst DeactivateBountyDenomination) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_DeactivateBountyDenomination,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst DeactivateBountyDenomination) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *DeactivateBountyDenomination) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Creator is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Mint is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Denomination is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *DeactivateBountyDenomination) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("DeactivateBountyDenomination")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("      creator", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("         mint", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta(" denomination", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj DeactivateBountyDenomination) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *DeactivateBountyDenomination) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewDeactivateBountyDenominationInstruction declares a new DeactivateBountyDenomination instruction with the provided parameters and accounts.
func NewDeactivateBountyDenominationInstruction(
	// Accounts:
	creator ag_solanago.PublicKey,
	mint ag_solanago.PublicKey,
	denomination ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *DeactivateBountyDenomination {
	return NewDeactivateBountyDenominationInstructionBuilder().
		SetCreatorAccount(creator).
		SetMintAccount(mint).
		SetDenominationAccount(denomination).
		SetSystemProgramAccount(systemProgram)
}
