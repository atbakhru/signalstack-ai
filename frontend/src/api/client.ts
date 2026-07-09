export const apiBaseUrl = import.meta.env.VITE_API_BASE_URL ?? 'http://localhost:8080'

export type ApiResponse<T> = T

export class ApiError extends Error {
	status: number

	constructor(message: string, status: number) {
		super(message)
		this.status = status
	}
}

export function getAuthToken() {
	return localStorage.getItem('signalstack_token')
}

export function setAuthToken(token: string | null) {
	if (token) {
		localStorage.setItem('signalstack_token', token)
		return
	}
	localStorage.removeItem('signalstack_token')
}

export async function request<T>(path: string, init: RequestInit = {}): Promise<T> {
	const headers = new Headers(init.headers)
	if (!headers.has('Content-Type') && init.body && !(init.body instanceof FormData)) {
		headers.set('Content-Type', 'application/json')
	}
	const token = getAuthToken()
	if (token) {
		headers.set('Authorization', `Bearer ${token}`)
	}

	const response = await fetch(`${apiBaseUrl}${path}`, {
		...init,
		headers,
	})

	const contentType = response.headers.get('content-type') ?? ''
	const payload = contentType.includes('application/json') ? await response.json() : null
	if (!response.ok) {
		throw new ApiError(payload?.error ?? response.statusText, response.status)
	}
	return payload as T
}
