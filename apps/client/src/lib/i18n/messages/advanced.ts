import type { Row } from '../row';

/** The advanced settings page. */
export const advanced = {
	'advanced.instance': ['Instance', 'Instance', 'Instancia', 'Instanz'],
	'advanced.instanceBody': [
		'The facts worth quoting when you file a bug against a self-hosted install.',
		'Les informations à citer quand vous signalez un bug sur une installation auto-hébergée.',
		'Los datos que conviene citar al reportar un error en una instalación autoalojada.',
		'Die Angaben, die du bei einem Fehlerbericht zu einer selbst gehosteten Installation nennen solltest.'
	],
	'advanced.apiBase': ['API base URL', 'URL de base de l’API', 'URL base de la API', 'Basis-URL der API'],
	'advanced.apiBaseBody': [
		'Same origin as this page in production.',
		'Même origine que cette page en production.',
		'Mismo origen que esta página en producción.',
		'In Produktion derselbe Ursprung wie diese Seite.'
	],
	'advanced.userId': ['Your user ID', 'Votre identifiant utilisateur', 'Tu ID de usuario', 'Deine Benutzer-ID'],
	'advanced.userIdBody': [
		'What the API knows you as.',
		'Sous quel identifiant l’API vous connaît.',
		'Como te conoce la API.',
		'Als wer die API dich kennt.'
	],
	'advanced.signIn': ['Sign-in', 'Connexion', 'Inicio de sesión', 'Anmeldung'],
	'advanced.signInBody': [
		'Everyone signs in with a Facile account, through Registre. Mail clients and scripts use API tokens.',
		'Tout le monde se connecte avec un compte Facile, via Registre. Les clients de messagerie et les scripts utilisent des jetons d’API.',
		'Todo el mundo inicia sesión con una cuenta de Facile, a través de Registre. Los clientes de correo y los scripts usan tokens de API.',
		'Alle melden sich mit einem Facile-Konto über Registre an. Mailclients und Skripte nutzen API-Tokens.'
	],
	'advanced.unknown': ['Unknown', 'Inconnu', 'Desconocido', 'Unbekannt'],
	'advanced.federated': ['Facile account (OIDC)', 'Compte Facile (OIDC)', 'Cuenta de Facile (OIDC)', 'Facile-Konto (OIDC)'],
	'advanced.notConfigured': ['Not configured', 'Non configuré', 'Sin configurar', 'Nicht konfiguriert'],
	'advanced.encryption': [
		'Credential encryption',
		'Chiffrement des identifiants',
		'Cifrado de credenciales',
		'Verschlüsselung der Zugangsdaten'
	],
	'advanced.encryptionBody': [
		'ENCRYPTION_KEY is what keeps your IMAP and SMTP passwords unreadable in the database.',
		'ENCRYPTION_KEY rend vos mots de passe IMAP et SMTP illisibles dans la base de données.',
		'ENCRYPTION_KEY es lo que mantiene ilegibles tus contraseñas IMAP y SMTP en la base de datos.',
		'ENCRYPTION_KEY sorgt dafür, dass deine IMAP- und SMTP-Passwörter in der Datenbank unlesbar sind.'
	],
	'advanced.keyConfigured': ['Key configured', 'Clé configurée', 'Clave configurada', 'Schlüssel konfiguriert'],
	'advanced.noKey': ['No key set', 'Aucune clé définie', 'Sin clave', 'Kein Schlüssel gesetzt'],
	'advanced.noKeyTitle': [
		'Mail credentials are not encrypted at rest',
		'Les identifiants de messagerie ne sont pas chiffrés au repos',
		'Las credenciales de correo no están cifradas en reposo',
		'Zugangsdaten sind nicht verschlüsselt gespeichert'
	],
	'advanced.noKeyBody': [
		'This instance is running without an ENCRYPTION_KEY. Set one in the server environment and restart before you add another mail account.',
		'Cette instance tourne sans ENCRYPTION_KEY. Définissez-en une dans l’environnement du serveur et redémarrez avant d’ajouter un autre compte.',
		'Esta instancia se ejecuta sin ENCRYPTION_KEY. Define una en el entorno del servidor y reinicia antes de añadir otra cuenta de correo.',
		'Diese Instanz läuft ohne ENCRYPTION_KEY. Setze einen in der Serverumgebung und starte neu, bevor du ein weiteres Konto hinzufügst.'
	],
	'advanced.danger': ['Danger zone', 'Zone de danger', 'Zona de peligro', 'Gefahrenzone'],
	'advanced.dangerBody': [
		'Irreversible, and nobody can undo it for you.',
		'Irréversible, et personne ne peut l’annuler à votre place.',
		'Irreversible, y nadie puede deshacerlo por ti.',
		'Nicht rückgängig zu machen, und niemand kann es für dich rückgängig machen.'
	],
	'advanced.disconnectAll': [
		'Disconnect every mail account',
		'Déconnecter tous les comptes de messagerie',
		'Desconectar todas las cuentas de correo',
		'Alle E-Mail-Konten trennen'
	],
	'advanced.disconnectNone': [
		'Nothing to disconnect: there are no mail accounts in this space.',
		'Rien à déconnecter : cet espace n’a aucun compte de messagerie.',
		'Nada que desconectar: este espacio no tiene cuentas de correo.',
		'Nichts zu trennen: In diesem Bereich gibt es keine E-Mail-Konten.'
	],
	'advanced.disconnectBody': [
		{
			one: 'Removes the account in this space, along with everything Courrier has cached from it.',
			other: 'Removes all {count} accounts in this space, along with everything Courrier has cached from them.'
		},
		{
			one: 'Supprime le compte de cet espace, avec tout ce que Courrier en a mis en cache.',
			other: 'Supprime les {count} comptes de cet espace, avec tout ce que Courrier en a mis en cache.'
		},
		{
			one: 'Elimina la cuenta de este espacio, junto con todo lo que Courrier tiene en caché de ella.',
			other: 'Elimina las {count} cuentas de este espacio, junto con todo lo que Courrier tiene en caché de ellas.'
		},
		{
			one: 'Entfernt das Konto in diesem Bereich samt allem, was Courrier davon zwischengespeichert hat.',
			other: 'Entfernt alle {count} Konten in diesem Bereich samt allem, was Courrier davon zwischengespeichert hat.'
		}
	],
	'advanced.disconnectAllButton': ['Disconnect all', 'Tout déconnecter', 'Desconectar todo', 'Alle trennen'],
	'advanced.disconnecting': ['Disconnecting…', 'Déconnexion…', 'Desconectando…', 'Wird getrennt…'],
	'advanced.wipeFailed': [
		'Could not remove every account',
		'Impossible de supprimer tous les comptes',
		'No se pudieron quitar todas las cuentas',
		'Nicht alle Konten konnten entfernt werden'
	],
	'advanced.stoppedAfter': [
		{ one: 'Stopped after {count} account', other: 'Stopped after {count} accounts' },
		{ one: 'Arrêt après {count} compte', other: 'Arrêt après {count} comptes' },
		{ one: 'Se detuvo tras {count} cuenta', other: 'Se detuvo tras {count} cuentas' },
		{ one: 'Nach {count} Konto gestoppt', other: 'Nach {count} Konten gestoppt' }
	],
	'advanced.disconnected': [
		{
			one: 'Disconnected {count} account. Your mail is untouched on the server.',
			other: 'Disconnected {count} accounts. Your mail is untouched on the server.'
		},
		{
			one: '{count} compte déconnecté. Votre courrier reste intact sur le serveur.',
			other: '{count} comptes déconnectés. Votre courrier reste intact sur le serveur.'
		},
		{
			one: 'Se desconectó {count} cuenta. Tu correo sigue intacto en el servidor.',
			other: 'Se desconectaron {count} cuentas. Tu correo sigue intacto en el servidor.'
		},
		{
			one: '{count} Konto getrennt. Deine Mails sind auf dem Server unverändert.',
			other: '{count} Konten getrennt. Deine Mails sind auf dem Server unverändert.'
		}
	],
	'advanced.wipeTitle': [
		'Disconnect every mail account?',
		'Déconnecter tous les comptes de messagerie ?',
		'¿Desconectar todas las cuentas de correo?',
		'Alle E-Mail-Konten trennen?'
	],
	'advanced.wipeBody': [
		{
			one: 'Courrier deletes the account in this space, its stored IMAP and SMTP credentials, and every folder, message and draft cached from it, so the mail views go empty. Your mail itself stays on the server, but adding the account back means re-entering its passwords and syncing from scratch.',
			other: 'Courrier deletes all {count} accounts in this space, their stored IMAP and SMTP credentials, and every folder, message and draft cached from them, so the mail views go empty. Your mail itself stays on the servers, but adding the accounts back means re-entering every password and syncing from scratch.'
		},
		{
			one: 'Courrier supprime le compte de cet espace, ses identifiants IMAP et SMTP enregistrés, ainsi que chaque dossier, message et brouillon mis en cache : les vues courrier se vident. Votre courrier lui-même reste sur le serveur, mais rajouter le compte oblige à ressaisir ses mots de passe et à tout resynchroniser.',
			other: 'Courrier supprime les {count} comptes de cet espace, leurs identifiants IMAP et SMTP enregistrés, ainsi que chaque dossier, message et brouillon mis en cache : les vues courrier se vident. Votre courrier lui-même reste sur les serveurs, mais rajouter les comptes oblige à ressaisir chaque mot de passe et à tout resynchroniser.'
		},
		{
			one: 'Courrier elimina la cuenta de este espacio, sus credenciales IMAP y SMTP guardadas y cada carpeta, mensaje y borrador en caché, así que las vistas de correo quedan vacías. Tu correo en sí permanece en el servidor, pero volver a añadir la cuenta obliga a reescribir sus contraseñas y sincronizar desde cero.',
			other: 'Courrier elimina las {count} cuentas de este espacio, sus credenciales IMAP y SMTP guardadas y cada carpeta, mensaje y borrador en caché, así que las vistas de correo quedan vacías. Tu correo en sí permanece en los servidores, pero volver a añadir las cuentas obliga a reescribir cada contraseña y sincronizar desde cero.'
		},
		{
			one: 'Courrier löscht das Konto in diesem Bereich, seine gespeicherten IMAP- und SMTP-Zugangsdaten sowie jeden zwischengespeicherten Ordner, jede Nachricht und jeden Entwurf, sodass die Mailansichten leer werden. Deine Mails selbst bleiben auf dem Server, aber beim erneuten Hinzufügen müssen die Passwörter neu eingegeben und alles von vorn synchronisiert werden.',
			other: 'Courrier löscht alle {count} Konten in diesem Bereich, ihre gespeicherten IMAP- und SMTP-Zugangsdaten sowie jeden zwischengespeicherten Ordner, jede Nachricht und jeden Entwurf, sodass die Mailansichten leer werden. Deine Mails selbst bleiben auf den Servern, aber beim erneuten Hinzufügen müssen alle Passwörter neu eingegeben und alles von vorn synchronisiert werden.'
		}
	],
	'advanced.leaveThem': ['Leave them alone', 'Ne rien toucher', 'Dejarlas como están', 'Nicht anfassen']
} satisfies Record<string, Row>;
