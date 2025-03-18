export async function onRequestGet(context) {
    const { env, request } = context;
    const url = new URL(request.url);
    const endpoint = url.pathname.split('/').pop();

    if (endpoint === 'mac-beta') {
        try {
            const data = await env.WARP_DIAG_CHECKER.get('mac-beta-version');
            return new Response(data, {
                status: 200,
                headers: {
                    'Content-Type': 'text/plain',
                    'Cache-Control': 'public, max-age=300'
                }
            });
        } catch (error) {
            return new Response('Error fetching data', { status: 500 });
        }
    } else {
        return new Response('Not found', { status: 404 });
    }
}
