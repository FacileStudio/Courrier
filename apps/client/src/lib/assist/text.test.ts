import { describe, expect, it } from 'vitest';
import { plainOf, textToHtml } from './text';

describe('assist text', () => {
	it('splits paragraphs and escapes', () => {
		expect(textToHtml('Bonjour <Lisa>,\n\nMerci\nYann')).toBe('<p>Bonjour &lt;Lisa&gt;,</p><p>Merci<br>Yann</p>');
	});

	it('strips html to text and keeps the tail when clipping', () => {
		expect(plainOf('<p>Hi&nbsp;there</p><p>bye</p>')).toBe('Hi there\nbye');
		expect(plainOf('abcdef', 3)).toBe('def');
	});
});
