import { web3 } from "@coral-xyz/anchor";
/**
 * sendAndConfirmTransaction is a simple wrapper around web3's sendAndConfirmTransaction
 * @param connection
 * @param transaction
 * @param latestBlockhash
 * @param signers
 * @returns
 */
export declare const sendAndConfirmTransaction: (connection: web3.Connection, transaction: web3.VersionedTransaction, latestBlockhash?: {
    blockhash: string;
    lastValidBlockHeight: number;
}, ...signers: web3.Signer[]) => Promise<web3.RpcResponseAndContext<web3.SignatureResult>>;
//# sourceMappingURL=utils.d.ts.map