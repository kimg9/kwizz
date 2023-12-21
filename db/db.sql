---------------------------------------------------------
---------------------USER setup--------------------------

create table users (
    user_id                int                  primary key generated always as identity,
    pseudonym              varchar(50)          unique not null,
    profil_pic             bytea                ,
    gallery_pic            bytea                
);

---------------------------------------------------------
-------------------QUIZZ setup---------------------------

create table user_score (
    id                      int             primary key generated always as identity,
    user_id                 int             references users(user_id),
    total_score             bigint          default '0'
);

create table categories (
    cat_id                  int                 primary key generated always as identity,
    cat_name                varchar(124)        unique not null,
    cat_short_name          varchar(124)        unique not null,
    cat_description         text                not null,
    cat_image               text                not null           
);

-- create table user_score_per_cat (
--     id      int             primary key generated always as identity,
--     cat_id  int             references categories(cat_id),
--     user_id int             references users(user_id),
--     score   int             default '0'
-- );

create table quizzes (
    quizz_id                int                 primary key generated always as identity,
    quizz_title             varchar(512)        not null,
    quizz_description       text                not null,
    created_at              timestamptz         not null default current_timestamp,
    cat_id                  INT                 references categories(cat_id)
);

create table questions (
    question_id             int                 primary key generated always as identity,
    quizz_id                int                 references quizzes(quizz_id) not null,
    question                text                not null,
    order_questions         smallint
);

create table quizz_sessions (
    session_id             int                  primary key generated always as identity,
    quizz_id               int                  references quizzes(quizz_id) not null,
    user_id                int                  references users(user_id) not null,
    created_at             timestamptz          not null default current_timestamp,
    finished               bool                 not null default false,
    score                  int                  not null default '0'
);

create table responses (
    response_id            int                  primary key generated always as identity,
    question_id            int                  references questions(question_id),
    answer                 text                 not null,
    isCorrect              bool                 not null default false
);

create table sess_resp (
    ID                     int                  primary key generated always as identity,
    session_id             int                  not null references quizz_sessions(session_id),
    response_id            int                  not null references responses(response_id)
);

---------------------------------------------------------
-----------------------TEST setup------------------------

INSERT INTO users(pseudonym)
VALUES ('Alain');

-- INSERT INTO users(pseudonym)
-- VALUES ('Nicolas');

INSERT INTO categories(cat_name, cat_short_name, cat_description, cat_image)
VALUES ('Culture Générale', 'culturegenerale', 'Entraine-toi sur des questions sur des sujets variés allant de l''histoire, à la politique, la géographie, etc.', '/public/cat_pic/culturegenerale.webp'),
('Mathématiques', 'mathematiques', 'Affronte les joies des petits problèmes de mathématiques.', '/public/cat_pic/mathematiques.webp'),
('Français', 'francais', 'Teste ta connaissance de la langue de Molière et ta mémoire sur ses auteurs et leurs oeuvres.', '/public/cat_pic/francais.webp'),
-- ('Mémoire', 'memoire', 'Des petits jeux de mémoire pour la remuscler !', '/public/cat_pic/memoire.webp'),
('Famille', 'famille', 'Pour t''aider à te souvenir de l''anniversaire de tes enfants !', '/public/cat_pic/famille.webp');

-- INSERT INTO quizzes(quizz_title, quizz_description, cat_id)
-- VALUES ('Les anniversaires', 'Seras-tu capable de te souvenir des anniversaires des membres de la famille ?', '4'),
-- ('Les voyages', 'Sauras-tu te souvenir de tous les lieux qu''on a visité ?', '4'),
-- ('Les déménagements', 'On est une famille de voyageurs, sauras-tu retrouver tous les déménagements qu''on a fait ?', '4');

-- INSERT INTO questions(quizz_id, question)
-- VALUES ('1', 'Quel est l''anniversaire de Kim ?'),
-- ('1', 'Quel est l''anniversaire de Véronique ?'),
-- ('1', 'Quel est l''anniversaire de Jérôme et Jenny ?');

-- INSERT INTO responses(question_id, answer, isCorrect)
-- VALUES ('1', '19 octobre', 'true'),
-- ('2', '17 mai', 'true'),
-- ('3', '25 janvier', 'true');

-- INSERT INTO responses(question_id, answer)
-- VALUES ('1', '17 octobre'),
-- ('1', '18 octobre'),
-- ('2', '19 mai'),
-- ('2', '18 mai'),
-- ('3', '26 janvier'),
-- ('3', '27 janvier');

insert into user_score(user_id, total_score) 
values ('1','0');

