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
	'pending.loadFailed': [
		'Could not load pending actions',
		'Impossible de charger les actions en attente',
		'No se pudieron cargar las acciones pendientes',
		'Ausstehende Aktionen konnten nicht geladen werden'
	]
} satisfies Record<string, Row>;
