import { describe, expect, it } from 'vitest';
import { routeFor } from './push-route';

describe('routeFor', () => {
	it('opens mail paths and falls back to the mailbox', () => {
		expect(routeFor('mail')).toBe('/mail');
		expect(routeFor('mail/inbox')).toBe('/mail/inbox');
		expect(routeFor(undefined)).toBe('/mail');
		expect(routeFor('//evil.example')).toBe('/mail');
		expect(routeFor('mail/../settings')).toBe('/mail');
	});
});
