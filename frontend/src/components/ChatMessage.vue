<template>
	<article class="chat-message" :class="role">
		<div class="meta">{{ roleLabel }}</div>
		<div class="content">{{ content }}</div>
	</article>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
	role: 'user' | 'assistant'
	content: string
}>()

const roleLabel = computed(() => (props.role === 'user' ? 'You' : 'SignalStack AI'))
</script>

<style scoped>
.chat-message {
	max-width: 84%;
	padding: 16px 18px;
	border-radius: 20px;
	border: 1px solid rgba(148, 163, 184, 0.14);
}

.chat-message.user {
	margin-left: auto;
	background: rgba(41, 211, 161, 0.12);
}

.chat-message.assistant {
	background: rgba(18, 29, 52, 0.9);
}

.meta {
	font-size: 0.8rem;
	color: #94a3b8;
	margin-bottom: 8px;
}

.content {
	line-height: 1.6;
	white-space: pre-wrap;
}
</style>
