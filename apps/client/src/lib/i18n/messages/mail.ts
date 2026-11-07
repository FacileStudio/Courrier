import type { Row } from '../row';

/** The mail list, the reading pane, the composer and the Bulle helper. */
export const mail = {
	'mail.inboxTitle': ['Inbox — Courrier', 'Boîte de réception — Courrier', 'Bandeja de entrada — Courrier', 'Posteingang — Courrier'],
	'mail.folderTitle': ['{folder} — Courrier', '{folder} — Courrier', '{folder} — Courrier', '{folder} — Courrier'],
	'mail.search': ['Search all mail…', 'Rechercher dans tout le courrier…', 'Buscar en todo el correo…', 'Alle Mails durchsuchen…'],
	'mail.searchLabel': ['Search all mail', 'Rechercher dans tout le courrier', 'Buscar en todo el correo', 'Alle Mails durchsuchen'],
	'mail.clearSearch': ['Clear search', 'Effacer la recherche', 'Borrar la búsqueda', 'Suche löschen'],
	'mail.searchHeading': ['Search', 'Recherche', 'Búsqueda', 'Suche'],
	'mail.results': [
		{ one: '{count} result', other: '{count} results' },
		{ one: '{count} résultat', other: '{count} résultats' },
		{ one: '{count} resultado', other: '{count} resultados' },
		{ one: '{count} Ergebnis', other: '{count} Ergebnisse' }
	],
	'mail.selectAll': ['Select all', 'Tout sélectionner', 'Seleccionar todo', 'Alle auswählen'],
	'mail.deselectAll': ['Deselect all', 'Tout désélectionner', 'Deseleccionar todo', 'Alle abwählen'],
	'mail.unread': ['Unread', 'Non lus', 'No leídos', 'Ungelesen'],
	'mail.refresh': ['Refresh', 'Actualiser', 'Actualizar', 'Aktualisieren'],
	'mail.loadingMore': [
		'Loading more messages',
		'Chargement de messages supplémentaires',
		'Cargando más mensajes',
		'Weitere Nachrichten werden geladen'
	],
	'mail.noMatchTitle': [
		'No messages match “{query}”',
		'Aucun message ne correspond à « {query} »',
		'Ningún mensaje coincide con «{query}»',
		'Keine Nachrichten passen zu „{query}“'
	],
	'mail.noMatchBody': [
		'Try a different sender, subject or word from the body.',
		'Essayez un autre expéditeur, un autre objet ou un autre mot du message.',
		'Prueba con otro remitente, asunto o palabra del cuerpo.',
		'Versuche einen anderen Absender, Betreff oder ein anderes Wort aus dem Text.'
	],
	'mail.emptyFolder': [
		'No email in {folder}',
		'Aucun message dans {folder}',
		'Ningún correo en {folder}',
		'Keine Mails in {folder}'
	],
	'mail.emptyUnread': [
		'No unread email in {folder}',
		'Aucun message non lu dans {folder}',
		'Ningún correo sin leer en {folder}',
		'Keine ungelesenen Mails in {folder}'
	],
	'mail.emptyUnreadBody': [
		'Everything here has been read.',
		'Tout ce qui est ici a été lu.',
		'Todo lo de aquí se ha leído.',
		'Alles hier wurde gelesen.'
	],
	'mail.emptyBody': [
		'New messages appear here after a sync.',
		'Les nouveaux messages apparaissent ici après une synchronisation.',
		'Los mensajes nuevos aparecen aquí tras una sincronización.',
		'Neue Nachrichten erscheinen hier nach einer Synchronisierung.'
	],
	'mail.nothingOpen': ['Nothing open', 'Rien d’ouvert', 'Nada abierto', 'Nichts geöffnet'],
	'mail.nothingOpenBody': [
		'Pick a message on the left to read it here.',
		'Choisissez un message à gauche pour le lire ici.',
		'Elige un mensaje a la izquierda para leerlo aquí.',
		'Wähle links eine Nachricht aus, um sie hier zu lesen.'
	],
	'mail.reply': ['Reply', 'Répondre', 'Responder', 'Antworten'],
	'mail.replyAll': ['Reply all', 'Répondre à tous', 'Responder a todos', 'Allen antworten'],
	'mail.forward': ['Forward', 'Transférer', 'Reenviar', 'Weiterleiten'],
	'mail.archive': ['Archive', 'Archiver', 'Archivar', 'Archivieren'],
	'mail.delete': ['Delete', 'Supprimer', 'Eliminar', 'Löschen'],
	'mail.markRead': ['Mark as read', 'Marquer comme lu', 'Marcar como leído', 'Als gelesen markieren'],
	'mail.markUnread': ['Mark as unread', 'Marquer comme non lu', 'Marcar como no leído', 'Als ungelesen markieren'],
	'mail.star': ['Star', 'Suivre', 'Destacar', 'Markieren'],
	'mail.unstar': ['Unstar', 'Ne plus suivre', 'Quitar destacado', 'Markierung entfernen'],
	'mail.clearSelection': ['Clear selection', 'Effacer la sélection', 'Borrar la selección', 'Auswahl aufheben'],
	'mail.selectConversation': ['Select conversation', 'Sélectionner la conversation', 'Seleccionar la conversación', 'Unterhaltung auswählen'],
	'mail.unreadMark': ['Unread', 'Non lu', 'Sin leer', 'Ungelesen'],
	'mail.hasAttachments': ['Has attachments', 'Contient des pièces jointes', 'Tiene archivos adjuntos', 'Enthält Anhänge'],
	'mail.to': ['to {list}', 'à {list}', 'para {list}', 'an {list}'],
	'mail.attachments': [
		{ one: '{count} attachment', other: '{count} attachments' },
		{ one: '{count} pièce jointe', other: '{count} pièces jointes' },
		{ one: '{count} archivo adjunto', other: '{count} archivos adjuntos' },
		{ one: '{count} Anhang', other: '{count} Anhänge' }
	],
	'mail.loadingMessage': ['Loading message…', 'Chargement du message…', 'Cargando el mensaje…', 'Nachricht wird geladen…'],
	'mail.emptyMessage': ['(empty message)', '(message vide)', '(mensaje vacío)', '(leere Nachricht)'],
	'mail.loadingConversation': [
		'Loading conversation',
		'Chargement de la conversation',
		'Cargando la conversación',
		'Unterhaltung wird geladen'
	],
	'mail.loadingConversationDots': [
		'Loading conversation…',
		'Chargement de la conversation…',
		'Cargando la conversación…',
		'Unterhaltung wird geladen…'
	],
	'mail.messageBody': ['Message body', 'Corps du message', 'Cuerpo del mensaje', 'Nachrichtentext'],
	'mail.deletePermTitle': [
		{ one: 'Permanently delete this email?', other: 'Permanently delete {count} emails?' },
		{ one: 'Supprimer définitivement ce message ?', other: 'Supprimer définitivement {count} messages ?' },
		{ one: '¿Eliminar definitivamente este correo?', other: '¿Eliminar definitivamente {count} correos?' },
		{ one: 'Diese Mail endgültig löschen?', other: '{count} Mails endgültig löschen?' }
	],
	'mail.deleteTitle': [
		{ one: 'Delete this email?', other: 'Delete {count} emails?' },
		{ one: 'Supprimer ce message ?', other: 'Supprimer {count} messages ?' },
		{ one: '¿Eliminar este correo?', other: '¿Eliminar {count} correos?' },
		{ one: 'Diese Mail löschen?', other: '{count} Mails löschen?' }
	],
	'mail.deletePermBody': [
		{
			one: 'This cannot be undone. This email will be permanently removed.',
			other: 'This cannot be undone. These emails will be permanently removed.'
		},
		{
			one: 'Cette action est irréversible. Ce message sera définitivement supprimé.',
			other: 'Cette action est irréversible. Ces messages seront définitivement supprimés.'
		},
		{
			one: 'No se puede deshacer. Este correo se eliminará definitivamente.',
			other: 'No se puede deshacer. Estos correos se eliminarán definitivamente.'
		},
		{
			one: 'Das lässt sich nicht rückgängig machen. Diese Mail wird endgültig entfernt.',
			other: 'Das lässt sich nicht rückgängig machen. Diese Mails werden endgültig entfernt.'
		}
	],
	'mail.deleteBody': [
		{
			one: 'This email will be moved to Trash. You can restore it from there.',
			other: 'These emails will be moved to Trash. You can restore them from there.'
		},
		{
			one: 'Ce message sera déplacé dans la corbeille. Vous pourrez l’y récupérer.',
			other: 'Ces messages seront déplacés dans la corbeille. Vous pourrez les y récupérer.'
		},
		{
			one: 'Este correo se moverá a la papelera. Podrás restaurarlo desde allí.',
			other: 'Estos correos se moverán a la papelera. Podrás restaurarlos desde allí.'
		},
		{
			one: 'Diese Mail wird in den Papierkorb verschoben. Du kannst sie dort wiederherstellen.',
			other: 'Diese Mails werden in den Papierkorb verschoben. Du kannst sie dort wiederherstellen.'
		}
	],
	'mail.deletedPermanently': [
		'{count} permanently deleted',
		{ one: '{count} supprimé définitivement', other: '{count} supprimés définitivement' },
		{ one: '{count} eliminado definitivamente', other: '{count} eliminados definitivamente' },
		'{count} endgültig gelöscht'
	],
	'mail.movedToTrash': [
		'{count} moved to trash',
		{ one: '{count} déplacé dans la corbeille', other: '{count} déplacés dans la corbeille' },
		{ one: '{count} movido a la papelera', other: '{count} movidos a la papelera' },
		'{count} in den Papierkorb verschoben'
	],
	'mail.archivedCount': [
		'{count} archived',
		{ one: '{count} archivé', other: '{count} archivés' },
		{ one: '{count} archivado', other: '{count} archivados' },
		'{count} archiviert'
	],
	'mail.archivedOne': ['Archived', 'Archivé', 'Archivado', 'Archiviert'],
	'mail.deleteFailed': ['Failed to delete', 'Échec de la suppression', 'No se pudo eliminar', 'Löschen fehlgeschlagen'],
	'mail.archiveFailed': ['Failed to archive', 'Échec de l’archivage', 'No se pudo archivar', 'Archivieren fehlgeschlagen'],
	'mail.markReadFailed': [
		'Failed to mark as read',
		'Impossible de marquer comme lu',
		'No se pudo marcar como leído',
		'Markieren als gelesen fehlgeschlagen'
	],
	'mail.markUnreadFailed': [
		'Failed to mark as unread',
		'Impossible de marquer comme non lu',
		'No se pudo marcar como no leído',
		'Markieren als ungelesen fehlgeschlagen'
	],

	'sender.label': [
		'Always file mail from this sender in',
		'Toujours classer le courrier de cet expéditeur dans',
		'Archivar siempre el correo de este remitente en',
		'Mails dieses Absenders immer ablegen in'
	],
	'sender.moveTo': ['Move sender to…', 'Déplacer l’expéditeur vers…', 'Mover remitente a…', 'Absender verschieben nach…'],
	'sender.auto': [
		'Let Courrier sort it',
		'Laisser Courrier trier',
		'Dejar que Courrier lo clasifique',
		'Courrier sortieren lassen'
	],
	'sender.autoDone': [
		'{address} is sorted automatically again',
		'{address} est de nouveau trié automatiquement',
		'{address} vuelve a clasificarse automáticamente',
		'{address} wird wieder automatisch sortiert'
	],
	'sender.pinned': [
		'Mail from {address} goes to {tab}',
		'Le courrier de {address} va dans {tab}',
		'El correo de {address} va a {tab}',
		'Mails von {address} landen in {tab}'
	],
	'sender.failed': [
		'Could not move this sender',
		'Impossible de déplacer cet expéditeur',
		'No se pudo mover este remitente',
		'Dieser Absender konnte nicht verschoben werden'
	],

	'assist.summarize': ['Summarize', 'Résumer', 'Resumir', 'Zusammenfassen'],
	'assist.askBulle': ['Ask Bulle', 'Demander à Bulle', 'Preguntar a Bulle', 'Bulle fragen'],
	'assist.summary': ['Thread summary', 'Résumé de la conversation', 'Resumen de la conversación', 'Zusammenfassung der Unterhaltung'],
	'assist.hide': ['Hide', 'Masquer', 'Ocultar', 'Ausblenden'],
	'assist.failed': [
		'Bulle could not answer',
		'Bulle n’a pas pu répondre',
		'Bulle no pudo responder',
		'Bulle konnte nicht antworten'
	],
	'assist.off': [
		'The Bulle helper is not set up on this server.',
		'L’assistant Bulle n’est pas configuré sur ce serveur.',
		'El asistente Bulle no está configurado en este servidor.',
		'Der Bulle-Assistent ist auf diesem Server nicht eingerichtet.'
	],
	'assist.notLinked': [
		'Sign in with your Facile account to use Bulle here.',
		'Connectez-vous avec votre compte Facile pour utiliser Bulle ici.',
		'Inicia sesión con tu cuenta de Facile para usar Bulle aquí.',
		'Melde dich mit deinem Facile-Konto an, um Bulle hier zu nutzen.'
	],
	'assist.noKey': [
		'Add a model key in Bulle to use the helper.',
		'Ajoutez une clé de modèle dans Bulle pour utiliser l’assistant.',
		'Añade una clave de modelo en Bulle para usar el asistente.',
		'Füge in Bulle einen Modellschlüssel hinzu, um den Assistenten zu nutzen.'
	],
	'assist.unsupported': [
		'Bulle’s model is not supported here yet.',
		'Le modèle de Bulle n’est pas encore pris en charge ici.',
		'El modelo de Bulle aún no es compatible aquí.',
		'Das Modell von Bulle wird hier noch nicht unterstützt.'
	],
	'assist.unreachable': [
		'Bulle is not answering right now.',
		'Bulle ne répond pas pour le moment.',
		'Bulle no responde en este momento.',
		'Bulle antwortet gerade nicht.'
	],
	'assist.checking': ['Checking Bulle…', 'Vérification de Bulle…', 'Comprobando Bulle…', 'Bulle wird geprüft…'],
	'assist.scopeSelection': ['the selection', 'la sélection', 'la selección', 'die Auswahl'],
	'assist.scopeAll': ['the whole message', 'tout le message', 'todo el mensaje', 'die ganze Nachricht'],
	'assist.askAbout': ['Ask Bulle about {scope}', 'Demander à Bulle à propos de {scope}', 'Preguntar a Bulle sobre {scope}', 'Bulle zu {scope} fragen'],
	'assist.worksOn': ['Works on {scope}', 'Agit sur {scope}', 'Actúa sobre {scope}', 'Wirkt auf {scope}'],
	'assist.draftReply': ['Draft a reply', 'Rédiger une réponse', 'Redactar una respuesta', 'Antwort entwerfen'],
	'assist.rewrite': ['Rewrite', 'Réécrire', 'Reescribir', 'Umschreiben'],
	'assist.shorten': ['Shorten', 'Raccourcir', 'Acortar', 'Kürzen'],
	'assist.expand': ['Make it longer', 'Allonger', 'Alargar', 'Verlängern'],
	'assist.fix': [
		'Fix spelling and grammar',
		'Corriger l’orthographe et la grammaire',
		'Corregir ortografía y gramática',
		'Rechtschreibung und Grammatik korrigieren'
	],
	'assist.formal': ['More formal', 'Plus formel', 'Más formal', 'Formeller'],
	'assist.friendly': ['Friendlier', 'Plus chaleureux', 'Más cercano', 'Freundlicher'],
	'assist.translate': ['Translate to…', 'Traduire en…', 'Traducir a…', 'Übersetzen in…'],
	'assist.langEnglish': ['English', 'Anglais', 'Inglés', 'Englisch'],
	'assist.langFrench': ['French', 'Français', 'Francés', 'Französisch'],
	'assist.langSpanish': ['Spanish', 'Espagnol', 'Español', 'Spanisch'],
	'assist.langGerman': ['German', 'Allemand', 'Alemán', 'Deutsch'],
	'assist.langItalian': ['Italian', 'Italien', 'Italiano', 'Italienisch'],
	'assist.instructions': [
		'Tell Bulle what to change',
		'Dites à Bulle quoi modifier',
		'Dile a Bulle qué cambiar',
		'Sag Bulle, was geändert werden soll'
	],
	'assist.suggestion': [
		'Bulle’s suggestion',
		'Suggestion de Bulle',
		'Sugerencia de Bulle',
		'Vorschlag von Bulle'
	],
	'assist.suggestionFor': [
		'Suggestion for {scope}',
		'Suggestion pour {scope}',
		'Sugerencia para {scope}',
		'Vorschlag für {scope}'
	],
	'assist.insert': ['Insert', 'Insérer', 'Insertar', 'Einfügen'],
	'assist.replace': ['Replace', 'Remplacer', 'Reemplazar', 'Ersetzen'],
	'assist.insertBelow': ['Insert below', 'Insérer en dessous', 'Insertar debajo', 'Darunter einfügen'],
	'assist.tryAgain': ['Try again', 'Réessayer', 'Reintentar', 'Erneut versuchen'],
	'assist.discard': ['Discard', 'Ignorer', 'Descartar', 'Verwerfen'],

	'editor.placeholder': [
		'Write your message…',
		'Écrivez votre message…',
		'Escribe tu mensaje…',
		'Schreibe deine Nachricht…'
	],
	'editor.linkPrompt': ['URL', 'URL', 'URL', 'URL'],
	'editor.bold': ['Bold', 'Gras', 'Negrita', 'Fett'],
	'editor.italic': ['Italic', 'Italique', 'Cursiva', 'Kursiv'],
	'editor.underline': ['Underline', 'Souligné', 'Subrayado', 'Unterstrichen'],
	'editor.strike': ['Strikethrough', 'Barré', 'Tachado', 'Durchgestrichen'],
	'editor.link': ['Link', 'Lien', 'Enlace', 'Link'],
	'editor.bulletList': ['Bullet list', 'Liste à puces', 'Lista con viñetas', 'Aufzählung'],
	'editor.orderedList': ['Numbered list', 'Liste numérotée', 'Lista numerada', 'Nummerierte Liste'],
	'editor.quote': ['Quote', 'Citation', 'Cita', 'Zitat'],
	'editor.code': ['Code', 'Code', 'Código', 'Code'],
	'address.remove': ['Remove {chip}', 'Retirer {chip}', 'Quitar {chip}', '{chip} entfernen'],

	'compose.title': ['Compose — Courrier', 'Écrire — Courrier', 'Redactar — Courrier', 'Verfassen — Courrier'],
	'compose.heading': ['New message', 'Nouveau message', 'Mensaje nuevo', 'Neue Nachricht'],
	'compose.draftSaved': ['Draft saved {time}', 'Brouillon enregistré à {time}', 'Borrador guardado a las {time}', 'Entwurf gespeichert um {time}'],
	'compose.useTemplate': ['Use template', 'Utiliser un modèle', 'Usar plantilla', 'Vorlage verwenden'],
	'compose.toRich': ['Switch to rich text', 'Passer au texte enrichi', 'Cambiar a texto enriquecido', 'Zu formatiertem Text wechseln'],
	'compose.toPlain': ['Switch to plain text', 'Passer au texte brut', 'Cambiar a texto sin formato', 'Zu reinem Text wechseln'],
	'compose.attach': ['Attach files', 'Joindre des fichiers', 'Adjuntar archivos', 'Dateien anhängen'],
	'compose.send': ['Send', 'Envoyer', 'Enviar', 'Senden'],
	'compose.sending': ['Sending…', 'Envoi…', 'Enviando…', 'Wird gesendet…'],
	'compose.noAccount': [
		'No mail account configured',
		'Aucun compte de messagerie configuré',
		'Ninguna cuenta de correo configurada',
		'Kein E-Mail-Konto eingerichtet'
	],
	'compose.noAccountBody': [
		'Add an IMAP/SMTP account before sending anything.',
		'Ajoutez un compte IMAP/SMTP avant d’envoyer quoi que ce soit.',
		'Añade una cuenta IMAP/SMTP antes de enviar nada.',
		'Füge ein IMAP/SMTP-Konto hinzu, bevor du etwas sendest.'
	],
	'compose.openSettings': ['Open settings', 'Ouvrir les paramètres', 'Abrir los ajustes', 'Einstellungen öffnen'],
	'compose.from': ['From', 'De', 'De', 'Von'],
	'compose.to': ['To', 'À', 'Para', 'An'],
	'compose.cc': ['Cc', 'Cc', 'Cc', 'Cc'],
	'compose.subject': ['Subject', 'Objet', 'Asunto', 'Betreff'],
	'compose.toPlaceholder': ['recipient@example.com', 'destinataire@exemple.fr', 'destinatario@ejemplo.com', 'empfaenger@beispiel.de'],
	'compose.ccPlaceholder': ['cc@example.com', 'cc@exemple.fr', 'cc@ejemplo.com', 'cc@beispiel.de'],
	'compose.removeFile': ['Remove {name}', 'Retirer {name}', 'Quitar {name}', '{name} entfernen'],
	'compose.bodyLabel': ['Message body', 'Corps du message', 'Cuerpo del mensaje', 'Nachrichtentext'],
	'compose.templates': ['Templates', 'Modèles', 'Plantillas', 'Vorlagen'],
	'compose.draftFailed': [
		'Failed to load draft',
		'Impossible de charger le brouillon',
		'No se pudo cargar el borrador',
		'Der Entwurf konnte nicht geladen werden'
	],
	'compose.originalFailed': [
		'Failed to load original email',
		'Impossible de charger le message d’origine',
		'No se pudo cargar el correo original',
		'Die ursprüngliche Mail konnte nicht geladen werden'
	],
	'compose.sent': ['Email sent', 'Message envoyé', 'Correo enviado', 'Mail gesendet'],
	'compose.sendFailed': ['Failed to send', 'Échec de l’envoi', 'No se pudo enviar', 'Senden fehlgeschlagen'],
	'compose.wrote': [
		'On {date}, {name} wrote:',
		'Le {date}, {name} a écrit :',
		'El {date}, {name} escribió:',
		'Am {date} schrieb {name}:'
	],
	'compose.forwardedHeader': [
		'---------- Forwarded message ----------',
		'---------- Message transféré ----------',
		'---------- Mensaje reenviado ----------',
		'---------- Weitergeleitete Nachricht ----------'
	],
	'compose.fwdFrom': ['From', 'De', 'De', 'Von'],
	'compose.fwdDate': ['Date', 'Date', 'Fecha', 'Datum'],
	'compose.fwdSubject': ['Subject', 'Objet', 'Asunto', 'Betreff'],
	'compose.fwdTo': ['To', 'À', 'Para', 'An']
} satisfies Record<string, Row>;
