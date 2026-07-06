import { describe, it, expect } from 'vitest';
import {
	passwordSchema,
	confirmResetSchema,
	formatZodErrors,
	formatBackendErrors,
	passwordRequirements
} from './validators.js';

const VALID = 'ValidPass1!';

describe('passwordSchema', () => {
	it('acepta una contraseña ASCII válida', () => {
		const r = passwordSchema.safeParse(VALID);
		expect(r.success).toBe(true);
	});

	it('acepta contraseñas con caracteres españoles (ñ, á, é, ü)', () => {
		expect(passwordSchema.safeParse('MiClaveÑoño2024!').success).toBe(true);
		expect(passwordSchema.safeParse('Mañá2024!Ché').success).toBe(true);
	});

	it('acepta una contraseña con minúsculas solo Unicode', () => {
		expect(passwordSchema.safeParse('ñandú2024!A').success).toBe(true);
	});

	it('acepta exactamente 8 caracteres con todas las clases', () => {
		expect(passwordSchema.safeParse('aA1!aA1!').success).toBe(true);
	});

	it('rechaza cadena vacía con mensaje "requerida"', () => {
		const r = passwordSchema.safeParse('');
		expect(r.success).toBe(false);
		if (!r.success) {
			expect(r.error.issues[0].message).toMatch(/requerida/i);
		}
	});

	it('rechaza contraseñas de menos de 8 caracteres', () => {
		expect(passwordSchema.safeParse('Ab1!').success).toBe(false);
	});

	it('rechaza sin minúscula', () => {
		expect(passwordSchema.safeParse('NOLOWER123!').success).toBe(false);
	});

	it('rechaza sin mayúscula', () => {
		expect(passwordSchema.safeParse('noupper123!').success).toBe(false);
	});

	it('rechaza sin dígito', () => {
		expect(passwordSchema.safeParse('NoDigit!!').success).toBe(false);
	});

	it('rechaza sin símbolo', () => {
		expect(passwordSchema.safeParse('NoSymbol123').success).toBe(false);
	});

	it('no cuenta el espacio como símbolo', () => {
		expect(passwordSchema.safeParse('Abc defg1').success).toBe(false);
	});

	it('no cuenta el tab como símbolo', () => {
		expect(passwordSchema.safeParse('Abc\tdefg1').success).toBe(false);
	});

	it('acepta la comilla doble como símbolo (está en el whitelist)', () => {
		expect(passwordSchema.safeParse(`Abc"defg1`).success).toBe(true);
	});

	it('acepta la barra invertida como símbolo', () => {
		expect(passwordSchema.safeParse('Abc\\defg1').success).toBe(true);
	});

	it('acepta el backtick como símbolo', () => {
		expect(passwordSchema.safeParse('Abc`defg1').success).toBe(true);
	});

	it('acepta la tilde como símbolo', () => {
		expect(passwordSchema.safeParse('Abc~defg1').success).toBe(true);
	});

	it('rechaza solo símbolos', () => {
		const r = passwordSchema.safeParse('!!!!!!!!');
		expect(r.success).toBe(false);
		if (!r.success) {
			const codes = r.error.issues.map((i) => i.message);
			expect(codes.some((m) => /minúscula/i.test(m))).toBe(true);
			expect(codes.some((m) => /mayúscula/i.test(m))).toBe(true);
			expect(codes.some((m) => /número/i.test(m))).toBe(true);
		}
	});

	it('cada issue tiene code y message no vacío', () => {
		const r = passwordSchema.safeParse('');
		if (!r.success) {
			expect(r.error.issues.length).toBeGreaterThan(0);
			for (const issue of r.error.issues) {
				expect(issue.code).toBeTypeOf('string');
				expect(issue.code).not.toBe('');
				expect(issue.message).toBeTypeOf('string');
				expect(issue.message).not.toBe('');
			}
		}
	});

	it('issues de passwordSchema validado dentro de un objeto llevan path=new_password', () => {
		const r = confirmResetSchema.safeParse({
			token: 'abc',
			new_password: '',
			confirmed_password: 'ValidPass1!'
		});
		if (!r.success) {
			const issuesForNewPassword = r.error.issues.filter((i) => i.path[0] === 'new_password');
			expect(issuesForNewPassword.length).toBeGreaterThan(0);
		}
	});
});

