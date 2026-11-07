import type { Row } from '../row';

/** Strings shared by the whole app: buttons, folders, the shell, sign-in and the first-run screens. */
export const shell = {
	'common.cancel': ['Cancel', 'Annuler', 'Cancelar', 'Abbrechen'],
	'common.saving': ['Saving…', 'Enregistrement…', 'Guardando…', 'Wird gespeichert…'],
	'common.delete': ['Delete', 'Supprimer', 'Eliminar', 'Löschen'],
	'common.edit': ['Edit', 'Modifier', 'Editar', 'Bearbeiten'],
	'common.open': ['Open', 'Ouvrir', 'Abrir', 'Öffnen'],
	'common.back': ['Back', 'Retour', 'Atrás', 'Zurück'],
	'common.next': ['Next', 'Suivant', 'Siguiente', 'Weiter'],
	'common.skip': ['Skip', 'Passer', 'Omitir', 'Überspringen'],
	'common.revoke': ['Revoke', 'Révoquer', 'Revocar', 'Widerrufen'],
	'common.email': ['Email', 'E-mail', 'Correo electrónico', 'E-Mail'],
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
	'error.uploadNetwork': [
		'Upload failed. The request never arrived.',
		'L’envoi a échoué. La requête n’est jamais arrivée.',
		'La subida falló. La solicitud nunca llegó.',
		'Upload fehlgeschlagen. Die Anfrage ist nie angekommen.'
	],
	'error.uploadCancelled': ['Upload cancelled', 'Envoi annulé', 'Subida cancelada', 'Upload abgebrochen'],
	'error.uploadNoProfile': [
		'Upload succeeded but the server sent no profile back',
		'L’envoi a réussi mais le serveur n’a renvoyé aucun profil',
		'La subida se completó pero el servidor no devolvió ningún perfil',
		'Upload erfolgreich, aber der Server hat kein Profil zurückgeschickt'
	],
	'error.uploadFailed': [
		'Upload failed with status {status}',
		'L’envoi a échoué (statut {status})',
		'La subida falló con el estado {status}',
		'Upload fehlgeschlagen (Status {status})'
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
	'shell.account': ['Account', 'Compte', 'Cuenta', 'Konto'],
	'shell.accountMenu': ['{name}, account menu', '{name}, menu du compte', '{name}, menú de la cuenta', '{name}, Kontomenü'],
	'shell.settings': ['Settings', 'Paramètres', 'Ajustes', 'Einstellungen'],
	'shell.desktopApp': ['Desktop app', 'Application de bureau', 'Aplicación de escritorio', 'Desktop-App'],
	'shell.signOut': ['Sign out', 'Se déconnecter', 'Cerrar sesión', 'Abmelden'],
	'shell.switchAccount': [
		'Switch account, current {email}',
		'Changer de compte, actuel : {email}',
		'Cambiar de cuenta, actual: {email}',
		'Konto wechseln, aktuell: {email}'
	],
	'shell.manageAccounts': ['Manage accounts', 'Gérer les comptes', 'Gestionar cuentas', 'Konten verwalten'],

	'login.title': ['Sign in — Courrier', 'Connexion — Courrier', 'Iniciar sesión — Courrier', 'Anmelden — Courrier'],
	'login.tagline': [
		'Your mail. Your server.',
		'Votre courrier. Votre serveur.',
		'Tu correo. Tu servidor.',
		'Deine Mails. Dein Server.'
	],
	'login.failed': [
		'Sign-in failed. Try again.',
		'La connexion a échoué. Réessayez.',
		'No se pudo iniciar sesión. Inténtalo de nuevo.',
		'Die Anmeldung ist fehlgeschlagen. Versuche es erneut.'
	],
	'login.nativeFailed': [
		'Could not open the sign-in page.',
		'Impossible d’ouvrir la page de connexion.',
		'No se pudo abrir la página de inicio de sesión.',
		'Die Anmeldeseite konnte nicht geöffnet werden.'
	],
	'login.nativeTitle': ['Sign in to Courrier', 'Connexion à Courrier', 'Inicia sesión en Courrier', 'Bei Courrier anmelden'],
	'login.nativeBody': [
		'Use your Facile account. Your mailboxes and spaces follow you from the web.',
		'Utilisez votre compte Facile. Vos boîtes mail et vos espaces vous suivent depuis le web.',
		'Usa tu cuenta de Facile. Tus buzones y espacios te siguen desde la web.',
		'Nutze dein Facile-Konto. Deine Postfächer und Bereiche folgen dir vom Web.'
	],
	'login.nativeHint': [
		'Sign-in opens in your browser, then brings you back here.',
		'La connexion s’ouvre dans votre navigateur, puis vous ramène ici.',
		'El inicio de sesión se abre en tu navegador y luego te trae de vuelta aquí.',
		'Die Anmeldung öffnet sich im Browser und bringt dich dann hierher zurück.'
	],
	'welcome.title': ['Welcome — Courrier', 'Bienvenue — Courrier', 'Te damos la bienvenida — Courrier', 'Willkommen — Courrier'],
	'welcome.unreachable': [
		'Courrier is not answering. Reload the page to try again.',
		'Courrier ne répond pas. Rechargez la page pour réessayer.',
		'Courrier no responde. Recarga la página para volver a intentarlo.',
		'Courrier antwortet nicht. Lade die Seite neu, um es erneut zu versuchen.'
	],

	'landing.title': [
		'Courrier — Self-hosted email for creative studios',
		'Courrier — La messagerie auto-hébergée des studios créatifs',
		'Courrier — Correo autoalojado para estudios creativos',
		'Courrier — Selbst gehostete E-Mail für Kreativstudios'
	],
	'landing.description': [
		'A self-hosted email client for creative studios. Connect your IMAP accounts, read and send from a clean interface.',
		'Un client de messagerie auto-hébergé pour les studios créatifs. Connectez vos comptes IMAP, lisez et envoyez depuis une interface claire.',
		'Un cliente de correo autoalojado para estudios creativos. Conecta tus cuentas IMAP, lee y envía desde una interfaz limpia.',
		'Ein selbst gehosteter E-Mail-Client für Kreativstudios. Verbinde deine IMAP-Konten, lies und sende in einer aufgeräumten Oberfläche.'
	],
	'landing.tagline': [
		'Self-hosted email for creative studios.',
		'La messagerie auto-hébergée des studios créatifs.',
		'Correo autoalojado para estudios creativos.',
		'Selbst gehostete E-Mail für Kreativstudios.'
	],
	'landing.heroTitle': [
		'Your mail. Your server.',
		'Votre courrier. Votre serveur.',
		'Tu correo. Tu servidor.',
		'Deine Mails. Dein Server.'
	],
	'landing.heroBody': [
		'Courrier is a self-hosted email client. Connect your IMAP accounts, read and send from a clean interface. No cloud, no tracking, no compromise.',
		'Courrier est un client de messagerie auto-hébergé. Connectez vos comptes IMAP, lisez et envoyez depuis une interface claire. Pas de cloud, pas de pistage, pas de compromis.',
		'Courrier es un cliente de correo autoalojado. Conecta tus cuentas IMAP, lee y envía desde una interfaz limpia. Sin nube, sin rastreo, sin concesiones.',
		'Courrier ist ein selbst gehosteter E-Mail-Client. Verbinde deine IMAP-Konten, lies und sende in einer aufgeräumten Oberfläche. Keine Cloud, kein Tracking, keine Kompromisse.'
	],
	'landing.downloadCta': [
		'Get the desktop and Android apps',
		'Obtenir les applications bureau et Android',
		'Descargar las apps de escritorio y Android',
		'Desktop- und Android-Apps herunterladen'
	],
	'landing.imapTitle': ['Any IMAP account', 'Tout compte IMAP', 'Cualquier cuenta IMAP', 'Jedes IMAP-Konto'],
	'landing.imapBody': [
		'Point Courrier at your IMAP and SMTP servers. Credentials are encrypted at rest.',
		'Indiquez à Courrier vos serveurs IMAP et SMTP. Les identifiants sont chiffrés au repos.',
		'Indica a Courrier tus servidores IMAP y SMTP. Las credenciales se cifran en reposo.',
		'Gib Courrier deine IMAP- und SMTP-Server an. Zugangsdaten werden verschlüsselt gespeichert.'
	],
	'landing.spacesTitle': ['Shared spaces', 'Espaces partagés', 'Espacios compartidos', 'Gemeinsame Bereiche'],
	'landing.spacesBody': [
		'Group accounts into spaces and invite the people who need to read the same inbox.',
		'Regroupez des comptes dans des espaces et invitez les personnes qui doivent lire la même boîte.',
		'Agrupa cuentas en espacios e invita a quienes necesitan leer la misma bandeja.',
		'Fasse Konten in Bereichen zusammen und lade die Personen ein, die denselben Posteingang lesen müssen.'
	],
	'landing.selfHostedTitle': ['Self-hosted', 'Auto-hébergé', 'Autoalojado', 'Selbst gehostet'],
	'landing.selfHostedBody': [
		'One binary and a Postgres database. Your mail stays on your server.',
		'Un seul binaire et une base Postgres. Votre courrier reste sur votre serveur.',
		'Un solo binario y una base de datos Postgres. Tu correo se queda en tu servidor.',
		'Eine Binärdatei und eine Postgres-Datenbank. Deine Mails bleiben auf deinem Server.'
	],
	'landing.pricing': ['Pricing', 'Tarifs', 'Precios', 'Preise'],
	'landing.ctaTitle': ['Ready to start?', 'Prêt à commencer ?', '¿Listo para empezar?', 'Bereit loszulegen?'],
	'landing.ctaBody': [
		'Sign in with your Facile account to open Courrier.',
		'Connectez-vous avec votre compte Facile pour ouvrir Courrier.',
		'Inicia sesión con tu cuenta de Facile para abrir Courrier.',
		'Melde dich mit deinem Facile-Konto an, um Courrier zu öffnen.'
	],
	'landing.navFeatures': ['Features', 'Fonctionnalités', 'Funciones', 'Funktionen'],
	'landing.navPricing': ['Pricing', 'Tarifs', 'Precios', 'Preise'],
	'landing.navApps': ['Get the apps', 'Les applications', 'Las apps', 'Die Apps'],
	'downloads.title': ['Get the app — Courrier', 'Obtenir l’application — Courrier', 'Descargar la app — Courrier', 'App herunterladen — Courrier'],
	'downloads.description': [
		'Download Courrier for Android, or use it in any browser.',
		'Téléchargez Courrier pour Android, ou utilisez-le dans n’importe quel navigateur.',
		'Descarga Courrier para Android o úsalo en cualquier navegador.',
		'Lade Courrier für Android herunter oder nutze es in jedem Browser.'
	],
	'downloads.home': ['Courrier home', 'Accueil Courrier', 'Inicio de Courrier', 'Courrier Startseite'],
	'downloads.open': ['Open Courrier', 'Ouvrir Courrier', 'Abrir Courrier', 'Courrier öffnen'],
	'downloads.heading': ['Get Courrier', 'Obtenir Courrier', 'Descargar Courrier', 'Courrier herunterladen'],
	'downloads.intro': [
		'The apps read the same mailboxes as the web. Sign in once with your Facile account.',
		'Les applications lisent les mêmes boîtes que le web. Connectez-vous une fois avec votre compte Facile.',
		'Las apps leen los mismos buzones que la web. Inicia sesión una vez con tu cuenta de Facile.',
		'Die Apps lesen dieselben Postfächer wie das Web. Melde dich einmal mit deinem Facile-Konto an.'
	],
	'downloads.desktop': ['Desktop', 'Bureau', 'Escritorio', 'Desktop'],
	'downloads.desktopBody': [
		'Its own window, signed installers, and it updates itself.',
		'Sa propre fenêtre, des installateurs signés, et des mises à jour automatiques.',
		'Su propia ventana, instaladores firmados y se actualiza solo.',
		'Ein eigenes Fenster, signierte Installationsprogramme und automatische Updates.'
	],
	'downloads.detailMacos': ['Apple silicon and Intel', 'Apple silicon et Intel', 'Apple silicon e Intel', 'Apple Silicon und Intel'],
	'downloads.detailWindows': [
		'Windows 10 or later, 64-bit',
		'Windows 10 ou plus récent, 64 bits',
		'Windows 10 o posterior, 64 bits',
		'Windows 10 oder neuer, 64 Bit'
	],
	'downloads.detailLinux': ['AppImage, x86-64', 'AppImage, x86-64', 'AppImage, x86-64', 'AppImage, x86-64'],
	'downloads.download': ['Download', 'Télécharger', 'Descargar', 'Herunterladen'],
	'downloads.debNote': [
		'On Debian or Ubuntu, the {link} installs Courrier system-wide. To check a file, compare it with the {checksums}.',
		'Sur Debian ou Ubuntu, le {link} installe Courrier pour tout le système. Pour vérifier un fichier, comparez-le aux {checksums}.',
		'En Debian o Ubuntu, el {link} instala Courrier para todo el sistema. Para comprobar un archivo, compáralo con las {checksums}.',
		'Unter Debian oder Ubuntu installiert das {link} Courrier systemweit. Um eine Datei zu prüfen, vergleiche sie mit den {checksums}.'
	],
	'downloads.debLink': ['.deb package', 'paquet .deb', 'paquete .deb', '.deb-Paket'],
	'downloads.checksumsLink': ['published checksums', 'sommes de contrôle publiées', 'sumas de verificación publicadas', 'veröffentlichten Prüfsummen'],
	'downloads.android': ['Android', 'Android', 'Android', 'Android'],
	'downloads.androidBody': [
		'Signed APK, Android 7 or later.',
		'APK signé, Android 7 ou plus récent.',
		'APK firmado, Android 7 o posterior.',
		'Signierte APK, Android 7 oder neuer.'
	],
	'downloads.apk': ['Download the APK', 'Télécharger l’APK', 'Descargar el APK', 'APK herunterladen'],
	'downloads.step1Title': ['Download the APK', 'Télécharger l’APK', 'Descargar el APK', 'APK herunterladen'],
	'downloads.step1Body': [
		'Tap the button above on your phone. The file is about 4 MB.',
		'Touchez le bouton ci-dessus sur votre téléphone. Le fichier pèse environ 4 Mo.',
		'Toca el botón de arriba en tu teléfono. El archivo pesa unos 4 MB.',
		'Tippe auf dem Telefon auf die Schaltfläche oben. Die Datei ist etwa 4 MB groß.'
	],
	'downloads.step2Title': ['Allow this source', 'Autoriser cette source', 'Permitir este origen', 'Diese Quelle erlauben'],
	'downloads.step2Body': [
		'Android asks once. Allow your browser to install unknown apps, then come back.',
		'Android le demande une seule fois. Autorisez votre navigateur à installer des applications inconnues, puis revenez.',
		'Android lo pide una sola vez. Permite que tu navegador instale apps desconocidas y vuelve.',
		'Android fragt nur einmal. Erlaube deinem Browser, unbekannte Apps zu installieren, und komm zurück.'
	],
	'downloads.step3Title': ['Open and sign in', 'Ouvrir et se connecter', 'Abrir e iniciar sesión', 'Öffnen und anmelden'],
	'downloads.step3Body': [
		'Open the downloaded file, install it, and sign in with your Facile account.',
		'Ouvrez le fichier téléchargé, installez-le et connectez-vous avec votre compte Facile.',
		'Abre el archivo descargado, instálalo e inicia sesión con tu cuenta de Facile.',
		'Öffne die heruntergeladene Datei, installiere sie und melde dich mit deinem Facile-Konto an.'
	],
	'downloads.apkNote': [
		'To check the file, compare {command} with the {checksum}. Older versions are on the {releases}.',
		'Pour vérifier le fichier, comparez {command} avec la {checksum}. Les versions précédentes sont sur la {releases}.',
		'Para comprobar el archivo, compara {command} con la {checksum}. Las versiones anteriores están en la {releases}.',
		'Um die Datei zu prüfen, vergleiche {command} mit der {checksum}. Ältere Versionen findest du auf der {releases}.'
	],
	'downloads.checksumLink': ['published checksum', 'somme de contrôle publiée', 'suma de verificación publicada', 'veröffentlichten Prüfsumme'],
	'downloads.releasesLink': ['releases page', 'page des versions', 'página de versiones', 'Seite mit den Versionen'],
	'downloads.browser': ['Any browser', 'N’importe quel navigateur', 'Cualquier navegador', 'Jeder Browser'],
	'downloads.browserBody': [
		'No install needed: Courrier works in any modern browser, on any computer or phone.',
		'Aucune installation : Courrier fonctionne dans tout navigateur récent, sur ordinateur comme sur téléphone.',
		'No hace falta instalar nada: Courrier funciona en cualquier navegador moderno, en ordenador o teléfono.',
		'Keine Installation nötig: Courrier läuft in jedem modernen Browser, auf Computer und Telefon.'
	],
	'downloads.openWeb': ['Open the web app', 'Ouvrir l’application web', 'Abrir la aplicación web', 'Web-App öffnen'],

	'intro.slide1Title': [
		'Your mail, in your pocket.',
		'Votre courrier, dans votre poche.',
		'Tu correo, en tu bolsillo.',
		'Deine Mails in der Hosentasche.'
	],
	'intro.slide1Body': [
		'Every IMAP mailbox you connected to Courrier, threaded, searchable and ready to answer.',
		'Chaque boîte IMAP connectée à Courrier, en conversations, consultable et prête à répondre.',
		'Cada buzón IMAP que conectaste a Courrier, en hilos, con búsqueda y listo para responder.',
		'Jedes IMAP-Postfach, das du mit Courrier verbunden hast: als Unterhaltung, durchsuchbar und bereit zum Antworten.'
	],
	'intro.slide2Title': [
		'Only the mail that matters.',
		'Seulement le courrier qui compte.',
		'Solo el correo que importa.',
		'Nur die Mails, die zählen.'
	],
	'intro.slide2Body': [
		'Newsletters and robots stay quiet. A notification means a person wrote to you.',
		'Les newsletters et les robots restent silencieux. Une notification signifie qu’une personne vous a écrit.',
		'Los boletines y los robots se quedan callados. Una notificación significa que una persona te escribió.',
		'Newsletter und Roboter bleiben still. Eine Benachrichtigung heißt: Eine Person hat dir geschrieben.'
	],
	'intro.slide3Title': [
		'Your server, your rules.',
		'Votre serveur, vos règles.',
		'Tu servidor, tus reglas.',
		'Dein Server, deine Regeln.'
	],
	'intro.slide3Body': [
		'Courrier is self-hosted. Your messages stay on your server, never on ours.',
		'Courrier est auto-hébergé. Vos messages restent sur votre serveur, jamais sur le nôtre.',
		'Courrier es autoalojado. Tus mensajes se quedan en tu servidor, nunca en el nuestro.',
		'Courrier wird selbst gehostet. Deine Nachrichten bleiben auf deinem Server, nie auf unserem.'
	],
	'intro.carousel': ['carousel', 'carrousel', 'carrusel', 'Karussell'],
	'intro.slide': ['slide', 'diapositive', 'diapositiva', 'Folie'],
	'intro.slides': ['Slides', 'Diapositives', 'Diapositivas', 'Folien'],
	'intro.slideOf': ['{index} of {total}', '{index} sur {total}', '{index} de {total}', '{index} von {total}'],
	'intro.goTo': ['Go to slide {index}', 'Aller à la diapositive {index}', 'Ir a la diapositiva {index}', 'Zu Folie {index}'],
	'intro.getStarted': ['Get started', 'Commencer', 'Empezar', 'Los geht’s'],

	'window.label': ['Window', 'Fenêtre', 'Ventana', 'Fenster'],
	'window.close': ['Close', 'Fermer', 'Cerrar', 'Schließen'],
	'window.minimize': ['Minimize', 'Réduire', 'Minimizar', 'Minimieren'],
	'window.maximize': ['Maximize', 'Agrandir', 'Maximizar', 'Maximieren'],

	'onboarding.welcomeTitle': ['Welcome, {name}.', 'Bienvenue, {name}.', 'Te damos la bienvenida, {name}.', 'Willkommen, {name}.'],
	'onboarding.welcomeBody': [
		'Three quick steps: your mailbox, your notifications, and you are reading mail.',
		'Trois étapes rapides : votre boîte mail, vos notifications, et vous lisez votre courrier.',
		'Tres pasos rápidos: tu buzón, tus notificaciones, y ya estás leyendo tu correo.',
		'Drei schnelle Schritte: dein Postfach, deine Benachrichtigungen, und schon liest du Mails.'
	],
	'onboarding.mailboxTitle': ['Your mailbox', 'Votre boîte mail', 'Tu buzón', 'Dein Postfach'],
	'onboarding.mailboxBody': [
		'Connect the address you want to read, or pick the one Courrier opens first.',
		'Connectez l’adresse que vous voulez lire, ou choisissez celle que Courrier ouvre en premier.',
		'Conecta la dirección que quieres leer o elige la que Courrier abre primero.',
		'Verbinde die Adresse, die du lesen willst, oder wähle die, die Courrier zuerst öffnet.'
	],
	'onboarding.notificationsTitle': [
		'Know when mail arrives',
		'Sachez quand un message arrive',
		'Entérate cuando llegue correo',
		'Erfahre, wenn Post ankommt'
	],
	'onboarding.notificationsBody': [
		'Courrier tells you when a person writes to you. Newsletters and robots stay silent.',
		'Courrier vous prévient quand une personne vous écrit. Les newsletters et les robots restent silencieux.',
		'Courrier te avisa cuando una persona te escribe. Los boletines y los robots se quedan callados.',
		'Courrier sagt dir Bescheid, wenn dir eine Person schreibt. Newsletter und Roboter bleiben still.'
	],
	'onboarding.doneTitle': ['You are all set.', 'Tout est prêt.', 'Todo listo.', 'Alles bereit.'],
	'onboarding.doneBody': [
		'Your inbox is ready. You can change any of this later in Settings.',
		'Votre boîte de réception est prête. Vous pourrez tout modifier plus tard dans les paramètres.',
		'Tu bandeja de entrada está lista. Puedes cambiar todo esto más tarde en Ajustes.',
		'Dein Posteingang ist bereit. Das kannst du später alles in den Einstellungen ändern.'
	],
	'onboarding.stepWelcome': ['Welcome', 'Bienvenue', 'Bienvenida', 'Willkommen'],
	'onboarding.stepMailbox': ['Mailbox', 'Boîte mail', 'Buzón', 'Postfach'],
	'onboarding.stepNotifications': ['Notifications', 'Notifications', 'Notificaciones', 'Benachrichtigungen'],
	'onboarding.stepDone': ['Ready', 'Prêt', 'Listo', 'Bereit'],
	'onboarding.stepOf': ['Step {index} of {total}', 'Étape {index} sur {total}', 'Paso {index} de {total}', 'Schritt {index} von {total}'],
	'onboarding.stepText': [
		'Step {index} of {total}: {label}',
		'Étape {index} sur {total} : {label}',
		'Paso {index} de {total}: {label}',
		'Schritt {index} von {total}: {label}'
	],
	'onboarding.progress': ['Setup progress', 'Progression de la configuration', 'Progreso de la configuración', 'Fortschritt der Einrichtung'],
	'onboarding.primaryWelcome': ['Set up my inbox', 'Configurer ma boîte', 'Configurar mi bandeja', 'Meinen Posteingang einrichten'],
	'onboarding.primaryAdd': ['Add a mail account', 'Ajouter un compte mail', 'Añadir una cuenta de correo', 'E-Mail-Konto hinzufügen'],
	'onboarding.primaryContinue': ['Continue', 'Continuer', 'Continuar', 'Weiter'],
	'onboarding.primaryDone': ['Open my inbox', 'Ouvrir ma boîte', 'Abrir mi bandeja', 'Meinen Posteingang öffnen'],
	'onboarding.skipSetup': ['Skip setup', 'Passer la configuration', 'Omitir la configuración', 'Einrichtung überspringen'],
	'onboarding.skipForNow': ['Skip for now', 'Plus tard', 'Ahora no', 'Jetzt nicht'],
	'onboarding.saveFailed': [
		'Could not save your choice',
		'Impossible d’enregistrer votre choix',
		'No se pudo guardar tu elección',
		'Deine Auswahl konnte nicht gespeichert werden'
	],
	'onboarding.noMailboxTitle': ['No mailbox yet', 'Aucune boîte mail pour l’instant', 'Aún no hay ningún buzón', 'Noch kein Postfach'],
	'onboarding.noMailboxBody': [
		'Add the IMAP and SMTP details of the address you want to read here. Courrier encrypts the passwords before it stores them.',
		'Ajoutez les paramètres IMAP et SMTP de l’adresse que vous voulez lire ici. Courrier chiffre les mots de passe avant de les enregistrer.',
		'Añade los datos IMAP y SMTP de la dirección que quieres leer aquí. Courrier cifra las contraseñas antes de guardarlas.',
		'Gib hier die IMAP- und SMTP-Daten der Adresse an, die du lesen willst. Courrier verschlüsselt die Passwörter vor dem Speichern.'
	],
	'onboarding.connected': ['Connected', 'Connecté', 'Conectado', 'Verbunden'],
	'onboarding.firstMailbox': [
		'Mailbox to open first',
		'Boîte mail à ouvrir en premier',
		'Buzón que se abre primero',
		'Postfach, das zuerst geöffnet wird'
	],
	'onboarding.addAnother': ['Add another mailbox', 'Ajouter une autre boîte mail', 'Añadir otro buzón', 'Weiteres Postfach hinzufügen'],

	'notify.on': ['Notifications are on', 'Les notifications sont activées', 'Las notificaciones están activadas', 'Benachrichtigungen sind aktiviert'],
	'notify.enableFailed': [
		'Could not turn on notifications',
		'Impossible d’activer les notifications',
		'No se pudieron activar las notificaciones',
		'Benachrichtigungen konnten nicht aktiviert werden'
	],
	'notify.testSent': ['Test notification sent', 'Notification de test envoyée', 'Notificación de prueba enviada', 'Testbenachrichtigung gesendet'],
	'notify.testFailed': [
		'Could not send the test',
		'Impossible d’envoyer le test',
		'No se pudo enviar la prueba',
		'Der Test konnte nicht gesendet werden'
	],
	'notify.notSetUpPhone': [
		'Phone notifications are not set up on this server yet.',
		'Les notifications du téléphone ne sont pas encore configurées sur ce serveur.',
		'Las notificaciones del teléfono aún no están configuradas en este servidor.',
		'Telefon-Benachrichtigungen sind auf diesem Server noch nicht eingerichtet.'
	],
	'notify.notSetUpBrowser': [
		'Browser notifications are not set up on this server yet.',
		'Les notifications du navigateur ne sont pas encore configurées sur ce serveur.',
		'Las notificaciones del navegador aún no están configuradas en este servidor.',
		'Browser-Benachrichtigungen sind auf diesem Server noch nicht eingerichtet.'
	],
	'notify.unsupported': [
		'This browser cannot show notifications.',
		'Ce navigateur ne peut pas afficher de notifications.',
		'Este navegador no puede mostrar notificaciones.',
		'Dieser Browser kann keine Benachrichtigungen anzeigen.'
	],
	'notify.blockedNative': [
		'Notifications are blocked. Turn them on in Settings > Apps > Courrier > Notifications.',
		'Les notifications sont bloquées. Activez-les dans Paramètres > Applications > Courrier > Notifications.',
		'Las notificaciones están bloqueadas. Actívalas en Ajustes > Aplicaciones > Courrier > Notificaciones.',
		'Benachrichtigungen sind blockiert. Aktiviere sie unter Einstellungen > Apps > Courrier > Benachrichtigungen.'
	],
	'notify.blockedDesktop': [
		'Notifications are blocked. Allow Courrier in your system notification settings.',
		'Les notifications sont bloquées. Autorisez Courrier dans les réglages de notification du système.',
		'Las notificaciones están bloqueadas. Permite Courrier en los ajustes de notificaciones del sistema.',
		'Benachrichtigungen sind blockiert. Erlaube Courrier in den Benachrichtigungseinstellungen des Systems.'
	],
	'notify.blockedBrowser': [
		'Notifications are blocked for this site. Allow them in your browser settings.',
		'Les notifications sont bloquées pour ce site. Autorisez-les dans les réglages de votre navigateur.',
		'Las notificaciones están bloqueadas para este sitio. Permítelas en los ajustes de tu navegador.',
		'Benachrichtigungen sind für diese Website blockiert. Erlaube sie in den Browsereinstellungen.'
	],
	'notify.turnOn': ['Turn on notifications', 'Activer les notifications', 'Activar las notificaciones', 'Benachrichtigungen aktivieren'],
	'notify.onThisDevice': ['On for this device', 'Activées sur cet appareil', 'Activadas en este dispositivo', 'Auf diesem Gerät aktiviert'],
	'notify.sendTest': ['Send a test', 'Envoyer un test', 'Enviar una prueba', 'Test senden'],

	'plan.reached': ['Plan limit reached', 'Limite du forfait atteinte', 'Límite del plan alcanzado', 'Limit des Tarifs erreicht'],
	'plan.fallback': [
		'Your plan limit is reached.',
		'La limite de votre forfait est atteinte.',
		'Has alcanzado el límite de tu plan.',
		'Das Limit deines Tarifs ist erreicht.'
	],
	'plan.upgrade': ['Upgrade the plan', 'Passer à un forfait supérieur', 'Mejorar el plan', 'Tarif upgraden'],
	'plan.askOwner': [
		'Ask a space owner to upgrade its plan in Registre.',
		'Demandez à un propriétaire de l’espace de changer son forfait dans Registre.',
		'Pide a un propietario del espacio que mejore su plan en Registre.',
		'Bitte einen Eigentümer des Bereichs, den Tarif in Registre zu upgraden.'
	],

	'settings.title': ['Settings — Courrier', 'Paramètres — Courrier', 'Ajustes — Courrier', 'Einstellungen — Courrier'],
	'settings.heading': ['Settings', 'Paramètres', 'Ajustes', 'Einstellungen'],
	'settings.intro': [
		'You, this browser, and the mail servers Courrier talks to on your behalf.',
		'Vous, ce navigateur, et les serveurs de messagerie auxquels Courrier se connecte à votre place.',
		'Tú, este navegador y los servidores de correo con los que Courrier habla en tu nombre.',
		'Du, dieser Browser und die Mailserver, mit denen Courrier in deinem Namen spricht.'
	],
	'settings.sections': ['Settings sections', 'Sections des paramètres', 'Secciones de ajustes', 'Einstellungsbereiche'],
	'settings.profile': ['Profile', 'Profil', 'Perfil', 'Profil'],
	'settings.appearance': ['Appearance', 'Apparence', 'Apariencia', 'Darstellung'],
	'settings.accounts': ['Accounts', 'Comptes', 'Cuentas', 'Konten'],
	'settings.notifications': ['Notifications', 'Notifications', 'Notificaciones', 'Benachrichtigungen'],
	'settings.templates': ['Templates', 'Modèles', 'Plantillas', 'Vorlagen'],
	'settings.api': ['API', 'API', 'API', 'API'],
	'settings.advanced': ['Advanced', 'Avancé', 'Avanzado', 'Erweitert'],
	'settings.updates': ['Updates', 'Mises à jour', 'Actualizaciones', 'Updates'],

	'appearance.theme': ['Theme', 'Thème', 'Tema', 'Design'],
	'appearance.themeBody': [
		'Stored in this browser. Every device you sign in from keeps its own choice.',
		'Enregistré dans ce navigateur. Chaque appareil depuis lequel vous vous connectez garde son propre choix.',
		'Se guarda en este navegador. Cada dispositivo desde el que inicias sesión conserva su propia elección.',
		'In diesem Browser gespeichert. Jedes Gerät, mit dem du dich anmeldest, behält seine eigene Wahl.'
	],
	'appearance.scheme': ['Colour scheme', 'Schéma de couleurs', 'Esquema de color', 'Farbschema'],
	'appearance.schemeBody': [
		'System follows your operating system, and keeps following it while the app is open.',
		'« Système » suit votre système d’exploitation, même pendant que l’application est ouverte.',
		'«Sistema» sigue tu sistema operativo y lo sigue mientras la app está abierta.',
		'„System“ folgt deinem Betriebssystem und folgt ihm auch, solange die App geöffnet ist.'
	],
	'appearance.system': ['System', 'Système', 'Sistema', 'System'],
	'appearance.light': ['Light', 'Clair', 'Claro', 'Hell'],
	'appearance.dark': ['Dark', 'Sombre', 'Oscuro', 'Dunkel'],
	'appearance.applied': ['Currently applied', 'Actuellement appliqué', 'Aplicado ahora', 'Derzeit angewendet'],
	'appearance.appliedBody': [
		'What System resolved to, or the scheme you picked.',
		'Ce que « Système » a donné, ou le schéma que vous avez choisi.',
		'Lo que resolvió «Sistema» o el esquema que elegiste.',
		'Was „System“ ergeben hat oder das Schema, das du gewählt hast.'
	],

	'notifySettings.title': ['Notifications', 'Notifications', 'Notificaciones', 'Benachrichtigungen'],
	'notifySettings.body': [
		'Courrier checks your inboxes every few minutes and tells you when a person writes to you. Newsletters, lists and robots never notify.',
		'Courrier vérifie vos boîtes de réception toutes les quelques minutes et vous prévient quand une personne vous écrit. Les newsletters, les listes et les robots ne notifient jamais.',
		'Courrier revisa tus bandejas de entrada cada pocos minutos y te avisa cuando una persona te escribe. Los boletines, las listas y los robots nunca notifican.',
		'Courrier prüft deine Posteingänge alle paar Minuten und meldet sich, wenn dir eine Person schreibt. Newsletter, Listen und Roboter benachrichtigen nie.'
	],
	'notifySettings.newMail': ['New mail', 'Nouveau courrier', 'Correo nuevo', 'Neue Mails'],
	'notifySettings.newMailBody': [
		'Stored per device. Turn it on for {device}; your other devices keep their own choice.',
		'Enregistré par appareil. Activez-le pour {device} ; vos autres appareils gardent leur propre choix.',
		'Se guarda por dispositivo. Actívalo para {device}; tus otros dispositivos conservan su propia elección.',
		'Pro Gerät gespeichert. Aktiviere es für {device}; deine anderen Geräte behalten ihre eigene Wahl.'
	],
	'notifySettings.devicePhone': ['this phone', 'ce téléphone', 'este teléfono', 'dieses Telefon'],
	'notifySettings.deviceComputer': ['this computer', 'cet ordinateur', 'este ordenador', 'diesen Computer'],
	'notifySettings.deviceBrowser': ['this browser', 'ce navigateur', 'este navegador', 'diesen Browser'],
	'notifySettings.shows': [
		'What a notification shows',
		'Ce qu’une notification affiche',
		'Qué muestra una notificación',
		'Was eine Benachrichtigung zeigt'
	],
	'notifySettings.showsBody': [
		'The sender’s name and the address it reached. Never the subject or the message.',
		'Le nom de l’expéditeur et l’adresse qui a reçu le message. Jamais l’objet ni le contenu.',
		'El nombre del remitente y la dirección que recibió el mensaje. Nunca el asunto ni el contenido.',
		'Der Name des Absenders und die Adresse, die die Nachricht erhalten hat. Nie der Betreff oder der Inhalt.'
	],

	'updates.title': ['Desktop updates', 'Mises à jour du bureau', 'Actualizaciones de escritorio', 'Desktop-Updates'],
	'updates.body': [
		'Keep Courrier desktop up to date with new features and fixes.',
		'Gardez Courrier pour le bureau à jour avec les nouveautés et les correctifs.',
		'Mantén Courrier de escritorio al día con nuevas funciones y correcciones.',
		'Halte Courrier für den Desktop mit neuen Funktionen und Korrekturen aktuell.'
	],
	'updates.version': ['Current version', 'Version actuelle', 'Versión actual', 'Aktuelle Version'],
	'updates.versionBody': [
		'The release build installed on this computer.',
		'La version publiée installée sur cet ordinateur.',
		'La versión publicada instalada en este ordenador.',
		'Die auf diesem Computer installierte Release-Version.'
	],
	'updates.status': ['Update status', 'État des mises à jour', 'Estado de las actualizaciones', 'Update-Status'],
	'updates.statusBody': [
		'The latest state from the Courrier release channel.',
		'Le dernier état du canal de publication de Courrier.',
		'El último estado del canal de versiones de Courrier.',
		'Der aktuelle Stand des Courrier-Release-Kanals.'
	],
	'updates.checking': ['Checking for updates…', 'Recherche de mises à jour…', 'Buscando actualizaciones…', 'Suche nach Updates…'],
	'updates.available': [
		'Version v{version} is available',
		'La version v{version} est disponible',
		'La versión v{version} está disponible',
		'Version v{version} ist verfügbar'
	],
	'updates.downloading': ['Downloading {progress}%', 'Téléchargement {progress} %', 'Descargando {progress} %', 'Download {progress} %'],
	'updates.ready': ['Ready to install', 'Prêt à installer', 'Listo para instalar', 'Bereit zur Installation'],
	'updates.failed': ['Update failed', 'Échec de la mise à jour', 'Error al actualizar', 'Update fehlgeschlagen'],
	'updates.upToDate': ['Courrier is up to date', 'Courrier est à jour', 'Courrier está al día', 'Courrier ist auf dem neuesten Stand'],
	'updates.idle': ['Ready to check', 'Prêt à vérifier', 'Listo para comprobar', 'Bereit zur Prüfung'],
	'updates.noneToInstall': [
		'No update available to install',
		'Aucune mise à jour à installer',
		'No hay ninguna actualización que instalar',
		'Kein Update zur Installation verfügbar'
	],
	'updates.progress': ['Download progress', 'Progression du téléchargement', 'Progreso de la descarga', 'Download-Fortschritt'],
	'updates.progressBody': [
		'Fetching update package from GitHub releases.',
		'Récupération du paquet de mise à jour depuis les versions GitHub.',
		'Obteniendo el paquete de actualización de las versiones de GitHub.',
		'Das Update-Paket wird von den GitHub-Releases geladen.'
	],
	'updates.actions': ['Actions', 'Actions', 'Acciones', 'Aktionen'],
	'updates.actionsBody': [
		'Manually check or install available updates.',
		'Vérifiez ou installez manuellement les mises à jour disponibles.',
		'Comprueba o instala manualmente las actualizaciones disponibles.',
		'Verfügbare Updates manuell prüfen oder installieren.'
	],
	'updates.install': ['Download & install', 'Télécharger et installer', 'Descargar e instalar', 'Herunterladen und installieren'],
	'updates.restart': ['Restart & install', 'Redémarrer et installer', 'Reiniciar e instalar', 'Neu starten und installieren'],
	'updates.check': ['Check for updates', 'Rechercher des mises à jour', 'Buscar actualizaciones', 'Nach Updates suchen'],
	'updates.checkingShort': ['Checking…', 'Recherche…', 'Buscando…', 'Suche…'],
	'updates.history': ['Release history', 'Historique des versions', 'Historial de versiones', 'Versionsverlauf'],
	'updates.historyBody': [
		'View changelog and previous release notes on GitHub.',
		'Consultez le journal des modifications et les notes des versions précédentes sur GitHub.',
		'Consulta el registro de cambios y las notas de versiones anteriores en GitHub.',
		'Änderungsprotokoll und frühere Versionshinweise auf GitHub ansehen.'
	],
	'updates.viewReleases': ['View releases', 'Voir les versions', 'Ver versiones', 'Versionen ansehen']
} satisfies Record<string, Row>;
