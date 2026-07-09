<template>
	<section class="stack">
		<div class="panel">
			<div class="panel-body stack">
				<div class="row">
					<div>
						<div class="eyebrow">Ingestion</div>
						<h1 class="title-xl">Automatic source ingestion</h1>
					</div>
					<div class="row">
						<button class="button secondary" type="button" :disabled="busy" @click="runSelected">Run selected source</button>
						<button class="button" type="button" :disabled="busy" @click="runAll">Run all sources</button>
					</div>
				</div>
				<div class="page-grid cols-3">
					<select v-model.number="selectedSourceId" class="select">
						<option :value="0">Select source</option>
						<option v-for="source in sources" :key="source.id" :value="source.id">{{ source.name }}</option>
					</select>
				</div>
				<p class="subtle">The backend ingests data directly from the public APIs. No uploads, no manual document imports.</p>
				<p v-if="error" class="subtle" style="color: #fca5a5">{{ error }}</p>
			</div>
		</div>

		<div class="page-grid cols-3">
			<SourceStatusCard v-for="source in sources" :key="source.name" :name="source.name" :api-name="source.api_name" :description="source.base_url" :enabled="source.enabled" />
		</div>

		<IngestionRunTable :runs="runRows" />
	</section>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import SourceStatusCard from '../components/SourceStatusCard.vue'
import IngestionRunTable from '../components/IngestionRunTable.vue'
import { ingestAll, ingestSource, listRuns, listSources, type IngestionRun, type Source } from '../api/ingestion'

const sources = ref<Source[]>([])
const runs = ref<IngestionRun[]>([])
const selectedSourceId = ref(0)
const busy = ref(false)
const error = ref('')

const runRows = computed(() => {
	const sourceMap = new Map(sources.value.map(source => [source.id, source.name]))
	return runs.value.slice(0, 20).map(run => ({
		id: run.id,
		source: sourceMap.get(run.source_id) ?? String(run.source_id),
		status: run.status,
		documents: run.documents_inserted,
		chunks: run.chunks_created,
		embeddings: run.embeddings_created,
	}))
})

async function refresh() {
	const [sourceResponse, runResponse] = await Promise.all([listSources(), listRuns()])
	sources.value = sourceResponse.sources
	runs.value = runResponse.runs
}

async function runAll() {
	busy.value = true
	error.value = ''
	try {
		await ingestAll()
		await refresh()
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Failed to run ingestion'
	} finally {
		busy.value = false
	}
}

async function runSelected() {
	if (!selectedSourceId.value) return
	busy.value = true
	error.value = ''
	try {
		const source = sources.value.find(item => item.id === selectedSourceId.value)
		if (source) {
			await ingestSource(source.name)
			await refresh()
		}
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Failed to run ingestion'
	} finally {
		busy.value = false
	}
}

onMounted(refresh)
</script>
