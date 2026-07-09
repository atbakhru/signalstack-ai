<template>
	<section class="stack">
		<div class="panel">
			<div class="panel-body stack">
				<div class="row">
					<div>
						<div class="eyebrow">Documents</div>
						<h1 class="title-xl">Indexed content browser</h1>
					</div>
					<div class="row" style="min-width: 320px">
						<input v-model="search" class="input" type="search" placeholder="Search by title, source, or date" @input="refresh" />
						<select v-model.number="sourceId" class="select" @change="refresh">
							<option :value="0">All sources</option>
							<option v-for="source in sources" :key="source.id" :value="source.id">{{ source.name }}</option>
						</select>
					</div>
				</div>
			</div>
		</div>

		<div class="panel">
			<div class="panel-body">
				<table class="table">
					<thead>
						<tr>
							<th>Title</th>
							<th>Source</th>
							<th>Published</th>
							<th>URL</th>
							<th></th>
						</tr>
					</thead>
					<tbody>
						<tr v-for="document in documents" :key="document.id">
							<td>{{ document.data?.title ?? document.title }}</td>
							<td>{{ document.source_name ?? document.source }}</td>
							<td>{{ formatDate(document.data?.published_at ?? document.published_at) }}</td>
							<td><a :href="document.data?.url ?? document.url" target="_blank" rel="noreferrer">Open</a></td>
							<td><button class="button secondary" type="button" @click="loadChunks(document.id)">Chunks</button></td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>

		<div class="page-grid cols-2">
			<div class="panel">
				<div class="panel-body stack">
					<div class="row"><h2 class="title-lg">Document chunks</h2><span class="pill">{{ selectedDocumentTitle || 'Select a document' }}</span></div>
					<div v-for="chunk in chunks" :key="chunk.id" class="panel" style="background: rgba(8, 17, 31, 0.6)">
						<div class="panel-body stack">
							<div class="row">
								<strong>Chunk {{ chunk.chunk_index }}</strong>
								<span class="pill">{{ chunk.token_count }} tokens</span>
							</div>
							<p class="subtle">{{ chunk.content }}</p>
						</div>
					</div>
				</div>
			</div>
		</div>
	</section>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { getDocumentChunks, listDocuments, type DocumentItem, type DocumentChunk } from '../api/documents'
import { listSources, type Source } from '../api/ingestion'

const documents = ref<DocumentItem[]>([])
const sources = ref<Source[]>([])
const chunks = ref<DocumentChunk[]>([])
const selectedDocumentTitle = ref('')
const search = ref('')
const sourceId = ref(0)

function formatDate(value?: string) {
	return value ? new Date(value).toLocaleDateString() : '—'
}

async function refresh() {
	const response = await listDocuments({
		search: search.value || undefined,
		sourceId: sourceId.value || undefined,
	})
	documents.value = response.documents
}

async function loadChunks(documentId: number) {
	const document = documents.value.find(item => item.id === documentId)
	selectedDocumentTitle.value = document?.data?.title ?? document?.title ?? ''
	const response = await getDocumentChunks(documentId)
	chunks.value = response.chunks
}

onMounted(async () => {
	const sourceResponse = await listSources()
	sources.value = sourceResponse.sources
	await refresh()
})
</script>
