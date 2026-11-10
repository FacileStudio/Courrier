import type { Row } from '../row';

/** The Ideas page: suggestions Bulle makes about the mail, run in place. */
export const ideas = {
	'ideas.title': ['Ideas — Courrier', 'Idées — Courrier', 'Ideas — Courrier', 'Ideen — Courrier'],
	'ideas.description': [
		'Bulle reads your inbox and suggests which replies to write, what to follow up and what to sort. Run an idea and it is prepared here.',
		'Bulle lit votre boîte de réception et propose les réponses à écrire, les relances à faire et le tri à mener. Lancez une idée, elle se prépare ici.',
		'Bulle lee tu bandeja de entrada y propone qué respuestas escribir, qué seguir y qué ordenar. Lanza una idea y se prepara aquí.',
		'Bulle liest Ihren Posteingang und schlägt vor, welche Antworten zu schreiben, was nachzufassen und was zu sortieren ist. Starten Sie eine Idee, sie wird hier vorbereitet.'
	],
	'ideas.approval.summary': [
		'Bulle wants to run "{tool}". Check it before it goes ahead.',
		'Bulle veut lancer « {tool} ». Vérifiez avant de continuer.',
		'Bulle quiere ejecutar «{tool}». Revísalo antes de seguir.',
		'Bulle möchte „{tool}“ ausführen. Prüfen Sie es, bevor es weitergeht.'
	],
	'ideas.tool.courrier_list_accounts': ['List mail accounts', 'Lister les comptes mail', 'Listar cuentas de correo', 'Postfächer auflisten'],
	'ideas.tool.courrier_search_mail': ['Search mail', 'Chercher des mails', 'Buscar correos', 'E-Mails suchen'],
	'ideas.tool.courrier_read_thread': ['Read a conversation', 'Lire une conversation', 'Leer una conversación', 'Unterhaltung lesen'],
	'ideas.tool.courrier_draft': ['Write a draft', 'Écrire un brouillon', 'Escribir un borrador', 'Entwurf schreiben'],
	'ideas.tool.courrier_send': ['Send mail', 'Envoyer un mail', 'Enviar correo', 'E-Mail senden'],
	'ideas.tool.courrier_organize': ['Organize mail', 'Ranger des mails', 'Organizar correos', 'E-Mails sortieren'],
	'ideas.tool.courrier_delete_draft': ['Delete a draft', 'Supprimer un brouillon', 'Eliminar un borrador', 'Entwurf löschen'],
	'ideas.error.load': [
		'Could not load the ideas.',
		'Impossible de charger les idées.',
		'No se pudieron cargar las ideas.',
		'Die Ideen konnten nicht geladen werden.'
	],
	'ideas.error.unreachable': [
		'Bulle is unreachable for now. Try again in a moment.',
		'Bulle est injoignable pour le moment. Réessayez dans un instant.',
		'Bulle no responde por ahora. Inténtalo en un momento.',
		'Bulle ist gerade nicht erreichbar. Versuchen Sie es gleich noch einmal.'
	],
	'ideas.error.notLinked': [
		'Reconnect your Facile account to use Bulle, or pick a space that Registre manages.',
		'Reconnectez votre compte Facile pour utiliser Bulle, ou choisissez un espace géré par Registre.',
		'Vuelve a conectar tu cuenta Facile para usar Bulle, o elige un espacio gestionado por Registre.',
		'Verbinden Sie Ihr Facile-Konto erneut oder wählen Sie einen von Registre verwalteten Bereich.'
	],
	'ideas.error.capped': [
		'You asked for ideas three times today. More come tomorrow.',
		'Vous avez demandé des idées trois fois aujourd’hui. D’autres arrivent demain.',
		'Ya pediste ideas tres veces hoy. Mañana llegan más.',
		'Sie haben heute dreimal nach Ideen gefragt. Morgen kommen neue.'
	],
	'ideas.error.busy': [
		'Bulle is already thinking of ideas for this space.',
		'Bulle cherche déjà des idées pour cet espace.',
		'Bulle ya está buscando ideas para este espacio.',
		'Bulle sucht bereits Ideen für diesen Bereich.'
	],
	'ideas.error.gone': [
		'This idea is no longer available.',
		'Cette idée n’est plus disponible.',
		'Esta idea ya no está disponible.',
		'Diese Idee ist nicht mehr verfügbar.'
	],
	'ideas.error.run': [
		'The run stopped before the end.',
		'L’exécution s’est arrêtée avant la fin.',
		'La ejecución se detuvo antes del final.',
		'Der Lauf wurde vor dem Ende abgebrochen.'
	],
	'ideas.example.reply.title': [
		'Answer Marc about the quote',
		'Répondre à Marc au sujet du devis',
		'Responder a Marc sobre el presupuesto',
		'Marc zum Angebot antworten'
	],
	'ideas.example.reply.why': [
		'He wrote on Monday and is still waiting. A short reply is drafted and ready to review.',
		'Il a écrit lundi et attend toujours. Une courte réponse est rédigée, prête à relire.',
		'Escribió el lunes y sigue esperando. Hay una respuesta breve redactada, lista para revisar.',
		'Er hat am Montag geschrieben und wartet noch. Eine kurze Antwort ist entworfen und bereit zur Durchsicht.'
	],
	'ideas.example.followup.title': [
		'Follow up on the unpaid invoice',
		'Relancer la facture impayée',
		'Reclamar la factura impagada',
		'Unbezahlte Rechnung nachfassen'
	],
	'ideas.example.followup.why': [
		'Invoice 0412 went out 30 days ago with no answer. A polite reminder is ready.',
		'La facture 0412 est partie il y a 30 jours sans réponse. Une relance polie est prête.',
		'La factura 0412 salió hace 30 días sin respuesta. Hay un recordatorio cortés listo.',
		'Rechnung 0412 ging vor 30 Tagen raus, ohne Antwort. Eine höfliche Erinnerung liegt bereit.'
	],
	'ideas.example.sort.title': [
		'Sort the newsletters',
		'Trier les newsletters',
		'Ordenar los boletines',
		'Newsletter sortieren'
	],
	'ideas.example.sort.why': [
		'Twenty-three newsletters crowd the inbox this week. They could move to their own tab.',
		'Vingt-trois newsletters encombrent la boîte cette semaine. Elles pourraient passer dans leur onglet.',
		'Veintitrés boletines saturan la bandeja esta semana. Podrían pasar a su propia pestaña.',
		'Dreiundzwanzig Newsletter füllen diese Woche den Posteingang. Sie könnten in ihren eigenen Tab wandern.'
	],
	'ideas.dismissed': ['Idea dismissed', 'Idée écartée', 'Idea descartada', 'Idee verworfen']
} as const satisfies Record<string, Row>;
