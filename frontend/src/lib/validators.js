import { z } from 'zod';

export const passwordSpecialRegex = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~`]/;

export const passwordSchema = z
	.string()
	.min(1, 'La contraseña es requerida.')
	.min(8, 'Debe tener al menos 8 caracteres.')
	.regex(/\p{Ll}/u, 'Debe incluir al menos una letra minúscula.')
	.regex(/\p{Lu}/u, 'Debe incluir al menos una letra mayúscula.')
	.regex(/\d/u, 'Debe incluir al menos un número.')
	.regex(passwordSpecialRegex, 'Debe incluir al menos un símbolo especial.');

export const confirmResetSchema = z
	.object({
		token: z.string().min(1, 'El enlace no es válido.'),
		new_password: passwordSchema,
		confirmed_password: z.string().min(1, 'Debes confirmar la contraseña.')
	})
	.refine((data) => data.new_password === data.confirmed_password, {
		message: 'Las contraseñas no coinciden.',
		path: ['confirmed_password']
	});

/**
 * Convierte un ZodError en un objeto { field: message } listo para bindear al estado.
 * Si el primer issue tiene `path` con elementos, usa el primero como field.
 */
export function formatZodErrors(error) {
	const fieldErrors = {};
	if (!error || !Array.isArray(error.issues)) return fieldErrors;
	for (const issue of error.issues) {
		const path = issue.path ?? [];
		const field = path[0];
		if (field && !fieldErrors[field]) {
			fieldErrors[field] = issue.message;
		}
	}
	return fieldErrors;
}

/**
 * Convierte el array `errors[]` que devuelve el backend (formato granulado)
 * en el mismo objeto { field: message } que formatZodErrors.
 * Acepta el array o undefined.
 */
export function formatBackendErrors(errors) {
	const fieldErrors = {};
	if (!Array.isArray(errors)) return fieldErrors;
	for (const err of errors) {
		if (err?.field && !fieldErrors[err.field]) {
			fieldErrors[err.field] = err.message ?? 'Valor inválido.';
		}
	}
	return fieldErrors;
}

export const passwordRequirements = [
	{ id: 'length', label: 'Mínimo 8 caracteres', test: (pw) => pw.length >= 8 },
	{ id: 'lower', label: 'Una letra minúscula', test: (pw) => /\p{Ll}/u.test(pw) },
	{ id: 'upper', label: 'Una letra mayúscula', test: (pw) => /\p{Lu}/u.test(pw) },
	{ id: 'digit', label: 'Un número', test: (pw) => /\d/u.test(pw) },
	{
		id: 'symbol',
		label: 'Un símbolo especial (!@#$%^&*…)',
		test: (pw) => passwordSpecialRegex.test(pw)
	}
];
