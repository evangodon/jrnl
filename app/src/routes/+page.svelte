<script lang="ts">
	import { api } from '$lib/api';

	import type { PageData } from './$types';

	export let data: PageData;
	export let numRows = data.daily.content.split('\n').length + 1;

	async function handleSubmit(event) {
		const formData = new FormData(this);

		try {
			// New entry
			if (data.daily.id == -1) {
				await api.post('daily/new', {
					body: JSON.stringify({
						id: data.daily.id,
						content: formData.get('content')
					})
				});

				return;
			}

			// Update entry
			await api.patch(`daily/${data.daily.id}`, {
				body: JSON.stringify({
					id: data.daily.id,
					content: formData.get('content')
				})
			});
		} catch (err) {
			console.error(err);
		}
	}
</script>

<form class="form" method="POST" on:submit|preventDefault={handleSubmit}>
	<textarea name="content" class="daily-input" value={data.daily.content} rows={numRows} />
	<button class="submit-btn" type="submit">{data.daily.id == -1 ? 'Create' : 'Update'}</button>
</form>

<style>
	.form {
		width: 100%;
		display: flex;
		align-items: center;
		flex-direction: column;
	}
	.daily-input {
		margin: 0 1rem;
		margin-top: 3rem;
		width: 95%;
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
		padding: 0.5rem 2rem;
		border-radius: 4px;
	}
</style>
