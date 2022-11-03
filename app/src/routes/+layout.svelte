<script>
	import { SvelteToast } from '@zerodevx/svelte-toast';
	import '../styles.css';

	export const toastOptions = {
		theme: {
			'--toastBarHeight': 0
		}
	};
</script>

<div class="app">
	<main class="app-container">
		<slot />
	</main>
	<SvelteToast options={toastOptions} />
</div>

<style>
	.app-container {
		margin: 0 auto;
		max-width: 38rem;
		display: flex;
		align-items: center;
		flex-direction: column;
	}
</style>
