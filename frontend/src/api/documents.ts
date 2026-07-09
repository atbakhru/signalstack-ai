import { request } from './client'

export type DocumentItem = {
	id: number
	data?: {
		id: number
		source_id: number
		external_id: string
		title: string
		summary: string
		url: string
		author: string
		published_at?: string
		raw_s3_key: string
		content_hash: string
		created_at: string
		updated_at: string
	}
	metadata?: unknown
	title?: string
	source_name?: string
	published_at?: string
	url?: string
	author?: string
	summary?: string
	source?: string
}

export type DocumentChunk = {
	id: number
	document_id: number
	source_id: number
	chunk_index: number
	content: string
	token_count: number
	created_at: string
}

export async function listDocuments(params: { sourceId?: number; search?: string } = {}) {
	const query = new URLSearchParams()
	if (params.sourceId) query.set('source_id', String(params.sourceId))
	if (params.search) query.set('search', params.search)
	const suffix = query.toString() ? `?${query.toString()}` : ''
	return request<{ documents: DocumentItem[] }>(`/api/documents${suffix}`)
}

export async function getDocument(id: number) {
	return request<{ document: { data: DocumentItem; source_name: string } }>(`/api/documents/${id}`)
}

export async function getDocumentChunks(id: number) {
	return request<{ chunks: DocumentChunk[] }>(`/api/documents/${id}/chunks`)
}

export async function deleteDocument(id: number) {
	return request<void>(`/api/documents/${id}`, { method: 'DELETE' })
}
