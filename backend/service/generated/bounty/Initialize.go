// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package bounty

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// initialize
//
// - Initializes the protocol
// - creates the bounty mint
type Initialize struct {

	// [0] = [WRITE, SIGNER] protocolOwner
	// ··········· creator is the owner of the protocol
	// ··········· should become a smart wallet over time
	//
	// [1] = [WRITE] protocol
	//
	// [2] = [WRITE] sandMint
	// ··········· mint to be used to distribute rewards
	//
	// [3] = [] tokenProgram
	//
	// [4] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewInitializeInstructionBuilder creates a new `Initialize` instruction builder.
func NewInitializeInstructionBuilder() *Initialize {
	nd := &Initialize{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 5),
	}
	return nd
}

// SetProtocolOwnerAccount sets the "protocolOwner" account.
// creator is the owner of the protocol
// should become a smart wallet over time
func (inst *Initialize) SetProtocolOwnerAccount(protocolOwner ag_solanago.PublicKey) *Initialize {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(protocolOwner).WRITE().SIGNER()
	return inst
}

// GetProtocolOwnerAccount gets the "protocolOwner" account.
// creator is the owner of the protocol
// should become a smart wallet over time
func (inst *Initialize) GetProtocolOwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetProtocolAccount sets the "protocol" account.
func (inst *Initialize) SetProtocolAccount(protocol ag_solanago.PublicKey) *Initialize {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(protocol).WRITE()
	return inst
}

// GetProtocolAccount gets the "protocol" account.
func (inst *Initialize) GetProtocolAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetSandMintAccount sets the "sandMint" account.
// mint to be used to distribute rewards
func (inst *Initialize) SetSandMintAccount(sandMint ag_solanago.PublicKey) *Initialize {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(sandMint).WRITE()
	return inst
}

// GetSandMintAccount gets the "sandMint" account.
// mint to be used to distribute rewards
func (inst *Initialize) GetSandMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *Initialize) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *Initialize {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *Initialize) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *Initialize) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *Initialize {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *Initialize) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

func (inst Initialize) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_Initialize,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Initialize) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Initialize) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.ProtocolOwner is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Protocol is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.SandMint is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *Initialize) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Initialize")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=5]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("protocolOwner", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("     protocol", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("     sandMint", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta(" tokenProgram", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(4)))
					})
				})
		})
}

func (obj Initialize) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *Initialize) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewInitializeInstruction declares a new Initialize instruction with the provided parameters and accounts.
func NewInitializeInstruction(
	// Accounts:
	protocolOwner ag_solanago.PublicKey,
	protocol ag_solanago.PublicKey,
	sandMint ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *Initialize {
	return NewInitializeInstructionBuilder().
		SetProtocolOwnerAccount(protocolOwner).
		SetProtocolAccount(protocol).
		SetSandMintAccount(sandMint).
		SetTokenProgramAccount(tokenProgram).
		SetSystemProgramAccount(systemProgram)
}
