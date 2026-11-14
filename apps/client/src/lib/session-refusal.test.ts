import { describe, expect, it } from 'vitest';
import { OWN_SESSION_REFUSAL, refusesOwnSession } from './session-refusal';

describe('refusesOwnSession', () => {
	it('matches porte refusing a session no browser sign-in opened', () => {
		expect(refusesOwnSession({ status: 403, code: 'permission_denied', message: OWN_SESSION_REFUSAL })).toBe(true);
	});

	it('leaves any other 403 to its own message', () => {
		expect(refusesOwnSession({ status: 403, code: 'permission_denied', message: 'access denied' })).toBe(false);
		expect(refusesOwnSession({ status: 403, code: 'forbidden', message: OWN_SESSION_REFUSAL })).toBe(false);
		expect(refusesOwnSession({ status: 401, code: 'permission_denied', message: OWN_SESSION_REFUSAL })).toBe(false);
		expect(refusesOwnSession(null)).toBe(false);
	});
});
