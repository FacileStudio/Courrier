import type { AssistStatus } from '$lib/backend';
import { t } from '$lib/i18n/index.svelte';
import type { AppKey } from '$lib/i18n/messages';
import { isSafeUpgradeUrl } from '$lib/plan-limit';

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

const REASONS: Record<NonNullable<AssistStatus['reason']>, AppKey> = {
	off: 'assist.off',
	plan_required: 'assist.planRequired',
	not_linked: 'assist.notLinked',
	unreachable: 'assist.unreachable'
};

/** Why the helper is unavailable, in words the user can act on. */
export function reasonText(status: AssistStatus | null): string {
	if (!status) return t('assist.checking');
	return status.reason ? t(REASONS[status.reason]) : '';
}

/** Where to get a plan with AI access, when the plan is why the helper is unavailable. */
export function upgradeLink(status: AssistStatus | null): string | undefined {
	if (status?.reason !== 'plan_required') return undefined;
	return isSafeUpgradeUrl(status.upgrade_url) ? status.upgrade_url : undefined;
}
