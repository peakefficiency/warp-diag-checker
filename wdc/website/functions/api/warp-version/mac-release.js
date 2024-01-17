

export async function onRequestGet(context) {
    const { env, request } = context;
    const url = new URL(request.url);
    const endpoint = url.pathname.split('/').pop();

    if (endpoint === 'mac-release') {
        try {
            // Replace 'MY_KV_NAMESPACE' with your actual KV namespace binding name
            const data = await env.WARP_DIAG_CHECKER.get('mac-release-version');
            return new Response(data, { status: 200 });
        } catch (error) {
            return new Response('Error fetching data', { status: 500 });
        }
    } else {
        return new Response('Not found', { status: 401 });
    }
}
