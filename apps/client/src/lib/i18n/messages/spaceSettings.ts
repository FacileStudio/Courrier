import type { Row } from '../row';

/** A space's settings and members pages. */
export const spaceSettings = {
	'spaceSettings.title': [
		'Settings — {name} — Courrier',
		'Paramètres — {name} — Courrier',
		'Ajustes — {name} — Courrier',
		'Einstellungen — {name} — Courrier'
	],
	'spaceSettings.heading': ['Space settings', 'Paramètres de l’espace', 'Ajustes del espacio', 'Bereichseinstellungen'],
	'spaceSettings.space': ['Space', 'Espace', 'Espacio', 'Bereich'],
	'spaceSettings.managedTitle': [
		'Managed in Registre',
		'Géré dans Registre',
		'Gestionado en Registre',
		'In Registre verwaltet'
	],
	'spaceSettings.managedBody': [
		'The name, the members and their roles are edited in Registre, for every Facile app at once.',
		'Le nom, les membres et leurs rôles se modifient dans Registre, pour toutes les apps Facile à la fois.',
		'El nombre, los miembros y sus roles se editan en Registre, para todas las apps de Facile a la vez.',
		'Name, Mitglieder und ihre Rollen werden in Registre bearbeitet, für alle Facile-Apps gleichzeitig.'
	],
	'spaceSettings.general': ['General', 'Général', 'General', 'Allgemein'],
	'spaceSettings.generalBody': [
		'What members see of this space.',
		'Ce que les membres voient de cet espace.',
		'Lo que los miembros ven de este espacio.',
		'Was Mitglieder von diesem Bereich sehen.'
	],
	'spaceSettings.save': ['Save', 'Enregistrer', 'Guardar', 'Speichern'],
	'spaceSettings.saving': ['Saving…', 'Enregistrement…', 'Guardando…', 'Wird gespeichert…'],
	'spaceSettings.saved': ['Settings saved', 'Paramètres enregistrés', 'Ajustes guardados', 'Einstellungen gespeichert'],
	'spaceSettings.saveFailed': [
		'Could not save',
		'Impossible d’enregistrer',
		'No se pudo guardar',
		'Speichern nicht möglich'
	],
	'spaceSettings.danger': ['Danger zone', 'Zone de danger', 'Zona de peligro', 'Gefahrenzone'],
	'spaceSettings.dangerBody': [
		'Deleting this space is immediate and final.',
		'Supprimer cet espace est immédiat et définitif.',
		'Eliminar este espacio es inmediato y definitivo.',
		'Das Löschen dieses Bereichs erfolgt sofort und ist endgültig.'
	],
	'spaceSettings.dangerDetail': [
		'Every member loses access, and the mail accounts attached to this space drop off the list: they are neither deleted nor made personal, they are simply attached to nothing anymore. The IMAP/SMTP credentials stay in the database.',
		'Tous les membres perdent leur accès, et les comptes mail rattachés à cet espace disparaissent de la liste : ils ne sont ni supprimés ni rendus personnels, ils ne sont simplement plus rattachés à rien. Les identifiants IMAP/SMTP restent en base.',
		'Todos los miembros pierden el acceso y las cuentas de correo asociadas a este espacio desaparecen de la lista: no se eliminan ni pasan a ser personales, simplemente dejan de estar asociadas a nada. Las credenciales IMAP/SMTP permanecen en la base de datos.',
		'Alle Mitglieder verlieren den Zugriff, und die diesem Bereich zugeordneten E-Mail-Konten verschwinden aus der Liste: Sie werden weder gelöscht noch persönlich, sie sind einfach keinem Bereich mehr zugeordnet. Die IMAP/SMTP-Zugangsdaten bleiben in der Datenbank.'
	],
	'spaceSettings.delete': ['Delete the space', 'Supprimer l’espace', 'Eliminar el espacio', 'Bereich löschen'],
	'spaceSettings.deleteTitle': [
		'Delete “{name}”?',
		'Supprimer « {name} » ?',
		'¿Eliminar «{name}»?',
		'„{name}“ löschen?'
	],
	'spaceSettings.deleteThis': ['this space', 'cet espace', 'este espacio', 'diesen Bereich'],
	'spaceSettings.deleteBody': [
		'Every member loses access, and the mail accounts attached to this space leave the list. They do not reappear in your personal accounts. This cannot be undone.',
		'Tous les membres perdent leur accès, et les comptes mail rattachés à cet espace quittent la liste : ils ne réapparaissent pas dans vos comptes personnels. Cette action est définitive.',
		'Todos los miembros pierden el acceso y las cuentas de correo asociadas a este espacio salen de la lista: no reaparecen en tus cuentas personales. Esta acción es definitiva.',
		'Alle Mitglieder verlieren den Zugriff, und die diesem Bereich zugeordneten E-Mail-Konten verschwinden aus der Liste. Sie erscheinen nicht in deinen persönlichen Konten wieder. Das lässt sich nicht rückgängig machen.'
	],
	'spaceSettings.deleteFailed': [
		'Could not delete the space',
		'Impossible de supprimer l’espace',
		'No se pudo eliminar el espacio',
		'Der Bereich konnte nicht gelöscht werden'
	],
	'spaceSettings.deleted': ['Space deleted', 'Espace supprimé', 'Espacio eliminado', 'Bereich gelöscht'],
	'members.title': [
		'Members — {name} — Courrier',
		'Membres — {name} — Courrier',
		'Miembros — {name} — Courrier',
		'Mitglieder — {name} — Courrier'
	],
	'members.heading': ['Members', 'Membres', 'Miembros', 'Mitglieder'],
	'members.addTitle': ['Add a member', 'Ajouter un membre', 'Añadir un miembro', 'Mitglied hinzufügen'],
	'members.addBody': [
		'A member sees the mail accounts attached to this space.',
		'Un membre voit les comptes mail rattachés à cet espace.',
		'Un miembro ve las cuentas de correo asociadas a este espacio.',
		'Ein Mitglied sieht die E-Mail-Konten, die diesem Bereich zugeordnet sind.'
	],
	'members.email': ['Email address', 'Adresse e-mail', 'Correo electrónico', 'E-Mail-Adresse'],
	'members.emailPlaceholder': ['name@example.com', 'nom@exemple.fr', 'nombre@ejemplo.es', 'name@beispiel.de'],
	'members.user': [
		'Or someone you share a space with',
		'Ou une personne d’un de vos espaces',
		'O alguien con quien compartes un espacio',
		'Oder jemand aus einem deiner Bereiche'
	],
	'members.pickUser': [
		'Select a person…',
		'Sélectionner une personne…',
		'Seleccionar una persona…',
		'Person auswählen…'
	],
	'members.pickUserError': [
		'Enter an email address or pick a person',
		'Saisissez une adresse e-mail ou choisissez une personne',
		'Introduce un correo o elige a una persona',
		'Gib eine E-Mail-Adresse ein oder wähle eine Person'
	],
	'members.noAccount': [
		'No account uses that email address yet.',
		'Aucun compte n’utilise encore cette adresse e-mail.',
		'Ninguna cuenta usa todavía esa dirección de correo.',
		'Noch kein Konto nutzt diese E-Mail-Adresse.'
	],
	'members.role': ['Role', 'Rôle', 'Rol', 'Rolle'],
	'members.roleOf': ['Role of {name}', 'Rôle de {name}', 'Rol de {name}', 'Rolle von {name}'],
	'members.add': ['Add', 'Ajouter', 'Añadir', 'Hinzufügen'],
	'members.adding': ['Adding…', 'Ajout…', 'Añadiendo…', 'Wird hinzugefügt…'],
	'members.addFailed': [
		'Could not add the member',
		'Impossible d’ajouter le membre',
		'No se pudo añadir al miembro',
		'Das Mitglied konnte nicht hinzugefügt werden'
	],
	'members.added': ['Member added', 'Membre ajouté', 'Miembro añadido', 'Mitglied hinzugefügt'],
	'members.count': [
		{ one: '{count} member', other: '{count} members' },
		{ one: '{count} membre', other: '{count} membres' },
		{ one: '{count} miembro', other: '{count} miembros' },
		{ one: '{count} Mitglied', other: '{count} Mitglieder' }
	],
	'members.roleUpdated': ['Role updated', 'Rôle mis à jour', 'Rol actualizado', 'Rolle aktualisiert'],
	'members.roleFailed': [
		'Could not update the role',
		'Impossible de mettre à jour le rôle',
		'No se pudo actualizar el rol',
		'Die Rolle konnte nicht aktualisiert werden'
	],
	'members.removed': ['Member removed', 'Membre retiré', 'Miembro quitado', 'Mitglied entfernt'],
	'members.removeFailed': [
		'Could not remove the member',
		'Impossible de retirer le membre',
		'No se pudo quitar al miembro',
		'Das Mitglied konnte nicht entfernt werden'
	],
	'members.nobodyBody': [
		'Add a teammate above and they will see this space immediately.',
		'Ajoutez un coéquipier au-dessus et il verra cet espace immédiatement.',
		'Añade a un compañero arriba y verá este espacio de inmediato.',
		'Füge oben eine Kollegin oder einen Kollegen hinzu, und sie oder er sieht diesen Bereich sofort.'
	],
	'members.colMember': ['Member', 'Membre', 'Miembro', 'Mitglied'],
	'members.colRole': ['Role', 'Rôle', 'Rol', 'Rolle'],
	'members.colActions': ['Actions', 'Actions', 'Acciones', 'Aktionen'],
	'members.remove': ['Remove', 'Retirer', 'Quitar', 'Entfernen'],
	'members.removeTitle': ['Remove {name}?', 'Retirer {name} ?', '¿Quitar a {name}?', '{name} entfernen?'],
	'members.removeThis': ['this member', 'ce membre', 'este miembro', 'dieses Mitglied'],
	'members.removeBody': [
		'This person loses access to the space and to the mail accounts attached to it. Their own Courrier account and personal mailboxes are not touched.',
		'Cette personne perd l’accès à l’espace et aux comptes mail qui y sont rattachés. Son propre compte Courrier et ses boîtes personnelles ne sont pas touchés.',
		'Esta persona pierde el acceso al espacio y a las cuentas de correo asociadas. Su propia cuenta de Courrier y sus buzones personales no se ven afectados.',
		'Diese Person verliert den Zugriff auf den Bereich und die zugeordneten E-Mail-Konten. Ihr eigenes Courrier-Konto und ihre persönlichen Postfächer bleiben unberührt.'
	],
	'members.removeConfirm': ['Remove the member', 'Retirer le membre', 'Quitar al miembro', 'Mitglied entfernen']
} satisfies Record<string, Row>;
