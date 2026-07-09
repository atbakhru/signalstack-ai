<template>
	<article class="metric-card panel">
		<div class="panel-body metric-card-body">
			<div class="metric-label">{{ label }}</div>
			<div class="metric-value">{{ value }}</div>
			<div class="metric-note">{{ note }}</div>
		</div>
	</article>
</template>

<script setup lang="ts">
defineProps<{
	label: string
	value: string
	note: string
}>()
</script>

<style scoped>
.metric-label {
	color: var(--muted);
	font-size: 0.88rem;
	margin-bottom: 10px;
}

.metric-value {
	font-size: 2rem;
	font-weight: 800;
	letter-spacing: -0.03em;
}

.metric-note {
	margin-top: 8px;
	color: #cbd5e1;
}
</style>
