import { createTheme } from '@facile/client/theme';

export type { Theme as ThemePreference } from '@facile/client/theme';

/**
 * The persisted theme. Courrier's `dark:` Tailwind utilities need an explicit `.dark` and muse needs
 * an explicit `.light` to escape its media block, so system mode resolves to one class as well.
 */
export const theme = createTheme('courrier.theme', { resolveSystem: true });
