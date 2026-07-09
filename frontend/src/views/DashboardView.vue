<template>
	<section class="stack">
		<div class="panel">
			<div class="panel-body stack">
				<div class="eyebrow">Overview</div>
				<div class="title-xl">SignalStack AI dashboard</div>
				<p class="subtle">Monitor ingestion volume, embeddings, source health, and retrieval latency across five public feeds.</p>
			</div>
		</div>

		<div class="page-grid cols-3">
			<MetricCard label="Documents ingested" :value="String(overview.documents_ingested)" note="Across all five sources" />
			<MetricCard label="Chunks generated" :value="String(overview.chunks_generated)" note="500-800 token windows" />
			<MetricCard label="Embeddings generated" :value="String(overview.embeddings_generated)" note="OpenAI vectors stored in pgvector" />
			<MetricCard label="Active sources" :value="String(overview.active_sources)" note="All public APIs enabled" />
			<MetricCard label="Avg retrieval latency" :value="formatMs(overview.average_retrieval_ms)" note="Top-k vector search" />
			<MetricCard label="Avg answer latency" :value="formatMs(overview.average_answer_ms)" note="Citation-grounded responses" />
		</div>

		<div class="page-grid cols-2">
			<section class="panel">
				<div class="panel-body stack">
					<div class="row">
						<h2 class="title-lg">Source health</h2>
						<span class="pill success">Live API status</span>
					</div>
					<SourceStatusCard
						v-for="source in sources"
						:key="source.name"
						:name="source.name"
						:api-name="source.api_name"
						:description="source.base_url"
						:enabled="source.enabled"
					/>
				</div>
			</section>

			<section class="panel">
				<div class="panel-body stack">
					<div class="row">
						<h2 class="title-lg">Recent ingestion runs</h2>
						<span class="pill">Latest status</span>
					</div>
					<IngestionRunTable :runs="recentRuns" />
				</div>
			</section>
		</div>
	</section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import MetricCard from '../components/MetricCard.vue'
import SourceStatusCard from '../components/SourceStatusCard.vue'
import IngestionRunTable from '../components/IngestionRunTable.vue'
import { getMetricsOverview } from '../api/metrics'
import { listSources, listRuns } from '../api/ingestion'

const overview = reactive({
	documents_ingested: 0,
	chunks_generated: 0,
	embeddings_generated: 0,
	active_sources: 0,
	average_retrieval_ms: 0,
	average_answer_ms: 0,
	token_usage: 0,
})

const sources = ref<Array<{ id: number; name: string; api_name: string; base_url: string; enabled: boolean }>>([])
const recentRuns = computed(() => runItems.value)
const runItems = ref<Array<{ id: number; source: string; status: string; documents: number; chunks: number; embeddings: number }>>([])

function formatMs(value: number) {
	return `${Math.round(value)} ms`
}

onMounted(async () => {
	const [overviewResponse, sourceResponse, runResponse] = await Promise.all([
		getMetricsOverview(),
		listSources(),
		listRuns(),
	])
	Object.assign(overview, overviewResponse)
	sources.value = sourceResponse.sources
	const sourceMap = new Map(sourceResponse.sources.map(source => [source.id, source.name]))
	runItems.value = runResponse.runs.slice(0, 6).map(run => ({
		id: run.id,
		source: sourceMap.get(run.source_id) ?? String(run.source_id),
		status: run.status,
		documents: run.documents_inserted,
		chunks: run.chunks_created,
		embeddings: run.embeddings_created,
	}))
})
</script>
