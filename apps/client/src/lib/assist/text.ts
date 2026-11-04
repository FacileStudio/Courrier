import type { AssistStatus } from '$lib/backend';

function escapeHtml(text: string): string {
	return text.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
}

/** Turns the helper's plain text into editor HTML: blank lines split paragraphs, single newlines break lines. */
export function textToHtml(text: string): string {
	return text
		.trim()
		.split(/\n{2,}/)
		.map((block) => `<p>${escapeHtml(block).replace(/\n/g, '<br>')}</p>`)
		.join('');
}

/** Plain text of an email body for the helper, HTML stripped and clipped. */
export function plainOf(html: string, max = 8000): string {
	const text = html
		.replace(/<(br|\/p|\/div|\/li)[^>]*>/gi, '\n')
		.replace(/<[^>]*>/g, '')
		.replace(/&nbsp;/g, ' ')
		.replace(/&lt;/g, '<')
		.replace(/&gt;/g, '>')
		.replace(/&amp;/g, '&')
		.replace(/\n{3,}/g, '\n\n')
		.trim();
	return text.length > max ? text.slice(text.length - max) : text;
}

const REASONS: Record<NonNullable<AssistStatus['reason']>, string> = {
	off: 'The Bulle helper is not set up on this server.',
	not_linked: 'Sign in with your Facile account to use Bulle here.',
	no_key: 'Add a model key in Bulle to use the helper.',
	unsupported: "Bulle's model is not supported here yet.",
	unreachable: 'Bulle is not answering right now.'
};

/** Why the helper is unavailable, in words the user can act on. */
export function reasonText(status: AssistStatus | null): string {
	if (!status) return 'Checking Bulle…';
	return status.reason ? REASONS[status.reason] : '';
}
