import { defineStore } from 'pinia'
import { ask, deleteSession as deleteSessionRequest, getSession as getSessionRequest, listSessions as listSessionsRequest, type ChatMessage, type ChatSession, type Citation, type RetrievedContext } from '../api/chat'

type ChatState = {
	sessions: ChatSession[]
	messages: ChatMessage[]
	citations: Citation[]
	retrievedContexts: RetrievedContext[]
	answer: string
	retrievalLatencyMs: number
	answerLatencyMs: number
	loading: boolean
	error: string | null
	activeSessionId: number | null
}

export const useChatStore = defineStore('chat', {
	state: (): ChatState => ({
		sessions: [],
		messages: [],
		citations: [],
		retrievedContexts: [],
		answer: '',
		retrievalLatencyMs: 0,
		answerLatencyMs: 0,
		loading: false,
		error: null,
		activeSessionId: null,
	}),
	actions: {
		async loadSessions() {
			const response = await listSessionsRequest()
			this.sessions = response.sessions
		},
		async loadSession(id: number) {
			const response = await getSessionRequest(id)
			this.activeSessionId = id
			this.messages = response.messages
		},
		async submitQuestion(question: string, topK = 5) {
			this.loading = true
			this.error = null
			try {
				const response = await ask({ question, session_id: this.activeSessionId ?? undefined, top_k: topK })
				this.activeSessionId = response.session.id
				this.answer = response.answer
				this.citations = response.citations
				this.retrievedContexts = response.retrieved_contexts
				this.retrievalLatencyMs = response.retrieval_latency_ms
				this.answerLatencyMs = response.answer_latency_ms
				this.messages = [response.user_message, response.assistant_message]
				await this.loadSessions()
			} catch (error) {
				this.error = error instanceof Error ? error.message : 'Failed to submit question'
				throw error
			} finally {
				this.loading = false
			}
		},
		async removeSession(id: number) {
			await deleteSessionRequest(id)
			this.sessions = this.sessions.filter(session => session.id !== id)
			if (this.activeSessionId === id) {
				this.activeSessionId = null
				this.messages = []
				this.citations = []
				this.retrievedContexts = []
			}
		},
	},
})
