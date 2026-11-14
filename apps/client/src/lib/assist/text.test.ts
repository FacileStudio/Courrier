import { describe, expect, it } from 'vitest';
import { plainOf, textToHtml, upgradeLink } from './text';

describe('assist text', () => {
	it('splits paragraphs and escapes', () => {
		expect(textToHtml('Bonjour <Lisa>,\n\nMerci\nYann')).toBe('<p>Bonjour &lt;Lisa&gt;,</p><p>Merci<br>Yann</p>');
	});

	it('strips html to text and keeps the tail when clipping', () => {
		expect(plainOf('<p>Hi&nbsp;there</p><p>bye</p>')).toBe('Hi there\nbye');
		expect(plainOf('abcdef', 3)).toBe('def');
	});

	it('links the upgrade only for a plan refusal with a safe URL', () => {
		const status = { available: false, actions: [], upgrade_url: 'https://sso.facile.studio/spaces' };
		expect(upgradeLink({ ...status, reason: 'plan_required' })).toBe('https://sso.facile.studio/spaces');
		expect(upgradeLink({ ...status, reason: 'unreachable' })).toBeUndefined();
		expect(upgradeLink({ ...status, reason: 'plan_required', upgrade_url: 'javascript:alert(1)' })).toBeUndefined();
	});
});
