<template>
	<section class="stack">
		<div class="panel">
			<div class="panel-body stack">
				<div class="eyebrow">Metrics</div>
				<h1 class="title-xl">System observability</h1>
				<p class="subtle">Track ingestion volume, embeddings, latency, and source distribution.</p>
			</div>
		</div>

		<div class="page-grid cols-3">
			<MetricCard label="Documents ingested" :value="String(overview.documents_ingested)" note="Rows in documents" />
			<MetricCard label="Embeddings generated" :value="String(overview.embeddings_generated)" note="Rows with non-null vectors" />
			<MetricCard label="Token usage" :value="String(overview.token_usage)" note="Documents and chat content" />
			<MetricCard label="Average retrieval latency" :value="formatMs(overview.average_retrieval_ms)" note="Benchmark results" />
			<MetricCard label="Average answer latency" :value="formatMs(overview.average_answer_ms)" note="Benchmark results" />
			<MetricCard label="Active sources" :value="String(overview.active_sources)" note="Enabled public APIs" />
		</div>

		<div class="page-grid cols-2">
			<div class="panel"><div class="panel-body">Ingestion volume: {{ overview.documents_ingested }} documents</div></div>
			<div class="panel"><div class="panel-body">Embeddings generated: {{ overview.embeddings_generated }}</div></div>
			<div class="panel"><div class="panel-body">Retrieval latency: {{ formatMs(overview.average_retrieval_ms) }}</div></div>
			<div class="panel"><div class="panel-body">Source distribution: {{ sources.length }} sources</div></div>
		</div>
	</section>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from 'vue'
import MetricCard from '../components/MetricCard.vue'
import { getMetricsOverview, getMetricsSources } from '../api/metrics'

const overview = reactive({
	documents_ingested: 0,
	chunks_generated: 0,
	embeddings_generated: 0,
	active_sources: 0,
	average_retrieval_ms: 0,
	average_answer_ms: 0,
	token_usage: 0,
})

const sources = ref<Array<unknown>>([])

function formatMs(value: number) {
	return `${Math.round(value)} ms`
}

onMounted(async () => {
	const [overviewResponse, sourcesResponse] = await Promise.all([getMetricsOverview(), getMetricsSources()])
	Object.assign(overview, overviewResponse)
	sources.value = sourcesResponse.sources
})
</script>
