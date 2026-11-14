import type { Row } from '../row';

/** The strings drawn inside the landing page's animated demos. */
export const landingDemo = {
	'landing.demoLabel': [
		'A message arrives in the Courrier inbox, Bulle sums up the long thread in three lines, and a reply goes out.',
		'Un message arrive dans la boîte Courrier, Bulle résume le long fil en trois lignes, puis une réponse part.',
		'Llega un mensaje a la bandeja de Courrier, Bulle resume el hilo largo en tres líneas y sale una respuesta.',
		'Eine Nachricht kommt im Courrier-Posteingang an, Bulle fasst den langen Verlauf in drei Zeilen zusammen, dann geht eine Antwort raus.'
	],
	'landing.demoInbox': ['Inbox', 'Réception', 'Recibidos', 'Posteingang'],
	'landing.demoPrimary': ['Primary', 'Principal', 'Principal', 'Allgemein'],
	'landing.demoNewsletters': ['Newsletters', 'Newsletters', 'Boletines', 'Newsletter'],
	'landing.demoNotifications': ['Notifications', 'Notifications', 'Avisos', 'Hinweise'],
	'landing.demoSubject': [
		'Feedback on the mockup',
		'Retours sur la maquette',
		'Comentarios sobre la maqueta',
		'Feedback zum Entwurf'
	],
	'landing.demoEarlier': ['4 earlier messages', '4 messages plus tôt', '4 mensajes anteriores', '4 frühere Nachrichten'],
	'landing.demoBody': [
		'The new layout works. Can we keep the darker header?',
		'La nouvelle mise en page fonctionne. On garde l’en-tête plus sombre ?',
		'El nuevo diseño funciona. ¿Mantenemos la cabecera más oscura?',
		'Das neue Layout passt. Behalten wir den dunkleren Kopfbereich?'
	],
	'landing.demoSummarize': ['Summarize', 'Résumer', 'Resumir', 'Zusammenfassen'],
	'landing.demoSummary1': [
		'The client approves the new layout.',
		'Le client valide la nouvelle mise en page.',
		'El cliente aprueba el nuevo diseño.',
		'Der Kunde nimmt das neue Layout an.'
	],
	'landing.demoSummary2': [
		'One question left: the darker header.',
		'Reste une question : l’en-tête sombre.',
		'Queda una duda: la cabecera oscura.',
		'Offen ist nur der dunkle Kopfbereich.'
	],
	'landing.demoSummary3': [
		'They expect a new version this week.',
		'Une nouvelle version est attendue cette semaine.',
		'Esperan una nueva versión esta semana.',
		'Eine neue Version wird diese Woche erwartet.'
	],
	'landing.demoReply': [
		'Yes, the darker one. New version tomorrow.',
		'Oui, la plus sombre. Nouvelle version demain.',
		'Sí, la oscura. Nueva versión mañana.',
		'Ja, den dunkleren. Neue Version morgen.'
	],
	'landing.demoYou': ['You', 'Vous', 'Tú', 'Sie'],
	'landing.demoSend': ['Send', 'Envoyer', 'Enviar', 'Senden'],
	'landing.demoSent': ['Sent', 'Envoyé', 'Enviado', 'Gesendet'],
	'landing.demoTeam': ['Atelier team', 'Équipe Atelier', 'Equipo Atelier', 'Team Atelier'],
	'landing.demoQuote': ['Quote request', 'Demande de devis', 'Solicitud de presupuesto', 'Angebotsanfrage'],
	'landing.demoQuoteReply': [
		'Thanks! Quote coming this afternoon.',
		'Merci ! Notre devis part cet après-midi.',
		'¡Gracias! El presupuesto sale hoy.',
		'Danke! Das Angebot kommt heute.'
	]
} satisfies Record<string, Row>;
