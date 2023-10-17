// Linker/link is responsible for validating the github login 
// and making a call to the identity service to link the github account to the wallet

import { error } from '@sveltejs/kit';
import jwt from 'jsonwebtoken';
import { Logger } from 'tslog';


export const POST = (async (event) => {
    const logger = new Logger();

    const ghJwt = event.cookies.get('ghJwt')
    if(!ghJwt) throw error(400,'No github jwt found')
    const jwtSecret = process.env.JWT_SECRET
    if(!jwtSecret) throw error(400,'No jwt secret found')
    const ghJwtDecoded = jwt.verify(ghJwt, jwtSecret)
    if(!ghJwtDecoded) throw error(400,'Invalid jwt')
    const ghAccessToken = ghJwtDecoded?.token
    logger.info(`linker.link.post: ghAccessToken=${ghAccessToken}, ghJwt=${ghJwt}, ghJwtDecoded=${JSON.stringify(ghJwtDecoded)}`);
    if(!ghAccessToken) throw error(400,'No github access token found')


    const getUserResponse = await fetch(`https://api.github.com/user`, {
		method: 'GET',
		headers: {
			Authorization: `Bearer ${ghAccessToken}`,
			accept: 'application/json'
		}
	});
	if (!getUserResponse.ok) throw error(400,`Failed to get user information: ${getUserResponse.status}`);
    const userData = await getUserResponse.json()
    const username = userData.login
    const userId = userData.id

    const requestBody = await event.request.json()
    logger.info(`linker.link.post: body = ${JSON.stringify(requestBody)}`);
    // wallet address
    const walletAddress = requestBody.walletAddress
    if(!walletAddress) throw error(400,'No wallet address found')

    // post data to rpc endpoint 
     
    return new Response(JSON.stringify({username,userId,walletAddress}))
})