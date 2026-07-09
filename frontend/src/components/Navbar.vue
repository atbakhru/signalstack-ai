<template>
	<header class="nav-shell">
		<div class="nav-inner">
			<div class="brand">
				<span class="brand-mark">S</span>
				<div>
					<div class="brand-name">SignalStack AI</div>
					<div class="brand-subtitle">Multi-source intelligence platform</div>
				</div>
			</div>

			<nav class="nav-links">
				<RouterLink to="/dashboard">Dashboard</RouterLink>
				<RouterLink to="/ingestion">Ingestion</RouterLink>
				<RouterLink to="/documents">Documents</RouterLink>
				<RouterLink to="/chat">Chat</RouterLink>
				<RouterLink to="/evaluation">Evaluation</RouterLink>
				<RouterLink to="/metrics">Metrics</RouterLink>
				<template v-if="authStore.isAuthenticated">
					<span class="nav-user">{{ authStore.user?.email ?? 'Signed in' }}</span>
					<button class="button secondary" type="button" @click="logout">Logout</button>
				</template>
			</nav>
		</div>
	</header>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/authStore'

const authStore = useAuthStore()
const router = useRouter()

async function logout() {
	authStore.logout()
	await router.push('/login')
}
</script>

<style scoped>
.nav-shell {
	position: sticky;
	top: 0;
	z-index: 20;
	backdrop-filter: blur(18px);
	background: rgba(6, 12, 22, 0.75);
	border-bottom: 1px solid rgba(148, 163, 184, 0.12);
}

.nav-inner {
	width: min(1240px, calc(100% - 32px));
	margin: 0 auto;
	min-height: 74px;
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: 20px;
}

.brand {
	display: flex;
	align-items: center;
	gap: 14px;
}

.brand-mark {
	width: 42px;
	height: 42px;
	border-radius: 14px;
	display: grid;
	place-items: center;
	font-weight: 800;
	background: linear-gradient(135deg, #29d3a1, #f97316);
	color: #041018;
}

.brand-name {
	font-weight: 800;
	letter-spacing: 0.01em;
}

.brand-subtitle {
	color: #94a3b8;
	font-size: 0.83rem;
}

.nav-links {
	display: flex;
	flex-wrap: wrap;
	gap: 14px;
	align-items: center;
}

.nav-links a {
	color: #cbd5e1;
	font-weight: 600;
}

.nav-links a.router-link-active {
	color: #29d3a1;
}

.nav-user {
	color: #94a3b8;
	font-size: 0.92rem;
}

@media (max-width: 860px) {
	.nav-inner {
		flex-direction: column;
		align-items: flex-start;
		padding: 14px 0;
	}
}
</style>
