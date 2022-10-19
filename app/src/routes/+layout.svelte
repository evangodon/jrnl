<script>
	import '../styles.css';
</script>

<div class="app">
<div class="top-bar"></div>

	<main class="app-container">
		<slot />
	</main>

</div>

<style>
  .top-bar {
    height: 1rem;
    background-color: var(--color-primary);
  }

  .app-container {
    margin: 0 auto;
    max-width: 38rem;
    display: flex;
    align-items: center;
    flex-direction: column;
  }
</style>
