import { t } from '$lib/i18n/index.svelte';

/** Registre's plans page: where people subscribe and manage a plan. */
export const SPACES_URL = 'https://sso.facile.studio/spaces';

/** A 402 refusal: the scope's plan has no room for what was being created. */
export type PlanLimit = {
	message: string;
	upgradeUrl?: string;
	key?: string;
	limit?: number;
};

/**
 * True for a same-origin path or an https URL with a host and no credentials. Anything
 * else, such as `javascript:` or a protocol-relative URL, never reaches an `href`.
 */
export function isSafeUpgradeUrl(raw: unknown): raw is string {
	if (typeof raw !== 'string' || raw === '' || /[\s\\\p{Cc}]/u.test(raw)) return false;
	if (raw.startsWith('/')) return !raw.startsWith('//');
	try {
		const url = new URL(raw);
		return url.protocol === 'https:' && url.host !== '' && !url.username && !url.password;
	} catch {
		return false;
	}
}

/**
 * Reads the suite error envelope of a 402 into a `PlanLimit`, or `undefined` for any other
 * status. An unsafe `upgrade_url` is dropped, never rendered.
 */
export function planLimitFrom(status: number, payload: unknown): PlanLimit | undefined {
	if (status !== 402) return undefined;
	const error = (payload as { error?: Record<string, unknown> } | undefined)?.error ?? {};
	return {
		message: typeof error.message === 'string' && error.message ? error.message : t('plan.fallback'),
		upgradeUrl: isSafeUpgradeUrl(error.upgrade_url) ? error.upgrade_url : undefined,
		key: typeof error.key === 'string' ? error.key : undefined,
		limit: typeof error.limit === 'number' && Number.isFinite(error.limit) ? error.limit : undefined
	};
}
