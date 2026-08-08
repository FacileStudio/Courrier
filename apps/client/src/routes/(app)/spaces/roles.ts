export type RoleTone = 'owner' | 'admin' | 'neutral';

const labels: Record<string, string> = {
	owner: 'Propriétaire',
	admin: 'Admin',
	member: 'Membre'
};

export function roleLabel(role: string): string {
	return labels[role] ?? role;
}

export function roleTone(role: string): RoleTone {
	if (role === 'owner') return 'owner';
	if (role === 'admin') return 'admin';
	return 'neutral';
}
