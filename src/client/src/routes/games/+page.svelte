<script lang="ts">
	import { onMount } from 'svelte';
	import type { Playlist } from '../../lib/types';
	import * as requests from '../../lib/axios';
	import { ensureToken } from '../../lib/auth';

	let playlists: Playlist[];
	let selectedPlaylistId: string;

	onMount(async () => {
		ensureToken();
		playlists = await requests.getPlaylists();
		selectedPlaylistId = playlists[0].id;
	});

	const createGame = async () => {
		const game = await requests.createGame(selectedPlaylistId);
		window.location.href = `/games/${game.id}`;
	};
</script>

<select>
	{#each playlists as playlist}
		<option onselect={() => (selectedPlaylistId = playlist.id)}>{playlist.name}</option>
	{/each}
</select>
<button onclick={() => createGame()}>Create game</button>
