import { error } from '@sveltejs/kit';

type BountyParams = {
    bountyId: string;
    tokenAddress: string
    bountyUIAmount: string 
    creatorAddress: string
};

/** @type {import('./$types').PageLoad} */
export async function load({ params,url }) {
    const bountyId= url.searchParams.get("bountyId");
    if(!bountyId) {
        throw error(404, 'Bounty not found');
    }
    const tokenAddress= url.searchParams.get("tokenAddress");
    if(!tokenAddress) {
        throw error(404, 'Token address not supplied');
    }
    const bountyUIAmount= url.searchParams.get("bountyUIAmount");
    if(!bountyUIAmount) {
        throw error(404, 'Bounty UI Amount not found');
    }
    const creatorAddress= url.searchParams.get("creatorAddress");
    if(!creatorAddress) {
        throw error(404, 'Creator not found');
    }

    const bountyParams: BountyParams = {
        bountyId ,
        tokenAddress,
        bountyUIAmount,
        creatorAddress 
    };
 
	return {
		bountyParams
	}
}