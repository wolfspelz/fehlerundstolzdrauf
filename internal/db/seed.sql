-- Seed data: Stories from the original index.html
INSERT INTO stories (year, title, text, status) VALUES
('2009', 'Die Bäckerei', 'Ich eröffnete eine Bäckerei. Ich konnte nicht backen. Das wusste ich. Ich hoffte, es wäre egal. Es war nicht egal. Nach sieben Monaten schloss ich. Ich schulde der Bank noch immer Geld. Aber ich bin Konditor geworden. Das habe ich von der Bäckerei gelernt: Man muss es selbst können.', 'approved'),
('2017', 'Vier Jahre Jura', 'Vier Jahre Jura. Ich bestand alle Prüfungen. Dann saß ich im ersten echten Büro und verstand: Das war nicht mein Leben. Ich verließ es. Heute mache ich Möbel aus altem Holz. Niemand fragt nach dem Abschluss.', 'approved'),
('2021', 'Achtzehn Absagen', 'Ich schrieb sechs Jahre an einem Roman. Achtzehn Verlage lehnten ab. Das Manuskript liegt in einer Schublade. Ich habe angefangen, ein zweites zu schreiben. Das ist besser. Vielleicht liegt das auch irgendwann in einer Schublade.', 'approved'),
('2012', 'Zu früh', 'Wir heirateten nach drei Monaten. Alle sagten, es sei zu früh. Wir sagten, alle anderen seien zu langsam. Nach zwei Jahren unterschrieben wir die Scheidung im selben Café, in dem wir uns kennengelernt hatten. Ich habe gelernt, dass Mut und Dummheit denselben Anfang haben. Nur ein anderes Ende.', 'approved'),
('2019', 'Niemand brauchte sie', 'Drei Freunde, ein Wohnzimmer, eine Idee. Wir bauten eine App, die niemand brauchte. Wir sagten uns, die Leute verstehen es noch nicht. Die Leute verstanden es. Sie brauchten es nicht. Das Geld war nach elf Monaten weg. Die Freundschaften auch. Einer von uns baut jetzt Software für Krankenhäuser. Etwas, das jemand braucht.', 'approved'),
('2015', 'Der Laden', 'Ich eröffnete einen Plattenladen. Im Jahr 2015. Alle sagten, das sei verrückt. Ich sagte, Vinyl sei zurück. Es stimmte. Aber nicht in meiner Straße. Nach vierzehn Monaten räumte ich die Regale leer. Die letzte Platte, die ich verkaufte, war "The End" von den Doors.', 'approved'),
('2008', 'Drei Semester Physik', 'Drei Semester Physik. Ich verstand alles in der Vorlesung. Dann kam die Klausur. Ich verstand nichts. Dreimal. Nach dem dritten Mal ging ich. Heute unterrichte ich Mathematik an einer Hauptschule. Die Schüler verstehen mich. Das ist mehr, als die Physik je von mir sagen konnte.', 'approved'),
('2020', 'Der Podcast', 'Wir starteten einen Podcast über Scheitern. Vierundzwanzig Folgen. Siebenunddreißig Hörer. Wir scheiterten am Scheitern. Die Ironie war uns bewusst. Wir machten trotzdem weiter. Bis wir aufhörten. Niemand bemerkte es.', 'approved'),
('2014', 'Das Restaurant', 'Ich kochte gut. Alle sagten das. Also eröffnete ich ein Restaurant. Kochen ist nicht Führen. Führen ist nicht Kochen. Ich konnte beides nicht gleichzeitig. Nach neun Monaten kochte ich wieder zu Hause. Für Freunde. Die sagen immer noch, ich koche gut.', 'approved'),
('2011', 'Der Umzug', 'Ich zog nach Barcelona. Ohne Spanisch. Ohne Job. Ohne Plan. Ich hatte Ersparnisse für drei Monate. Nach zwei Monaten waren sie weg. Ich kam zurück. Kleiner. Aber ich hatte das Meer gesehen. Jeden Morgen. Das war es wert.', 'approved'),
('2022', 'Die Galerie', 'Ich malte zehn Jahre. Dann mietete ich eine Galerie für eine Ausstellung. Drei Wochen. Zweitausend Euro. Vierzehn Besucher. Zwei davon meine Eltern. Ich male immer noch. Aber für mich. Das hätte ich von Anfang an tun sollen.', 'approved'),
('2016', 'Zwei Jahre China', 'Zwei Jahre Mandarin gelernt. Jeden Tag. Dann flog ich nach Peking. Niemand verstand mich. Mein Lehrer hatte einen Dialekt unterrichtet, den nur sein Dorf sprach. Ich lachte. Dann weinte ich. Dann bestellte ich auf Englisch.', 'approved'),
('2013', 'Das Erbe', 'Mein Onkel hinterließ mir fünfzigtausend Euro. Ich investierte alles in Bitcoin. Im Jahr 2013. Ich verkaufte bei zweihundert Euro. Panik. Heute wäre es ein Vermögen. Ich habe gelernt, dass Geduld kein Talent ist. Es ist eine Entscheidung. Ich habe falsch entschieden.', 'approved'),
('2018', 'Der Marathon', 'Ich trainierte acht Monate für den Berlin-Marathon. Am Tag des Rennens stand ich an der Startlinie. Kilometer dreiundzwanzig: mein Knie gab nach. Ich ging die letzten neunzehn Kilometer. Zu Fuß. Es dauerte vier Stunden. Aber ich kam an. Manchmal ist Ankommen genug.', 'approved'),
('2023', 'Die Bewerbung', 'Ich bewarb mich auf meinen Traumjob. Dreimal. Beim ersten Mal kam keine Antwort. Beim zweiten eine Absage. Beim dritten Mal ein Gespräch. Sie sagten, ich sei überqualifiziert. Ich arbeite jetzt woanders. Es ist nicht mein Traum. Aber es ist gut. Manchmal reicht gut.', 'approved');

