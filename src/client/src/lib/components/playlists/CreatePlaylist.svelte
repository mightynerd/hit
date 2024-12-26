<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { createPlaylist } from '../../axios';

	let name: string;
	let spotifyId: string;

	let waiting = false;
	let progress = 0;
	export let modalOpen = false;

	export let onCreated = () => {};

	const onSubmit = async (event: SubmitEvent) => {
		event.preventDefault();
		waiting = true;

		const interval = setInterval(() => (progress = Math.random() * 100), 1000);

		await createPlaylist(name, spotifyId);

		clearInterval(interval);
		progress = 100;
		waiting = false;
		modalOpen = false;
		onCreated();
	};
</script>

<div>
	<input type="checkbox" id="modal-control" class="modal" bind:checked={modalOpen} />

	{#if modalOpen}
		<div role="dialog" aria-labelledby="dialog-title">
			<div class="card">
				<label for="modal-control" class="modal-close" onclick={() => (modalOpen = false)}></label>

				<h3 class="section" id="dialog-title">Create playlist</h3>
				<form onsubmit={onSubmit}>
					<label for="name">Name</label>
					<input type="text" id="name" placeholder="My playlist" bind:value={name} required />

					<label for="spotify_id">Spotify ID</label>
					<input
						type="text"
						id="spotify_id"
						placeholder="2Ak3dxncNLLGQW2A9l3tI1"
						bind:value={spotifyId}
						required
					/>

					<input type="submit" value="Create" />
				</form>

				{#if waiting}
					<progress value={progress} max="100"></progress>
				{/if}
			</div>
		</div>
	{/if}
</div>
