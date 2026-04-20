<script>
	import { page } from '$app/stores';

	const token = $derived($page.params.token);

	let newPassword = $state('');
	let confirmedPassword = $state('');
	let loading = $state(false);
	let validating = $state(true);
	let tokenValid = $state(false);
	let expiresAt = $state('');
	let done = $state(false);
	let error = $state('');

	$effect(() => {
		if (token) validateToken();
	});

	async function validateToken() {
		validating = true;
		error = '';
		try {
			const res = await fetch(`/api/v1/password-reset/${token}`);
			if (res.ok) {
				const body = await res.json();
				tokenValid = true;
				expiresAt = body.data?.expires_at ?? '';
			} else if (res.status === 410) {
				error = 'Este enlace ya fue utilizado o ha vencido.';
			} else if (res.status === 404) {
				error = 'El enlace no es válido.';
			} else {
				error = 'Error al verificar el enlace. Intenta más tarde.';
			}
		} catch {
			error = 'No se pudo conectar con el servidor.';
		} finally {
			validating = false;
		}
	}

	async function handleSubmit(e) {
		e.preventDefault();
		error = '';

		if (newPassword !== confirmedPassword) {
			error = 'Las contraseñas no coinciden.';
			return;
		}

		loading = true;
		try {
			const res = await fetch('/api/v1/password-reset/confirm', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					token,
					new_password: newPassword,
					confirmed_password: confirmedPassword
				})
			});

			if (res.ok) {
				done = true;
				return;
			}

			const body = await res.json().catch(() => ({}));

			if (res.status === 400) {
				error = 'Los datos enviados no son válidos.';
			} else if (res.status === 410) {
				error = 'El enlace ha vencido. Solicita uno nuevo.';
				tokenValid = false;
			} else if (res.status === 422) {
				error =
					body.detail ??
					'La contraseña no cumple los requisitos de seguridad. Usa al menos 8 caracteres con mayúsculas, minúsculas y números.';
			} else if (res.status === 429) {
				error = 'Demasiados intentos. Espera unos minutos.';
			} else {
				error = 'Error interno. Intenta más tarde.';
			}
		} catch {
			error = 'No se pudo conectar con el servidor.';
		} finally {
			loading = false;
		}
	}

	function formatExpiry(iso) {
		if (!iso) return '';
		return new Date(iso).toLocaleTimeString('es-CL', { hour: '2-digit', minute: '2-digit' });
	}
</script>

<div class="card">
	<h1>Nueva contraseña</h1>

	{#if validating}
		<p class="subtitle">Verificando enlace...</p>
	{:else if done}
		<div class="alert alert-success">
			¡Contraseña actualizada correctamente! Ya puedes iniciar sesión con tu nueva contraseña.
		</div>
	{:else if !tokenValid}
		<div class="alert alert-error">{error}</div>
		<a href="/" class="btn-back">Solicitar nuevo enlace</a>
	{:else}
		<p class="subtitle">
			Ingresa tu nueva contraseña.
			{#if expiresAt}
				El enlace vence a las {formatExpiry(expiresAt)}.
			{/if}
		</p>

		{#if error}
			<div class="alert alert-error">{error}</div>
		{/if}

		<form onsubmit={handleSubmit}>
			<label for="new-password">Nueva contraseña</label>
			<input
				id="new-password"
				type="password"
				bind:value={newPassword}
				placeholder="Mínimo 8 caracteres"
				required
				minlength="8"
				disabled={loading}
				autocomplete="new-password"
			/>

			<label for="confirm-password">Confirmar contraseña</label>
			<input
				id="confirm-password"
				type="password"
				bind:value={confirmedPassword}
				placeholder="Repite la contraseña"
				required
				minlength="8"
				disabled={loading}
				autocomplete="new-password"
			/>

			<button type="submit" disabled={loading || !newPassword || !confirmedPassword}>
				{loading ? 'Guardando...' : 'Cambiar contraseña'}
			</button>
		</form>
	{/if}
</div>

<style>
	.btn-back {
		display: block;
		text-align: center;
		margin-top: 1rem;
		color: #6366f1;
		font-size: 0.875rem;
		text-decoration: underline;
	}
</style>
