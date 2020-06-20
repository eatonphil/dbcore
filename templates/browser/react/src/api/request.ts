// @ts-ignore
const apiPrefix = window.DBCORE_API_PREFIX; // Substituted by build process

export async function request(
  endpoint: string,
  body?: object,
  method?: 'POST' | 'GET' | 'DELETE' | 'PUT',
) {
  const req = await window.fetch(apiPrefix + '/{{api.router_prefix}}'+endpoint, {
    method: !method ? (body ? 'POST' : 'GET') : method,
    body: body ? JSON.stringify(body) : undefined,
    headers: body ? {
      'content-type': 'application/json',
    } : undefined,
    credentials: 'include',
  });

  return req.json();
}
