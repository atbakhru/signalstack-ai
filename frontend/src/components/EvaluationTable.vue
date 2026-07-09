<template>
	<div class="panel">
		<div class="panel-body">
			<table class="table">
				<thead>
					<tr>
						<th>Question</th>
						<th>Expected Source</th>
						<th>Top-k Accuracy</th>
						<th>Retrieval MS</th>
						<th>Answer MS</th>
					</tr>
				</thead>
				<tbody>
					<tr v-for="item in results" :key="item.question">
						<td>{{ item.question }}</td>
						<td>{{ item.expectedSource }}</td>
						<td>{{ item.topKAccuracy }}%</td>
						<td>{{ item.retrievalLatencyMs }}</td>
						<td>{{ item.answerLatencyMs }}</td>
					</tr>
				</tbody>
			</table>
		</div>
	</div>
</template>

<script setup lang="ts">
defineProps<{
	results: Array<{
		question: string
		expectedSource: string
		topKAccuracy: number
		retrievalLatencyMs: number
		answerLatencyMs: number
	}>
}>()
</script>
