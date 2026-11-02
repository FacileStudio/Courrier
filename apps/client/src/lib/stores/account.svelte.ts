const PREFIX = 'courrier.active-account.';

function key(spaceId: string | null | undefined): string {
	return PREFIX + (spaceId ?? 'personal');
}

export const activeAccountStore = {
	get(spaceId: string | null | undefined): number | null {
		try {
			const raw = localStorage.getItem(key(spaceId));
			if (!raw) return null;
			const id = Number(raw);
			return Number.isInteger(id) ? id : null;
		} catch {
			return null;
		}
	},

	set(spaceId: string | null | undefined, accountId: number) {
		try {
			localStorage.setItem(key(spaceId), String(accountId));
		} catch {}
	}
};
