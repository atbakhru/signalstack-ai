<template>
	<div class="panel">
		<div class="panel-body">
			<table class="table">
				<thead>
					<tr>
						<th>Title</th>
						<th>Source</th>
						<th>Published</th>
						<th>URL</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="document in documents" :key="document.id">
						<td>{{ document.title }}</td>
						<td>{{ document.source }}</td>
						<td>{{ document.publishedAt }}</td>
						<td>{{ document.url }}</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</template>

<script setup lang="ts">
defineProps<{
	documents: Array<{
		id: number
		title: string
		source: string
		publishedAt: string
		url: string
	}>
}>()
</script>
