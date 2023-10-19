import * as proto from '$lib/index_pb';
import { error } from '@sveltejs/kit';


/** @type {import('./$types').PageLoad} */
export async function load({ params,url }) {
    let bountyIdProto= url.searchParams.get("bountyId");
    if(!bountyIdProto) {
        throw error(404, 'Bounty not found');
    }
    const bountyId = parseInt(bountyIdProto)
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
    const installationId = url.searchParams.get("installationId");
    if(!installationId) {
        throw error(404, 'Installation Id not found');
    }

    const bountyParams = new proto.BountyMessage({
        BountySignStatus: proto.BountySignStatus.CREATED,
        Bountyid:BigInt(bountyId),
        BountyUIAmount: bountyUIAmount,
        TokenAddress: tokenAddress,
        CreatorAddress: creatorAddress,
        InstallationId: BigInt(installationId)
    })


    
 
	return {
		bountyParams: bountyParams.toJson()
	}
}