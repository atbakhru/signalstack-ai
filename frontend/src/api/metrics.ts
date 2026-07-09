import { request } from './client'
import type { Source } from './ingestion'

export type MetricsOverview = {
	documents_ingested: number
	chunks_generated: number
	embeddings_generated: number
	active_sources: number
	average_retrieval_ms: number
	average_answer_ms: number
	token_usage: number
}

export async function getMetricsOverview() {
	return request<MetricsOverview>('/api/metrics/overview')
}

export async function getMetricsSources() {
	return request<{ sources: Source[] }>('/api/metrics/sources')
}

export async function getMetricsLatency() {
	return request<{ summary: { average_retrieval_ms: number; average_answer_ms: number }; results: unknown[] }>('/api/metrics/latency')
}

export async function getMetricsEmbeddings() {
	return request<{ embeddings_generated: number }>('/api/metrics/embeddings')
}
