

import { error } from "console"
import * as identity from "../sdk-ts/dist/cjs/index"
import * as web3 from "@solana/web3.js"
import { sendAndConfirmTransaction } from "../sdk-ts/dist/cjs/utils"
import * as bs58 from "bs58";
import { BN } from "bn.js";


async function main() {
    // load wallet from env 
    const userId = new BN(47750504)
    const social = "github"

    const identityPDA = await identity.getIdentityPDA(social, userId)
    console.log(`Identity PDA: ${identityPDA[0].toBase58()}`)

}


main().catch(err => {
    console.error(err)
    process.exit(1)
})