import { defineStore } from 'pinia'
import { login as loginRequest, me as meRequest, register as registerRequest, type User } from '../api/auth'
import { getAuthToken, setAuthToken } from '../api/client'

type AuthState = {
	user: User | null
	token: string | null
	initialized: boolean
}

export const useAuthStore = defineStore('auth', {
	state: (): AuthState => ({
		user: null,
		token: getAuthToken(),
		initialized: false,
	}),
	getters: {
		isAuthenticated: state => Boolean(state.token),
	},
	actions: {
		async initialize() {
			if (this.initialized) return
			this.initialized = true
			if (!this.token) return
			try {
				const response = await meRequest()
				this.user = response.user
			} catch {
				this.logout()
			}
		},
		setSession(user: User, token: string) {
			this.user = user
			this.token = token
			setAuthToken(token)
		},
		logout() {
			this.user = null
			this.token = null
			setAuthToken(null)
		},
		async login(email: string, password: string) {
			const response = await loginRequest({ email, password })
			this.setSession(response.user, response.token)
			return response
		},
		async register(name: string, email: string, password: string) {
			const response = await registerRequest({ name, email, password })
			this.setSession(response.user, response.token)
			return response
		},
	},
})
