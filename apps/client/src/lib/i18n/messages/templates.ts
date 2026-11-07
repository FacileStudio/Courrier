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
	'api.replace': ['Replace token', 'Remplacer le jeton', 'Reemplazar token', 'Token ersetzen'],
	'api.new': ['New token', 'Nouveau jeton', 'Nuevo token', 'Neues Token'],
	'api.none': [
		'No token yet. Nothing outside the browser can talk to this instance as you until there is one.',
		'Aucun jeton pour l’instant. Rien en dehors du navigateur ne peut parler à cette instance en votre nom tant qu’il n’y en a pas.',
		'Aún no hay token. Nada fuera del navegador puede hablar con esta instancia como tú hasta que exista uno.',
		'Noch kein Token. Bis eines existiert, kann nichts außerhalb des Browsers als du mit dieser Instanz sprechen.'
	],
	'api.colName': ['Name', 'Nom', 'Nombre', 'Name'],
	'api.colToken': ['Token', 'Jeton', 'Token', 'Token'],
	'api.colCreated': ['Created', 'Créé', 'Creado', 'Erstellt'],
	'api.colActions': ['Actions', 'Actions', 'Acciones', 'Aktionen'],
	'api.storedHashed': ['stored hashed', 'stocké haché', 'guardado con hash', 'als Hash gespeichert'],
	'api.readFailed': [
		'Could not read your token',
		'Impossible de lire votre jeton',
		'No se pudo leer tu token',
		'Dein Token konnte nicht gelesen werden'
	],
	'api.createFailed': [
		'Could not create the token',
		'Impossible de créer le jeton',
		'No se pudo crear el token',
		'Das Token konnte nicht erstellt werden'
	],
	'api.revokeFailed': [
		'Could not revoke the token',
		'Impossible de révoquer le jeton',
		'No se pudo revocar el token',
		'Das Token konnte nicht widerrufen werden'
	],
	'api.drawerReplace': ['Replace API token', 'Remplacer le jeton d’API', 'Reemplazar token de API', 'API-Token ersetzen'],
	'api.drawerNew': ['New API token', 'Nouveau jeton d’API', 'Nuevo token de API', 'Neues API-Token'],
	'api.copyNow': ['Copy it now', 'Copiez-le maintenant', 'Cópialo ahora', 'Jetzt kopieren'],
	'api.copyNowBody': [
		'Courrier stores only a hash of this token, so this is the one and only time it is shown. Lose it and the fix is issuing another one.',
		'Courrier ne garde qu’une empreinte de ce jeton : c’est la seule fois où il est affiché. Si vous le perdez, il faudra en émettre un autre.',
		'Courrier guarda solo un hash de este token, así que esta es la única vez que se muestra. Si lo pierdes, la solución es emitir otro.',
		'Courrier speichert nur einen Hash dieses Tokens, deshalb wird es nur dieses eine Mal angezeigt. Geht es verloren, musst du ein neues ausstellen.'
	],
	'api.tokenLabel': ['Token', 'Jeton', 'Token', 'Token'],
	'api.tokenHelp': [
		'Put it in your password manager or your CI secret store.',
		'Rangez-le dans votre gestionnaire de mots de passe ou dans le coffre de secrets de votre CI.',
		'Guárdalo en tu gestor de contraseñas o en el almacén de secretos de tu CI.',
		'Lege es in deinem Passwortmanager oder im Secret-Speicher deiner CI ab.'
	],
	'api.done': ['Done', 'Terminé', 'Hecho', 'Fertig'],
	'api.replacesTitle': [
		'This replaces your current token',
		'Ceci remplace votre jeton actuel',
		'Esto reemplaza tu token actual',
		'Dies ersetzt dein aktuelles Token'
	],
	'api.replacesBody': [
		'The token named "{name}" stops working the moment the new one is issued, and anything still using it starts failing.',
		'Le jeton nommé « {name} » cesse de fonctionner dès que le nouveau est émis, et tout ce qui l’utilise encore se met à échouer.',
		'El token llamado «{name}» deja de funcionar en cuanto se emite el nuevo, y todo lo que lo siga usando empieza a fallar.',
		'Das Token „{name}“ funktioniert nicht mehr, sobald das neue ausgestellt wird, und alles, was es noch nutzt, schlägt fehl.'
	],
	'api.nameLabel': ['Name', 'Nom', 'Nombre', 'Name'],
	'api.nameHelp': [
		'Where the token will live: a machine, a cron job, a script.',
		'Où vivra le jeton : une machine, une tâche cron, un script.',
		'Dónde vivirá el token: una máquina, una tarea cron, un script.',
		'Wo das Token genutzt wird: eine Maschine, ein Cronjob, ein Skript.'
	],
	'api.namePlaceholder': ['laptop-cli', 'laptop-cli', 'laptop-cli', 'laptop-cli'],
	'api.create': ['Create token', 'Créer le jeton', 'Crear token', 'Token erstellen'],
	'api.creating': ['Creating…', 'Création…', 'Creando…', 'Wird erstellt…'],
	'api.revokeTitle': ['Revoke this token?', 'Révoquer ce jeton ?', '¿Revocar este token?', 'Dieses Token widerrufen?'],
	'api.revokeBody': [
		'Every script, cron job and CLI still holding it starts getting 401s immediately, and it cannot be un-revoked. The only way back is issuing a new token and updating whatever used the old one.',
		'Chaque script, tâche cron et CLI qui le détient reçoit aussitôt des 401, et la révocation est irréversible. Il faut émettre un nouveau jeton et mettre à jour tout ce qui utilisait l’ancien.',
		'Todo script, tarea cron y CLI que lo conserve empieza a recibir 401 de inmediato, y la revocación no se puede deshacer. La única vía es emitir un token nuevo y actualizar lo que usaba el anterior.',
		'Jedes Skript, jeder Cronjob und jede CLI, die es noch hält, erhält sofort 401-Fehler, und der Widerruf lässt sich nicht rückgängig machen. Der einzige Weg zurück ist ein neues Token und die Aktualisierung von allem, was das alte nutzte.'
	],
	'api.revokeConfirm': ['Revoke token', 'Révoquer le jeton', 'Revocar token', 'Token widerrufen']
} satisfies Record<string, Row>;
