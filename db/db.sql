---------------------------------------------------------
---------------------USER setup--------------------------

create table users (
    user_id                int                  primary key generated always as identity,
    pseudonym              varchar(50)          unique not null,
    profil_pic             bytea                ,
    gallery_pic            bytea                
);

create table user_score (
    user_id                 int                 primary key generated always as identity,
    total_score             bigint              ,
    total_score_per_cat     int                 
);

---------------------------------------------------------
-------------------QUIZZ setup---------------------------

create table categories (
    cat_id                  int                 primary key generated always as identity,
    cat_name                varchar(124)        unique not null,
    cat_short_name          varchar(124)        unique not null,
    cat_description         text                not null,
    cat_image               text                not null           
);

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

INSERT INTO quizzes(quizz_title, quizz_description, cat_id)
VALUES ('Les anniversaires', 'Seras-tu capable de te souvenir des anniversaires des membres de la famille ?', '4'),
('Les voyages', 'Sauras-tu te souvenir de tous les lieux qu''on a visité ?', '4'),
('Les déménagements', 'On est une famille de voyageurs, sauras-tu retrouver tous les déménagements qu''on a fait ?', '4');

INSERT INTO questions(quizz_id, question)
VALUES ('1', 'Quel est l''anniversaire de Kim ?'),
('1', 'Quel est l''anniversaire de Véronique ?'),
('1', 'Quel est l''anniversaire de Jérôme et Jenny ?');

INSERT INTO responses(question_id, answer, isCorrect)
VALUES ('1', '19 octobre', 'true'),
('2', '17 mai', 'true'),
('3', '25 janvier', 'true');

INSERT INTO responses(question_id, answer)
VALUES ('1', '17 octobre'),
('1', '18 octobre'),
('2', '19 mai'),
('2', '18 mai'),
('3', '26 janvier'),
('3', '27 janvier');

-- create view finished_quizz as
-- select responses.responses_id, question_id, answer, isCorrect
-- from  responses
-- inner join sess_resp
-- on responses.responses_id = sess_resp.responses_id;

