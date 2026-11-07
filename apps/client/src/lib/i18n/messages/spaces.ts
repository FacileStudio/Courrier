import type { Row } from '../row';

/** The spaces list and a space's page. */
export const spaces = {
	'spaces.title': ['Spaces — Courrier', 'Espaces — Courrier', 'Espacios — Courrier', 'Bereiche — Courrier'],
	'spaces.heading': ['Spaces', 'Espaces', 'Espacios', 'Bereiche'],
	'spaces.intro': [
		'Work with your team in shared spaces: each space has its own members and mail accounts.',
		'Collaborez avec votre équipe dans des espaces partagés : chaque espace a ses membres et ses comptes mail.',
		'Colabora con tu equipo en espacios compartidos: cada espacio tiene sus miembros y sus cuentas de correo.',
		'Arbeite mit deinem Team in gemeinsamen Bereichen: Jeder Bereich hat eigene Mitglieder und E-Mail-Konten.'
	],
	'spaces.settings': ['Settings', 'Réglages', 'Ajustes', 'Einstellungen'],
	'spaces.manageInRegistre': ['Manage in Registre', 'Gérer dans Registre', 'Gestionar en Registre', 'In Registre verwalten'],
	'spaces.new': ['New space', 'Nouvel espace', 'Nuevo espacio', 'Neuer Bereich'],
	'spaces.none': ['No spaces', 'Aucun espace', 'Ningún espacio', 'Keine Bereiche'],
	'spaces.noneRegistre': [
		'Spaces are created in Registre and appear here within seconds.',
		'Les espaces se créent dans Registre et apparaissent ici en quelques secondes.',
		'Los espacios se crean en Registre y aparecen aquí en unos segundos.',
		'Bereiche werden in Registre erstellt und erscheinen hier innerhalb weniger Sekunden.'
	],
	'spaces.noneLocal': [
		'Create a space to share mail accounts and folders with your team.',
		'Créez un espace pour partager des comptes mail et des dossiers avec votre équipe.',
		'Crea un espacio para compartir cuentas de correo y carpetas con tu equipo.',
		'Erstelle einen Bereich, um E-Mail-Konten und Ordner mit deinem Team zu teilen.'
	],
	'spaces.createFirst': [
		'Create your first space',
		'Créer votre premier espace',
		'Crear tu primer espacio',
		'Ersten Bereich erstellen'
	],
	'spaces.noDescription': ['No description.', 'Pas de description.', 'Sin descripción.', 'Keine Beschreibung.'],
	'spaces.createTitle': ['New space', 'Nouvel espace', 'Nuevo espacio', 'Neuer Bereich'],
	'spaces.nameLabel': ['Name', 'Nom', 'Nombre', 'Name'],
	'spaces.namePlaceholder': ['e.g. Marketing team', 'ex. Équipe Marketing', 'p. ej. Equipo de marketing', 'z. B. Marketing-Team'],
	'spaces.descriptionLabel': ['Description', 'Description', 'Descripción', 'Beschreibung'],
	'spaces.descriptionHelp': [
		'Optional. What is this space for?',
		'Optionnel : à quoi sert cet espace ?',
		'Opcional. ¿Para qué sirve este espacio?',
		'Optional. Wofür ist dieser Bereich?'
	],
	'spaces.descriptionPlaceholder': [
		'The team’s shared mailboxes.',
		'Boîtes partagées de l’équipe.',
		'Buzones compartidos del equipo.',
		'Die gemeinsamen Postfächer des Teams.'
	],
	'spaces.nameRequired': ['The name is required', 'Le nom est requis', 'El nombre es obligatorio', 'Der Name ist erforderlich'],
	'spaces.createFailed': [
		'Could not create the space',
		'Impossible de créer l’espace',
		'No se pudo crear el espacio',
		'Der Bereich konnte nicht erstellt werden'
	],
	'spaces.created': [
		'Space “{name}” created',
		'Espace « {name} » créé',
		'Espacio «{name}» creado',
		'Bereich „{name}“ erstellt'
	],
	'spaces.cancel': ['Cancel', 'Annuler', 'Cancelar', 'Abbrechen'],
	'spaces.createSubmit': ['Create the space', 'Créer l’espace', 'Crear el espacio', 'Bereich erstellen'],
	'spaces.creating': ['Creating…', 'Création…', 'Creando…', 'Wird erstellt…'],
	'spaces.roleOwner': ['Owner', 'Propriétaire', 'Propietario', 'Eigentümer'],
	'spaces.roleAdmin': ['Admin', 'Admin', 'Admin', 'Admin'],
	'spaces.roleMember': ['Member', 'Membre', 'Miembro', 'Mitglied'],
	'space.fallbackName': ['Space', 'Espace', 'Espacio', 'Bereich'],
	'space.notFound': ['Space not found', 'Espace introuvable', 'Espacio no encontrado', 'Bereich nicht gefunden'],
	'space.notFoundBody': [
		'It may have been deleted, or your access was removed.',
		'Il a peut-être été supprimé, ou votre accès a été retiré.',
		'Puede que se haya eliminado o que te hayan quitado el acceso.',
		'Er wurde möglicherweise gelöscht oder dein Zugriff wurde entfernt.'
	],
	'space.backToSpaces': ['Back to spaces', 'Retour aux espaces', 'Volver a los espacios', 'Zurück zu den Bereichen'],
	'space.spaces': ['Spaces', 'Espaces', 'Espacios', 'Bereiche'],
	'space.switch': ['Switch to this space', 'Basculer vers cet espace', 'Cambiar a este espacio', 'Zu diesem Bereich wechseln'],
	'space.active': ['Active space', 'Espace actif', 'Espacio activo', 'Aktiver Bereich'],
	'space.activeToast': ['Active space: {name}', 'Espace actif : {name}', 'Espacio activo: {name}', 'Aktiver Bereich: {name}'],
	'space.settings': ['Settings', 'Paramètres', 'Ajustes', 'Einstellungen'],
	'space.members': ['Members', 'Membres', 'Miembros', 'Mitglieder'],
	'space.createdOn': ['Created on', 'Créé le', 'Creado el', 'Erstellt am'],
	'space.yourRole': ['Your role: {role}', 'Votre rôle : {role}', 'Tu rol: {role}', 'Deine Rolle: {role}'],
	'space.membersIntro': [
		'Everyone here sees the mail accounts attached to this space.',
		'Tout le monde ici voit les comptes mail rattachés à cet espace.',
		'Todas las personas de aquí ven las cuentas de correo asociadas a este espacio.',
		'Alle hier sehen die E-Mail-Konten, die diesem Bereich zugeordnet sind.'
	],
	'space.manageMembers': ['Manage members', 'Gérer les membres', 'Gestionar miembros', 'Mitglieder verwalten'],
	'space.nobody': [
		'Nobody here yet',
		'Personne ici pour l’instant',
		'Todavía no hay nadie aquí',
		'Noch niemand hier'
	],
	'space.nobodyRegistre': [
		'Add a teammate in Registre and they will see this space within seconds.',
		'Ajoutez un coéquipier dans Registre et il verra cet espace en quelques secondes.',
		'Añade a un compañero en Registre y verá este espacio en unos segundos.',
		'Füge in Registre eine Kollegin oder einen Kollegen hinzu, und sie oder er sieht diesen Bereich innerhalb weniger Sekunden.'
	],
	'space.nobodyLocal': [
		'Add a teammate and they will see this space immediately.',
		'Ajoutez un coéquipier et il verra cet espace immédiatement.',
		'Añade a un compañero y verá este espacio de inmediato.',
		'Füge eine Kollegin oder einen Kollegen hinzu, und sie oder er sieht diesen Bereich sofort.'
	]
} satisfies Record<string, Row>;
