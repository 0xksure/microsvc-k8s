<script lang="ts">
    import * as proto from "$lib/index_pb";
    import { WalletMultiButton } from "@svelte-on-solana/wallet-adapter-ui";
    import { walletStore } from "@svelte-on-solana/wallet-adapter-core";

    import * as bounty from "bounty-sdk";
    import { Connection, PublicKey } from "@solana/web3.js";
    export let data: {
        bountyParams: proto.BountyMessage;
    };
    let referrer = data.referrer;

    let REDIRECT_IN_SECONDS = 5;
    let startedRedirect = 0;
    let redirectIn = REDIRECT_IN_SECONDS;

    function redirectWithTimeout() {
        startedRedirect = Math.floor(Date.now() / 1000);
        setTimeout(() => {
            window.location.assign(referrer);
        }, REDIRECT_IN_SECONDS * 1000);

        setInterval(() => {
            redirectIn =
                startedRedirect +
                REDIRECT_IN_SECONDS -
                Math.floor(Date.now() / 1000);
        }, 1000);
    }

    async function createBounty() {
        console.log("create bounty");
        const rpcUrl =
            process.env.RPC_URL ?? "https://api.mainnet-beta.solana.com";
        if (!rpcUrl) throw new Error("RPC_URL is not defined");
        const connection = new Connection(rpcUrl, "confirmed");
        // sign transaction and send
        if (
            !walletStore ||
            !$walletStore.publicKey ||
            !$walletStore.signTransaction
        ) {
            return;
        }

        const bountySDK = new bounty.BountySdk(
            $walletStore.publicKey,
            connection
        );
        const createBounty = await bountySDK.createBounty({
            id: data.bountyParams.Bountyid.toString(),
            bountyAmount: data.bountyParams.BountyUIAmount,
            bountyCreator: $walletStore.publicKey,
            mint: new PublicKey(data.bountyParams.TokenAddress),
            platform: data.bountyParams.platform,
            organization: data.bountyParams.organization,
            team: data.bountyParams.team,
            domainType: data.bountyParams.domainType,
        });
        await $walletStore?.signTransaction(await createBounty.vtx);
        await bounty.utils.sendAndConfirmTransaction(
            connection,
            await createBounty.vtx
        );

        // call backend with info to create bounty
        fetch("http://localhost:3030/bounty/create", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                bountySignStatus: proto.BountySignStatus.FAILED_TO_SIGN,
                bountyId: data.bountyParams.Bountyid,
                bountyUIAmount: data.bountyParams.BountyUIAmount,
                tokenAddress: data.bountyParams.TokenAddress,
                creatorAddress: data.bountyParams.CreatorAddress,
                installationId: data.bountyParams.InstallationId,
            }),
        })
            .then((response) => response.json())
            .then((data) => {
                console.log("Success:", data);
            })
            .catch((error) => {
                console.error("Error:", error);
            });
    }
</script>

<div class="flex flex-col items-center">
    <div class="address">
        <WalletMultiButton />
    </div>
    <div>
        <h2 class="text-2xl text-white">Bounty</h2>
    </div>
    <div>
        <p class="text-white">Bounty ID: {data.bountyParams.Bountyid}</p>
        <p class="text-white">Amount: {data.bountyParams.BountyUIAmount}</p>
        <p class="text-white">
            Token address: {data.bountyParams.TokenAddress}
        </p>
        <p class="text-white">Creator: {data.bountyParams.CreatorAddress}</p>
    </div>
    <button
        class="m-2 border-md bg-blue-200 p-2 rounded-md"
        on:click={createBounty}
    >
        Create Bounty
    </button>
</div>
