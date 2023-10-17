import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import jwt from 'jsonwebtoken';

export const GET = (async ({ url, cookies }) => {
	const code = String(url.searchParams.get('code'));
	const jwtSecret = process.env.JWT_SECRET;
	if (!jwtSecret) throw Error('No jwt secret found');

	// authorize code
	const uri = `https://github.com/login/oauth/access_token?code=${code}&client_id=${process.env.GITHUB_CLIENT_ID}&client_secret=${process.env.GITHUB_CLIENT_SECRET}`;
	const resp = await fetch(uri, {
		method: 'GET',
		headers: {
			accept: 'application/json'
		}
	});
	if (!resp.ok) throw Error(`Failed to authorize code: ${resp.status}`);
	const text = JSON.parse(await resp.text());
	const accessToken = text.access_token;
	if (!accessToken) throw redirect(300, '/error');

	// get information about the user
	const getUserResponse = await fetch(`https://api.github.com/user`, {
		method: 'GET',
		headers: {
			Authorization: `Bearer ${accessToken}`,
			accept: 'application/json'
		}
	});
	if (!getUserResponse.ok) throw Error(`Failed to get user information: ${getUserResponse.status}`);

	const jwtToken = jwt.sign({ token: accessToken }, jwtSecret, { expiresIn: '1h' });
	cookies.set('ghJwt', jwtToken, {
		path: '/',
		maxAge: 60 * 60 * 24 * 7,
		sameSite: 'strict',
		httpOnly: true,
		secure: true
	});

	const githubUser = await getUserResponse.json();
	throw redirect(307, `${url.origin}/linker/?ghLogin=${githubUser.login}`);
}) satisfies RequestHandler;