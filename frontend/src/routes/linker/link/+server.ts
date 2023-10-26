// Linker/link is responsible for validating the github login 
// and making a call to the identity service to link the github account to the wallet

import { error } from '@sveltejs/kit';
import * as anchor from "@coral-xyz/anchor"
import jwt from 'jsonwebtoken';
import { Logger } from 'tslog';
import { Kafka } from 'kafkajs';
import * as proto from '$lib/index_pb';
import { Keypair, VersionedTransaction } from '@solana/web3.js';
import { bs58 } from '@coral-xyz/anchor/dist/cjs/utils/bytes/index.js';
import * as identity from "idlinker-sdk";


export interface LinkerResponse {
    username: string,
    userId: number,
    walletAddress: string
    vtx: VersionedTransaction
}

export const POST = (async (event) => {
    const logger = new Logger();

    const ghJwt = event.cookies.get('ghJwt')
    if (!ghJwt) throw error(400, 'No github jwt found')
    const jwtSecret = process.env.JWT_SECRET
    if (!jwtSecret) throw error(400, 'No jwt secret found')
    const ghJwtDecoded = jwt.verify(ghJwt, jwtSecret)
    if (!ghJwtDecoded) throw error(400, 'Invalid jwt')
    const ghAccessToken = (ghJwtDecoded as any)?.token
    logger.info(`linker.link.post: ghAccessToken=${ghAccessToken}, ghJwt=${ghJwt}, ghJwtDecoded=${JSON.stringify(ghJwtDecoded)}`);
    if (!ghAccessToken) throw error(400, 'No github access token found')


    const getUserResponse = await fetch(`https://api.github.com/user`, {
        method: 'GET',
        headers: {
            Authorization: `Bearer ${ghAccessToken}`,
            accept: 'application/json'
        }
    });
    if (!getUserResponse.ok) throw error(400, `Failed to get user information: ${getUserResponse.status}`);
    const userData = await getUserResponse.json()
    const username = userData.login
    const userId = userData.id

    const requestBody = await event.request.json()
    logger.info(`linker.link.post: body = ${JSON.stringify(requestBody)}`);
    // wallet address
    const walletAddress = requestBody.walletAddress
    if (!walletAddress) throw error(400, 'No wallet address found')
    const identityOwner = requestBody.identityOwner
    if (!identityOwner) throw error(400, 'No identity owner found')

    // load wallet from json file 
    const secretKey = process.env.WALLET_SECRET_KEY
    if (!secretKey) throw error(400, 'No wallet secret key found')
    const wallet = Keypair.fromSecretKey(bs58.decode(secretKey))

    // create identity transaction
    const identitySdk = new identity.IdentitySdk(wallet.publicKey);
    const createIdentity = await identitySdk.createIdentity({
        social: "github",
        username,
        userId,
        identityOwner: identityOwner,
        protocolOwner: wallet.publicKey,
    });
    const vtx = createIdentity.vtx
    vtx.sign([wallet])

    // post data to kafka 
    const kafkaPwd = process.env.KAFKA_PASSWORD
    if (!kafkaPwd) throw error(400, 'No kafka password found')
    const kafka = new Kafka({
        clientId: 'my-app',
        brokers: ["kafka-controller-0.kafka-controller-headless.default.svc.cluster.local:9092",
            "kafka-controller-1.kafka-controller-headless.default.svc.cluster.local:9092",
            "kafka-controller-2.kafka-controller-headless.default.svc.cluster.local:9092"],
        ssl: false,
        sasl: {
            mechanism: 'scram-sha-256',
            username: 'user1',
            password: kafkaPwd,
        }
    })

    const linkerMessage = new proto.LinkerMessage({
        Username: username,
        UserId: userId,
        WalletAddress: walletAddress
    })

    const producer = kafka.producer()
    await producer.connect()
    await producer.send({
        topic: 'linker',
        messages: [
            { value: Buffer.from(linkerMessage.toBinary()), partition: 0 },
        ],
    })

    const linkerResponse: LinkerResponse = {
        username,
        userId,
        walletAddress,
        vtx
    }
    return new Response(JSON.stringify(linkerResponse))
})