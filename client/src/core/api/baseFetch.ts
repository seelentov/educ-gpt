import { API_URL } from "./api";

export async function baseFetch(
    url: string,
    body: string | undefined | null = null,
    method: string = "GET",
    token: string = "",
    headers: { [key: string]: string } | null | undefined = null
): Promise<any> {
    try {
        const init: RequestInit = {
            method,
            headers: {
                'Content-Type': 'application/json',
            },
            mode: "no-cors"
        };

        if (headers) {
            init.headers = {
                ...init.headers,
                ...headers,
            };
        }

        if (token !== "") {
            init.headers = {
                ...init.headers,
                "Authorization": "Bearer " + token
            }
        }

        if (body) {
            init.body = body;
        }

        const response = await fetch(`${API_URL}${url}`, init);

        const contentLength = response.headers.get('Content-Length');
        if (response.status >= 200 && response.status < 300 && (contentLength === '0' || !contentLength)) {
            return { status: response.status };
        }

        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Ошибка:', error);

        if (error instanceof Error) {
            throw error;
        }

        throw new Error(`Неизвестная ошибка: ${error}`);
    }
}