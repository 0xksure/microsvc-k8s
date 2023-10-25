import * as anchor from "@coral-xyz/anchor";
import { Identity } from "./idl/identity";
export * as utils from "./utils";
export { IdentitySdk };
declare class IdentitySdk {
    readonly signer: anchor.web3.PublicKey;
    readonly connection?: anchor.web3.Connection;
    program: anchor.Program<Identity>;
    constructor(signer: anchor.web3.PublicKey, connection?: anchor.web3.Connection);
    private static setUpProgram;
    /**
     * createVersionedTransaction takes a list of instructions and creates a versioned transaction
     *
     * @param ixs: instructions
     * @returns
     */
    createVersionedTransaction: (ixs: anchor.web3.TransactionInstruction[], payer?: anchor.web3.PublicKey) => Promise<anchor.web3.VersionedTransaction>;
    initializeProtocol: ({ signer }?: {
        signer?: anchor.web3.PublicKey;
    }) => Promise<{
        vtx: anchor.web3.VersionedTransaction;
        ix: anchor.web3.TransactionInstruction;
        identityProgramPDA: [anchor.web3.PublicKey, number];
    }>;
    createIdentity: ({ social, userId, username, identityOwner, protocolOwner }: {
        social: string;
        userId: number;
        username: string;
        identityOwner: anchor.web3.PublicKey;
        protocolOwner?: anchor.web3.PublicKey;
    }) => Promise<{
        vtx: anchor.web3.VersionedTransaction;
        ix: anchor.web3.TransactionInstruction;
        identityPDA: [anchor.web3.PublicKey, number];
    }>;
    updateUsername: ({ username, identityOwner, social, userId }: {
        username: string;
        social: string;
        userId: number;
        identityOwner: anchor.web3.PublicKey;
    }) => Promise<{
        vtx: anchor.web3.VersionedTransaction;
        ix: anchor.web3.TransactionInstruction;
        identityPDA: [anchor.web3.PublicKey, number];
    }>;
    transferOwnership: ({ identityOwner, newIdentityOwner, social, userId }: {
        identityOwner: anchor.web3.PublicKey;
        newIdentityOwner: anchor.web3.PublicKey;
        social: string;
        userId: number;
    }) => Promise<{
        vtx: anchor.web3.VersionedTransaction;
        ix: anchor.web3.TransactionInstruction;
        identityPDA: [anchor.web3.PublicKey, number];
    }>;
    deleteIdentity: ({ identityOwner, social, userId }: {
        identityOwner: anchor.web3.PublicKey;
        social: string;
        userId: number;
    }) => Promise<{
        vtx: anchor.web3.VersionedTransaction;
        ix: anchor.web3.TransactionInstruction;
        identityPDA: [anchor.web3.PublicKey, number];
    }>;
}
//# sourceMappingURL=index.d.ts.map