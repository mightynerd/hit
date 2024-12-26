<script lang="ts">
	import { onMount } from 'svelte';
	import { ensureToken } from '../../lib/auth';
	import * as requests from '../../lib/axios';
	import type { Playlist } from '../../lib/types';
	import CreatePlaylist from '../../lib/components/playlists/CreatePlaylist.svelte';

	let playlists: Playlist[];
	let isCreateModelOpen = false;

	const fetchPlaylists = async () => {
		const result = await requests.getPlaylists();
		playlists = result;
	};

	onMount(() => {
		ensureToken();
		fetchPlaylists();
	});

	const onPlaylistCreated = async () => {
		await fetchPlaylists();
	};

	const deletePlaylist = async (playlistId: string) => {
		await requests.deletePlaylist(playlistId);
		await fetchPlaylists();
	};
</script>

<div class="container">
	<div class="row">
		<div class="col-sm-9 col-md-10 col-lg-11"></div>
		<div class="col-sm-3 col-md-2 col-lg-1">
			<label class="button doc" onclick={() => (isCreateModelOpen = true)}> Create playlist </label>
			<CreatePlaylist bind:modalOpen={isCreateModelOpen} onCreated={() => onPlaylistCreated()} />
		</div>
	</div>
	<div class="row">
		<div class="col-sm">
			<table style="max-height: fit-content; overflow-x: hidden;">
				<thead>
					<tr>
						<th>Name</th>
						<th>Created</th>
						<th></th>
					</tr>
				</thead>

				<tbody>
					{#each playlists as playlist}
						<tr>
							<td data-label="Name"><a href="/playlists/{playlist.id}">{playlist.name}</a></td>
							<td data-label="Created">{playlist.created_at}</td>
							<td data-label=""
								><button class="secondary" onclick={() => deletePlaylist(playlist.id)}
									>Delete</button
								></td
							>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
</div>
