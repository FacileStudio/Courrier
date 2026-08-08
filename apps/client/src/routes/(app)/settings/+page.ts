import { redirect } from '@sveltejs/kit';

export const prerender = false;

/** /settings is not a section of its own — it opens on the first one. */
export function load() {
	redirect(307, '/settings/profile');
}
