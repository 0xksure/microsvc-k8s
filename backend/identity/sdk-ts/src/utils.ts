import { web3 } from "@coral-xyz/anchor";

/**
 * sendAndConfirmTransaction is a simple wrapper around web3's sendAndConfirmTransaction
 * @param connection 
 * @param transaction 
 * @param latestBlockhash 
 * @param signers 
 * @returns 
 */
export const sendAndConfirmTransaction = async (
    connection: web3.Connection,
    transaction: web3.VersionedTransaction,
    signers: web3.Signer[],
    latestBlockhash?: {
        blockhash: string;
        lastValidBlockHeight: number;
    },
) => {
    try {
        if (!latestBlockhash) {
            latestBlockhash = await connection.getLatestBlockhash();
        }
        if (signers.length !== 0) {
            transaction.sign(signers);
        }
        const signature = await connection.sendTransaction(transaction, {
            skipPreflight: false,
        });
        const confirmation = await connection.confirmTransaction({
            signature: signature,
            ...latestBlockhash
        });
        return confirmation
    } catch (err) {
        console.log("err", err)
        throw err
    }
}