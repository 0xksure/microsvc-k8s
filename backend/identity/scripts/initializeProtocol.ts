import { error } from "console"
import * as identity from "../sdk-ts/src/index"
import * as web3 from "@solana/web3.js"
import { sendAndConfirmTransaction } from "../sdk-ts/src/utils"
import * as bs58 from "bs58";


async function main() {
    // load wallet from env 
    const secretKey = process.env.WALLET_SECRET_KEY
    if (!secretKey) throw error(400, 'No wallet secret key found')

    const wallet = web3.Keypair.fromSecretKey(bs58.decode(secretKey))

    // setup connection from env rpc url
    const rpcUrl = process.env.RPC_URL
    if (!rpcUrl) throw error(400, 'No rpc url found')
    const connection = new web3.Connection(rpcUrl)
    const latestBlockhash = await connection.getLatestBlockhash();
    console.log(`Latest blockhash: ${latestBlockhash.blockhash}`)
    // create initialize identity transaction
    const identitySdk = new identity.IdentitySdk(wallet.publicKey, connection);
    const initializeIdentity = await identitySdk.initializeProtocol()
    await sendAndConfirmTransaction(connection, initializeIdentity.vtx, [wallet], latestBlockhash)
    console.log(`Identity program initialized at ${initializeIdentity.identityProgramPDA[0].toBase58()}`)
}

main().catch(err => {
    console.error(err)
    process.exit(1)
})