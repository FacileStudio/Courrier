import { backend, type AssistStatus } from '$lib/backend';

let status = $state<AssistStatus | null>(null);
let pending: Promise<void> | null = null;

/** The helper's availability, fetched once per page load and shared by every blob. */
export const assistStatus = {
	get current() {
		return status;
	},
	load(): Promise<void> {
		pending ??= backend.assistStatus().then(
			(next) => {
				status = next;
			},
			() => {
				status = { available: false, reason: 'unreachable', actions: [] };
			}
		);
		return pending;
	}
};
