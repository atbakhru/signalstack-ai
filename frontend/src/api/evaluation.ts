import { request } from './client'

export type EvaluationQuestion = {
	id: number
	question: string
	expected_source: string
	expected_keywords: string[]
	created_at: string
}

export type EvaluationResult = {
	id: number
	question_id: number
	retrieved_document_ids: number[]
	top_k_accuracy: number
	answer_latency_ms: number
	retrieval_latency_ms: number
	citation_count: number
	created_at: string
}

export async function runEvaluation() {
	return request<{ summary: unknown; results: EvaluationResult[] }>('/api/evaluate/run', { method: 'POST' })
}

export async function listEvaluationResults() {
	return request<{ summary: unknown; results: EvaluationResult[] }>('/api/evaluate/results')
}

export async function listEvaluationQuestions() {
	return request<{ questions: EvaluationQuestion[] }>('/api/evaluate/questions')
}