-- Seed data: Featured story
INSERT INTO featured (year_range, title, intro, quote, outro) VALUES
('1998 – 2003', 'Ich habe eine Firma gegen die Wand gefahren. Achtzehn Menschen verloren ihren Job.', 'Es war meine erste Firma. Ich war dreiundzwanzig. Ich dachte, Enthusiasmus ersetzt Erfahrung. Das tut er nicht.', 'Ich habe jedem einzeln angerufen. Es hat zwei Tage gedauert. Das war das Richtige. Es hat nichts besser gemacht.', 'Heute führe ich eine andere Firma. Wir sind klein. Ich weiß, was jeder verdient. Ich weiß, was jeder braucht. Das habe ich damals nicht gewusst. Jetzt weiß ich es. Das war teuer.'),
('2005 – 2010', 'Fünf Jahre lang habe ich ein Buch geschrieben, das niemand lesen wollte.', 'Jeden Morgen um fünf aufgestanden. Vor der Arbeit geschrieben. Nach der Arbeit geschrieben. Am Wochenende geschrieben. Meine Frau sagte, ich sei besessen. Sie hatte recht.', 'Das Manuskript hatte achthundert Seiten. Der Lektor sagte: Die Geschichte beginnt auf Seite dreihundert.', 'Ich habe die ersten dreihundert Seiten gelöscht. Dann die nächsten fünfhundert. Übrig blieb eine Kurzgeschichte. Zwölf Seiten. Sie wurde veröffentlicht. In einer Zeitschrift, die niemand liest. Aber sie wurde veröffentlicht.'),
('2010 – 2015', 'Wir wollten die Bildung revolutionieren. Wir haben nur Geld verbrannt.', 'Vier Gründer, eine Vision. Jedes Kind verdient die beste Bildung. Kostenlos. Digital. Wir hatten alles: Investoren, ein Büro in Kreuzberg, zwanzig Mitarbeiter. Alles außer Schüler.', 'Unser Produkt war perfekt. Für niemanden. Wir hatten nie einen Lehrer gefragt, was er braucht.', 'Die Firma gibt es nicht mehr. Zwei der Gründer reden nicht mehr miteinander. Einer ist jetzt Lehrer. Er sagt, er lernt mehr von seinen Schülern als wir je von unseren Nutzern gelernt haben.');