-- insert into user_score_per_cat(cat_id, user_id, score) 
-- values ('1','1','0'),
-- ('2','1','0'),
-- ('3','1','0'),
-- ('4','1','0');

create view v_selected_questions as
select
    s.session_id,
    q.question_id,
    q.question,
    array_agg(r.response_id) response_ids,
    array_agg(r.answer) answers,
    array_agg(r.isCorrect) isCorrects,
    array_agg(sr.response_id notnull) selected
from 
    quizz_sessions s
    join questions q on (q.quizz_id = s.quizz_id)
    join responses r on (r.question_id = q.question_id)
    left join sess_resp sr on (sr.session_id = s.session_id and sr.response_id = r.response_id)
group by
    q.question_id, s.session_id
order by
    s.session_id, q.question_id
;

-- create view v_score_per_user as
-- select
--     u.user_id,
--     s.total_score,
--     c.cat_id,
--     c.cat_name,
--     sc.score
-- from
--     users u
--     join user_score s on (u.user_id = s.user_id)
--     join user_score_per_cat sc on (u.user_id = sc.user_id)
--     join categories c on (sc.cat_id = c.cat_id)
-- group by 
--     c.cat_name
-- ;

-- create view v_score_per_user as
-- select
--     u.user_id,
--     sum(s.total_score),
--     -- sc.score
-- from
--     users u
--     join user_score s on (u.user_id = s.user_id)
--     -- join user_score_per_cat sc on (u.user_id = sc.user_id)
-- ;ore_per_user as
-- select 
--     qs.user_id,
--     c.cat_id,
--     c.cat_name,
--     coalesce(sum(qs.score), 0)
-- from
--     categories c
--     left join quizzes q on (q.cat_id = c.cat_id)
--     left join quizz_sessions qs on (qs.quizz_id = q.quizz_id)
-- group by c.cat_id, qs.user_id
-- ;

-- create view v_score_per_user as
-- select 
--     qs.user_id,
--     c.cat_id,
--     c.cat_name,
--     coalesce(sum(qs.score), 0)
-- from
--     categories c
--     left join quizzes q on (q.cat_id = c.cat_id)
--     left join quizz_sessions qs on (qs.quizz_id = q.quizz_id)
-- group by c.cat_id, qs.user_id
-- ;

INSERT INTO quizzes(quizz_title, quizz_description, cat_id)
VALUES ('Les films et cinéma', 'Connaissez-vous les grands classiques du cinéma d''hier et d''aujourd''hui ?', '1');

INSERT INTO questions(quizz_id, question)
VALUES ('1', 'Qui a réalisé "La bon, la brute et le truand" avec Clint Eastwood et Lee Van Cleef en 1966 ?'),
('1', 'Avatar" est un long-métrage mis en scène par Steven Spielberg.'),
('1', 'Quel personnage interprété par Louis De Funès crie "ma cassette, ma cassette" ?'),
('1', 'Pour quel film Jean Dujardin a-t-il reçu l''Oscar du meilleur acteur en 2012 ?'),
('1', 'Lequel de ces films n''a pas été realisé par Luc Besson ?'),
('1', 'Quel acteur incarne Nino Quincampoix dans le film Le Fabuleux Destin d’Amélie Poulain ?'),
('1', 'Dans "Les Tontons flingueurs" quel est le surnom de Louis ?'),
('1', 'Quel acteur incarne Pierrot le fou dans le film de Jean-Luc Godard ?'),
('1', 'Qui est la réalisatrice du film Polisse, sorti au cinéma en 2011 ?'),
('1', 'Quel est le nom du personnage joué par Francis Huster dans le film “Le Dîner de Cons” ?');

INSERT INTO responses(question_id, answer, isCorrect)
VALUES ('1', 'Clint Eastwood', 'false'),
('1', 'Sergio Corbucci', 'false'),
('1', 'Sergio Leone', 'true'),
('2', 'Vrai', 'false'),
('2', 'Faux', 'true'),
('3', 'Harpagon (L''Avare)', 'true'),
('3', 'Stanislas Lefort (La Grande Vadrouille)', 'false'),
('3', 'Léopold Saroyan (Le Corniaud)', 'false'),
('4', 'Brice de Nice', 'false'),
('4', 'The Artist', 'true'),
('4', 'Les petits mouchoirs', 'false'),
('5', 'Taxi', 'true'),
('5', 'Le Grand Bleu', 'false'),
('5', 'Lucy', 'false'),
('6', 'Mathieu Kassovitz', 'true'),
('6', 'Guillaume Canet', 'false'),
('6', 'Vincent Cassel', 'false'),
('7', 'Le Filou', 'false'),
('7', 'Le Menteur', 'false'),
('7', 'Le Mexicain', 'true'),
('8', 'Alain Delon', 'false'),
('8', 'Jean Gabin', 'false'),
('8', 'Jean-Paul Belmondo', 'true'),
('9', 'Marina Foïs', 'false'),
('9', 'Maïwenn', 'true'),
('9', 'Emmanuelle Bercot', 'false'),
('10', 'Lucien Cheval', 'false'),
('10', 'Pierre Brochant', 'false'),
('10', 'Juste Leblanc', 'true');


