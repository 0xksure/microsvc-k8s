import { BN, web3, Program, AnchorProvider } from "@coral-xyz/anchor";
import { Identity, IDL } from "./idl/identity.js"
import bs58 from "bs58"
import NodeWallet from "@coral-xyz/anchor/dist/cjs/nodewallet.js";
export * as utils from "./utils.js"
const IDENTITY_SEED = "identity"
const IDENTITY_PROGRAM_ID = new web3.PublicKey("3Nt1tyTJ6VBf4APaPPWixXFJr6DtfGvvTwHY1aGYT4Ws")

export const getIdentityProgramPDA = () => {
    return web3.PublicKey.findProgramAddressSync(
        [Buffer.from(IDENTITY_SEED)],
        IDENTITY_PROGRAM_ID
    );
}

export const getIdentityPDA = (social: string, userId: BN) => {
    return web3.PublicKey.findProgramAddressSync(
        [Buffer.from(IDENTITY_SEED), Buffer.from(social), new BN(userId).toBuffer("le", 8)],
        IDENTITY_PROGRAM_ID
    )
}


/**
 * convertByteArrayToString takes a an array of bytes and 
 * strips the trailing zeros and converts it to a string
 * @param bytes 
 * @returns 
 */
export const convertByteArrayToString = (bytes: Uint8Array) => {
    const lastSignificantByte = bytes.findIndex((byte) => byte === 0)
    const byteSubset = bytes.slice(0, lastSignificantByte)
    return Buffer.from(byteSubset).toString()
}

export const convertStringOfSizeToString = (string: string) => {
    var enc = new TextEncoder();
    return convertByteArrayToString(enc.encode(string))
}


export class IdentitySdk {

    public program: Program<Identity>;
    constructor(
        readonly signer: web3.PublicKey,
        readonly connection?: web3.Connection,
    ) {
        this.program = IdentitySdk.setUpProgram({
            keypair: web3.Keypair.generate(),
            connection: connection
        });
    }

    private static setUpProgram({
        keypair,
        connection
    }:
        {
            keypair: web3.Keypair,
            connection?: web3.Connection
        }) {
        const provider = new AnchorProvider(connection ?? new web3.Connection("https://api.solana.com"), new NodeWallet(keypair), {
            preflightCommitment: "recent",
            commitment: "confirmed",
        })
        return new Program<Identity>(IDL, IDENTITY_PROGRAM_ID, provider);
    }

    /**
     * createVersionedTransaction takes a list of instructions and creates a versioned transaction
     *
     * @param ixs: instructions
     * @returns
     */
    createVersionedTransaction = async (
        ixs: web3.TransactionInstruction[], payer: web3.PublicKey = this.signer
    ) => {
        const txMessage = await new web3.TransactionMessage({
            instructions: ixs,
            recentBlockhash: (
                await this.program.provider.connection.getLatestBlockhash()
            ).blockhash,
            payerKey: payer,
        }).compileToV0Message();

        return new web3.VersionedTransaction(txMessage);
    };


    /**
     * getIdentityAccount uses the standard fetch method to get the account data
     * It also deals with the last mile deserialization and standardization of the data
     * @param param0 
     * @returns 
     */
    getIdentityAccount = async ({
        social,
        userId
    }: {
        social: string,
        userId: BN
    }) => {
        const identityPDA = getIdentityPDA(social, userId)
        const identity = await this.program.account.identity.fetch(identityPDA[0])
        identity.username = convertStringOfSizeToString(identity.username)
        identity.social = convertStringOfSizeToString(identity.social)
        return identity
    }

    getIdentityAccountFromAddress = async ({
        address,
    }: {
        address: web3.PublicKey
    }) => {
        const identity = await this.program.account.identity.fetch(address)
        identity.username = convertStringOfSizeToString(identity.username)
        identity.social = convertStringOfSizeToString(identity.social)
        return identity
    }

    /**
     * getIdentityProgramAccount is a simple helper method
     * @returns 
     */
    getIdentityProgramAccount = async () => {
        const identityProgramPDA = getIdentityProgramPDA()
        return this.program.account.identityProgram.fetch(identityProgramPDA[0])
    }

    initializeProtocol = async ({
        signer
    }: {
        signer?: web3.PublicKey
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
        userId: BN,
        username: string,
        identityOwner: web3.PublicKey,
        protocolOwner?: web3.PublicKey
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
        userId: BN,
        identityOwner: web3.PublicKey
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
        identityOwner: web3.PublicKey,
        newIdentityOwner: web3.PublicKey,
        social: string,
        userId: BN
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
        identityOwner: web3.PublicKey,
        social: string,
        userId: BN
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

    /**
     * getIdentityFromAddress allows you to get all identities associated with an address
     * @param address 
     * @returns 
     */
    getIdentityFromAddress = async ({ address }: { address: web3.PublicKey }) => {
        const memcmpFilters = [
            {
                memcmp: {
                    offset: 8,
                    bytes: address.toBase58()
                }
            }
        ]
        const identityAccounts = await this.program.account.identity.all(memcmpFilters)
        return identityAccounts.map((identity) => {
            identity.account.username = convertStringOfSizeToString(identity.account.username)
            identity.account.social = convertStringOfSizeToString(identity.account.social)
            return identity
        })
    }

    /**
     * getIdentityFromUsername allows you to get all identities associated with a username
     * @param param0 
     * @returns 
     */
    getIdentityFromUsername = async ({
        username
    }: { username: string }) => {
        const memcmpFilters = [
            {
                memcmp: {
                    offset: 8 + 32 + 4 + 32 + 8 + 4,
                    bytes: bs58.encode(Buffer.from(username))
                }
            }
        ]

        const identityAccount = await this.program.account.identity.all(memcmpFilters)
        return identityAccount.map((identity) => {
            identity.account.username = convertStringOfSizeToString(identity.account.username)
            identity.account.social = convertStringOfSizeToString(identity.account.social)
            return identity
        })
    }

}