-- Seed data: Quotes
INSERT INTO quotes (text, attribution) VALUES
('Ich habe nicht versagt. Ich habe 10.000 Wege gefunden, die nicht funktionieren.', 'Thomas Edison'),
('Erfolg ist die Fähigkeit, von einem Misserfolg zum nächsten zu gehen, ohne seine Begeisterung zu verlieren.', 'Winston Churchill'),
('Wer einen Fehler gemacht hat und ihn nicht korrigiert, begeht einen zweiten.', 'Konfuzius'),
('Unsere größte Schwäche liegt im Aufgeben. Der sicherste Weg zum Erfolg ist immer, es noch einmal zu versuchen.', 'Thomas Edison'),
('Es ist unmöglich, zu leben, ohne an etwas zu scheitern, es sei denn, man lebt so vorsichtig, dass man genauso gut gar nicht gelebt haben könnte.', 'J.K. Rowling'),
('Fehler sind die Portale der Entdeckung.', 'James Joyce'),
('Nur wer nichts tut, macht keine Fehler.', 'Theodor Fontane'),
('Ich bin dankbar für alle, die Nein gesagt haben. Ihretwegen habe ich es selbst gemacht.', 'Albert Einstein'),
('Das Scheitern ist der Grundstein des Erfolges.', 'Lao Tse'),
('Man muss noch Chaos in sich haben, um einen tanzenden Stern gebären zu können.', 'Friedrich Nietzsche'),
('Der größte Ruhm im Leben liegt nicht darin, nie zu fallen, sondern jedes Mal wieder aufzustehen.', 'Nelson Mandela'),
('Aus Fehlern lernt man. Aus großen Fehlern lernt man Großes.', 'Unbekannt'),
('Wer immer tut, was er schon kann, bleibt immer das, was er schon ist.', 'Henry Ford'),
('Ein Experte ist ein Mann, der alle Fehler gemacht hat, die man in einem sehr engen Gebiet machen kann.', 'Niels Bohr'),
('Der einzige wirkliche Fehler ist der, von dem wir nichts lernen.', 'Henry Ford'),
('Perfektion ist nicht erreichbar. Aber wenn wir Perfektion anstreben, können wir Exzellenz erreichen.', 'Vince Lombardi'),
('Scheitern ist einfach die Gelegenheit, wieder von vorn anzufangen, diesmal intelligenter.', 'Henry Ford'),
('Das Leben ist wie Fahrradfahren. Um die Balance zu halten, muss man in Bewegung bleiben.', 'Albert Einstein'),
('Nicht weil es schwer ist, wagen wir es nicht, sondern weil wir es nicht wagen, ist es schwer.', 'Seneca'),
('Jeder Misserfolg ist ein Schritt zum Erfolg.', 'William Whewell'),
('Die reinste Form des Wahnsinns ist es, alles beim Alten zu lassen und gleichzeitig zu hoffen, dass sich etwas ändert.', 'Albert Einstein'),
('Wer aufhört, Fehler zu machen, lernt nichts mehr dazu.', 'Theodor Fontane'),
('Hindernisse sind jene furchtbaren Dinge, die man sieht, wenn man die Augen von seinem Ziel abwendet.', 'Henry Ford'),
('Ein Schiff im Hafen ist sicher, aber dafür werden Schiffe nicht gebaut.', 'John A. Shedd'),
('Der Weg zum Erfolg und der Weg zum Misserfolg sind fast genau derselbe.', 'Colin R. Davis'),
('Mut steht am Anfang des Handelns, Glück am Ende.', 'Demokrit'),
('Wer kämpft, kann verlieren. Wer nicht kämpft, hat schon verloren.', 'Bertolt Brecht'),
('Es ist nicht genug, zu wissen, man muss auch anwenden. Es ist nicht genug, zu wollen, man muss auch tun.', 'Johann Wolfgang von Goethe');

