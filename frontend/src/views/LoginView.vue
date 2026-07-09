<template>
	<section class="page-grid cols-2">
		<div class="panel">
			<div class="panel-body stack">
				<div class="eyebrow">Authentication</div>
				<h1 class="title-xl">Sign in to SignalStack AI</h1>
				<p class="subtle">Use JWT-backed auth to access the dashboard, chat, and evaluation pages.</p>
			</div>
		</div>

		<form class="panel" @submit.prevent="submit">
			<div class="panel-body stack">
				<input v-model="email" class="input" type="email" placeholder="Email" />
				<input v-model="password" class="input" type="password" placeholder="Password" />
				<p v-if="error" class="subtle" style="color: #fca5a5">{{ error }}</p>
				<button class="button" type="submit" :disabled="loading">{{ loading ? 'Signing in...' : 'Sign in' }}</button>
			</div>
		</form>
	</section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

const router = useRouter()
const authStore = useAuthStore()
const email = ref('demo@signalstack.local')
const password = ref('demo-password')
const loading = ref(false)
const error = ref('')

async function submit() {
	loading.value = true
	error.value = ''
	try {
		await authStore.login(email.value, password.value)
		await router.push('/dashboard')
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Failed to sign in'
	} finally {
		loading.value = false
	}
}
</script>
