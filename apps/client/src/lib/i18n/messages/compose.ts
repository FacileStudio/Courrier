import type { Row } from '../row';

/** The composer. */
export const compose = {
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
