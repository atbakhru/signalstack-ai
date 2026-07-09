import { request, setAuthToken } from './client'

export type User = {
	id: number
	name: string
	email: string
	created_at: string
}

export type AuthResponse = {
	user: User
	token: string
}

export async function register(payload: { name: string; email: string; password: string }) {
	const response = await request<AuthResponse>('/api/auth/register', {
		method: 'POST',
		body: JSON.stringify(payload),
	})
	setAuthToken(response.token)
	return response
}

export async function login(payload: { email: string; password: string }) {
	const response = await request<AuthResponse>('/api/auth/login', {
		method: 'POST',
		body: JSON.stringify(payload),
	})
	setAuthToken(response.token)
	return response
}

export async function me() {
	return request<{ user: User }>('/api/auth/me')
}
