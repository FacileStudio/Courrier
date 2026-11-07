import type { Row } from '../row';

/** Words and errors shared by the whole app, the folders and the shell. */
export const common = {
	'common.cancel': ['Cancel', 'Annuler', 'Cancelar', 'Abbrechen'],
	'common.saving': ['Saving…', 'Enregistrement…', 'Guardando…', 'Wird gespeichert…'],
	'common.delete': ['Delete', 'Supprimer', 'Eliminar', 'Löschen'],
	'common.edit': ['Edit', 'Modifier', 'Editar', 'Bearbeiten'],
	'common.open': ['Open', 'Ouvrir', 'Abrir', 'Öffnen'],
	'common.back': ['Back', 'Retour', 'Atrás', 'Zurück'],
	'common.next': ['Next', 'Suivant', 'Siguiente', 'Weiter'],
	'common.skip': ['Skip', 'Passer', 'Omitir', 'Überspringen'],
	'common.noSubject': ['(no subject)', '(sans objet)', '(sin asunto)', '(kein Betreff)'],
	'common.you': ['You', 'Vous', 'Tú', 'Du'],
	'error.requestFailed': [
		'Request failed with status {status}',
		'La requête a échoué (statut {status})',
		'La solicitud falló con el estado {status}',
		'Anfrage fehlgeschlagen (Status {status})'
	],
	'error.downloadFailed': [
		'Download failed with status {status}',
		'Le téléchargement a échoué (statut {status})',
		'La descarga falló con el estado {status}',
		'Download fehlgeschlagen (Status {status})'
	],
	'error.signInState': [
		'Sign-in state did not match. Try again.',
		'L’état de connexion ne correspond pas. Réessayez.',
		'El estado del inicio de sesión no coincide. Inténtalo de nuevo.',
		'Der Anmeldestatus stimmt nicht überein. Versuche es erneut.'
	],
	'error.signInIncomplete': [
		'Sign-in could not be completed. Try again.',
		'La connexion n’a pas pu aboutir. Réessayez.',
		'No se pudo completar el inicio de sesión. Inténtalo de nuevo.',
		'Die Anmeldung konnte nicht abgeschlossen werden. Versuche es erneut.'
	],
	'error.signInFailed': ['Sign-in failed.', 'La connexion a échoué.', 'No se pudo iniciar sesión.', 'Die Anmeldung ist fehlgeschlagen.'],
	'folder.inbox': ['Inbox', 'Boîte de réception', 'Bandeja de entrada', 'Posteingang'],
	'folder.sent': ['Sent', 'Envoyés', 'Enviados', 'Gesendet'],
	'folder.drafts': ['Drafts', 'Brouillons', 'Borradores', 'Entwürfe'],
	'folder.archive': ['Archive', 'Archives', 'Archivo', 'Archiv'],
	'folder.junk': ['Junk', 'Indésirables', 'Correo no deseado', 'Spam'],
	'folder.trash': ['Trash', 'Corbeille', 'Papelera', 'Papierkorb'],
	'tabs.label': ['Inbox tabs', 'Onglets de la boîte de réception', 'Pestañas de la bandeja de entrada', 'Posteingang-Tabs'],
	'tabs.primary': ['Primary', 'Principal', 'Principal', 'Primär'],
	'tabs.newsletters': ['Newsletters', 'Newsletters', 'Boletines', 'Newsletter'],
	'tabs.notifications': ['Notifications', 'Notifications', 'Notificaciones', 'Benachrichtigungen'],
	'shell.compose': ['Compose', 'Écrire', 'Redactar', 'Verfassen'],
	'shell.mail': ['Mail', 'Courrier', 'Correo', 'Mail'],
	'shell.desktopApp': ['Desktop app', 'Application de bureau', 'Aplicación de escritorio', 'Desktop-App'],
	'shell.switchAccount': [
		'Switch account, current {email}',
		'Changer de compte, actuel : {email}',
		'Cambiar de cuenta, actual: {email}',
		'Konto wechseln, aktuell: {email}'
	],
	'shell.manageAccounts': ['Manage accounts', 'Gérer les comptes', 'Gestionar cuentas', 'Konten verwalten']
} satisfies Record<string, Row>;
