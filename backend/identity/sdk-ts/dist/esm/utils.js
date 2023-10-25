/**
 * sendAndConfirmTransaction is a simple wrapper around web3's sendAndConfirmTransaction
 * @param connection
 * @param transaction
 * @param latestBlockhash
 * @param signers
 * @returns
 */
export const sendAndConfirmTransaction = async (connection, transaction, latestBlockhash, ...signers) => {
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
        return confirmation;
    }
    catch (err) {
        console.log("err", err);
        throw err;
    }
};
//# sourceMappingURL=utils.js.map