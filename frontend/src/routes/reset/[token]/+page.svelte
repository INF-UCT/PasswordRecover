<script>
	import { page } from '$app/stores';
	import {
		confirmResetSchema,
		formatZodErrors,
		formatBackendErrors,
		passwordRequirements
	} from '$lib/validators.js';

	const token = $derived($page.params.token);

	let newPassword = $state('');
	let confirmedPassword = $state('');
	let loading = $state(false);
	let validating = $state(true);
	let tokenValid = $state(false);
	let expiresAt = $state('');
	let done = $state(false);
	let fieldErrors = $state({});
	let generalError = $state('');

	const requirementsState = $derived(
		passwordRequirements.map((req) => ({ ...req, met: req.test(newPassword) }))
	);

	$effect(() => {
		if (token) validateToken();
	});

	async function validateToken() {
		validating = true;
		generalError = '';
		try {
			const res = await fetch(`/api/v1/password-reset/${token}`);
			if (res.ok) {
				const body = await res.json();
				tokenValid = true;
				expiresAt = body.data?.expires_at ?? '';
			} else if (res.status === 410) {
				generalError = 'Este enlace ya fue utilizado o ha vencido.';
			} else if (res.status === 404) {
				generalError = 'El enlace no es válido.';
			} else {
				generalError = 'Error al verificar el enlace. Intenta más tarde.';
			}
		} catch {
			generalError = 'No se pudo conectar con el servidor.';
		} finally {
			validating = false;
		}
	}

	async function handleSubmit(e) {
		e.preventDefault();
		fieldErrors = {};
		generalError = '';

		const parsed = confirmResetSchema.safeParse({
			token,
			new_password: newPassword,
			confirmed_password: confirmedPassword
		});

		if (!parsed.success) {
			fieldErrors = formatZodErrors(parsed.error);
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
				const backendErrors = formatBackendErrors(body.errors);
				if (Object.keys(backendErrors).length > 0) {
					fieldErrors = backendErrors;
				} else {
					generalError = 'Los datos enviados no son válidos.';
				}
			} else if (res.status === 404) {
				generalError = 'El enlace no es válido o ya fue utilizado.';
			} else if (res.status === 410) {
				generalError = 'El enlace ha vencido. Solicita uno nuevo.';
				tokenValid = false;
			} else if (res.status === 422) {
				generalError =
					body.detail ??
					'La contraseña no cumple los requisitos de seguridad.';
			} else if (res.status === 429) {
				generalError = 'Demasiados intentos. Espera unos minutos.';
			} else {
				generalError = 'Error interno. Intenta más tarde.';
			}
		} catch {
			generalError = 'No se pudo conectar con el servidor.';
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
		<div class="alert alert-error">{generalError}</div>
		<a href="/" class="btn-back">Solicitar nuevo enlace</a>
	{:else}
		<p class="subtitle">
			Ingresa tu nueva contraseña.
			{#if expiresAt}
				El enlace vence a las {formatExpiry(expiresAt)}.
			{/if}
		</p>

		{#if generalError}
			<div class="alert alert-error">{generalError}</div>
		{/if}

		<form onsubmit={handleSubmit} novalidate>
			<label for="new-password">Nueva contraseña</label>
			<input
				id="new-password"
				type="password"
				bind:value={newPassword}
				placeholder="Ingresa tu nueva contraseña"
				disabled={loading}
				autocomplete="new-password"
				aria-invalid={fieldErrors.new_password ? 'true' : 'false'}
				aria-describedby={fieldErrors.new_password ? 'new-password-error' : 'password-requirements'}
			/>
			{#if fieldErrors.new_password}
				<small id="new-password-error" class="field-error">{fieldErrors.new_password}</small>
			{:else}
				<ul id="password-requirements" class="requirements">
					{#each requirementsState as req (req.id)}
						<li class:met={req.met} class:unmet={!req.met}>
							<span class="icon">{req.met ? '✓' : '○'}</span>
							{req.label}
						</li>
					{/each}
				</ul>
			{/if}

			<label for="confirm-password">Confirmar contraseña</label>
			<input
				id="confirm-password"
				type="password"
				bind:value={confirmedPassword}
				placeholder="Repite la contraseña"
				disabled={loading}
				autocomplete="new-password"
				aria-invalid={fieldErrors.confirmed_password ? 'true' : 'false'}
				aria-describedby={fieldErrors.confirmed_password ? 'confirm-password-error' : undefined}
			/>
			{#if fieldErrors.confirmed_password}
				<small id="confirm-password-error" class="field-error"
					>{fieldErrors.confirmed_password}</small
				>
			{/if}

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

	.field-error {
		display: block;
		margin-top: 0.25rem;
		color: #dc2626;
		font-size: 0.8125rem;
	}

	.requirements {
		list-style: none;
		padding: 0;
		margin: 0.5rem 0 1rem;
		font-size: 0.8125rem;
		display: grid;
		gap: 0.25rem;
	}

	.requirements li {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		color: #6b7280;
	}

	.requirements li.met {
		color: #16a34a;
	}

	.requirements li.unmet {
		color: #9ca3af;
	}

	.requirements .icon {
		font-weight: 700;
		width: 1em;
		display: inline-block;
		text-align: center;
	}
</style>
