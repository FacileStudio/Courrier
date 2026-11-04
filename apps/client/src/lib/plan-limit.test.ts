import { describe, expect, it } from 'vitest';
import { isSafeUpgradeUrl, planLimitFrom } from './plan-limit';

describe('planLimitFrom', () => {
	it('reads the 402 envelope', () => {
		const body = {
			error: {
				code: 'payment_required',
				message: 'limit reached',
				upgrade_url: 'https://sso.facile.studio/spaces/spc_1',
				key: 'courrier.mailboxes',
				limit: 1
			}
		};
		expect(planLimitFrom(402, body)).toEqual({
			message: 'limit reached',
			upgradeUrl: 'https://sso.facile.studio/spaces/spc_1',
			key: 'courrier.mailboxes',
			limit: 1
		});
	});

	it('ignores every other status', () => {
		expect(planLimitFrom(403, { error: { message: 'no' } })).toBeUndefined();
	});

	it('drops an unsafe upgrade link', () => {
		for (const link of ['javascript:alert(1)', 'http://x.test', '//evil.test', 'https://u:p@x.test', ' https://x.test']) {
			expect(planLimitFrom(402, { error: { upgrade_url: link } })?.upgradeUrl).toBeUndefined();
		}
	});
});

describe('isSafeUpgradeUrl', () => {
	it('accepts https and same-origin paths only', () => {
		expect(isSafeUpgradeUrl('https://sso.facile.studio/spaces')).toBe(true);
		expect(isSafeUpgradeUrl('/spaces')).toBe(true);
		expect(isSafeUpgradeUrl('data:text/html,x')).toBe(false);
		expect(isSafeUpgradeUrl(undefined)).toBe(false);
	});
});
