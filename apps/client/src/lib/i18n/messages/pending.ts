import type { Row } from '../row';

/** The page where the person confirms what an AI client asked to do. */
export const pending = {
	'pending.title': ['Pending actions', 'Actions en attente', 'Acciones pendientes', 'Ausstehende Aktionen'],
	'pending.intro': [
		'An AI client asked Courrier to do these. Nothing happens until you confirm.',
		'Un client IA a demandé ces actions à Courrier. Rien ne se passe tant que vous ne confirmez pas.',
		'Un cliente de IA pidió estas acciones a Courrier. No ocurre nada hasta que confirmes.',
		'Ein KI-Client hat Courrier um diese Aktionen gebeten. Nichts passiert, bevor du bestätigst.'
	],
	'pending.none': ['Nothing is waiting', 'Rien en attente', 'Nada pendiente', 'Nichts ausstehend'],
	'pending.noneBody': [
		'Actions an AI client asks for appear here until you confirm or dismiss them.',
		'Les actions demandées par un client IA apparaissent ici jusqu’à ce que vous les confirmiez ou les ignoriez.',
		'Las acciones que pide un cliente de IA aparecen aquí hasta que las confirmes o descartes.',
		'Aktionen, um die ein KI-Client bittet, erscheinen hier, bis du sie bestätigst oder verwirfst.'
	],
	'pending.send': ['Send email', 'Envoyer l’e-mail', 'Enviar correo', 'E-Mail senden'],
	'pending.deleteDraft': ['Delete draft', 'Supprimer le brouillon', 'Eliminar borrador', 'Entwurf löschen'],
	'pending.from': ['From', 'De', 'De', 'Von'],
	'pending.to': ['To', 'À', 'Para', 'An'],
	'pending.cc': ['Cc', 'Cc', 'Cc', 'Cc'],
	'pending.bcc': ['Bcc', 'Cci', 'Cco', 'Bcc'],
	'pending.askedBy': ['Asked by {client}', 'Demandé par {client}', 'Pedido por {client}', 'Angefragt von {client}'],
	'pending.confirmSend': ['Send', 'Envoyer', 'Enviar', 'Senden'],
	'pending.confirmDelete': ['Delete', 'Supprimer', 'Eliminar', 'Löschen'],
	'pending.dismiss': ['Dismiss', 'Ignorer', 'Descartar', 'Verwerfen'],
	'pending.openDraft': ['Open draft', 'Ouvrir le brouillon', 'Abrir borrador', 'Entwurf öffnen'],
	'pending.changed': [
		'The draft changed after the request. Ask the AI client again.',
		'Le brouillon a changé après la demande. Redemandez au client IA.',
		'El borrador cambió después de la solicitud. Vuelve a pedírselo al cliente de IA.',
		'Der Entwurf wurde nach der Anfrage geändert. Frag den KI-Client erneut.'
	],
	'pending.missing': [
		'The draft no longer exists.',
		'Le brouillon n’existe plus.',
		'El borrador ya no existe.',
		'Der Entwurf existiert nicht mehr.'
	],
	'pending.sent': ['Email sent', 'E-mail envoyé', 'Correo enviado', 'E-Mail gesendet'],
	'pending.deleted': ['Draft deleted', 'Brouillon supprimé', 'Borrador eliminado', 'Entwurf gelöscht'],
	'pending.dismissed': ['Dismissed. The draft is kept.', 'Ignoré. Le brouillon est conservé.', 'Descartado. Se conserva el borrador.', 'Verworfen. Der Entwurf bleibt erhalten.'],
	'pending.failed': ['That did not work', 'Cela n’a pas marché', 'No ha funcionado', 'Das hat nicht geklappt'],
	'pending.signInTitle': [
		'Sign in again to confirm',
		'Reconnectez-vous pour confirmer',
		'Vuelve a iniciar sesión para confirmar',
		'Melde dich erneut an, um zu bestätigen'
	],
	'pending.signInBody': [
		'Only a fresh sign-in in your browser can confirm what an AI client asked for. Sign in again, then come back to this page.',
		'Seule une connexion récente dans votre navigateur peut confirmer ce qu’un client IA a demandé. Reconnectez-vous, puis revenez sur cette page.',
		'Solo un inicio de sesión reciente en tu navegador puede confirmar lo que pidió un cliente de IA. Vuelve a iniciar sesión y regresa a esta página.',
		'Nur eine frische Anmeldung in deinem Browser kann bestätigen, worum ein KI-Client gebeten hat. Melde dich erneut an und komm dann auf diese Seite zurück.'
	],
	'pending.appCannotConfirm': [
		'The mobile and desktop apps cannot confirm what an AI client asked for. Open Courrier in your browser, sign in, and confirm there.',
		'Les applications mobile et bureau ne peuvent pas confirmer ce qu’un client IA a demandé. Ouvrez Courrier dans votre navigateur, connectez-vous et confirmez-y.',
		'Las aplicaciones móvil y de escritorio no pueden confirmar lo que pidió un cliente de IA. Abre Courrier en tu navegador, inicia sesión y confirma allí.',
		'Die Mobil- und Desktop-Apps können nicht bestätigen, worum ein KI-Client gebeten hat. Öffne Courrier im Browser, melde dich an und bestätige dort.'
	],
	'pending.signInAgain': ['Sign in again', 'Se reconnecter', 'Volver a iniciar sesión', 'Erneut anmelden'],
	'pending.loadFailed': [
		'Could not load pending actions',
		'Impossible de charger les actions en attente',
		'No se pudieron cargar las acciones pendientes',
		'Ausstehende Aktionen konnten nicht geladen werden'
	]
} satisfies Record<string, Row>;
