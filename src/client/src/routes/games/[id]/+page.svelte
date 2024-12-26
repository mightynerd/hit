<script lang="ts">
	import { onMount } from 'svelte';
	import { advanceGame } from '../../../lib/axios';
	import type { Track } from '../../../lib/types.js';
	import { ensureToken } from '../../../lib/auth';

	export let data;
	const { params } = data;
	let gameId = params.id;

	let state: 'init' | 'hidden' | 'visible' = 'init';
	let currentTrack: Track;

	onMount(() => {
		ensureToken();
	});

	const advance = async () => {
		state = 'hidden';
		currentTrack = await advanceGame(gameId);
	};

	const show = async () => {
		state = 'visible';
	};
</script>

<div class="container">
	<div class="row">
		<div class="card" style="height: 30em; width: 20em; vertical-align: middle;">
			{#if state === 'init'}
				<button onclick={() => advance()}>Start</button>
			{/if}

			{#if state === 'hidden'}
				<img onclick={() => show()} src="/robot.gif" alt="robot" />
				<button onclick={() => show()}>Show</button>
			{/if}

			{#if state === 'visible'}
				<table class="horizontal">
					<thead>
						<tr>
							<th>Title</th>
							<th>Artist</th>
							<th>Year</th>
						</tr>
					</thead>
					<tbody>
						<tr>
							<td>{currentTrack.title}</td>
							<td>{currentTrack.artist}</td>
							<td>{currentTrack.year}</td>
						</tr>
					</tbody>
				</table>

				<button onclick={() => advance()}>Play next</button>
			{/if}
		</div>
	</div>
</div>
