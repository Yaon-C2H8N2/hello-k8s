import Cookies from "js-cookie";

export const fetchAPI = async (
    input: RequestInfo | URL,
    init?: RequestInit | undefined,
) => {
    const token = Cookies.get("token");
    const wrappedHeaders = new Headers(init?.headers) || new Headers();
    (wrappedHeaders as Headers).append("Authorization", `Bearer ${token}`);

    init = {
        ...init,
        headers: wrappedHeaders,
    };
    return fetch(input, init);
}
