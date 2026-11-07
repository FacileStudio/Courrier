import type { AppKey } from '$lib/i18n/messages';

export type RoleTone = 'owner' | 'admin' | 'neutral';

const labels: Record<string, AppKey> = {
	owner: 'spaces.roleOwner',
	admin: 'spaces.roleAdmin',
	member: 'spaces.roleMember'
};

/** The message key of a role's name, or null for a role this app does not name. */
export function roleLabelKey(role: string): AppKey | null {
	return labels[role] ?? null;
}

export function roleTone(role: string): RoleTone {
	if (role === 'owner') return 'owner';
	if (role === 'admin') return 'admin';
	return 'neutral';
}
