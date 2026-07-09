<template>
	<section class="stack">
		<div class="panel">
			<div class="panel-body stack">
				<div class="eyebrow">Evaluation</div>
				<h1 class="title-xl">Benchmark retrieval quality</h1>
				<div class="row">
					<span class="pill success">Top-k accuracy: {{ summary.top_k_accuracy.toFixed(2) }}</span>
					<span class="pill">Average answer latency: {{ summary.average_answer_ms }} ms</span>
					<button class="button" type="button" :disabled="loading" @click="runBenchmark">{{ loading ? 'Running...' : 'Run benchmark' }}</button>
				</div>
				<p v-if="error" class="subtle" style="color: #fca5a5">{{ error }}</p>
			</div>
		</div>

		<EvaluationTable :results="tableResults" />
	</section>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import EvaluationTable from '../components/EvaluationTable.vue'
import { listEvaluationResults, runEvaluation, type EvaluationResult } from '../api/evaluation'

const loading = ref(false)
const error = ref('')
const results = ref<EvaluationResult[]>([])
const summary = reactive({ top_k_accuracy: 0, average_answer_ms: 0, average_retrieval_ms: 0, average_citations: 0, failed_retrievals: 0 })

const tableResults = computed(() => results.value.map(result => ({
	question: String(result.question_id),
	expectedSource: 'Indexed source',
	topKAccuracy: Math.round(result.top_k_accuracy * 100),
	retrievalLatencyMs: result.retrieval_latency_ms,
	answerLatencyMs: result.answer_latency_ms,
})))

async function refresh() {
	const response = await listEvaluationResults()
	results.value = response.results
	Object.assign(summary, response.summary)
}

async function runBenchmark() {
	loading.value = true
	error.value = ''
	try {
		const response = await runEvaluation()
		results.value = response.results
		Object.assign(summary, response.summary)
	} catch (err) {
		error.value = err instanceof Error ? err.message : 'Failed to run benchmark'
	} finally {
		loading.value = false
	}
}

onMounted(refresh)
</script>
