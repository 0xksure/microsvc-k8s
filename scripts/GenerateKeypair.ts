import { Keypair } from "@solana/web3.js";
import bs58 from "bs58";

function main() {
    const neeKeypair = Keypair.generate();
    const base58PK = bs58.encode(neeKeypair.secretKey);
    console.log("\n Private key (base58): ", base58PK, "\n");
    console.log("Public key (base58): ", neeKeypair.publicKey.toBase58());

}

main()