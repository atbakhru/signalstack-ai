import { request } from './client'

export type Source = {
	id: number
	name: string
	api_name: string
	base_url: string
	enabled: boolean
	created_at: string
}

export type IngestionRun = {
	id: number
	source_id: number
	status: string
	documents_fetched: number
	documents_inserted: number
	chunks_created: number
	embeddings_created: number
	started_at: string
	completed_at?: string
	error_message?: string
}

export async function ingestAll() {
	return request<{ results: unknown[] }>('/api/ingest/all', { method: 'POST' })
}

export async function ingestSource(source: string) {
	return request(`/api/ingest/${encodeURIComponent(source)}`, { method: 'POST' })
}

export async function listRuns() {
	return request<{ runs: IngestionRun[] }>('/api/ingest/runs')
}

export async function getRun(id: number) {
	return request<{ run: IngestionRun }>(`/api/ingest/runs/${id}`)
}

export async function listSources() {
	return request<{ sources: Source[] }>('/api/ingest/sources')
}

export async function toggleSource(id: number) {
	return request<{ source: Source }>(`/api/ingest/sources/${id}/toggle`, { method: 'PATCH' })
}
