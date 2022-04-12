<script>
	import { createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';

	export let title;
	export let formFields;

	const dispatch = createEventDispatcher();

	const onSubmit = (event) => {
		const formData = new FormData(event.target);
		const data = {};
		for (const [key, value] of formData.entries()) {
			data[key] = value;
		}
		dispatch('formSubmitted', data);
	};

	const onCancel = () => {
		dispatch('formCancelled', {});
	};
</script>

<div transition:fade={{ duration: 300 }} class="overlay">
	<div class="modal">
		<div class="title">
			{title}
		</div>
		<form on:submit|preventDefault={onSubmit}>
			<div class="form-group">
				{#each formFields as field}
					<label for={field.name}>{field.label}</label>
					<input
						type={field.type}
						id={field.name}
						name={field.name}
						placeholder={field.placeholder}
					/>
				{/each}
				<button type="submit">Submit</button>
				<p class="form-cancel" on:click={onCancel}>Cancel</p>
			</div>
		</form>
	</div>
</div>

<style>
	.overlay {
		position: absolute;
		top: 0;
		left: 0;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		width: 100%;
		height: 100vh;
		background-color: rgba(0, 0, 0, 0.5);
		z-index: 1000;
	}

	.modal {
		display: flex;

		padding: 10px 20px;
		border-radius: 1em;
		background-color: gray;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		max-width: 600px;
	}
	.title {
		font-size: 1.5em;
		font-weight: bold;
		margin-bottom: 20px;
	}

	.form-group {
		display: flex;
		flex-direction: column;
		align-items: flex-start;
		justify-content: flex-start;
		margin-bottom: 10px;
	}

	.form-group label {
		font-size: 1em;
		font-weight: bold;
		margin-bottom: 10px;
	}

	.form-group input {
		background-color: transparent;
		color: white;
		border: none;
		border-bottom: rgba(0, 0, 0, 0.5) solid 1px;
		font-size: 1em;
		margin-bottom: 20px;
	}

	.form-group button {
		background-color: #0091b1;
		color: white;
		border: none;
		padding: 10px 20px;
		font-size: 1em;
		margin-bottom: 5px;
		margin: 10px auto;
	}

	.form-group button:hover {
		background-color: #41b4ce;
		cursor: pointer;
	}

	::placeholder {
		color: rgba(255, 255, 255, 0.5);
	}

	.form-cancel {
		font-size: 1em;
		margin: auto;
		margin-top: -5px;
		margin-bottom: 0;
	}

	.form-cancel:hover {
		cursor: pointer;
		font-weight: bold;
		color: rgb(255, 108, 108);
		text-decoration: underline;
	}

	textarea:foxus,
	input:focus {
		outline: none;
	}
</style>
