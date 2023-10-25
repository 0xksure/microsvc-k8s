"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.sendAndConfirmTransaction = void 0;
/**
 * sendAndConfirmTransaction is a simple wrapper around web3's sendAndConfirmTransaction
 * @param connection
 * @param transaction
 * @param latestBlockhash
 * @param signers
 * @returns
 */
const sendAndConfirmTransaction = (connection, transaction, latestBlockhash, ...signers) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        if (!latestBlockhash) {
            latestBlockhash = yield connection.getLatestBlockhash();
        }
        if (signers.length !== 0) {
            transaction.sign(signers);
        }
        const signature = yield connection.sendTransaction(transaction, {
            skipPreflight: false,
        });
        const confirmation = yield connection.confirmTransaction(Object.assign({ signature: signature }, latestBlockhash));
        return confirmation;
    }
    catch (err) {
        console.log("err", err);
        throw err;
    }
});
exports.sendAndConfirmTransaction = sendAndConfirmTransaction;
//# sourceMappingURL=utils.js.map