-- Seed data: Historical
INSERT INTO historical (year, title, text) VALUES
('1968', 'Die Klebstoff-Katastrophe', 'Spencer Silver wollte bei 3M einen superstarken Klebstoff entwickeln. Heraus kam das Gegenteil: ein Kleber, der kaum haftete. Sechs Jahre lang interessierte sich niemand dafür. Dann brauchte sein Kollege Art Fry ein Lesezeichen, das nicht rausfiel. Das Post-it war geboren.'),
('1492', 'Der falsche Kontinent', 'Kolumbus rechnete den Erdumfang um ein Drittel zu klein. Deshalb glaubte er, Indien sei erreichbar. Jeder Kartograph in Europa wusste, dass er falsch lag. Niemand konnte ihn überzeugen. Er starb in dem Glauben, in Asien gewesen zu sein. Sein Fehler veränderte die Weltkarte.'),
('1928', 'Die vergessene Petrischale', 'Alexander Fleming ging in den Urlaub und vergaß, seine Petrischalen zu reinigen. Als er zurückkam, war eine davon mit Schimmel bedeckt. Um den Schimmel herum waren alle Bakterien tot. Er hatte gerade Penicillin entdeckt. Durch Schlamperei.'),
('1895', 'Die unsichtbaren Strahlen', 'Wilhelm Conrad Röntgen experimentierte mit Kathodenstrahlen. Ein Bildschirm in der Nähe leuchtete auf, obwohl er es nicht sollte. Statt den Fehler zu ignorieren, untersuchte er ihn. Sechs Wochen später hatte er die Röntgenstrahlen entdeckt und die Medizin für immer verändert.'),
('1853', 'Die zu dünnen Kartoffeln', 'Ein Gast im Moon Lake Lodge beschwerte sich, die Bratkartoffeln seien zu dick. Koch George Crum schnitt sie aus Trotz hauchdünn und frittierte sie knusprig. Der Gast war begeistert. Die Kartoffelchips waren erfunden. Aus Ärger.'),
('1945', 'Die geschmolzene Schokolade', 'Percy Spencer stand vor einem Magnetron und bemerkte, dass der Schokoriegel in seiner Tasche geschmolzen war. Statt sich zu ärgern, stellte er Popcorn-Körner vor das Gerät. Sie platzten. Am nächsten Tag versuchte er es mit einem Ei. Es explodierte. Die Mikrowelle war geboren.'),
('1907', 'Der falsche Kunststoff', 'Leo Baekeland wollte einen Ersatz für Schellack herstellen. Das Experiment ging schief. Heraus kam ein Material, das sich nicht mehr verformen ließ, sobald es einmal erhitzt wurde. Er nannte es Bakelit. Es wurde der erste vollsynthetische Kunststoff und revolutionierte die Industrie.'),
('1826', 'Acht Stunden Belichtung', 'Joseph Nicéphore Niépce versuchte jahrelang, Bilder festzuhalten. Sein erstes brauchbares Foto benötigte acht Stunden Belichtungszeit. Das Ergebnis war kaum erkennbar. Aber es war ein Foto. Das erste der Welt. Alles, was danach kam – von Portraits bis Instagram – begann mit diesem unscharfen Bild.'),
('1974', 'Der gescheiterte Klebstoff wird Kunst', 'Art Fry sang im Kirchenchor. Seine Lesezeichen fielen ständig aus dem Gesangbuch. Er erinnerte sich an Spencer Silvers nutzlosen Klebstoff von 3M. Er bestrich Papierstreifen damit. Sie hafteten und ließen sich ablösen. Die Firma sagte: Das braucht niemand. Zwölf Jahre später machte 3M mit Post-its Milliarden.'),
('1956', 'Das falsche Medikament', 'Forscher testeten Iproniazid als Tuberkulose-Medikament. Es half kaum gegen TB. Aber die Patienten wurden auffallend fröhlich. Manche tanzten durch die Station. Die Forscher hatten versehentlich das erste Antidepressivum entdeckt. Ein gescheitertes Medikament heilte eine andere Krankheit.'),
('1839', 'Der vergessene Gummi', 'Charles Goodyear versuchte jahrelang, Gummi haltbar zu machen. Alles scheiterte. Dann ließ er versehentlich eine Mischung aus Gummi und Schwefel auf den heißen Ofen fallen. Statt zu schmelzen, wurde sie elastisch und fest. Vulkanisation. Zufall und Verzweiflung, in der richtigen Reihenfolge.'),
('1980', 'Die Teppich-Maschine', 'James Dyson hasste seinen Staubsauger. Der Beutel verstopfte ständig. Er baute einen Prototyp ohne Beutel. Dann noch einen. Und noch einen. 5.127 Prototypen und fünfzehn Jahre später funktionierte er. Kein Hersteller wollte ihn lizenzieren. Also verkaufte Dyson ihn selbst. Heute ist er Milliardär.');
