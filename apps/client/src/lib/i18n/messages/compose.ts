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
	'compose.fwdTo': ['To', 'À', 'Para', 'An'],
	'compose.capsule': ['Add a capsule', 'Ajouter une capsule', 'Añadir una cápsula', 'Kapsel hinzufügen'],
	'compose.capsuleTitle': ['Add a secret link', 'Ajouter un lien secret', 'Añadir un enlace secreto', 'Geheimen Link hinzufügen'],
	'compose.capsuleBody': [
		'The text is encrypted in your browser. The link opens once, then the secret is destroyed.',
		'Le texte est chiffré dans votre navigateur. Le lien s’ouvre une fois, puis le secret est détruit.',
		'El texto se cifra en tu navegador. El enlace se abre una vez y luego el secreto se destruye.',
		'Der Text wird in Ihrem Browser verschlüsselt. Der Link öffnet sich einmal, dann wird das Geheimnis gelöscht.'
	],
	'compose.capsuleSecret': ['Secret', 'Secret', 'Secreto', 'Geheimnis'],
	'compose.capsulePlaceholder': ['A password, a code, an address…', 'Un mot de passe, un code, une adresse…', 'Una contraseña, un código, una dirección…', 'Ein Passwort, ein Code, eine Adresse…'],
	'compose.capsuleExpiry': ['Expires if unopened after', 'Expire sans ouverture après', 'Caduca sin abrir después de', 'Läuft ungeöffnet ab nach'],
	'compose.capsule1h': ['1 hour', '1 heure', '1 hora', '1 Stunde'],
	'compose.capsule24h': ['1 day', '1 jour', '1 día', '1 Tag'],
	'compose.capsule7d': ['7 days', '7 jours', '7 días', '7 Tage'],
	'compose.capsule30d': ['30 days', '30 jours', '30 días', '30 Tage'],
	'compose.capsuleInsert': ['Seal and insert', 'Sceller et insérer', 'Sellar e insertar', 'Versiegeln und einfügen'],
	'compose.capsuleSealing': ['Sealing…', 'Scellement…', 'Sellando…', 'Wird versiegelt…'],
	'compose.capsuleLine': ['Secret link, readable once:', 'Lien secret, lisible une seule fois :', 'Enlace secreto, legible una sola vez:', 'Geheimer Link, nur einmal lesbar:'],
	'compose.capsuleAdded': ['Secret link added', 'Lien secret ajouté', 'Enlace secreto añadido', 'Geheimer Link hinzugefügt'],
	'compose.capsuleFailed': ['Could not seal the capsule', 'Impossible de sceller la capsule', 'No se pudo sellar la cápsula', 'Kapsel konnte nicht versiegelt werden']
} satisfies Record<string, Row>;
