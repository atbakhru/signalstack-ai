<template>
	<Teleport to="body">
		<div v-if="open" class="modal-backdrop" @click.self="emit('close')">
			<div class="modal-panel" role="dialog" aria-modal="true">
				<div class="modal-header row">
					<div>
						<div class="eyebrow">Source chunk</div>
						<div class="title-lg">{{ context?.document_title ?? 'Unknown document' }}</div>
					</div>
					<button class="button secondary" type="button" @click="emit('close')">Close</button>
				</div>

				<div class="modal-meta row">
					<span class="pill">{{ context?.source_name }}</span>
					<span class="pill">Similarity {{ context?.similarity_score.toFixed(3) }}</span>
					<span class="pill">{{ context?.chunk.token_count }} tokens</span>
				</div>

				<div class="chunk-text">{{ context?.chunk.content ?? 'No content available.' }}</div>

				<div class="modal-footer row">
					<a
						v-if="context?.url"
						class="button"
						:href="context.url"
						target="_blank"
						rel="noreferrer"
					>Open original source</a>
				</div>
			</div>
		</div>
	</Teleport>
</template>

<script setup lang="ts">
import type { RetrievedContext } from '../api/chat'

defineProps<{
	open: boolean
	context: RetrievedContext | null
}>()

const emit = defineEmits<{
	close: []
}>()
</script>

<style scoped>
.modal-backdrop {
	position: fixed;
	inset: 0;
	z-index: 1000;
	background: rgba(0, 0, 0, 0.6);
	backdrop-filter: blur(4px);
	display: flex;
	align-items: center;
	justify-content: center;
	padding: 24px;
}

.modal-panel {
	background: var(--bg-card);
	border: 1px solid var(--border);
	border-radius: 16px;
	box-shadow: var(--shadow);
	width: 100%;
	max-width: 720px;
	max-height: 80vh;
	display: flex;
	flex-direction: column;
	gap: 16px;
	padding: 28px;
	overflow: hidden;
}

.modal-header {
	align-items: flex-start;
}

.modal-meta {
	flex-wrap: wrap;
	gap: 8px;
}

.chunk-text {
	flex: 1;
	overflow-y: auto;
	white-space: pre-wrap;
	line-height: 1.7;
	font-size: 0.9rem;
	padding: 16px;
	background: rgba(8, 17, 31, 0.5);
	border: 1px solid var(--border);
	border-radius: 10px;
	color: var(--text);
}

.modal-footer {
	justify-content: flex-end;
}
</style>

