<script >
    import { WalletMultiButton } from "@svelte-on-solana/wallet-adapter-ui";
    import { walletStore } from '@svelte-on-solana/wallet-adapter-core';

    /** @type {import('./$types').PageData} */
	export let data;

    function linkWalletAndProfile() {
        console.log("linking bounty");
        const res = fetch('/linker/link', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                walletAddress: $walletStore.publicKey,
                ghLogin: data.ghLogin
            })
        }).then(res => res.json()).then(res => {
            console.log(res);
        })
    }
</script>

<div class="flex flex-col items-center gap-4">
    <div class="address">
        <WalletMultiButton />
    </div>
    <div>
            <h2 class="text-2xl text-white">Linker</h2>
            <p class="text-white">Link your github profile and wallet</p>
    </div>
    
    <div>
        <a
        class={`z-0 w-40 border-2 rounded-md border-black p-1 shadow-md text-center flex items-center ${
            data.ghLogin ? 'bg-green-800' :'bg-indigo-200'
        }`}
        href={'/login/github'}
    >
        <p class={`mx-auto font-pixel ${data.ghLogin ? 'text-white' :  'text-slate-800'}`}>
            Login with Github
        </p>
    </a>
    </div>

    <div class="flex flex-col items-center gap-2">
        {#if data.ghLogin &&  $walletStore.publicKey}

        <div>
            <h2 class="text-xl text-white">Ready to link accounts</h2>
        </div>
        <div class="flex flex-col justify-center items-center">
            <p  class="text-white">
                {$walletStore.publicKey}
        </p>
            <div class="text-green-800">
                {"<->"}
            </div>
            <p  class="text-white">
                {data.ghLogin}
        </p>
        </div>
        <button
        class={`z-0 w-40 border-2 rounded-md border-black p-1 shadow-md text-center flex items-center ${
             ' bg-indigo-200'
        }`}
        on:click={linkWalletAndProfile}
    >
        <p class={`mx-auto font-pixel ${ 'text-slate-800'}`}>
            Link
        </p>
    </button>
        {/if}

    </div>
</div>