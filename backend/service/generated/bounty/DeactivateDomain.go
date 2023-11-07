// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package bounty

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// deactivate domain
type DeactivateDomain struct {

	// [0] = [SIGNER] signer
	//
	// [1] = [WRITE] domain
	//
	// [2] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewDeactivateDomainInstructionBuilder creates a new `DeactivateDomain` instruction builder.
func NewDeactivateDomainInstructionBuilder() *DeactivateDomain {
	nd := &DeactivateDomain{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// SetSignerAccount sets the "signer" account.
func (inst *DeactivateDomain) SetSignerAccount(signer ag_solanago.PublicKey) *DeactivateDomain {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(signer).SIGNER()
	return inst
}

// GetSignerAccount gets the "signer" account.
func (inst *DeactivateDomain) GetSignerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetDomainAccount sets the "domain" account.
func (inst *DeactivateDomain) SetDomainAccount(domain ag_solanago.PublicKey) *DeactivateDomain {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(domain).WRITE()
	return inst
}

// GetDomainAccount gets the "domain" account.
func (inst *DeactivateDomain) GetDomainAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *DeactivateDomain) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *DeactivateDomain {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *DeactivateDomain) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

func (inst DeactivateDomain) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_DeactivateDomain,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst DeactivateDomain) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *DeactivateDomain) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Signer is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Domain is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *DeactivateDomain) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("DeactivateDomain")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=3]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("       signer", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("       domain", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(2)))
					})
				})
		})
}

func (obj DeactivateDomain) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *DeactivateDomain) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewDeactivateDomainInstruction declares a new DeactivateDomain instruction with the provided parameters and accounts.
func NewDeactivateDomainInstruction(
	// Accounts:
	signer ag_solanago.PublicKey,
	domain ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *DeactivateDomain {
	return NewDeactivateDomainInstructionBuilder().
		SetSignerAccount(signer).
		SetDomainAccount(domain).
		SetSystemProgramAccount(systemProgram)
}
