import { error } from '@sveltejs/kit';
import { Logger } from 'tslog';
import { Kafka } from 'kafkajs';
import * as proto from '$lib/index_pb';


export const POST = (async (event) => {
    const logger = new Logger();

    // 
    const requestBody = await event.request.json()
    const {
        bountySignStatus,
        bountyId,
        bountyUIAmount,
        tokenAddress,
        creatorAddress,
        installationId
    } = requestBody

    logger.info(`linker.link.post: body = ${JSON.stringify(requestBody)}`);
    // wallet address

    // post data to kafka 
    const kafkaPwd = process.env.KAFKA_PASSWORD
    if (kafkaPwd) {
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

        // create bounty message
        const kafkaPayload = new proto.BountyMessage({
            BountySignStatus: bountySignStatus,
            Bountyid: bountyId,
            BountyUIAmount: bountyUIAmount,
            TokenAddress: tokenAddress,
            CreatorAddress: creatorAddress,
            InstallationId: installationId
        })


        const producer = kafka.producer()
        await producer.connect()
        await producer.send({
            topic: 'bounty',
            messages: [
                { value: Buffer.from(kafkaPayload.toBinary()), partition: 0 },
            ],
        })

        return new Response(JSON.stringify(kafkaPayload))
    } else {
        return new Response(JSON.stringify({
            status: "ok",
            subStatus: "Did not post data to kafka"
        }))
    }

})