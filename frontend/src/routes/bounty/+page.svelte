<script >
    import { WalletMultiButton } from "@svelte-on-solana/wallet-adapter-ui";
    /** @type {import('./$types').PageData} */
	export let data;

    function createBounty() {
        console.log("create bounty");
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
        <p class="text-white">Bounty ID: {data.bountyParams.bountyId}</p>
        <p class="text-white">Amount: {data.bountyParams.bountyUIAmount}</p>
        <p class="text-white">Token address: {data.bountyParams.tokenAddress}</p>
        <p class="text-white">Creator: {data.bountyParams.creatorAddress}</p>
    </div>
    <button class="m-2 border-md bg-blue-200 p-2 rounded-md" on:click={createBounty}>
        Create Bounty
    </button>
</div>