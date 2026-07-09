<template>
	<section class="page-grid cols-2">
		<div class="panel">
			<div class="panel-body stack">
				<div class="eyebrow">Create account</div>
				<h1 class="title-xl">Register for SignalStack AI</h1>
				<p class="subtle">Users can save chat sessions and revisit cited answers later.</p>
			</div>
		</div>

		<form class="panel" @submit.prevent="submit">
			<div class="panel-body stack">
				<input v-model="name" class="input" type="text" placeholder="Name" />
				<input v-model="email" class="input" type="email" placeholder="Email" />
				<input v-model="password" class="input" type="password" placeholder="Password" />
				<p v-if="error" class="subtle" style="color: #fca5a5">{{ error }}</p>
				<button class="button" type="submit" :disabled="loading">{{ loading ? 'Creating...' : 'Create account' }}</button>
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
const name = ref('SignalStack User')
const email = ref('user@signalstack.local')
const password = ref('demo-password')
const loading = ref(false)
const error = ref('')

async function submit() {
	loading.value = true
	error.value = ''
	try {
		await authStore.register(name.value, email.value, password.value)
		await router.push('/dashboard')
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Failed to register'
	} finally {
		loading.value = false
	}
}
</script>
