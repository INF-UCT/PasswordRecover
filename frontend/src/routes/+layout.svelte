<script>
	let { children } = $props();
</script>

<svelte:head>
	<title>Cambio de contraseña</title>
</svelte:head>

<main>
	{@render children()}
</main>

<style>
	:global(*, *::before, *::after) {
		box-sizing: border-box;
		margin: 0;
		padding: 0;
	}

	:global(body) {
		font-family: system-ui, -apple-system, sans-serif;
		background: #f1f5f9;
		color: #1e293b;
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
	}

	main {
		width: 100%;
		display: flex;
		justify-content: center;
		padding: 1rem;
	}

	:global(.card) {
		background: white;
		border-radius: 0.75rem;
		box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
		padding: 2.5rem;
		width: 100%;
		max-width: 420px;
	}

	:global(h1) {
		font-size: 1.5rem;
		font-weight: 700;
		margin-bottom: 0.5rem;
	}

	:global(.subtitle) {
		color: #64748b;
		font-size: 0.9rem;
		margin-bottom: 2rem;
	}

	:global(label) {
		display: block;
		font-size: 0.875rem;
		font-weight: 500;
		margin-bottom: 0.375rem;
		color: #374151;
	}

	:global(input) {
		width: 100%;
		padding: 0.625rem 0.875rem;
		border: 1px solid #d1d5db;
		border-radius: 0.5rem;
		font-size: 1rem;
		transition: border-color 0.15s;
		outline: none;
		margin-bottom: 1.25rem;
	}

	:global(input:focus) {
		border-color: #6366f1;
		box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
	}

	:global(input:disabled) {
		background: #f9fafb;
		cursor: not-allowed;
	}

	:global(button[type='submit']) {
		width: 100%;
		padding: 0.75rem;
		background: #6366f1;
		color: white;
		border: none;
		border-radius: 0.5rem;
		font-size: 1rem;
		font-weight: 600;
		cursor: pointer;
		transition: background 0.15s;
	}

	:global(button[type='submit']:hover:not(:disabled)) {
		background: #4f46e5;
	}

	:global(button[type='submit']:disabled) {
		background: #a5b4fc;
		cursor: not-allowed;
	}

	:global(.alert) {
		padding: 0.75rem 1rem;
		border-radius: 0.5rem;
		font-size: 0.875rem;
		margin-bottom: 1.25rem;
	}

	:global(.alert-error) {
		background: #fef2f2;
		border: 1px solid #fecaca;
		color: #991b1b;
	}

	:global(.alert-success) {
		background: #f0fdf4;
		border: 1px solid #bbf7d0;
		color: #166534;
	}

	:global(.alert-info) {
		background: #eff6ff;
		border: 1px solid #bfdbfe;
		color: #1e40af;
	}
</style>
