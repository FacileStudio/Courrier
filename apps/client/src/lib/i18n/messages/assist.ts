import type { Row } from '../row';

/** The Bulle helper and the editor toolbar. */
export const assist = {
	'assist.summarize': ['Summarize', 'Résumer', 'Resumir', 'Zusammenfassen'],
	'assist.askBulle': ['Ask Bulle', 'Demander à Bulle', 'Preguntar a Bulle', 'Bulle fragen'],
	'assist.summary': ['Thread summary', 'Résumé de la conversation', 'Resumen de la conversación', 'Zusammenfassung der Unterhaltung'],
	'assist.hide': ['Hide', 'Masquer', 'Ocultar', 'Ausblenden'],
	'assist.failed': [
		'Bulle could not answer',
		'Bulle n’a pas pu répondre',
		'Bulle no pudo responder',
		'Bulle konnte nicht antworten'
	],
	'assist.off': [
		'The Bulle helper is not set up on this server.',
		'L’assistant Bulle n’est pas configuré sur ce serveur.',
		'El asistente Bulle no está configurado en este servidor.',
		'Der Bulle-Assistent ist auf diesem Server nicht eingerichtet.'
	],
	'assist.notLinked': [
		'Sign in with your Facile account to use Bulle here.',
		'Connectez-vous avec votre compte Facile pour utiliser Bulle ici.',
		'Inicia sesión con tu cuenta de Facile para usar Bulle aquí.',
		'Melde dich mit deinem Facile-Konto an, um Bulle hier zu nutzen.'
	],
	'assist.noKey': [
		'Add a model key in Bulle to use the helper.',
		'Ajoutez une clé de modèle dans Bulle pour utiliser l’assistant.',
		'Añade una clave de modelo en Bulle para usar el asistente.',
		'Füge in Bulle einen Modellschlüssel hinzu, um den Assistenten zu nutzen.'
	],
	'assist.unsupported': [
		'Bulle’s model is not supported here yet.',
		'Le modèle de Bulle n’est pas encore pris en charge ici.',
		'El modelo de Bulle aún no es compatible aquí.',
		'Das Modell von Bulle wird hier noch nicht unterstützt.'
	],
	'assist.unreachable': [
		'Bulle is not answering right now.',
		'Bulle ne répond pas pour le moment.',
		'Bulle no responde en este momento.',
		'Bulle antwortet gerade nicht.'
	],
	'assist.checking': ['Checking Bulle…', 'Vérification de Bulle…', 'Comprobando Bulle…', 'Bulle wird geprüft…'],
	'assist.scopeSelection': ['the selection', 'la sélection', 'la selección', 'die Auswahl'],
	'assist.scopeAll': ['the whole message', 'tout le message', 'todo el mensaje', 'die ganze Nachricht'],
	'assist.askAbout': ['Ask Bulle about {scope}', 'Demander à Bulle à propos de {scope}', 'Preguntar a Bulle sobre {scope}', 'Bulle zu {scope} fragen'],
	'assist.worksOn': ['Works on {scope}', 'Agit sur {scope}', 'Actúa sobre {scope}', 'Wirkt auf {scope}'],
	'assist.draftReply': ['Draft a reply', 'Rédiger une réponse', 'Redactar una respuesta', 'Antwort entwerfen'],
	'assist.rewrite': ['Rewrite', 'Réécrire', 'Reescribir', 'Umschreiben'],
	'assist.shorten': ['Shorten', 'Raccourcir', 'Acortar', 'Kürzen'],
	'assist.expand': ['Make it longer', 'Allonger', 'Alargar', 'Verlängern'],
	'assist.fix': [
		'Fix spelling and grammar',
		'Corriger l’orthographe et la grammaire',
		'Corregir ortografía y gramática',
		'Rechtschreibung und Grammatik korrigieren'
	],
	'assist.formal': ['More formal', 'Plus formel', 'Más formal', 'Formeller'],
	'assist.friendly': ['Friendlier', 'Plus chaleureux', 'Más cercano', 'Freundlicher'],
	'assist.translate': ['Translate to…', 'Traduire en…', 'Traducir a…', 'Übersetzen in…'],
	'assist.langEnglish': ['English', 'Anglais', 'Inglés', 'Englisch'],
	'assist.langFrench': ['French', 'Français', 'Francés', 'Französisch'],
	'assist.langSpanish': ['Spanish', 'Espagnol', 'Español', 'Spanisch'],
	'assist.langGerman': ['German', 'Allemand', 'Alemán', 'Deutsch'],
	'assist.langItalian': ['Italian', 'Italien', 'Italiano', 'Italienisch'],
	'assist.instructions': [
		'Tell Bulle what to change',
		'Dites à Bulle quoi modifier',
		'Dile a Bulle qué cambiar',
		'Sag Bulle, was geändert werden soll'
	],
	'assist.suggestion': [
		'Bulle’s suggestion',
		'Suggestion de Bulle',
		'Sugerencia de Bulle',
		'Vorschlag von Bulle'
	],
	'assist.suggestionFor': [
		'Suggestion for {scope}',
		'Suggestion pour {scope}',
		'Sugerencia para {scope}',
		'Vorschlag für {scope}'
	],
	'assist.insert': ['Insert', 'Insérer', 'Insertar', 'Einfügen'],
	'assist.replace': ['Replace', 'Remplacer', 'Reemplazar', 'Ersetzen'],
	'assist.insertBelow': ['Insert below', 'Insérer en dessous', 'Insertar debajo', 'Darunter einfügen'],
	'assist.tryAgain': ['Try again', 'Réessayer', 'Reintentar', 'Erneut versuchen'],
	'assist.discard': ['Discard', 'Ignorer', 'Descartar', 'Verwerfen'],
	'editor.placeholder': [
		'Write your message…',
		'Écrivez votre message…',
		'Escribe tu mensaje…',
		'Schreibe deine Nachricht…'
	],
	'editor.linkPrompt': ['URL', 'URL', 'URL', 'URL'],
	'editor.bold': ['Bold', 'Gras', 'Negrita', 'Fett'],
	'editor.italic': ['Italic', 'Italique', 'Cursiva', 'Kursiv'],
	'editor.underline': ['Underline', 'Souligné', 'Subrayado', 'Unterstrichen'],
	'editor.strike': ['Strikethrough', 'Barré', 'Tachado', 'Durchgestrichen'],
	'editor.link': ['Link', 'Lien', 'Enlace', 'Link'],
	'editor.bulletList': ['Bullet list', 'Liste à puces', 'Lista con viñetas', 'Aufzählung'],
	'editor.orderedList': ['Numbered list', 'Liste numérotée', 'Lista numerada', 'Nummerierte Liste'],
	'editor.quote': ['Quote', 'Citation', 'Cita', 'Zitat'],
	'editor.code': ['Code', 'Code', 'Código', 'Code'],
	'address.remove': ['Remove {chip}', 'Retirer {chip}', 'Quitar {chip}', '{chip} entfernen']
} satisfies Record<string, Row>;
