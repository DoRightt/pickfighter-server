SELECT 'CREATE DATABASE fighters_db' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'fighters_db');

\c fighters_db;

--- pf_fighters table

CREATE TABLE IF NOT EXISTS public.pf_fighters (
    fighter_id serial NOT NULL,
    name character varying(255) NOT NULL,
    nickname character varying(255) DEFAULT ''::character varying,
    division integer NOT NULL,
    status character varying(50) NOT NULL,
    hometown character varying(100) DEFAULT ''::character varying,
    trains_at character varying(100) DEFAULT ''::character varying,
    fighting_style character varying(100) DEFAULT ''::character varying,
    age integer NOT NULL,
    height double precision,
    weight double precision,
    octagon_debut character varying(50) DEFAULT ''::character varying,
    debut_timestamp bigint NOT NULL,
    reach integer,
    leg_reach integer,
    fighter_url character varying(255) NOT NULL,
    image_url text,
    wins integer DEFAULT 0 NOT NULL,
    loses integer DEFAULT 0 NOT NULL,
    draw integer DEFAULT 0 NOT NULL
);

ALTER TABLE ONLY public.pf_fighters
    ADD CONSTRAINT pf_fighters_name_debut_timestamp_key UNIQUE (name, debut_timestamp, fighter_url);

ALTER TABLE ONLY public.pf_fighters
    ADD CONSTRAINT pf_fighters_pk PRIMARY KEY (fighter_id);

CREATE UNIQUE INDEX pf_fighters_fighter_url_uindex ON public.pf_fighters USING btree (fighter_url);

INSERT INTO public.pf_fighters (fighter_id, name, nickname, division, status, hometown, trains_at, fighting_style, age, height, weight, octagon_debut, debut_timestamp, reach, leg_reach, fighter_url, image_url, wins, loses, draw) VALUES (57918, 'Rostem Akman', '', 4, 'Not Fighting', '', '', '', 31, 70, 171, 'Jun. 1, 2019', 1559347200, 72, 38, 'https://www.ufc.com/athlete/rostam-akman', 'https://dmxg5wxfqgb4u.cloudfront.net/styles/athlete_bio_full_body/s3/image/ufc-fighter-container/71542/profile-galery/fullbodyleft-picture/AKMAN_ROSTAM_L.png?VersionId=s0Xyj_DSjzTjrVVAvaeImvkXyz9WVs3Z&itok=sOszamHM', 0, 2, 0);
INSERT INTO public.pf_fighters (fighter_id, name, nickname, division, status, hometown, trains_at, fighting_style, age, height, weight, octagon_debut, debut_timestamp, reach, leg_reach, fighter_url, image_url, wins, loses, draw) VALUES (57919, 'Razak Al-Hassan', '"Razor"', 6, 'Not Fighting', '', '', '', 41, 74, 205, 'Dec. 10, 2008', 1228867200, 0, 0, 'https://www.ufc.com/athlete/razak-al-hassan', '', 7, 2, 0);
INSERT INTO public.pf_fighters (fighter_id, name, nickname, division, status, hometown, trains_at, fighting_style, age, height, weight, octagon_debut, debut_timestamp, reach, leg_reach, fighter_url, image_url, wins, loses, draw) VALUES (57901, 'Tank Abbott', '"Tank"', 7, 'Not Fighting', '', '', '', 0, 72, 253, 'Jul. 14, 1995', 805680000, 0, 0, 'https://www.ufc.com/athlete/tank-abbott', '', 8, 10, 0);
INSERT INTO public.pf_fighters (fighter_id, name, nickname, division, status, hometown, trains_at, fighting_style, age, height, weight, octagon_debut, debut_timestamp, reach, leg_reach, fighter_url, image_url, wins, loses, draw) VALUES (57904, 'Daichi Abe', '', 4, 'Not Fighting', '', '', '', 31, 69, 170.5, 'Sep. 22, 2017', 1506038400, 71, 42, 'https://www.ufc.com/athlete/daichi-abe', 'https://dmxg5wxfqgb4u.cloudfront.net/styles/athlete_bio_full_body/s3/2022-03/8c026dcd-60a4-457e-9cb5-b2c20969cd8f%252FDaichi-Abe_635302_LeftFullBodyImage.png?itok=5Gozk7xe', 6, 1, 0);
INSERT INTO public.pf_fighters (fighter_id, name, nickname, division, status, hometown, trains_at, fighting_style, age, height, weight, octagon_debut, debut_timestamp, reach, leg_reach, fighter_url, image_url, wins, loses, draw) VALUES (57905, 'Papy Abedi', '"Makambo"', 5, 'Not Fighting', '', '', '', 44, 71, 184.5, 'Nov. 5, 2011', 1320451200, 74, 0, 'https://www.ufc.com/athlete/papy-abedi', 'https://dmxg5wxfqgb4u.cloudfront.net/styles/athlete_bio_full_body/s3/2022-03/b234b354-8110-42f1-8dc3-e12c64426cce%252FPapy-Abedi_205746_LeftFullBodyImage.png?itok=cxcjAQ75', 9, 3, 0);
INSERT INTO public.pf_fighters (fighter_id, name, nickname, division, status, hometown, trains_at, fighting_style, age, height, weight, octagon_debut, debut_timestamp, reach, leg_reach, fighter_url, image_url, wins, loses, draw) VALUES (57913, 'Fabio Agu', '', 5, 'Active', '', '', '', 35, 0, 0, 'May. 16, 2024', 1715817600, 0, 0, 'https://www.ufc.com/athlete/fabio-agu', '', 0, 0, 0);


