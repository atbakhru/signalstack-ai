<template>
	<article class="panel citation-card">
		<div class="panel-body stack">
			<div class="row">
				<div>
					<div class="title-lg">{{ citation.title }}</div>
					<div class="subtle">{{ citation.source_name }} · relevance {{ citation.relevance_score.toFixed(3) }}</div>
				</div>
				<div class="row" style="gap: 8px">
					<button class="button secondary" type="button" @click="emit('preview', citation)">View chunk</button>
					<a class="button secondary" :href="citation.url" target="_blank" rel="noreferrer">Open source</a>
				</div>
			</div>
			<div class="subtle citation-url">{{ citation.url }}</div>
		</div>
	</article>
</template>

<script setup lang="ts">
import type { Citation } from '../api/chat'

defineProps<{
	citation: Citation
}>()

const emit = defineEmits<{
	preview: [citation: Citation]
}>()
</script>

