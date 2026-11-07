import type { Row } from '../row';

/** The mail accounts page and its drawer. */
export const accounts = {
	'accounts.title': ['Mail accounts', 'Comptes de messagerie', 'Cuentas de correo', 'E-Mail-Konten'],
	'accounts.bodySpace': [
		'The IMAP and SMTP servers Courrier reads and sends through in {space}.',
		'Les serveurs IMAP et SMTP par lesquels Courrier lit et envoie dans {space}.',
		'Los servidores IMAP y SMTP con los que Courrier lee y envía en {space}.',
		'Die IMAP- und SMTP-Server, über die Courrier in {space} liest und sendet.'
	],
	'accounts.body': [
		'The IMAP and SMTP servers Courrier reads and sends through.',
		'Les serveurs IMAP et SMTP par lesquels Courrier lit et envoie.',
		'Los servidores IMAP y SMTP con los que Courrier lee y envía.',
		'Die IMAP- und SMTP-Server, über die Courrier liest und sendet.'
	],
	'accounts.add': ['Add account', 'Ajouter un compte', 'Añadir cuenta', 'Konto hinzufügen'],
	'accounts.empty': [
		'No mail accounts yet',
		'Aucun compte de messagerie',
		'Aún no hay cuentas de correo',
		'Noch keine E-Mail-Konten'
	],
	'accounts.emptyBody': [
		'Courrier has nothing to read until you point it at an IMAP server. Credentials are encrypted at rest.',
		'Courrier n’a rien à lire tant que vous ne l’avez pas relié à un serveur IMAP. Les identifiants sont chiffrés au repos.',
		'Courrier no tiene nada que leer hasta que lo conectes a un servidor IMAP. Las credenciales se cifran en reposo.',
		'Courrier hat nichts zu lesen, bis du es mit einem IMAP-Server verbindest. Zugangsdaten werden verschlüsselt gespeichert.'
	],
	'accounts.addFirst': ['Add your first account', 'Ajouter votre premier compte', 'Añade tu primera cuenta', 'Erstes Konto hinzufügen'],
	'accounts.default': ['Default', 'Par défaut', 'Predeterminada', 'Standard'],
	'accounts.checkRunning': ['Checking…', 'Vérification…', 'Comprobando…', 'Wird geprüft…'],
	'accounts.notChecked': [
		'Not checked: {reason}',
		'Non vérifié : {reason}',
		'Sin comprobar: {reason}',
		'Nicht geprüft: {reason}'
	],
	'accounts.checkFailedToRun': [
		'Courrier could not run the check',
		'Courrier n’a pas pu lancer la vérification',
		'Courrier no pudo ejecutar la comprobación',
		'Courrier konnte die Prüfung nicht ausführen'
	],
	'accounts.noServer': [
		'No {protocol} server set',
		'Aucun serveur {protocol} défini',
		'Ningún servidor {protocol} configurado',
		'Kein {protocol}-Server eingerichtet'
	],
	'accounts.legNotChecked': [
		'{protocol} not checked yet',
		'{protocol} pas encore vérifié',
		'{protocol} aún sin comprobar',
		'{protocol} noch nicht geprüft'
	],
	'accounts.legConnected': ['{protocol} connected', '{protocol} connecté', '{protocol} conectado', '{protocol} verbunden'],
	'accounts.legRejected': [
		'{protocol} rejected the credentials',
		'{protocol} a refusé les identifiants',
		'{protocol} rechazó las credenciales',
		'{protocol} hat die Zugangsdaten abgelehnt'
	],
	'accounts.legUnreachable': [
		'{protocol} unreachable',
		'{protocol} injoignable',
		'{protocol} inaccesible',
		'{protocol} nicht erreichbar'
	],
	'accounts.retryIn': ['retry in {seconds}s', 'réessayez dans {seconds} s', 'reintenta en {seconds} s', 'erneut versuchen in {seconds} s'],
	'accounts.retryNow': ['you can retry now', 'vous pouvez réessayer', 'ya puedes reintentar', 'du kannst es jetzt erneut versuchen'],
	'accounts.checkRunFailed': [
		'The check could not be run',
		'La vérification n’a pas pu être lancée',
		'No se pudo ejecutar la comprobación',
		'Die Prüfung konnte nicht ausgeführt werden'
	],
	'accounts.loadFailed': [
		'Could not load your mail accounts',
		'Impossible de charger vos comptes de messagerie',
		'No se pudieron cargar tus cuentas de correo',
		'Deine E-Mail-Konten konnten nicht geladen werden'
	],
	'accounts.defaultFailed': [
		'Could not change the default account',
		'Impossible de changer le compte par défaut',
		'No se pudo cambiar la cuenta predeterminada',
		'Das Standardkonto konnte nicht geändert werden'
	],
	'accounts.removeFailed': [
		'Could not remove the account',
		'Impossible de supprimer le compte',
		'No se pudo quitar la cuenta',
		'Das Konto konnte nicht entfernt werden'
	],
	'accounts.protocolFailed': ['{protocol} failed', '{protocol} a échoué', '{protocol} falló', '{protocol} fehlgeschlagen'],
	'accounts.couldNotCheck': [
		'Could not be checked',
		'Vérification impossible',
		'No se pudo comprobar',
		'Konnte nicht geprüft werden'
	],
	'accounts.checkConnection': ['Check connection', 'Vérifier la connexion', 'Comprobar la conexión', 'Verbindung prüfen'],
	'accounts.checkHint': [
		'Signs in to IMAP and SMTP with the stored passwords, then hangs up. No mail is read.',
		'Se connecte à IMAP et SMTP avec les mots de passe enregistrés, puis raccroche. Aucun message n’est lu.',
		'Inicia sesión en IMAP y SMTP con las contraseñas guardadas y cuelga. No se lee ningún correo.',
		'Meldet sich mit den gespeicherten Passwörtern bei IMAP und SMTP an und trennt wieder. Es werden keine Mails gelesen.'
	],
	'accounts.makeDefault': ['Set as default', 'Définir par défaut', 'Establecer como predeterminada', 'Als Standard festlegen'],
	'accounts.makeDefaultHint': [
		'The default account is what the CLI and new devices open first.',
		'Le compte par défaut est celui que la CLI et les nouveaux appareils ouvrent en premier.',
		'La cuenta predeterminada es la que abren primero la CLI y los dispositivos nuevos.',
		'Das Standardkonto öffnen die CLI und neue Geräte zuerst.'
	],
	'accounts.defaultNote': [
		'The default account is what the CLI and new devices open first. The mail view has its own switcher.',
		'Le compte par défaut est celui que la CLI et les nouveaux appareils ouvrent en premier. La vue courrier a son propre sélecteur.',
		'La cuenta predeterminada es la que abren primero la CLI y los dispositivos nuevos. La vista de correo tiene su propio selector.',
		'Das Standardkonto öffnen die CLI und neue Geräte zuerst. Die Mailansicht hat einen eigenen Umschalter.'
	],
	'accounts.removeTitle': [
		'Remove this mail account?',
		'Supprimer ce compte de messagerie ?',
		'¿Quitar esta cuenta de correo?',
		'Dieses E-Mail-Konto entfernen?'
	],
	'accounts.removeBody': [
		'Courrier deletes {name}, its stored credentials, and every folder, message and draft it has cached locally. Your mail itself stays on {host}, but re-adding the account means syncing it all again from scratch.',
		'Courrier supprime {name}, ses identifiants enregistrés, ainsi que chaque dossier, message et brouillon mis en cache localement. Votre courrier lui-même reste sur {host}, mais rajouter le compte oblige à tout resynchroniser depuis le début.',
		'Courrier elimina {name}, sus credenciales guardadas y cada carpeta, mensaje y borrador que tiene en caché local. Tu correo en sí permanece en {host}, pero volver a añadir la cuenta obliga a sincronizarlo todo desde cero.',
		'Courrier löscht {name}, die gespeicherten Zugangsdaten sowie jeden lokal zwischengespeicherten Ordner, jede Nachricht und jeden Entwurf. Deine Mails selbst bleiben auf {host}, aber beim erneuten Hinzufügen des Kontos muss alles von vorn synchronisiert werden.'
	],
	'accounts.thisAccount': ['this account', 'ce compte', 'esta cuenta', 'dieses Konto'],
	'accounts.theServer': ['the server', 'le serveur', 'el servidor', 'dem Server'],
	'accounts.removeConfirm': ['Remove account', 'Supprimer le compte', 'Quitar cuenta', 'Konto entfernen'],
	'accounts.keepIt': ['Keep it', 'Le garder', 'Conservarla', 'Behalten'],
	'drawer.editTitle': ['Edit {name}', 'Modifier {name}', 'Editar {name}', '{name} bearbeiten'],
	'drawer.addTitle': ['Add a mail account', 'Ajouter un compte de messagerie', 'Añadir una cuenta de correo', 'E-Mail-Konto hinzufügen'],
	'drawer.displayName': ['Display name', 'Nom affiché', 'Nombre visible', 'Anzeigename'],
	'drawer.displayNameHelp': [
		'What this inbox is called inside Courrier.',
		'Le nom de cette boîte dans Courrier.',
		'Cómo se llama esta bandeja dentro de Courrier.',
		'So heißt dieser Posteingang in Courrier.'
	],
	'drawer.namePlaceholder': ['Work', 'Travail', 'Trabajo', 'Arbeit'],
	'drawer.email': ['Email address', 'Adresse e-mail', 'Dirección de correo', 'E-Mail-Adresse'],
	'drawer.emailHelp': [
		'The address messages are sent from.',
		'L’adresse depuis laquelle les messages sont envoyés.',
		'La dirección desde la que se envían los mensajes.',
		'Die Adresse, von der Nachrichten gesendet werden.'
	],
	'drawer.imap': ['IMAP, reading mail', 'IMAP, lecture du courrier', 'IMAP, lectura del correo', 'IMAP, Mails lesen'],
	'drawer.smtp': ['SMTP, sending mail', 'SMTP, envoi du courrier', 'SMTP, envío del correo', 'SMTP, Mails senden'],
	'drawer.host': ['Host', 'Serveur', 'Servidor', 'Host'],
	'drawer.port': ['Port', 'Port', 'Puerto', 'Port'],
	'drawer.username': ['Username', 'Nom d’utilisateur', 'Usuario', 'Benutzername'],
	'drawer.usernameHelp': [
		'Left empty, the email address is used.',
		'Si le champ est vide, l’adresse e-mail est utilisée.',
		'Si se deja vacío, se usa la dirección de correo.',
		'Bleibt das Feld leer, wird die E-Mail-Adresse verwendet.'
	],
	'drawer.imapPassword': ['IMAP password', 'Mot de passe IMAP', 'Contraseña IMAP', 'IMAP-Passwort'],
	'drawer.smtpPassword': ['SMTP password', 'Mot de passe SMTP', 'Contraseña SMTP', 'SMTP-Passwort'],
	'drawer.storedHelp': [
		'Already stored, encrypted at rest. Type over it to replace it; leave it to keep it.',
		'Déjà enregistré, chiffré au repos. Saisissez-en un autre pour le remplacer ; laissez-le tel quel pour le garder.',
		'Ya guardada, cifrada en reposo. Escribe encima para reemplazarla; déjala para conservarla.',
		'Bereits gespeichert und verschlüsselt. Überschreibe es, um es zu ersetzen; lass es stehen, um es zu behalten.'
	],
	'drawer.imapPasswordHelp': [
		'Encrypted at rest with this instance’s ENCRYPTION_KEY.',
		'Chiffré au repos avec la ENCRYPTION_KEY de cette instance.',
		'Cifrada en reposo con la ENCRYPTION_KEY de esta instancia.',
		'Verschlüsselt gespeichert mit dem ENCRYPTION_KEY dieser Instanz.'
	],
	'drawer.smtpPasswordHelp': [
		'Often the same as the IMAP password, but not always.',
		'Souvent identique au mot de passe IMAP, mais pas toujours.',
		'A menudo es la misma que la contraseña IMAP, pero no siempre.',
		'Oft identisch mit dem IMAP-Passwort, aber nicht immer.'
	],
	'drawer.signature': ['Signature', 'Signature', 'Firma', 'Signatur'],
	'drawer.signatureHelp': [
		'Appended to messages you send from this account.',
		'Ajoutée aux messages que vous envoyez depuis ce compte.',
		'Se añade a los mensajes que envías desde esta cuenta.',
		'Wird an Nachrichten angehängt, die du von diesem Konto sendest.'
	],
	'drawer.signaturePlaceholder': ['Your name', 'Votre nom', 'Tu nombre', 'Dein Name'],
	'drawer.useDefault': [
		'Use as the default account',
		'Utiliser comme compte par défaut',
		'Usar como cuenta predeterminada',
		'Als Standardkonto verwenden'
	],
	'drawer.useDefaultHelp': [
		'The account Courrier opens on and composes from unless you pick another.',
		'Le compte que Courrier ouvre et utilise pour écrire, sauf si vous en choisissez un autre.',
		'La cuenta que Courrier abre y usa para redactar, salvo que elijas otra.',
		'Das Konto, das Courrier öffnet und zum Schreiben verwendet, solange du kein anderes wählst.'
	],
	'drawer.test': ['Test connection', 'Tester la connexion', 'Probar la conexión', 'Verbindung testen'],
	'drawer.testing': ['Testing…', 'Test en cours…', 'Probando…', 'Wird getestet…'],
	'drawer.testConnecting': [
		'Connecting to both servers…',
		'Connexion aux deux serveurs…',
		'Conectando con ambos servidores…',
		'Verbindung zu beiden Servern…'
	],
	'drawer.testPassed': [
		'IMAP and SMTP both answered',
		'IMAP et SMTP ont tous deux répondu',
		'IMAP y SMTP respondieron',
		'IMAP und SMTP haben beide geantwortet'
	],
	'drawer.testFailedShort': ['Test failed', 'Échec du test', 'La prueba falló', 'Test fehlgeschlagen'],
	'drawer.testNeedsPasswords': [
		'Testing sends both passwords to the servers, so it needs them typed in. For an account already saved, use Check on its row instead: that one uses the stored credentials.',
		'Le test envoie les deux mots de passe aux serveurs, il faut donc les saisir. Pour un compte déjà enregistré, utilisez plutôt « Vérifier » sur sa ligne : cette vérification utilise les identifiants enregistrés.',
		'La prueba envía ambas contraseñas a los servidores, así que hay que escribirlas. Para una cuenta ya guardada, usa mejor «Comprobar» en su fila: esa usa las credenciales guardadas.',
		'Der Test sendet beide Passwörter an die Server, sie müssen also eingegeben werden. Für ein bereits gespeichertes Konto nutze stattdessen „Prüfen“ in seiner Zeile: Das verwendet die gespeicherten Zugangsdaten.'
	],
	'drawer.testFailedTitle': [
		'Connection test failed',
		'Échec du test de connexion',
		'La prueba de conexión falló',
		'Verbindungstest fehlgeschlagen'
	],
	'drawer.testError': [
		'Connection test failed',
		'Échec du test de connexion',
		'La prueba de conexión falló',
		'Verbindungstest fehlgeschlagen'
	],
	'drawer.saveFailed': [
		'Could not save the account',
		'Impossible d’enregistrer le compte',
		'No se pudo guardar la cuenta',
		'Das Konto konnte nicht gespeichert werden'
	],
	'drawer.saveAccount': ['Save account', 'Enregistrer le compte', 'Guardar cuenta', 'Konto speichern'],
	'drawer.addAccount': ['Add account', 'Ajouter le compte', 'Añadir cuenta', 'Konto hinzufügen'],
	'drawer.removeAccount': ['Remove account', 'Supprimer le compte', 'Quitar cuenta', 'Konto entfernen']
} satisfies Record<string, Row>;
