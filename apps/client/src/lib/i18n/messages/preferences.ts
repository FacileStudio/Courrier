import type { Row } from '../row';

/** The settings shell, appearance, notifications and desktop updates. */
export const preferences = {
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
	'settings.preferences': ['Preferences', 'Préférences', 'Preferencias', 'Einstellungen'],
	'settings.accounts': ['Accounts', 'Comptes', 'Cuentas', 'Konten'],
	'settings.notifications': ['Notifications', 'Notifications', 'Notificaciones', 'Benachrichtigungen'],
	'settings.templates': ['Templates', 'Modèles', 'Plantillas', 'Vorlagen'],
	'settings.tokens': ['API tokens', 'Jetons d’API', 'Tokens de API', 'API-Tokens'],
	'settings.advanced': ['Advanced', 'Avancé', 'Avanzado', 'Erweitert'],
	'settings.updates': ['Updates', 'Mises à jour', 'Actualizaciones', 'Updates'],
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
