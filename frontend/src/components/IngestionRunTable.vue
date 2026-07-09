<template>
	<div class="panel">
		<div class="panel-body">
			<table class="table">
				<thead>
					<tr>
						<th>Source</th>
						<th>Status</th>
						<th>Documents</th>
						<th>Chunks</th>
						<th>Embeddings</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="run in runs" :key="run.id">
						<td>{{ run.source }}</td>
						<td>{{ run.status }}</td>
						<td>{{ run.documents }}</td>
						<td>{{ run.chunks }}</td>
						<td>{{ run.embeddings }}</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</template>

<script setup lang="ts">
defineProps<{
	runs: Array<{
		id: number
		source: string
		status: string
		documents: number
		chunks: number
		embeddings: number
	}>
}>()
</script>
