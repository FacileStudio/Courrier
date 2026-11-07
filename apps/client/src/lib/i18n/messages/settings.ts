import type { Row } from '../row';

/** The profile, mail accounts, templates, API token and advanced settings pages. */
export const settings = {
	'profile.memberSince': ['Member since', 'Membre depuis', 'Miembro desde', 'Mitglied seit'],
	'profile.activeSpace': ['Active space', 'Espace actif', 'Espacio activo', 'Aktiver Bereich'],
	'profile.signIn': ['Sign-in', 'Connexion', 'Inicio de sesión', 'Anmeldung'],
	'profile.signInValue': ['Facile account', 'Compte Facile', 'Cuenta de Facile', 'Facile-Konto'],
	'profile.rejectType': [
		'that file is not a PNG, JPEG, GIF or WebP image',
		'ce fichier n’est pas une image PNG, JPEG, GIF ou WebP',
		'ese archivo no es una imagen PNG, JPEG, GIF o WebP',
		'diese Datei ist kein PNG-, JPEG-, GIF- oder WebP-Bild'
	],
	'profile.rejectSize': [
		'that file is larger than 5 MB',
		'ce fichier dépasse 5 Mo',
		'ese archivo pesa más de 5 MB',
		'diese Datei ist größer als 5 MB'
	],
	'profile.rejectCount': [
		'only one photo at a time',
		'une seule photo à la fois',
		'solo una foto a la vez',
		'nur ein Foto auf einmal'
	],
	'profile.rejectOther': [
		'that file was rejected',
		'ce fichier a été refusé',
		'ese archivo fue rechazado',
		'diese Datei wurde abgelehnt'
	],
	'profile.notUploaded': [
		'Not uploaded: {reason}.',
		'Non envoyé : {reason}.',
		'No se subió: {reason}.',
		'Nicht hochgeladen: {reason}.'
	],
	'profile.saveFailed': [
		'Could not save your details',
		'Impossible d’enregistrer vos informations',
		'No se pudieron guardar tus datos',
		'Deine Angaben konnten nicht gespeichert werden'
	],
	'profile.uploadFailed': ['Upload failed', 'Échec de l’envoi', 'Error al subir', 'Upload fehlgeschlagen'],
	'profile.removeFailed': [
		'Could not remove the photo',
		'Impossible de retirer la photo',
		'No se pudo quitar la foto',
		'Das Foto konnte nicht entfernt werden'
	],
	'profile.identity': ['Identity', 'Identité', 'Identidad', 'Identität'],
	'profile.identityBody': [
		'How Courrier names you when you send mail.',
		'Comment Courrier vous nomme quand vous envoyez du courrier.',
		'Cómo te nombra Courrier cuando envías correo.',
		'Wie Courrier dich nennt, wenn du Mails sendest.'
	],
	'profile.displayName': ['Display name', 'Nom affiché', 'Nombre visible', 'Anzeigename'],
	'profile.displayNameBody': [
		'Used as the From name on messages you send and on your sidebar card.',
		'Utilisé comme nom d’expéditeur sur vos messages et sur la carte de la barre latérale.',
		'Se usa como nombre del remitente en los mensajes que envías y en la tarjeta de la barra lateral.',
		'Wird als Absendername deiner Nachrichten und auf der Karte in der Seitenleiste verwendet.'
	],
	'profile.namePlaceholder': ['Your name', 'Votre nom', 'Tu nombre', 'Dein Name'],
	'profile.emailBody': [
		'Managed by your Facile account. Change it at sso.facile.studio and it syncs back here.',
		'Géré par votre compte Facile. Modifiez-le sur sso.facile.studio et il se synchronise ici.',
		'Lo gestiona tu cuenta de Facile. Cámbialo en sso.facile.studio y se sincroniza aquí.',
		'Wird von deinem Facile-Konto verwaltet. Ändere es auf sso.facile.studio und es wird hierher synchronisiert.'
	],
	'profile.saved': ['Saved.', 'Enregistré.', 'Guardado.', 'Gespeichert.'],
	'profile.save': ['Save changes', 'Enregistrer les modifications', 'Guardar cambios', 'Änderungen speichern'],
	'profile.photo': ['Photo', 'Photo', 'Foto', 'Foto'],
	'profile.photoBody': [
		'Shown on your sidebar card. Square images crop best.',
		'Affichée sur la carte de la barre latérale. Les images carrées se recadrent le mieux.',
		'Se muestra en la tarjeta de la barra lateral. Las imágenes cuadradas se recortan mejor.',
		'Wird auf der Karte in der Seitenleiste angezeigt. Quadratische Bilder lassen sich am besten zuschneiden.'
	],
	'profile.photoFromFacile': [
		'Your photo comes from your Facile account',
		'Votre photo vient de votre compte Facile',
		'Tu foto viene de tu cuenta de Facile',
		'Dein Foto stammt aus deinem Facile-Konto'
	],
	'profile.photoFromFacileBody': [
		'Change it in your {account} and it shows here on the next page load. Courrier keeps no copy, so there is nothing to upload or delete on this side.',
		'Modifiez-la dans votre {account} et elle apparaît ici au prochain chargement de la page. Courrier n’en garde aucune copie, il n’y a donc rien à envoyer ni à supprimer de ce côté.',
		'Cámbiala en tu {account} y aparecerá aquí en la próxima carga de la página. Courrier no guarda ninguna copia, así que no hay nada que subir ni borrar de este lado.',
		'Ändere es in deinem {account} und es erscheint hier beim nächsten Laden der Seite. Courrier speichert keine Kopie, daher gibt es hier nichts hochzuladen oder zu löschen.'
	],
	'profile.facileAccount': ['Facile account', 'compte Facile', 'cuenta de Facile', 'Facile-Konto'],
	'profile.drop': ['Drop a photo here', 'Déposez une photo ici', 'Suelta una foto aquí', 'Foto hier ablegen'],
	'profile.dropHint': [
		'PNG, JPEG, GIF or WebP, up to 5 MB',
		'PNG, JPEG, GIF ou WebP, jusqu’à 5 Mo',
		'PNG, JPEG, GIF o WebP, hasta 5 MB',
		'PNG, JPEG, GIF oder WebP, bis zu 5 MB'
	],
	'profile.removePhoto': ['Remove photo', 'Retirer la photo', 'Quitar la foto', 'Foto entfernen'],
	'profile.removing': ['Removing…', 'Suppression…', 'Quitando…', 'Wird entfernt…'],
	'profile.session': ['Session', 'Session', 'Sesión', 'Sitzung'],
	'profile.sessionBody': [
		'Signed in through your Facile account at sso.facile.studio.',
		'Connecté avec votre compte Facile sur sso.facile.studio.',
		'Sesión iniciada con tu cuenta de Facile en sso.facile.studio.',
		'Angemeldet über dein Facile-Konto auf sso.facile.studio.'
	],
	'profile.signOutBody': [
		'Ends this session in this browser. Other browsers stay signed in.',
		'Termine la session dans ce navigateur. Les autres navigateurs restent connectés.',
		'Cierra la sesión en este navegador. Los demás navegadores siguen con la sesión iniciada.',
		'Beendet die Sitzung in diesem Browser. Andere Browser bleiben angemeldet.'
	],
	'profile.signingOut': ['Signing out…', 'Déconnexion…', 'Cerrando sesión…', 'Abmeldung…'],

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
	'drawer.removeAccount': ['Remove account', 'Supprimer le compte', 'Quitar cuenta', 'Konto entfernen'],

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
	'api.revokeConfirm': ['Revoke token', 'Révoquer le jeton', 'Revocar token', 'Token widerrufen'],

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
