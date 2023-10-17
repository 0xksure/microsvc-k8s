
/** @type {import('./$types').PageLoad} */
export async function load({url }) {
	const ghLogin = url.searchParams.get('ghLogin');
	const ghId = url.searchParams.get('ghId');
	return { ghLogin,ghId };
}