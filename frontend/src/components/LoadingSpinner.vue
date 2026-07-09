<template>
	<div class="spinner"></div>
</template>

<style scoped>
.spinner {
	width: 22px;
	height: 22px;
	border-radius: 50%;
	border: 2px solid rgba(148, 163, 184, 0.25);
	border-top-color: #29d3a1;
	animation: spin 0.8s linear infinite;
}

@keyframes spin {
	to {
		transform: rotate(360deg);
	}
}
</style>