describe('confirmResetSchema', () => {
	const base = { token: 'abc123', new_password: VALID, confirmed_password: VALID };

	it('pasa con datos válidos', () => {
		const r = confirmResetSchema.safeParse(base);
		expect(r.success).toBe(true);
	});

	it('falla con error en confirmed_password cuando no coinciden', () => {
		const r = confirmResetSchema.safeParse({
			...base,
			new_password: 'ValidPass1!',
			confirmed_password: 'ValidPass2!'
		});
		expect(r.success).toBe(false);
		if (!r.success) {
			const errs = formatZodErrors(r.error);
			expect(errs.confirmed_password).toMatch(/coinciden/i);
		}
	});

	it('falla con token vacío', () => {
		const r = confirmResetSchema.safeParse({ ...base, token: '' });
		expect(r.success).toBe(false);
		if (!r.success) {
			const errs = formatZodErrors(r.error);
			expect(errs.token).toBeDefined();
		}
	});

	it('falla con confirmed_password vacía', () => {
		const r = confirmResetSchema.safeParse({ ...base, confirmed_password: '' });
		expect(r.success).toBe(false);
	});
});

describe('formatZodErrors', () => {
	it('mapea el primer issue por campo', () => {
		const r = confirmResetSchema.safeParse({
			token: 'abc',
			new_password: 'short',
			confirmed_password: ''
		});
		if (!r.success) {
			const errs = formatZodErrors(r.error);
			expect(errs.new_password).toBeDefined();
			expect(errs.confirmed_password).toBeDefined();
		}
	});

	it('ignora issues con path vacío', () => {
		const errs = formatZodErrors({ issues: [{ message: 'orphan', path: [] }] });
		expect(errs).toEqual({});
	});

	it('devuelve {} con un error malformado (sin issues)', () => {
		expect(formatZodErrors(null)).toEqual({});
		expect(formatZodErrors(undefined)).toEqual({});
		expect(formatZodErrors({})).toEqual({});
	});
});

describe('formatBackendErrors', () => {
	it('mapea el array del backend a { field: message }', () => {
		const arr = [
			{ field: 'new_password', code: 'too_short', message: 'Muy corta' },
			{ field: 'confirmed_password', code: 'password_mismatch', message: 'No coinciden' }
		];
		const errs = formatBackendErrors(arr);
		expect(errs.new_password).toBe('Muy corta');
		expect(errs.confirmed_password).toBe('No coinciden');
	});

	it('devuelve {} con undefined', () => {
		expect(formatBackendErrors(undefined)).toEqual({});
	});

	it('devuelve {} con null', () => {
		expect(formatBackendErrors(null)).toEqual({});
	});

	it('devuelve {} con un valor que no es array', () => {
		expect(formatBackendErrors('not an array')).toEqual({});
	});

	it('usa fallback "Valor inválido." si el item no tiene message', () => {
		const errs = formatBackendErrors([{ field: 'new_password', code: 'x' }]);
		expect(errs.new_password).toBe('Valor inválido.');
	});

	it('ignora items sin field', () => {
		const errs = formatBackendErrors([{ code: 'x', message: 'y' }]);
		expect(errs).toEqual({});
	});

	it('toma solo el primer mensaje por campo', () => {
		const arr = [
			{ field: 'new_password', code: 'a', message: 'primero' },
			{ field: 'new_password', code: 'b', message: 'segundo' }
		];
		expect(formatBackendErrors(arr).new_password).toBe('primero');
	});
});

describe('passwordRequirements', () => {
	const find = (id) => passwordRequirements.find((r) => r.id === id);

	it('todas las reglas se cumplen para ValidPass1!', () => {
		expect(passwordRequirements.every((r) => r.test(VALID))).toBe(true);
	});

	it('length falla con menos de 8', () => {
		expect(find('length').test('Ab1!')).toBe(false);
		expect(find('length').test('12345678')).toBe(true);
	});

	it('lower acepta ñ como minúscula', () => {
		expect(find('lower').test('ñandú2024!A')).toBe(true);
	});

	it('upper acepta Ñ como mayúscula', () => {
		expect(find('upper').test('MiClaveÑoño2024!')).toBe(true);
	});

	it('digit acepta dígitos ASCII', () => {
		expect(find('digit').test('Ab1')).toBe(true);
	});

	it('symbol falla con espacio', () => {
		expect(find('symbol').test('Abc defg1')).toBe(false);
	});

	it('symbol acepta todos los caracteres del whitelist', () => {
		const sample = '!@#$%^&*()_+-=[]{};\':"\\|,.<>/?~`';
		expect(find('symbol').test(sample)).toBe(true);
	});

	it('cada requisito tiene id, label y test', () => {
		for (const r of passwordRequirements) {
			expect(r.id).toBeTypeOf('string');
			expect(r.label).toBeTypeOf('string');
			expect(r.test).toBeTypeOf('function');
		}
	});
});
