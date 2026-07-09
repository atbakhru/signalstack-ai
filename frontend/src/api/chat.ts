import { request } from './client'

export type Citation = {
	id: number
	message_id: number
	document_id: number
	chunk_id: number
	source_name: string
	title: string
	url: string
	relevance_score: number
	created_at: string
}

export type ChatMessage = {
	id: number
	session_id: number
	user_id: number
	role: 'user' | 'assistant'
	content: string
	created_at: string
}

export type ChatSession = {
	id: number
	user_id: number
	title: string
	created_at: string
}

export type RetrievedContext = {
	chunk: { id: number; document_id: number; source_id: number; chunk_index: number; content: string; token_count: number; created_at: string }
	document_title: string
	source_name: string
	url: string
	similarity_score: number
}

export type ChatAnswer = {
	session: ChatSession
	user_message: ChatMessage
	assistant_message: ChatMessage
	citations: Citation[]
	retrieved_contexts: RetrievedContext[]
	answer: string
	used_context: boolean
	insufficient: boolean
	retrieval_latency_ms: number
	answer_latency_ms: number
}

export async function ask(payload: { question: string; session_id?: number; top_k?: number; source_ids?: number[] }) {
	return request<ChatAnswer>('/api/chat/ask', {
		method: 'POST',
		body: JSON.stringify(payload),
	})
}

export async function listSessions() {
	return request<{ sessions: ChatSession[] }>('/api/chat/sessions')
}

export async function getSession(id: number) {
	return request<{ session: ChatSession; messages: ChatMessage[] }>(`/api/chat/sessions/${id}`)
}

export async function deleteSession(id: number) {
	return request<void>(`/api/chat/sessions/${id}`, { method: 'DELETE' })
}
