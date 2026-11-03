/**
 * Maps a server notification path such as "mail" or "mail/inbox" to an app route. Anything
 * else lands on the mailbox, so a crafted payload cannot send the app to an arbitrary URL.
 */
export function routeFor(path: string | undefined): string {
	return /^mail(\/[a-z]+)?$/.test(path ?? '') ? `/${path}` : '/mail';
}
