import { ApiError, backend, type AssistStatus } from '$lib/backend';

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
	},
	/** Records a helper run refused for the plan, so every blob shows the upgrade link. */
	refused(err: unknown): void {
		if (!(err instanceof ApiError) || !err.planLimit) return;
		status = { available: false, reason: 'plan_required', actions: status?.actions ?? [], bulle_url: status?.bulle_url, upgrade_url: err.planLimit.upgradeUrl };
	}
};
