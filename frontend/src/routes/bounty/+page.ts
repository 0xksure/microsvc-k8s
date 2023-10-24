import * as proto from '$lib/index_pb';
import { error } from '@sveltejs/kit';


/** @type {import('./$types').PageLoad} */
export async function load(data) {
    let bountyIdProto = data.url.searchParams.get("bountyId");
    if (!bountyIdProto) {
        throw error(404, 'Bounty not found');
    }
    const bountyId = parseInt(bountyIdProto)
    const tokenAddress = data.url.searchParams.get("tokenAddress");
    if (!tokenAddress) {
        throw error(404, 'Token address not supplied');
    }
    const bountyUIAmount = data.url.searchParams.get("bountyUIAmount");
    if (!bountyUIAmount) {
        throw error(404, 'Bounty UI Amount not found');
    }
    const creatorAddress = data.url.searchParams.get("creatorAddress");
    if (!creatorAddress) {
        throw error(404, 'Creator not found');
    }
    const installationId = data.url.searchParams.get("installationId");
    if (!installationId) {
        throw error(404, 'Installation Id not found');
    }

    const platform = data.url.searchParams.get("platform");
    if (!platform) {
        throw error(404, 'Platform not found');
    }
    const organization = data.url.searchParams.get("organization");
    if (!organization) {
        throw error(404, 'Organization not found');
    }
    const team = data.url.searchParams.get("team");
    if (!team) {
        throw error(404, 'Team not found');
    }
    const domainType = data.url.searchParams.get("domainType");
    if (!domainType) {
        throw error(404, 'Domain type not found');
    }

    const bountyParams = new proto.BountyMessage({
        BountySignStatus: proto.BountySignStatus.CREATED,
        Bountyid: BigInt(bountyId),
        BountyUIAmount: bountyUIAmount,
        TokenAddress: tokenAddress,
        CreatorAddress: creatorAddress,
        InstallationId: BigInt(installationId),
        platform: platform,
        organization: organization,
        team: team,
        domainType: domainType
    })

    const referrer = data.url.searchParams.get('referrer');


    return {
        referrer: referrer,
        bountyParams: bountyParams.toJson()
    }
}