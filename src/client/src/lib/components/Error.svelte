<script lang="ts">
	import { onDestroy } from 'svelte';
	import { errorStore } from '../axios';

	let errorMessage: string | null = null;
	let modelOpen = false;

	const unsubscribe = errorStore.subscribe((v) => {
		if (!v) {
			return;
		}

		errorMessage = v;
		modelOpen = true;
	});

	onDestroy(() => {
		unsubscribe();
	});
</script>

<input type="checkbox" id="error-modal-control" class="modal" bind:checked={modelOpen} />

{#if modelOpen}
	<div>
		<div class="card">
			<label for="error-modal-control" class="modal-close" onclick={() => (modelOpen = false)}
			></label>
			<h3 class="section">Modal</h3>
			<p class="section">{errorMessage}</p>
		</div>
	</div>
{/if}
