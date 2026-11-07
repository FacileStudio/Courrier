import type { Row } from '../row';

/** The profile page. */
export const profile = {
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
	'profile.signingOut': ['Signing out…', 'Déconnexion…', 'Cerrando sesión…', 'Abmeldung…']
} satisfies Record<string, Row>;
