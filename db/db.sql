---------------------------------------------------------
---------------------USER setup--------------------------

create table users (
    user_id                int                  primary key generated always as identity,
    pseudonym              VARCHAR(50)          unique not null,
    profil_pic             bytea                ,
    gallery_pic            bytea                
);

create table user_score (
    user_id                 int                 primary key generated always as identity,
    total_score             bigint              ,
    total_score_per_cat     INT                 
);

---------------------------------------------------------
-------------------QUIZZ setup---------------------------

create table categories (
    cat_id                  int                 primary key generated always as identity,
    cat_name                VARCHAR(124)        not null,
    cat_short_name          VARCHAR(124)        not null,
    cat_description         text                not null,
    cat_image               text                not null           
);

create table quizzes (
    quizz_id                int                 primary key generated always as identity,
    quizz_title             VARCHAR(512)        not null,
    quizz_description       text                not null,
    created_at              timestamptz         not null default current_timestamp,
    cat_id                  INT                 references categories(cat_id)
);

create table questions (
    question_id             int                 primary key generated always as identity,
    quizz_id                int                 references quizzes(quizz_id) not null,
    question                text                not null,
    response                text                not null,
    order_questions         smallint
);

create table quizz_sessions (
    session_id             INT                  primary key generated always as identity,
    quizz_id               INT                  references quizzes(quizz_id) not null,
    user_id                INT                  references users(user_id) not null,
    finished               bool                 not null default false,
    score                  INT                  not null default '0'
);

create table reponses (
    responses_id           int                  primary key generated always as identity,
    question_id            INT                  references questions(question_id),
    session_id             INT                  references quizz_sessions(session_id),
    isCorrect              bool
);

---------------------------------------------------------
---------------------INITIAL setup-----------------------

INSERT INTO users(pseudonym)
VALUES ('Alain');

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