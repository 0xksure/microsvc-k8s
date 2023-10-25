"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || function (mod) {
    if (mod && mod.__esModule) return mod;
    var result = {};
    if (mod != null) for (var k in mod) if (k !== "default" && Object.prototype.hasOwnProperty.call(mod, k)) __createBinding(result, mod, k);
    __setModuleDefault(result, mod);
    return result;
};
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
exports.IdentitySdk = exports.utils = void 0;
const anchor = __importStar(require("@coral-xyz/anchor"));
const identity_1 = require("./idl/identity");
exports.utils = __importStar(require("./utils"));
const IDENTITY_SEED = "identity";
const IDENTITY_PROGRAM_ID = new anchor.web3.PublicKey("3rQketG7pSopHE1APQKZu1BQofanqbCBP7spZ4CBGrUm");
const getIdentityProgramPDA = () => {
    return anchor.web3.PublicKey.findProgramAddressSync([Buffer.from(IDENTITY_SEED)], IDENTITY_PROGRAM_ID);
};
const getIdentityPDA = (social, userId) => {
    return anchor.web3.PublicKey.findProgramAddressSync([Buffer.from(IDENTITY_SEED), Buffer.from(social), Buffer.from(userId.toString())], IDENTITY_PROGRAM_ID);
};
class IdentitySdk {
    constructor(signer, connection) {
        this.signer = signer;
        this.connection = connection;
        /**
         * createVersionedTransaction takes a list of instructions and creates a versioned transaction
         *
         * @param ixs: instructions
         * @returns
         */
        this.createVersionedTransaction = (ixs, payer = this.signer) => __awaiter(this, void 0, void 0, function* () {
            const txMessage = yield new anchor.web3.TransactionMessage({
                instructions: ixs,
                recentBlockhash: (yield this.program.provider.connection.getLatestBlockhash()).blockhash,
                payerKey: payer,
            }).compileToV0Message();
            return new anchor.web3.VersionedTransaction(txMessage);
        });
        this.initializeProtocol = ({ signer } = {}) => __awaiter(this, void 0, void 0, function* () {
            const identityProgramPDA = getIdentityProgramPDA();
            const ix = yield this.program.methods.initialize().accounts({
                protocolOwner: signer !== null && signer !== void 0 ? signer : this.signer,
                identityProgram: identityProgramPDA[0]
            }).instruction();
            return {
                vtx: yield this.createVersionedTransaction([ix], signer !== null && signer !== void 0 ? signer : this.signer),
                ix,
                identityProgramPDA
            };
        });
        this.createIdentity = ({ social, userId, username, identityOwner, protocolOwner }) => __awaiter(this, void 0, void 0, function* () {
            const identityPDA = getIdentityPDA(social, userId);
            const protocolPDA = getIdentityProgramPDA();
            const ix = yield this.program.methods.createIdentity(social, username, userId).accounts({
                accountHolder: identityOwner,
                protocolOwner: protocolOwner !== null && protocolOwner !== void 0 ? protocolOwner : this.signer,
                identityProgram: protocolPDA[0],
                identity: identityPDA[0],
            }).instruction();
            return {
                vtx: yield this.createVersionedTransaction([ix], identityOwner),
                ix,
                identityPDA
            };
        });
        this.updateUsername = ({ username, identityOwner, social, userId }) => __awaiter(this, void 0, void 0, function* () {
            const identityPDA = getIdentityPDA(social, userId);
            const ix = yield this.program.methods.updateUsername(username).accounts({
                accountHolder: identityOwner,
                identity: identityPDA[0],
            }).instruction();
            const signers = [identityOwner];
            return {
                vtx: yield this.createVersionedTransaction([ix], ...signers),
                ix,
                identityPDA
            };
        });
        this.transferOwnership = ({ identityOwner, newIdentityOwner, social, userId }) => __awaiter(this, void 0, void 0, function* () {
            const identityPDA = getIdentityPDA(social, userId);
            const ix = yield this.program.methods.transferOwnership().accounts({
                accountHolderCurr: identityOwner,
                accountHolderNew: newIdentityOwner,
                identity: identityPDA[0],
            }).instruction();
            const signers = [identityOwner, newIdentityOwner];
            return {
                vtx: yield this.createVersionedTransaction([ix], ...signers),
                ix,
                identityPDA
            };
        });
        this.deleteIdentity = ({ identityOwner, social, userId }) => __awaiter(this, void 0, void 0, function* () {
            const identityPDA = getIdentityPDA(social, userId);
            const ix = yield this.program.methods.deleteIdentity().accounts({
                accountHolder: identityOwner,
                identity: identityPDA[0],
            }).instruction();
            const signers = [identityOwner];
            return {
                vtx: yield this.createVersionedTransaction([ix], ...signers),
                ix,
                identityPDA
            };
        });
        this.program = IdentitySdk.setUpProgram({
            keypair: anchor.web3.Keypair.generate(),
            connection: connection
        });
    }
    static setUpProgram({ keypair, connection }) {
        const provider = new anchor.AnchorProvider(connection !== null && connection !== void 0 ? connection : new anchor.web3.Connection("https://api.solana.com"), new anchor.Wallet(keypair), {
            preflightCommitment: "recent",
            commitment: "confirmed",
        });
        return new anchor.Program(identity_1.IDL, IDENTITY_PROGRAM_ID, provider);
    }
}
exports.IdentitySdk = IdentitySdk;
//# sourceMappingURL=index.js.map