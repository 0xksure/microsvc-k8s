// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package bounty

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// RemoveRelayer is the `removeRelayer` instruction.
type RemoveRelayer struct {

	// [0] = [WRITE, SIGNER] signer
	//
	// [1] = [] protocol
	//
	// [2] = [WRITE] relayer
	//
	// [3] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewRemoveRelayerInstructionBuilder creates a new `RemoveRelayer` instruction builder.
func NewRemoveRelayerInstructionBuilder() *RemoveRelayer {
	nd := &RemoveRelayer{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetSignerAccount sets the "signer" account.
func (inst *RemoveRelayer) SetSignerAccount(signer ag_solanago.PublicKey) *RemoveRelayer {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(signer).WRITE().SIGNER()
	return inst
}

// GetSignerAccount gets the "signer" account.
func (inst *RemoveRelayer) GetSignerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetProtocolAccount sets the "protocol" account.
func (inst *RemoveRelayer) SetProtocolAccount(protocol ag_solanago.PublicKey) *RemoveRelayer {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(protocol)
	return inst
}

// GetProtocolAccount gets the "protocol" account.
func (inst *RemoveRelayer) GetProtocolAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetRelayerAccount sets the "relayer" account.
func (inst *RemoveRelayer) SetRelayerAccount(relayer ag_solanago.PublicKey) *RemoveRelayer {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(relayer).WRITE()
	return inst
}

// GetRelayerAccount gets the "relayer" account.
func (inst *RemoveRelayer) GetRelayerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *RemoveRelayer) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *RemoveRelayer {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *RemoveRelayer) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst RemoveRelayer) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_RemoveRelayer,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst RemoveRelayer) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *RemoveRelayer) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Signer is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Protocol is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Relayer is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *RemoveRelayer) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("RemoveRelayer")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("       signer", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("     protocol", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("      relayer", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj RemoveRelayer) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *RemoveRelayer) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewRemoveRelayerInstruction declares a new RemoveRelayer instruction with the provided parameters and accounts.
func NewRemoveRelayerInstruction(
	// Accounts:
	signer ag_solanago.PublicKey,
	protocol ag_solanago.PublicKey,
	relayer ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *RemoveRelayer {
	return NewRemoveRelayerInstructionBuilder().
		SetSignerAccount(signer).
		SetProtocolAccount(protocol).
		SetRelayerAccount(relayer).
		SetSystemProgramAccount(systemProgram)
}