<template>
	<section class="page-grid cols-2">
		<!-- Left: chat interaction panel -->
		<ChatWindow @preview="openPreview" />

		<!-- Right: sessions sidebar -->
		<div class="stack">
			<div class="panel">
				<div class="panel-body stack">
					<div class="row">
						<h2 class="title-lg">Sessions</h2>
						<button class="button secondary" type="button" @click="refreshSessions">Refresh</button>
					</div>
					<p v-if="!chatStore.sessions.length" class="subtle">No sessions yet. Ask a question to start one.</p>
					<div
						v-for="session in chatStore.sessions"
						:key="session.id"
						class="panel"
						style="background: rgba(8, 17, 31, 0.6)"
					>
						<div class="panel-body row">
							<div>
								<strong>{{ session.title }}</strong>
								<div class="subtle">Session {{ session.id }}</div>
							</div>
							<div class="row" style="gap: 8px">
								<button class="button secondary" type="button" @click="openSession(session.id)">Open</button>
								<button class="button secondary" type="button" style="color: #fca5a5" @click="removeSession(session.id)">Delete</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>

	<!-- Source chunk preview modal -->
	<SourceModal :open="modalOpen" :context="modalContext" @close="closePreview" />
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import ChatWindow from '../components/ChatWindow.vue'
import SourceModal from '../components/SourceModal.vue'
import { useChatStore } from '../stores/chatStore'
import type { Citation, RetrievedContext } from '../api/chat'

const chatStore = useChatStore()
const modalOpen = ref(false)
const modalContext = ref<RetrievedContext | null>(null)

function openPreview(citation: Citation) {
	const ctx = chatStore.retrievedContexts.find(rc => rc.chunk.id === citation.chunk_id) ?? null
	modalContext.value = ctx
	modalOpen.value = true
}

function closePreview() {
	modalOpen.value = false
	modalContext.value = null
}

async function refreshSessions() {
	await chatStore.loadSessions()
}

async function openSession(id: number) {
	await chatStore.loadSession(id)
}

async function removeSession(id: number) {
	await chatStore.removeSession(id)
}

onMounted(async () => {
	await chatStore.loadSessions()
})
</script>

