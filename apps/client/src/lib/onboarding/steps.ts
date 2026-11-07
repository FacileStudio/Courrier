import type { AppKey } from '$lib/i18n/messages';

export type StepId = 'welcome' | 'mailbox' | 'notifications' | 'done';

export type Step = { id: StepId; labelKey: AppKey; skippable: boolean };

export const STEPS: readonly Step[] = [
	{ id: 'welcome', labelKey: 'onboarding.stepWelcome', skippable: false },
	{ id: 'mailbox', labelKey: 'onboarding.stepMailbox', skippable: true },
	{ id: 'notifications', labelKey: 'onboarding.stepNotifications', skippable: true },
	{ id: 'done', labelKey: 'onboarding.stepDone', skippable: false }
];

export const DONE_KEY = 'courrier.onboarded';

const LAST = STEPS.length - 1;

/** Keeps a step index inside the flow. */
export function clampStep(index: number): number {
	if (!Number.isFinite(index)) return 0;
	return Math.min(LAST, Math.max(0, Math.trunc(index)));
}

/** True on the final step. */
export function isLastStep(index: number): boolean {
	return clampStep(index) === LAST;
}

/** The message key of the primary button on a step, given how many mailboxes are connected. */
export function primaryLabel(index: number, accounts: number): AppKey {
	switch (STEPS[clampStep(index)].id) {
		case 'welcome':
			return 'onboarding.primaryWelcome';
		case 'mailbox':
			return accounts === 0 ? 'onboarding.primaryAdd' : 'onboarding.primaryContinue';
		case 'notifications':
			return 'onboarding.primaryContinue';
		default:
			return 'onboarding.primaryDone';
	}
}

/**
 * Says whether to run the first-run flow: once per device in the Android app, and on the web
 * only for someone with no mailbox yet, so existing web users are never sent through it.
 */
export function needsOnboarding(state: { native: boolean; seen: boolean; accounts: number }): boolean {
	if (state.seen) return false;
	return state.native || state.accounts === 0;
}

/** Reads whether this device finished or skipped the flow. */
export function onboardingSeen(store: Pick<Storage, 'getItem'> | null = safeStore()): boolean {
	try {
		return store?.getItem(DONE_KEY) === '1';
	} catch {
		return false;
	}
}

/** Remembers on this device that the flow is done. */
export function markOnboardingSeen(store: Pick<Storage, 'setItem'> | null = safeStore()): void {
	try {
		store?.setItem(DONE_KEY, '1');
	} catch {
		return;
	}
}

function safeStore(): Storage | null {
	try {
		return globalThis.localStorage ?? null;
	} catch {
		return null;
	}
}
