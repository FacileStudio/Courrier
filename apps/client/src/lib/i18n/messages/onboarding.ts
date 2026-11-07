import type { Row } from '../row';

/** The first-run intro, the setup steps and the notification control. */
export const onboarding = {
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
	'notify.sendTest': ['Send a test', 'Envoyer un test', 'Enviar una prueba', 'Test senden']
} satisfies Record<string, Row>;
