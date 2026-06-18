const SPACE_KEY = 'courrier.active-space';

export type SpaceContext = {
	id: string;
	name: string;
	role: string;
};

function loadActiveSpace(): SpaceContext | null {
	try {
		const raw = localStorage.getItem(SPACE_KEY);
		if (!raw) return null;
		return JSON.parse(raw) as SpaceContext;
	} catch {
		return null;
	}
}

function saveActiveSpace(space: SpaceContext | null) {
	if (space) {
		localStorage.setItem(SPACE_KEY, JSON.stringify(space));
	} else {
		localStorage.removeItem(SPACE_KEY);
	}
}

let current = $state<SpaceContext | null>(loadActiveSpace());

export const spaceStore = {
	get active() {
		return current;
	},

	set(space: SpaceContext | null) {
		current = space;
		saveActiveSpace(space);
	},

	clear() {
		current = null;
		saveActiveSpace(null);
	}
};