--- pf_fighters_stats table

CREATE TABLE IF NOT EXISTS public.pf_fighter_stats (
    stat_id serial NOT NULL,
    fighter_id integer,
    total_sig_str_landed integer,
    total_sig_str_attempted integer,
    str_accuracy integer,
    total_tkd_landed integer,
    total_tkd_attempted integer,
    tkd_accuracy integer,
    sig_str_landed double precision,
    sig_str_absorbed double precision,
    sig_str_defense integer,
    takedown_defense integer,
    takedown_avg double precision,
    submission_avg double precision,
    knockdown_avg double precision,
    avg_fight_time character varying(50),
    win_by_ko integer,
    win_by_sub integer,
    win_by_dec integer
);

INSERT INTO public.pf_fighter_stats (stat_id, fighter_id, total_sig_str_landed, total_sig_str_attempted, str_accuracy, total_tkd_landed, total_tkd_attempted, tkd_accuracy, sig_str_landed, sig_str_absorbed, sig_str_defense, takedown_defense, takedown_avg, submission_avg, knockdown_avg, avg_fight_time, win_by_ko, win_by_sub, win_by_dec) VALUES (29333, 57918, 54, 210, 25, 2, 5, 40, 1.7999999523162842, 3.9700000286102295, 58, 88, 1, 0, 0, '15:00', 5, 1, 0);
INSERT INTO public.pf_fighter_stats (stat_id, fighter_id, total_sig_str_landed, total_sig_str_attempted, str_accuracy, total_tkd_landed, total_tkd_attempted, tkd_accuracy, sig_str_landed, sig_str_absorbed, sig_str_defense, takedown_defense, takedown_avg, submission_avg, knockdown_avg, avg_fight_time, win_by_ko, win_by_sub, win_by_dec) VALUES (29334, 57919, 44, 135, 32, 0, 5, 0, 2.309999942779541, 3.7799999713897705, 60, 57, 0.7900000214576721, 0, 0, '09:32', 0, 0, 0);
INSERT INTO public.pf_fighter_stats (stat_id, fighter_id, total_sig_str_landed, total_sig_str_attempted, str_accuracy, total_tkd_landed, total_tkd_attempted, tkd_accuracy, sig_str_landed, sig_str_absorbed, sig_str_defense, takedown_defense, takedown_avg, submission_avg, knockdown_avg, avg_fight_time, win_by_ko, win_by_sub, win_by_dec) VALUES (29316, 57901, 12, 31, 38, 0, 0, 0, 2.4100000858306885, 10.029999732971191, 38, 67, 0, 0, 0, '01:40', 0, 0, 0);
INSERT INTO public.pf_fighter_stats (stat_id, fighter_id, total_sig_str_landed, total_sig_str_attempted, str_accuracy, total_tkd_landed, total_tkd_attempted, tkd_accuracy, sig_str_landed, sig_str_absorbed, sig_str_defense, takedown_defense, takedown_avg, submission_avg, knockdown_avg, avg_fight_time, win_by_ko, win_by_sub, win_by_dec) VALUES (29319, 57904, 171, 508, 33, 1, 2, 50, 3.799999952316284, 4.489999771118164, 57, 0, 0.33000001311302185, 0, 0.33000001311302185, '15:00', 4, 0, 2);
INSERT INTO public.pf_fighter_stats (stat_id, fighter_id, total_sig_str_landed, total_sig_str_attempted, str_accuracy, total_tkd_landed, total_tkd_attempted, tkd_accuracy, sig_str_landed, sig_str_absorbed, sig_str_defense, takedown_defense, takedown_avg, submission_avg, knockdown_avg, avg_fight_time, win_by_ko, win_by_sub, win_by_dec) VALUES (29320, 57905, 97, 176, 55, 0, 14, 0, 2.799999952316284, 3.1500000953674316, 49, 50, 3.4700000286102295, 1.2999999523162842, 0, '08:39', 5, 2, 2);
INSERT INTO public.pf_fighter_stats (stat_id, fighter_id, total_sig_str_landed, total_sig_str_attempted, str_accuracy, total_tkd_landed, total_tkd_attempted, tkd_accuracy, sig_str_landed, sig_str_absorbed, sig_str_defense, takedown_defense, takedown_avg, submission_avg, knockdown_avg, avg_fight_time, win_by_ko, win_by_sub, win_by_dec) VALUES (29328, 57913, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '00:00', 0, 0, 0);

ALTER TABLE ONLY public.pf_fighter_stats
    ADD CONSTRAINT pf_fighter_stats_pkey PRIMARY KEY (stat_id);

ALTER TABLE ONLY public.pf_fighter_stats
    ADD CONSTRAINT pf_fighter_stats_fighter_id_fkey FOREIGN KEY (fighter_id) REFERENCES public.pf_fighters(fighter_id);