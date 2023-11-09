<script lang="ts">
    import * as proto from "$lib/index_pb";
    import { WalletMultiButton } from "@svelte-on-solana/wallet-adapter-ui";
    import { walletStore } from "@svelte-on-solana/wallet-adapter-core";
    import BN from "bn.js";
    import * as bounty from "bounty-sdk/dist/cjs/index.js";
    import { Connection, PublicKey } from "@solana/web3.js";
    export let data: {
        bountyParams: proto.BountyMessage;
    };
    let referrer = data.referrer;

    let inputRef = null;

    let localState = {
        err: "",
    };

    let bountyId = data.bountyParams.Bountyid;
    let bountyAmount = data.bountyParams.BountyUIAmount;
    let tokenAddress = data.bountyParams.TokenAddress;
    let creatorAddress = data.bountyParams.CreatorAddress;

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

    async function createBounty(event: Event) {
        try {
            const mint = new PublicKey(tokenAddress);
        } catch (e) {
            console.log("Invalid mint address");
            console.log(inputRef);
            document
                .getElementById("token-address")
                .classList.add("border-2", "border-red-500");
            localState.err = "Invalid mint address. Please try again.";
            return;
        }

        const rpcUrl = process?.env?.RPC_URL ?? "https://api.devnet.solana.com";
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
        const createBountyArgs = {
            id: bountyId.toString(),
            bountyAmount: bountyAmount,
            bountyCreator: $walletStore.publicKey,
            mint: new PublicKey(tokenAddress),
            platform: data.bountyParams.platform,
            organization: data.bountyParams.organization,
            team: data.bountyParams.team,
            domainType: data.bountyParams.domainType,
        };
        console.log(createBountyArgs);
        const createBounty = await bountySDK.createBounty({
            id: bountyId.toString(),
            bountyAmount: new BN(bountyAmount),
            bountyCreator: $walletStore.publicKey,
            mint: new PublicKey(tokenAddress),
            platform: data.bountyParams.platform,
            organization: data.bountyParams.organization,
            team: data.bountyParams.team,
            domainType: data.bountyParams.domainType,
        });
        try {
            const vtx = await $walletStore?.signTransaction(
                await createBounty.vtx
            );
            await bounty.utils.sendAndConfirmTransaction(connection, vtx, []);
            // call backend with info to create bounty
            fetch("/bounty/create", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    bountySignStatus: proto.BountySignStatus.SIGNED,
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
                    localState.err = "Error storing bounty. Please try again.";
                });
        } catch (e) {
            console.log("error: ", e);
            localState.err =
                "Error creating bounty. Please try again. cause" + e?.message;
            throw e;
        }
    }
</script>

<div class="flex flex-col items-center gap-2">
    <div class="address">
        <WalletMultiButton />
    </div>
    <div>
        <h2 class="text-2xl text-white">Bounty from url</h2>
    </div>
    <form
        class="flex flex-col justify-center gap-2"
        on:submit|preventDefault={createBounty}
    >
        <label class="flex flex-row gap-2 items-center justify-between">
            <p class="text-white">Bounty ID</p>
            <input
                id="bountyId"
                name="BountyId"
                type="text"
                class=" text-white border-2 border-gray-500 rounded-md p-1 bg-transparent"
                placeholder={data.bountyParams.Bountyid.toString()}
                bind:value={bountyId}
            />
        </label>
        <label class="flex flex-row gap-2 items-center justify-between">
            <p class="text-white">Bounty Amount</p>
            <input
                id="bountyAmount"
                name="AmountUI"
                type="text"
                class="text-white border-2 border-gray-500 rounded-md p-1 bg-transparent"
                bind:value={bountyAmount}
            />
        </label>

        <label class="flex flex-row gap-2 items-center justify-between">
            <p class="text-white">Token address</p>
            <input
                id="token-address"
                name="tokenAddress"
                type="text"
                bind:value={tokenAddress}
                class=" text-white border-2 border-gray-500 rounded-md p-1 bg-transparent"
            />
        </label>

        <label class="flex flex-row gap-2 items-center justify-between">
            <p class="text-white">Creator address</p>
            <input
                id="creator-address"
                name="creatorAddress"
                type="text"
                placeholder="Enter creator address"
                bind:value={creatorAddress}
                class="text-white border-2 border-gray-500 rounded-md p-1 bg-transparent"
            />
        </label>
        <button class="m-2 border-md bg-blue-200 p-2 rounded-md">
            Create Bounty
        </button>
    </form>
    {#if localState.err != ""}
        <p class="text-red-500">{localState.err}</p>
    {/if}
</div>
