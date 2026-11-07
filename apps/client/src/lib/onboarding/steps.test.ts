import { describe, expect, it } from 'vitest';
import {
	STEPS,
	clampStep,
	isLastStep,
	markOnboardingSeen,
	needsOnboarding,
	onboardingSeen,
	primaryLabel
} from './steps';

describe('needsOnboarding', () => {
	it('runs once per device in the app and only for empty accounts on the web', () => {
		expect(needsOnboarding({ native: true, seen: false, accounts: 3 })).toBe(true);
		expect(needsOnboarding({ native: true, seen: true, accounts: 0 })).toBe(false);
		expect(needsOnboarding({ native: false, seen: false, accounts: 0 })).toBe(true);
		expect(needsOnboarding({ native: false, seen: false, accounts: 1 })).toBe(false);
	});
});

describe('steps', () => {
	it('clamps and ends on the ready step', () => {
		expect(clampStep(-4)).toBe(0);
		expect(clampStep(Number.NaN)).toBe(0);
		expect(clampStep(99)).toBe(STEPS.length - 1);
		expect(isLastStep(STEPS.length - 1)).toBe(true);
		expect(primaryLabel(1, 0)).toBe('onboarding.primaryAdd');
		expect(primaryLabel(1, 2)).toBe('onboarding.primaryContinue');
	});

	it('remembers the flow in storage and survives a throwing store', () => {
		const data = new Map<string, string>();
		const store = { getItem: (k: string) => data.get(k) ?? null, setItem: (k: string, v: string) => void data.set(k, v) };
		expect(onboardingSeen(store)).toBe(false);
		markOnboardingSeen(store);
		expect(onboardingSeen(store)).toBe(true);
		const broken = { getItem: () => { throw new Error('denied'); } };
		expect(onboardingSeen(broken)).toBe(false);
	});
});