INSERT INTO quizzes(quizz_title, quizz_description, cat_id)
VALUES ('Gastronomie et nourriture française', 'Avez-vous une bonne connaissance des plats régionaux français, des traditions culinaires, de ses fromages …? Alors tentez votre chance !', '1');

INSERT INTO questions(quizz_id, question)
VALUES ('2', 'Quel fromage est traditionnellement rôti au four dans sa boîte avant d''être dégusté ?'),
('2', 'De quel coin de France est originaire l''Aligot ?'),
('2', 'D''après une ancienne tradition provençale, combien de desserts sont servis lors de la veillée de la fête de Noël ?'),
('2', 'Quel est le rôle d''un fusil dans une cuisine ?'),
('2', 'Quel légume est appelé chicon dans le Nord de la France ?'),
('2', 'Que sont les panisses, spécialités qui se mangent de Nice à Marseille ?'),
('2', 'La tapenade est une spécialité niçoise et provençale composée d''olives broyées, d''anchois et...'),
('2', 'Avec quel aliment est mélangé la graisse d''oie pour faire la garbure béarnaise ?'),
('2', 'Avec quoi sont faits les grattons lyonnais ?'),
('2', 'Combien le guide Michelin compte-t-il de femmes cheffes trois fois étoilées en France ?');

INSERT INTO responses(question_id, answer, isCorrect)
VALUES ('11', 'Maroilles', 'false'),
('11', 'Mont d''or', 'true'),
('11', 'Tête de Moine', 'false'),
('12', 'Bourgogne', 'false'),
('12', 'Bretagne', 'false'),
('12', 'Auvergne', 'true'),
('13', '8', 'false'),
('13', '11', 'false'),
('13', '13', 'true'),
('14', 'A aiguiser des couteaux', 'false'),
('14', 'A épépiner les fruits', 'true'),
('14', 'A écaler les oeufs', 'false'),
('15', 'Le poireau', 'false'),
('15', 'Le radis', 'false'),
('15', 'L''endive', 'true'),
('16', 'Des beignets anisés', 'false'),
('16', 'Des préparations à base de pois chiches et frite dans l''huile', 'false'),
('16', 'Des pâtisseries à base d''anis', 'false'),
('17', 'De câpres', 'true'),
('17', 'D''oignons cuits', 'false'),
('17', 'D''ail pilé', 'false'),
('18', 'Les pruneaux', 'false'),
('18', 'Les cardes', 'false'),
('18', 'Le chou vert', 'true'),
('19', 'Avec de la graisse et de la viande de porc', 'true'),
('19', 'Avec des abats de boeuf', 'false'),
('19', 'Avec des pieds de cochon', 'false'),
('20', '1', 'true'),
('20', '10', 'false'),
('20', '20', 'false');

INSERT INTO quizzes(quizz_title, quizz_description, cat_id)
VALUES ('Football', 'Passionné de foot, vous vous croyez incollable sur le football ? ', '1');

INSERT INTO questions(quizz_id, question)
VALUES ('3', 'Dans quel club de foot, Karim Benzema est-il resté durant 15 ans ?'),
('3', 'Combien de buts Zinédine Zidane a-t-il marqués tout au long de sa carrière dans l''Équipe de France de Foot ?'),
('3', 'En quelle année Antoine Griezmann est-il sélectionné pour la première fois en équipe de France de foot ?'),
('3', 'Le recours à l''assistance vidéo à l''arbitrage est autorisé au football. Vrai ou faux ?'),
('3', 'Dans quel pays se déroule la toute première Coupe du Monde de football en 1930 ?'),
('3', 'Avec quel score la France s''impose-t-elle en finale de la Coupe du Monde 2018 face à la Croatie ?'),
('3', 'Lequel de ces joueurs a marqué exactement 200 buts lors de sa carrière au PSG ?'),
('3', 'Quelle chaîne de télévision a été pendant plusieurs années propriétaire du PSG ?'),
('3', 'Quel homme d''affaires a pris la tête du club de foot de l''Olympique lyonnais en 1987 ?'),
('3', 'Parmi ces grands noms de l''Histoire du football, lequel n''a-t-il jamais joué au FC Barcelone ?');

