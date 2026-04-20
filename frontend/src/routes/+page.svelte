<script>
	let email = $state('');
	let loading = $state(false);
	let sent = $state(false);
	let error = $state('');

	async function handleSubmit(e) {
		e.preventDefault();
		error = '';
		loading = true;

		try {
			const res = await fetch('/api/v1/password-reset/request', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ email })
			});

			if (res.status === 400) {
				error = 'Ingresa un correo electrónico válido.';
				return;
			}
			if (res.status === 429) {
				error = 'Demasiados intentos. Espera unos minutos antes de volver a intentarlo.';
				return;
			}
			if (res.status === 500) {
				error = 'Error interno del servidor. Intenta más tarde.';
				return;
			}

			sent = true;
		} catch {
			error = 'No se pudo conectar con el servidor. Verifica tu conexión.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="card">
	<h1>Recuperar contraseña</h1>
	<p class="subtitle">
		Ingresa tu correo y te enviaremos un enlace para restablecer tu contraseña.
	</p>

	{#if sent}
		<div class="alert alert-success">
			Si el correo existe en el sistema, recibirás un enlace para cambiar tu contraseña en los
			próximos minutos.
		</div>
		<button type="button" class="btn-link" onclick={() => { sent = false; email = ''; }}>
			Enviar a otro correo
		</button>
	{:else}
		{#if error}
			<div class="alert alert-error">{error}</div>
		{/if}

		<form onsubmit={handleSubmit}>
			<label for="email">Correo electrónico</label>
			<input
				id="email"
				type="email"
				bind:value={email}
				placeholder="usuario@dominio.cl"
				required
				disabled={loading}
				autocomplete="email"
			/>

			<button type="submit" disabled={loading}>
				{loading ? 'Enviando...' : 'Enviar enlace'}
			</button>
		</form>
	{/if}
</div>

<style>
	.btn-link {
		width: 100%;
		background: none;
		border: none;
		color: #6366f1;
		font-size: 0.875rem;
		cursor: pointer;
		text-decoration: underline;
		padding: 0.5rem 0;
	}
</style>
