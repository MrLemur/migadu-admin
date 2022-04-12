<script>
	import { onMount } from 'svelte';
	import { activeDomain } from '../components/stores.js';
	import 'carbon-components-svelte/css/g80.css';
	import { loop_guard } from 'svelte/internal';

	let domains = [];

	const getDomains = async () => {
		const response = await fetch('http://localhost:5000/api/domains');
		domains = await response.json();
	};

	onMount(async () => {
		await getDomains();

		$activeDomain = $activeDomain === '' ? domains[0] : $activeDomain;
	});
</script>

<div class="container">
	<div class="title">Migadu Admin</div>
	<div class="domains-list">
		<select bind:value={$activeDomain}>
			{#each domains as domain}
				<option value={domain}>
					{domain}
				</option>
			{/each}
		</select>
	</div>
	<div class="navbar">
		<div class="navitem">
			<a href="/mailboxes">Mailboxes</a>
			<a href="/aliases">Aliases</a>
			<a href="/aliases">Redirections</a>
		</div>
	</div>
	<slot />
</div>

<style>
	.title {
		align-self: top;
		font-size: 30pt;
		font-weight: bold;
		align-items: center;
		margin-left: auto;
		margin-right: auto;
	}

	.domains-list {
		align-self: top;
		padding-top: 10px;
		font-size: 15pt;
		align-items: center;
		margin-left: auto;
		margin-right: auto;
	}

	.container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100vh;

		padding: 20px 30px;
	}

	.navbar {
		margin: 20px;

		display: flex;
		flex-direction: row;
		justify-content: space-between;
		font-size: 20pt;
	}

	.navbar > .navitem > a,
	a:visited,
	a:active {
		text-decoration: none;
	}
</style>