INSERT INTO responses(question_id, answer, isCorrect)
VALUES ('21', 'Real Madrid', 'true'),
('21', 'FC Barcelone', 'false'),
('21', 'Atlético de Madrid', 'false'),
('22', '20', 'false'),
('22', '31', 'true'),
('22', '46', 'false'),
('23', '2012', 'false'),
('23', '2014', 'true'),
('23', '2016', 'false'),
('24', 'Vrai', 'true'),
('24', 'Faux', 'false'),
('25', 'France', 'false'),
('25', 'Mexique', 'false'),
('25', 'Uruguay', 'true'),
('26', '4-1', 'false'),
('26', '4-2', 'true'),
('26', '3-2','false'),
('27', 'Zlatan Ibrahimovic', 'false'),
('27', 'Edinson Cavani', 'true'),
('27', 'Angel Di Maria', 'false'),
('28', 'TF1', 'false'),
('28', 'La Cinq', 'false'),
('28', 'Canal+', 'true'),
('29', 'Jean-Michel Aulas', 'true'),
('29', 'Noël le Graët', 'false'),
('29', 'Louis Nicollin', 'false'),
('30', 'Neymar', 'false'),
('30', 'Kaká', 'true'),
('30', 'Deco', 'false');

INSERT INTO quizzes(quizz_title, quizz_description, cat_id)
VALUES ('Problèmes de Maths', 'Saurez-vous résoudre ces petits problèmes de mathématiques ? ', '2');

INSERT INTO questions(quizz_id, question)
VALUES ('4', 'Je profite d''une promo de 20% sur un achat à 40 €. Quel est le montant que je vais économiser ?'),
('4', 'Si l''on divise une minute en 4, combien y a-t-il de secondes dans chaque part ?'),
('4', 'Sophie dépense 11 € par mois pour son abonnement à la bibliothèque. Quel montant paye-t-elle à l''année ?'),
('4', 'Un escargot est au fond d''un puits de 10 mètres. Le jour, il monte de 3 mètres et la nuit il redescend de 2 mètres. Quand arrivera-t-il hors du puits ?'),
('4', '89, 106, 113, 118, 128, ? Quel est le prochain nombre de cette suite logique ?'),
('4', 'Il y a plusieurs livres sur une étagère. Si un livre est le cinquième en partant de la gauche et le cinquième en partant de la droite, combien y a t-il de livres sur cette étagère ?'),
('4', 'J''assiste à un match de rugby qui se termine dans 8 minutes. Une mi-temps dure 40 minutes et la pause entre les deux mi-temps dure un quart d''heure. Depuis combien de minutes le match a-t-il commencé ?'),
('4', 'Le cross du collège est parti à 13h13. J''ai franchi la ligne d''arrivée à 14h25. Quel est le temps que j''ai réalisé en minutes ?'),
('4', 'Voici une suite logique de nombres : 7 ; 21 ; 18 ; 72 ; 68… Quel est le nombre suivant ?'),
('4', 'Toutes mes économies, soit 270 000 € sont placées à la banque sur trois comptes A, B et C. La répartition est la suivante : les deux-tiers sur le compte A, 15% sur le compte B et le reste sur le compte C. 
Combien y a-t-il en euros sur le compte C ?');

INSERT INTO responses(question_id, answer, isCorrect)
VALUES ('31', '5€', 'false'),
('31', '8€', 'true'),
('31', '12€', 'false'),
('32', '20', 'false'),
('32', '12,5', 'false'),
('32', '15', 'true'),
('33', '121', 'false'),
('33', '124', 'false'),
('33', '132', 'true'),
('34', '5 jours', 'false'),
('34', '7 jours', 'false'),
('34', '8 jours', 'true'),
('35', '138', 'false'),
('35', '139', 'true'),
('35', '143', 'false'),
('36', '9', 'true'),
('36', '10', 'false'),
('36', '11','false'),
('37', '82', 'false'),
('37', '87', 'true'),
('37', '88', 'false'),
('38', '1h12', 'true'),
('38', '1h15', 'false'),
('38', '1h18', 'true'),
('39', '144', 'false'),
('39', '216', 'false'),
('39', '340', 'true'),
('40', '49 500€', 'true'),
('40', '51 500€', 'false'),
('40', '55 000€', 'false');