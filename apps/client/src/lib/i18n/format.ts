import { locale } from './index.svelte';

/** Formats a date or instant with `Intl.DateTimeFormat` in the current locale. */
export function formatDate(value: Date | string | number, options: Intl.DateTimeFormatOptions = {}): string {
	return new Intl.DateTimeFormat(locale(), options).format(new Date(value));
}

/** Formats the clock time of an instant in the current locale. */
export function formatTime(value: Date | string | number): string {
	return formatDate(value, { hour: '2-digit', minute: '2-digit' });
}

/** Formats a count in the current locale. */
export function formatNumber(value: number): string {
	return new Intl.NumberFormat(locale()).format(value);
}

/** Formats a byte count as B, KB or MB in the current locale. */
export function formatBytes(bytes: number): string {
	const units = ['B', 'KB', 'MB', 'GB'];
	let size = bytes;
	let unit = 0;
	while (size >= 1024 && unit < units.length - 1) {
		size /= 1024;
		unit += 1;
	}
	const digits = unit === 0 ? 0 : 1;
	return `${new Intl.NumberFormat(locale(), { maximumFractionDigits: digits }).format(size)} ${units[unit]}`;
}

/** Splits a translated sentence at its `{name}` placeholders so a component can put a link or tag in each one. */
export function slots(text: string): (string | { slot: string })[] {
	return text
		.split(/\{(\w+)\}/)
		.map((part, index) => (index % 2 === 1 ? { slot: part } : part))
		.filter((part) => part !== '');
}

const DAY = 1000 * 60 * 60 * 24;

/** A message date for a list row: the time today, the weekday this week, else day and month. */
export function formatListDate(value: string): string {
	const date = new Date(value);
	const age = (Date.now() - date.getTime()) / DAY;
	if (age < 1) return formatTime(date);
	if (age < 7) return formatDate(date, { weekday: 'short' });
	return formatDate(date, { day: 'numeric', month: 'short' });
}

/** A message date for a conversation header: like a list row, with more detail. */
export function formatThreadDate(value: string): string {
	const date = new Date(value);
	const age = (Date.now() - date.getTime()) / DAY;
	if (age < 1) return formatTime(date);
	if (age < 7) return formatDate(date, { weekday: 'short', hour: '2-digit', minute: '2-digit' });
	return formatDate(date, { day: 'numeric', month: 'short', year: 'numeric' });
}
