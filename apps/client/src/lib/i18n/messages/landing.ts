import type { Row } from '../row';

/** The public landing page and the downloads page. */
export const landing = {
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
	'downloads.openWeb': ['Open the web app', 'Ouvrir l’application web', 'Abrir la aplicación web', 'Web-App öffnen']
} satisfies Record<string, Row>;
