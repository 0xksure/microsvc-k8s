import * as anchor from "@coral-xyz/anchor";
import { Identity, IDL } from "./idl/identity"
import { encode } from "@coral-xyz/anchor/dist/cjs/utils/bytes/bs58";
export * as utils from "./utils"
export { IdentitySdk }

const IDENTITY_SEED = "identity"
const IDENTITY_PROGRAM_ID = new anchor.web3.PublicKey("3rQketG7pSopHE1APQKZu1BQofanqbCBP7spZ4CBGrUm")

const getIdentityProgramPDA = () => {
    return anchor.web3.PublicKey.findProgramAddressSync(
        [Buffer.from(IDENTITY_SEED)],
        IDENTITY_PROGRAM_ID
    );
}

const getIdentityPDA = (social: string, userId: number) => {
    return anchor.web3.PublicKey.findProgramAddressSync(
        [Buffer.from(IDENTITY_SEED), Buffer.from(social), new anchor.BN(userId).toBuffer("le", 4)],
        IDENTITY_PROGRAM_ID
    )
}



class IdentitySdk {

    public program: anchor.Program<Identity>;
    constructor(
        readonly signer: anchor.web3.PublicKey,
        readonly connection?: anchor.web3.Connection,
    ) {
        this.program = IdentitySdk.setUpProgram({
            keypair: anchor.web3.Keypair.generate(),
            connection: connection
        });
    }

    private static setUpProgram({
        keypair,
        connection
    }:
        {
            keypair: anchor.web3.Keypair,
            connection?: anchor.web3.Connection
        }) {
        const provider = new anchor.AnchorProvider(connection ?? new anchor.web3.Connection("https://api.solana.com"), new anchor.Wallet(keypair), {
            preflightCommitment: "recent",
            commitment: "confirmed",
        })
        return new anchor.Program<Identity>(IDL, IDENTITY_PROGRAM_ID, provider);
    }

    /**
     * createVersionedTransaction takes a list of instructions and creates a versioned transaction
     *
     * @param ixs: instructions
     * @returns
     */
    createVersionedTransaction = async (
        ixs: anchor.web3.TransactionInstruction[], payer: anchor.web3.PublicKey = this.signer
    ) => {
        const txMessage = await new anchor.web3.TransactionMessage({
            instructions: ixs,
            recentBlockhash: (
                await this.program.provider.connection.getLatestBlockhash()
            ).blockhash,
            payerKey: payer,
        }).compileToV0Message();

        return new anchor.web3.VersionedTransaction(txMessage);
    };

    initializeProtocol = async ({
        signer
    }: {
        signer?: anchor.web3.PublicKey
    } = {
        }) => {
        const identityProgramPDA = getIdentityProgramPDA()
        const ix = await this.program.methods.initialize().accounts({
            protocolOwner: signer ?? this.signer,
            identityProgram: identityProgramPDA[0]
        }).instruction()

        return {
            vtx: await this.createVersionedTransaction([ix], signer ?? this.signer),
            ix,
            identityProgramPDA
        }
    }

    createIdentity = async ({
        social,
        userId,
        username,
        identityOwner,
        protocolOwner
    }: {
        social: string,
        userId: number,
        username: string,
        identityOwner: anchor.web3.PublicKey,
        protocolOwner?: anchor.web3.PublicKey
    }) => {
        const identityPDA = getIdentityPDA(social, userId)
        const protocolPDA = getIdentityProgramPDA()
        const ix = await this.program.methods.createIdentity(social, username, userId).accounts({
            accountHolder: identityOwner,
            protocolOwner: protocolOwner ?? this.signer,
            identityProgram: protocolPDA[0],
            identity: identityPDA[0],
        }).instruction()

        return {
            vtx: await this.createVersionedTransaction([ix], identityOwner),
            ix,
            identityPDA
        }
    }

    updateUsername = async ({
        username,
        identityOwner,
        social,
        userId
    }: {
        username: string,
        social: string,
        userId: number,
        identityOwner: anchor.web3.PublicKey
    }) => {
        const identityPDA = getIdentityPDA(social, userId)
        const ix = await this.program.methods.updateUsername(username).accounts({
            accountHolder: identityOwner,
            identity: identityPDA[0],
        }).instruction()

        const signers = [identityOwner]
        return {
            vtx: await this.createVersionedTransaction([ix], ...signers),
            ix,
            identityPDA
        }
    }

    transferOwnership = async ({
        identityOwner,
        newIdentityOwner,
        social,
        userId
    }: {
        identityOwner: anchor.web3.PublicKey,
        newIdentityOwner: anchor.web3.PublicKey,
        social: string,
        userId: number
    }) => {
        const identityPDA = getIdentityPDA(social, userId)
        const ix = await this.program.methods.transferOwnership().accounts({
            accountHolderCurr: identityOwner,
            accountHolderNew: newIdentityOwner,
            identity: identityPDA[0],
        }).instruction()

        const signers = [identityOwner, newIdentityOwner]
        return {
            vtx: await this.createVersionedTransaction([ix], ...signers),
            ix,
            identityPDA
        }
    }

    deleteIdentity = async ({
        identityOwner,
        social,
        userId
    }: {
        identityOwner: anchor.web3.PublicKey,
        social: string,
        userId: number
    }) => {
        const identityPDA = getIdentityPDA(social, userId)
        const ix = await this.program.methods.deleteIdentity().accounts({
            accountHolder: identityOwner,
            identity: identityPDA[0],
        }).instruction()

        const signers = [identityOwner]
        return {
            vtx: await this.createVersionedTransaction([ix], ...signers),
            ix,
            identityPDA
        }
    }


}