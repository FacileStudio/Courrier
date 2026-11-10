import type { Catalogs, Message } from '@facile/muse';
import { accounts } from './accounts';
import { advanced } from './advanced';
import { assist } from './assist';
import { auth } from './auth';
import { common } from './common';
import { compose } from './compose';
import { ideas } from './ideas';
import { landing } from './landing';
import { mail } from './mail';
import { onboarding } from './onboarding';
import { pending } from './pending';
import { preferences } from './preferences';
import { spaceSettings } from './spaceSettings';
import { spaces } from './spaces';
import { templates } from './templates';
import type { Row } from '../row';

const rows = {
	...common,
	...auth,
	...landing,
	...onboarding,
	...preferences,
	...accounts,
	...templates,
	...advanced,
	...mail,
	...assist,
	...compose,
	...spaces,
	...spaceSettings,
	...pending,
	...ideas
};

export type AppKey = keyof typeof rows;
type AppCatalog = Record<AppKey, Message>;

function column(index: number): AppCatalog {
	const entries = Object.entries(rows as Record<string, Row>).map(([key, row]) => [key, row[index]]);
	return Object.fromEntries(entries) as AppCatalog;
}

/** The app's catalogs: English is the typed base, and every row carries all four locales. */
export const catalogs: Catalogs<AppCatalog> = {
	en: column(0),
	fr: column(1),
	es: column(2),
	de: column(3)
};
