import type { Row } from '../row';

/** The sign-in screens, the welcome page and the plan limit. */
export const auth = {
	'login.title': ['Sign in — Courrier', 'Connexion — Courrier', 'Iniciar sesión — Courrier', 'Anmelden — Courrier'],
	'login.tagline': [
		'Your mail. Your server.',
		'Votre courrier. Votre serveur.',
		'Tu correo. Tu servidor.',
		'Deine Mails. Dein Server.'
	],
	'login.failed': [
		'Sign-in failed. Try again.',
		'La connexion a échoué. Réessayez.',
		'No se pudo iniciar sesión. Inténtalo de nuevo.',
		'Die Anmeldung ist fehlgeschlagen. Versuche es erneut.'
	],
	'login.nativeFailed': [
		'Could not open the sign-in page.',
		'Impossible d’ouvrir la page de connexion.',
		'No se pudo abrir la página de inicio de sesión.',
		'Die Anmeldeseite konnte nicht geöffnet werden.'
	],
	'login.nativeTitle': ['Sign in to Courrier', 'Connexion à Courrier', 'Inicia sesión en Courrier', 'Bei Courrier anmelden'],
	'login.nativeBody': [
		'Use your Facile account. Your mailboxes and spaces follow you from the web.',
		'Utilisez votre compte Facile. Vos boîtes mail et vos espaces vous suivent depuis le web.',
		'Usa tu cuenta de Facile. Tus buzones y espacios te siguen desde la web.',
		'Nutze dein Facile-Konto. Deine Postfächer und Bereiche folgen dir vom Web.'
	],
	'login.nativeHint': [
		'Sign-in opens in your browser, then brings you back here.',
		'La connexion s’ouvre dans votre navigateur, puis vous ramène ici.',
		'El inicio de sesión se abre en tu navegador y luego te trae de vuelta aquí.',
		'Die Anmeldung öffnet sich im Browser und bringt dich dann hierher zurück.'
	],
	'welcome.title': ['Welcome — Courrier', 'Bienvenue — Courrier', 'Te damos la bienvenida — Courrier', 'Willkommen — Courrier'],
	'welcome.unreachable': [
		'Courrier is not answering. Reload the page to try again.',
		'Courrier ne répond pas. Rechargez la page pour réessayer.',
		'Courrier no responde. Recarga la página para volver a intentarlo.',
		'Courrier antwortet nicht. Lade die Seite neu, um es erneut zu versuchen.'
	],
	'plan.reached': ['Plan limit reached', 'Limite du forfait atteinte', 'Límite del plan alcanzado', 'Limit des Tarifs erreicht'],
	'plan.fallback': [
		'Your plan limit is reached.',
		'La limite de votre forfait est atteinte.',
		'Has alcanzado el límite de tu plan.',
		'Das Limit deines Tarifs ist erreicht.'
	],
	'plan.upgrade': ['Upgrade the plan', 'Passer à un forfait supérieur', 'Mejorar el plan', 'Tarif upgraden'],
	'plan.askOwner': [
		'Ask a space owner to upgrade its plan in Registre.',
		'Demandez à un propriétaire de l’espace de changer son forfait dans Registre.',
		'Pide a un propietario del espacio que mejore su plan en Registre.',
		'Bitte einen Eigentümer des Bereichs, den Tarif in Registre zu upgraden.'
	]
} satisfies Record<string, Row>;
