<template>
	<div class="chat-window stack">
		<div class="panel">
			<div class="panel-body stack">
				<div class="eyebrow">Chat</div>
				<h1 class="title-xl">Ask across all indexed sources</h1>
				<div class="stack">
					<input
						v-model="question"
						class="input"
						type="text"
						placeholder="What are recent trends in AI infrastructure?"
						:disabled="chatStore.loading"
						@keyup.enter="submit"
					/>
					<div class="answer-box" :class="{ empty: !chatStore.answer }">
						{{ chatStore.answer || 'Retrieved answers will appear here with citations.' }}
					</div>
					<div class="row">
						<span class="pill warn">Citations required</span>
						<button class="button" type="button" :disabled="chatStore.loading" @click="submit">
							{{ chatStore.loading ? 'Thinking…' : 'Ask SignalStack AI' }}
						</button>
					</div>
					<p v-if="chatStore.error" class="error-text">{{ chatStore.error }}</p>
					<div v-if="chatStore.answer" class="row">
						<span class="pill">Retrieval {{ chatStore.retrievalLatencyMs }} ms</span>
						<span class="pill">Answer {{ chatStore.answerLatencyMs }} ms</span>
					</div>
				</div>
			</div>
		</div>

		<template v-if="chatStore.messages.length">
			<ChatMessage
				v-for="message in chatStore.messages"
				:key="message.id"
				:role="message.role"
				:content="message.content"
			/>
		</template>

		<template v-if="chatStore.citations.length">
			<div class="eyebrow" style="padding: 4px 0">Sources</div>
			<CitationCard
				v-for="citation in chatStore.citations"
				:key="citation.id"
				:citation="citation"
				@preview="emit('preview', citation)"
			/>
		</template>
	</div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ChatMessage from './ChatMessage.vue'
import CitationCard from './CitationCard.vue'
import { useChatStore } from '../stores/chatStore'
import type { Citation } from '../api/chat'

const emit = defineEmits<{
	preview: [citation: Citation]
}>()

const chatStore = useChatStore()
const question = ref('')

async function submit() {
	if (!question.value.trim() || chatStore.loading) return
	const q = question.value
	question.value = ''
	await chatStore.submitQuestion(q)
}
</script>

<style scoped>
.answer-box {
	padding: 14px 16px;
	border: 1px solid var(--border);
	border-radius: 10px;
	background: rgba(8, 17, 31, 0.5);
	min-height: 120px;
	white-space: pre-wrap;
	line-height: 1.65;
	font-size: 0.95rem;
}

.answer-box.empty {
	color: var(--muted);
}

.error-text {
	margin: 0;
	color: #fca5a5;
	font-size: 0.875rem;
}
</style>

