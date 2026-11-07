import type { Row } from '../row';

/** The templates and API token pages. */
export const templates = {
	'templates.title': ['Email templates', 'Modèles d’e-mail', 'Plantillas de correo', 'E-Mail-Vorlagen'],
	'templates.body': [
		'Canned messages you can drop into a compose window instead of retyping them.',
		'Des messages types à insérer dans une fenêtre de rédaction plutôt que de les retaper.',
		'Mensajes prediseñados que puedes insertar en una ventana de redacción en lugar de volver a escribirlos.',
		'Vorgefertigte Nachrichten, die du in ein Schreibfenster einfügen kannst, statt sie neu zu tippen.'
	],
	'templates.new': ['New template', 'Nouveau modèle', 'Nueva plantilla', 'Neue Vorlage'],
	'templates.empty': ['No templates yet', 'Aucun modèle', 'Aún no hay plantillas', 'Noch keine Vorlagen'],
	'templates.emptyBody': [
		'Anything you type more than twice, an onboarding reply or a quote follow-up, belongs here.',
		'Tout ce que vous tapez plus de deux fois, une réponse d’accueil ou une relance de devis, a sa place ici.',
		'Todo lo que escribes más de dos veces, una respuesta de bienvenida o un seguimiento de presupuesto, tiene su lugar aquí.',
		'Alles, was du mehr als zweimal tippst, etwa eine Willkommensantwort oder eine Angebotsnachfrage, gehört hierher.'
	],
	'templates.writeFirst': ['Write your first one', 'Écrire le premier', 'Escribe la primera', 'Erste Vorlage schreiben'],
	'templates.colName': ['Name', 'Nom', 'Nombre', 'Name'],
	'templates.colSubject': ['Subject', 'Objet', 'Asunto', 'Betreff'],
	'templates.colUpdated': ['Updated', 'Modifié', 'Actualizada', 'Aktualisiert'],
	'templates.colActions': ['Actions', 'Actions', 'Acciones', 'Aktionen'],
	'templates.loadFailed': [
		'Could not load your templates',
		'Impossible de charger vos modèles',
		'No se pudieron cargar tus plantillas',
		'Deine Vorlagen konnten nicht geladen werden'
	],
	'templates.saveFailed': [
		'Could not save the template',
		'Impossible d’enregistrer le modèle',
		'No se pudo guardar la plantilla',
		'Die Vorlage konnte nicht gespeichert werden'
	],
	'templates.deleteFailed': [
		'Could not delete the template',
		'Impossible de supprimer le modèle',
		'No se pudo eliminar la plantilla',
		'Die Vorlage konnte nicht gelöscht werden'
	],
	'templates.editTitle': ['Edit {name}', 'Modifier {name}', 'Editar {name}', '{name} bearbeiten'],
	'templates.name': ['Name', 'Nom', 'Nombre', 'Name'],
	'templates.nameHelp': [
		'How you will find it in the compose window.',
		'Comment vous le retrouverez dans la fenêtre de rédaction.',
		'Cómo la encontrarás en la ventana de redacción.',
		'So findest du sie im Schreibfenster.'
	],
	'templates.namePlaceholder': ['Quote follow-up', 'Relance de devis', 'Seguimiento de presupuesto', 'Angebotsnachfrage'],
	'templates.subject': ['Subject line', 'Objet', 'Asunto', 'Betreffzeile'],
	'templates.subjectHelp': [
		'Optional. Leave it empty to keep whatever you typed.',
		'Facultatif. Laissez vide pour garder ce que vous avez saisi.',
		'Opcional. Déjalo vacío para conservar lo que hayas escrito.',
		'Optional. Lass es leer, um zu behalten, was du getippt hast.'
	],
	'templates.subjectPlaceholder': [
		'Following up on our quote',
		'Suite à notre devis',
		'Seguimiento de nuestro presupuesto',
		'Zu unserem Angebot'
	],
	'templates.bodyLabel': ['Body', 'Message', 'Cuerpo', 'Text'],
	'templates.bodyPlaceholder': ['Hi there,', 'Bonjour,', 'Hola,', 'Hallo,'],
	'templates.save': ['Save template', 'Enregistrer le modèle', 'Guardar plantilla', 'Vorlage speichern'],
	'templates.create': ['Create template', 'Créer le modèle', 'Crear plantilla', 'Vorlage erstellen'],
	'templates.deleteTitle': [
		'Delete this template?',
		'Supprimer ce modèle ?',
		'¿Eliminar esta plantilla?',
		'Diese Vorlage löschen?'
	],
	'templates.deleteBody': [
		'"{name}" disappears from every compose window in this space, for everyone who shares it. Messages already sent from it are untouched, and there is no undo.',
		'« {name} » disparaît de toutes les fenêtres de rédaction de cet espace, pour toutes les personnes qui le partagent. Les messages déjà envoyés ne changent pas, et il n’y a pas d’annulation.',
		'«{name}» desaparece de todas las ventanas de redacción de este espacio, para todas las personas que lo comparten. Los mensajes ya enviados no cambian y no hay forma de deshacerlo.',
		'„{name}“ verschwindet aus jedem Schreibfenster in diesem Bereich, für alle, die ihn teilen. Bereits gesendete Nachrichten bleiben unverändert, und es gibt kein Rückgängigmachen.'
	],
	'templates.deleteConfirm': ['Delete template', 'Supprimer le modèle', 'Eliminar plantilla', 'Vorlage löschen'],
	'api.endpoint': ['Endpoint', 'Point d’accès', 'Punto de acceso', 'Endpunkt'],
	'api.endpointBody': [
		'Point a script or an HTTP client here. Not a secret. The token is.',
		'Pointez un script ou un client HTTP ici. Ce n’est pas un secret. Le jeton, si.',
		'Apunta un script o un cliente HTTP aquí. No es un secreto. El token sí.',
		'Richte ein Skript oder einen HTTP-Client hierhin. Das ist kein Geheimnis. Das Token schon.'
	],
	'api.baseUrl': ['Base URL', 'URL de base', 'URL base', 'Basis-URL'],
	'api.tokenTitle': ['Personal API token', 'Jeton d’API personnel', 'Token de API personal', 'Persönliches API-Token'],
	'api.tokenBody': [
		'Authenticates a script, the CLI or a mail client as you. Courrier keeps one per person, so issuing a new token retires the old one on the spot.',
		'Authentifie un script, la CLI ou un client de messagerie en votre nom. Courrier en garde un par personne : émettre un nouveau jeton retire l’ancien aussitôt.',
		'Autentica un script, la CLI o un cliente de correo como tú. Courrier guarda uno por persona, así que emitir un token nuevo retira el anterior al instante.',
		'Authentifiziert ein Skript, die CLI oder einen Mailclient als dich. Courrier speichert eines pro Person, ein neues Token macht das alte sofort ungültig.'
	],
	'api.storedHashed': ['stored hashed', 'stocké haché', 'guardado con hash', 'als Hash gespeichert'],
	'api.readFailed': [
		'Could not read your token',
		'Impossible de lire votre jeton',
		'No se pudo leer tu token',
		'Dein Token konnte nicht gelesen werden'
	],
} satisfies Record<string, Row>;
