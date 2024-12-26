<script lang="ts">
	import { onMount } from 'svelte';
	import { ensureToken } from '../lib/auth';

	onMount(() => {
		ensureToken();
	});
</script>

<head>
	<link
		href="https://unpkg.com/material-components-web@latest/dist/material-components-web.min.css"
		rel="stylesheet"
	/>
	<script
		src="https://unpkg.com/material-components-web@latest/dist/material-components-web.min.js"
	></script>
</head>

<p>Hello</p>
