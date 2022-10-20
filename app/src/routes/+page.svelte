<script lang="ts">
	import { toast } from '@zerodevx/svelte-toast';
	import { api } from '$lib/api';

	import type { PageData } from './$types';
	import { invalidate } from '$app/navigation';

	export let data: PageData;
	export let newContent = data.daily.content;
	export let numRows = data.daily.content.split('\n').length + 1;

	async function handleSubmit() {
		const formData = new FormData(this);

		try {
			// New entry
			if (data.daily.id === '-1') {
				await api.post('daily/new', {
					body: JSON.stringify({
						content: formData.get('content')
					})
				});

				toast.push(' ✅ Entry created for today');
				await invalidate('daily-entry:new');

				return;
			}

			// Update entry
			await api.patch(`daily/${data.daily.id}`, {
				body: JSON.stringify({
					id: data.daily.id,
					content: formData.get('content')
				})
			});
			toast.push(' ✅ Entry updated');
		} catch (err) {
			toast.push(' ❗ An unexpected error occured');
			console.error(err);
		}
	}
</script>

<form class="form" method="POST" on:submit|preventDefault={handleSubmit}>
	<textarea name="content" class="daily-input" bind:value={newContent} rows={numRows} />
	<button class="submit-btn" type="submit" disabled={data.daily.content === newContent}>
		{data.daily.id === '-1' ? 'Create' : 'Update'}
	</button>
</form>

<style>
	.form {
		width: 95%;
		display: flex;
		align-items: center;
		flex-direction: column;
	}
	.daily-input {
		margin: 0 1rem;
		margin-top: 3rem;
		width: 100%;
		min-height: 50vh;
		padding: 0.8rem;
		border: 1px solid transparent;
		background-color: #eff1f5;
		flex: 1;
		border-radius: 4px;
		font-size: 18px;
		resize: none;
	}

	.submit-btn {
		margin-top: 1rem;
		border: 2px solid transparent;
		background-color: var(--color-primary);
		color: white;
		font-weight: 600;
		padding: 0.5rem 2rem;
		border-radius: 4px;
	}
	.submit-btn:disabled {
		opacity: 75%;
	}

	@media only screen and (max-width: 600px) {
		.submit-btn {
			font-size: 22px;
			width: 100%;
		}
	}
</style>
