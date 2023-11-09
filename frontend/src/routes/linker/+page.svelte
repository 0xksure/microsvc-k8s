<script lang="ts">
    import { WalletMultiButton } from "@svelte-on-solana/wallet-adapter-ui";
    import { walletStore } from "@svelte-on-solana/wallet-adapter-core";
    import type { LinkerResponse } from "./link/+server";
    import {
        Connection,
        PublicKey,
        VersionedTransaction,
    } from "@solana/web3.js";
    import * as linker from "idlinker-sdk/dist/cjs/index.js";
    import { onMount } from "svelte";

    /** @type {import('./$types').PageData} */
    export let data;
    let linked = false;
    let linkerSdk: linker.IdentitySdk | null;
    let connection: Connection | null;
    let identityAccount: {
        account: { social: string; username: string; address: PublicKey };
    } | null;

    async function isLinked(
        linkerSdk: linker.IdentitySdk,
        ghLogin?: string,
        address?: string
    ) {
        if (!ghLogin || !address) return false;
        if (data && data.ghLogin) {
            const identityAccounts = await linkerSdk.getIdentityFromUsername({
                username: data.ghLogin,
            });
            if (identityAccounts.length > 0) {
                identityAccount = identityAccounts[0];
            }
            console.log(identityAccount);
        }
    }

    onMount(async () => {
        const rpcUrl = process.env.RPC_URL;
        if (!rpcUrl) throw new Error("RPC_URL is not defined");
        const publicKey = $walletStore?.publicKey ?? PublicKey.default;
        connection = new Connection(rpcUrl, "confirmed");
        linkerSdk = new linker.IdentitySdk(publicKey, connection);
        await isLinked(linkerSdk, data.ghLogin, publicKey.toString());
    });

    $: async () => {
        console.log("getting identity account");
        if (linkerSdk) {
            console.log("getting identity account");
            if (data.ghLogin) {
                identityAccount = await linkerSdk.getIdentityFromUsername(
                    data.ghLogin
                );
                console.log(identityAccount);
            }
        }
    };

    async function linkWalletAndProfile() {
        // sign link transaction
        if (
            !$walletStore ||
            !$walletStore.publicKey ||
            !$walletStore.signTransaction
        ) {
            return;
        }

        console.log("linking profiles");
        const res = await fetch("/linker/link", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                walletAddress: $walletStore.publicKey,
                ghLogin: data.ghLogin,
                identityOwner: $walletStore.publicKey,
            }),
        });
        if (!res.ok) {
            console.log("error linking wallet and profile");
            return;
        }

        const linkerResponse: LinkerResponse = await res.json();
        const vtx = VersionedTransaction.deserialize(
            Buffer.from(linkerResponse.vtx)
        );
        const signedVtx = await $walletStore?.signTransaction(vtx);
        if (!signedVtx) {
            console.log("error signing link transaction");
            return;
        }
        const rpcUrl = process.env?.RPC_URL ?? "https://api.devnet.solana.com";
        if (!rpcUrl) throw new Error("RPC_URL is not defined");

        try {
            const linkerSignature = await connection.sendTransaction(signedVtx);
            console.log("linked wallet and profile");
            const latestBlockhash = await connection.getLatestBlockhash();
            await connection.confirmTransaction({
                signature: linkerSignature,
                ...latestBlockhash,
            });
            linked = true;
        } catch (e) {
            console.log(e);
        }

        // update state

        //
    }
</script>

<div class="flex flex-col items-center gap-4">
    <div class="address">
        <WalletMultiButton />
    </div>
    <div>
        <h1 class="text-3xl text-white">Linker</h1>
        <p class="text-white">Link your github profile and wallet</p>
    </div>
    {#if identityAccount}
        <div class="border-2 border-green-800 p-2">
            <h2 class="text-center text-xl text-white">
                {`Accounts are already linked!`}
            </h2>
            <div class="flex flex-col justify-center items-center">
                <p class="text-white">
                    {identityAccount.account.username}
                </p>
                <div class="text-green-800">
                    {"<->"}
                </div>
                <p class="text-white">
                    {identityAccount.account.address
                        .toString()
                        .slice(0, 5)}...{identityAccount.account.address
                        .toString()
                        .slice(-5)}
                </p>
            </div>
            <h2 class="text-center text-xl text-white">
                {`You are ready to earn bounties!`}
            </h2>
        </div>
    {:else}
        <div>
            {#if data.ghLogin}
                <div class="border-2 border-green-800 p-2">
                    <h2 class="text-xl text-white">
                        {`Logged in with ${data.ghLogin}`}
                    </h2>
                </div>
            {:else}
                <a
                    class={`z-0 w-40 border-2 rounded-md border-black p-1 shadow-md text-center flex items-center ${
                        data.ghLogin ? "bg-green-800" : "bg-indigo-200"
                    }`}
                    href={"/login/github"}
                >
                    <p
                        class={`mx-auto font-pixel ${
                            data.ghLogin ? "text-white" : "text-slate-800"
                        }`}
                    >
                        Login with Github
                    </p>
                </a>
            {/if}
        </div>

        <div>
            {#if $walletStore.publicKey}
                <div class="border-2 border-green-800 p-2">
                    <h2 class="text-xl text-white">
                        {`Wallet connected: ${$walletStore.publicKey
                            .toString()
                            .slice(0, 5)}...${$walletStore.publicKey
                            .toString()
                            .slice(-5)}`}
                    </h2>
                </div>
            {:else}
                <WalletMultiButton />
            {/if}
        </div>

        {#if data.ghLogin && $walletStore.publicKey && !linked}
            <div
                class="flex flex-col items-center gap-2 border-2 border-yellow-500 p-2"
            >
                <div>
                    <h2 class="text-xl text-white">Ready to link accounts</h2>
                </div>
                <div class="flex flex-col justify-center items-center">
                    <p class="text-white">
                        {$walletStore.publicKey}
                    </p>
                    <div class="text-green-800">
                        {"<->"}
                    </div>
                    <p class="text-white">
                        {data.ghLogin}
                    </p>
                </div>
                <button
                    class={`z-0 w-40 border-2 rounded-md border-black p-1 shadow-md text-center flex items-center ${" bg-indigo-200"}`}
                    on:click={linkWalletAndProfile}
                >
                    <p class={`mx-auto font-pixel ${"text-slate-800"}`}>Link</p>
                </button>
            </div>
        {/if}
    {/if}
</div>
