import { ApiError, backend, type AssistStatus } from '$lib/backend';

let status = $state<AssistStatus | null>(null);
let pending: Promise<void> | null = null;
let loadedFor: number | undefined;

/**
 * The helper's availability for one mailbox, shared by every blob: a space
 * mailbox may get AI from the space's plan. Fetched again only when the
 * mailbox changes.
 */
export const assistStatus = {
	get current() {
		return status;
	},
	load(accountId?: number): Promise<void> {
		if (pending && loadedFor === accountId) return pending;
		loadedFor = accountId;
		const next = backend.assistStatus(accountId).then(
			(answer) => {
				if (loadedFor === accountId) status = answer;
			},
			() => {
				if (loadedFor === accountId) status = { available: false, reason: 'unreachable', actions: [] };
			}
		);
		pending = next;
		return next;
	},
	/** Records a helper run refused for the plan, so every blob shows the upgrade link. */
	refused(err: unknown): void {
		if (!(err instanceof ApiError) || !err.planLimit) return;
		status = { available: false, reason: 'plan_required', actions: status?.actions ?? [], bulle_url: status?.bulle_url, upgrade_url: err.planLimit.upgradeUrl };
	}
};
