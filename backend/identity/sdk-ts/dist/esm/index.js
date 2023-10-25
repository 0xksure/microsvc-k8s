import * as anchor from "@coral-xyz/anchor";
import { IDL } from "./idl/identity";
export * as utils from "./utils";
export { IdentitySdk };
const IDENTITY_SEED = "identity";
const IDENTITY_PROGRAM_ID = new anchor.web3.PublicKey("3rQketG7pSopHE1APQKZu1BQofanqbCBP7spZ4CBGrUm");
const getIdentityProgramPDA = () => {
    return anchor.web3.PublicKey.findProgramAddressSync([Buffer.from(IDENTITY_SEED)], IDENTITY_PROGRAM_ID);
};
const getIdentityPDA = (social, userId) => {
    return anchor.web3.PublicKey.findProgramAddressSync([Buffer.from(IDENTITY_SEED), Buffer.from(social), Buffer.from(userId.toString())], IDENTITY_PROGRAM_ID);
};
class IdentitySdk {
    signer;
    connection;
    program;
    constructor(signer, connection) {
        this.signer = signer;
        this.connection = connection;
        this.program = IdentitySdk.setUpProgram({
            keypair: anchor.web3.Keypair.generate(),
            connection: connection
        });
    }
    static setUpProgram({ keypair, connection }) {
        const provider = new anchor.AnchorProvider(connection ?? new anchor.web3.Connection("https://api.solana.com"), new anchor.Wallet(keypair), {
            preflightCommitment: "recent",
            commitment: "confirmed",
        });
        return new anchor.Program(IDL, IDENTITY_PROGRAM_ID, provider);
    }
    /**
     * createVersionedTransaction takes a list of instructions and creates a versioned transaction
     *
     * @param ixs: instructions
     * @returns
     */
    createVersionedTransaction = async (ixs, payer = this.signer) => {
        const txMessage = await new anchor.web3.TransactionMessage({
            instructions: ixs,
            recentBlockhash: (await this.program.provider.connection.getLatestBlockhash()).blockhash,
            payerKey: payer,
        }).compileToV0Message();
        return new anchor.web3.VersionedTransaction(txMessage);
    };
    initializeProtocol = async ({ signer } = {}) => {
        const identityProgramPDA = getIdentityProgramPDA();
        const ix = await this.program.methods.initialize().accounts({
            protocolOwner: signer ?? this.signer,
            identityProgram: identityProgramPDA[0]
        }).instruction();
        return {
            vtx: await this.createVersionedTransaction([ix], signer ?? this.signer),
            ix,
            identityProgramPDA
        };
    };
    createIdentity = async ({ social, userId, username, identityOwner, protocolOwner }) => {
        const identityPDA = getIdentityPDA(social, userId);
        const protocolPDA = getIdentityProgramPDA();
        const ix = await this.program.methods.createIdentity(social, username, userId).accounts({
            accountHolder: identityOwner,
            protocolOwner: protocolOwner ?? this.signer,
            identityProgram: protocolPDA[0],
            identity: identityPDA[0],
        }).instruction();
        return {
            vtx: await this.createVersionedTransaction([ix], identityOwner),
            ix,
            identityPDA
        };
    };
    updateUsername = async ({ username, identityOwner, social, userId }) => {
        const identityPDA = getIdentityPDA(social, userId);
        const ix = await this.program.methods.updateUsername(username).accounts({
            accountHolder: identityOwner,
            identity: identityPDA[0],
        }).instruction();
        const signers = [identityOwner];
        return {
            vtx: await this.createVersionedTransaction([ix], ...signers),
            ix,
            identityPDA
        };
    };
    transferOwnership = async ({ identityOwner, newIdentityOwner, social, userId }) => {
        const identityPDA = getIdentityPDA(social, userId);
        const ix = await this.program.methods.transferOwnership().accounts({
            accountHolderCurr: identityOwner,
            accountHolderNew: newIdentityOwner,
            identity: identityPDA[0],
        }).instruction();
        const signers = [identityOwner, newIdentityOwner];
        return {
            vtx: await this.createVersionedTransaction([ix], ...signers),
            ix,
            identityPDA
        };
    };
    deleteIdentity = async ({ identityOwner, social, userId }) => {
        const identityPDA = getIdentityPDA(social, userId);
        const ix = await this.program.methods.deleteIdentity().accounts({
            accountHolder: identityOwner,
            identity: identityPDA[0],
        }).instruction();
        const signers = [identityOwner];
        return {
            vtx: await this.createVersionedTransaction([ix], ...signers),
            ix,
            identityPDA
        };
    };
}
//# sourceMappingURL=index.js.map