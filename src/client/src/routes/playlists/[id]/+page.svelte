<script lang="ts">
	import { onMount } from 'svelte';
	import type { Track } from '../../../lib/types.js';
	import { ensureToken } from '../../../lib/auth';
	import * as requests from '../../../lib/axios.js';

	export let data;
	const { params } = data;
	let tracks: Track[] = [];
	let playlistId = params.id;
	const size = 50;
	let page = 0;

	const updateTracks = async () => {
		tracks = await requests.getTracks(playlistId, page, size).catch(() => (tracks = []));
	};

	const nextPage = async () => {
		page++;
		await updateTracks();
	};

	const previousPage = async () => {
		page--;
		await updateTracks();
	};

	const deleteTrack = async (trackId: string) => {
		await requests.deleteTrack(playlistId, trackId);
		await updateTracks();
	};

	onMount(async () => {
		ensureToken();
		await updateTracks();
	});
</script>

<div class="container">
	<div class="row">
		<div class="col-sm">
			<table style="max-height: fit-content; overflow-x: hidden;">
				<caption>Tracks</caption>
				<thead>
					<tr>
						<th>Title</th>
						<th>Artist</th>
						<th>Year</th>
						<th></th>
					</tr>
				</thead>
				<tbody>
					{#each tracks as track}
						<tr>
							<td>{track.title}</td>
							<td>{track.artist}</td>
							<td>{track.year}</td>
							<td data-label=""
								><button class="secondary" onclick={() => deleteTrack(track.id)}>Delete</button></td
							>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	</div>
	<div class="row">
		<div class="col-sm-6 col-md-9 col-lg-10"></div>
		<div class="col-sm-5 col-md-3 col-lg-2">
			<button onclick={() => previousPage()}>Page {page}</button>
			<button onclick={() => nextPage()}>Page {page + 1}</button>
		</div>
	</div>
</div>
