<template>
	<article class="panel">
		<div class="panel-body stack">
			<div class="row">
				<div>
					<div class="title-lg">{{ name }}</div>
					<div class="subtle">{{ apiName }}</div>
				</div>
				<span class="pill" :class="enabled ? 'success' : 'danger'">
					{{ enabled ? 'Enabled' : 'Disabled' }}
				</span>
			</div>
			<div class="subtle">{{ description }}</div>
		</div>
	</article>
</template>

<script setup lang="ts">
defineProps<{
	name: string
	apiName: string
	description: string
	enabled: boolean
}>()
</script>
