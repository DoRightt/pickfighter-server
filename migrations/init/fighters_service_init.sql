\c pickfighter;

-- Create tables
CREATE TABLE IF NOT EXISTS fighters.divisions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    value VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS fighters.fighters (
    fighter_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    nickname VARCHAR(255) DEFAULT '',
    division_id INT REFERENCES fighters.divisions(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL,
    hometown VARCHAR(100) DEFAULT '',
    trains_at VARCHAR(100) DEFAULT '',
    fighting_style VARCHAR(100) DEFAULT '',
    age INT NOT NULL,
    height DOUBLE PRECISION,
    weight DOUBLE PRECISION,
    octagon_debut VARCHAR(50) DEFAULT '',
    debut_timestamp BIGINT NOT NULL,
    reach INT,
    leg_reach INT,
    fighter_url VARCHAR(255) NOT NULL,
    image_url TEXT,
    wins INT DEFAULT 0 NOT NULL,
    loses INT DEFAULT 0 NOT NULL,
    draw INT DEFAULT 0 NOT NULL
);

CREATE TABLE IF NOT EXISTS fighters.fighter_stats (
    stat_id SERIAL PRIMARY KEY,
    fighter_id INT REFERENCES fighters.fighters(fighter_id),
    total_sig_str_landed INT,
    total_sig_str_attempted INT,
    str_accuracy INT,
    total_tkd_landed INT,
    total_tkd_attempted INT,
    tkd_accuracy INT,
    sig_str_landed DOUBLE PRECISION,
    sig_str_absorbed DOUBLE PRECISION,
    sig_str_defense INT,
    takedown_defense INT,
    takedown_avg DOUBLE PRECISION,
    submission_avg DOUBLE PRECISION,
    knockdown_avg DOUBLE PRECISION,
    avg_fight_time VARCHAR(50),
    win_by_ko INT,
    win_by_sub INT,
    win_by_dec INT
);

-- Create sequences
CREATE SEQUENCE IF NOT EXISTS fighters.divisions_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS fighters.pf_fighter_stats_stat_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;
CREATE SEQUENCE IF NOT EXISTS fighters.pf_fighters_fighter_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;

-- Set ownership of sequences
ALTER SEQUENCE fighters.divisions_id_seq OWNED BY fighters.divisions.id;
ALTER SEQUENCE fighters.pf_fighter_stats_stat_id_seq OWNED BY fighters.fighter_stats.stat_id;
ALTER SEQUENCE fighters.pf_fighters_fighter_id_seq OWNED BY fighters.fighters.fighter_id;

-- Set default values for columns linked to sequences
ALTER TABLE ONLY fighters.divisions ALTER COLUMN id SET DEFAULT nextval('fighters.divisions_id_seq'::regclass);
ALTER TABLE ONLY fighters.fighter_stats ALTER COLUMN stat_id SET DEFAULT nextval('fighters.pf_fighter_stats_stat_id_seq'::regclass);
ALTER TABLE ONLY fighters.fighters ALTER COLUMN fighter_id SET DEFAULT nextval('fighters.pf_fighters_fighter_id_seq'::regclass);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_fighter_name ON fighters.fighters(name);
CREATE INDEX IF NOT EXISTS idx_fighter_status ON fighters.fighters(status);
CREATE INDEX IF NOT EXISTS idx_fighter_division_id ON fighters.fighters(division_id);
CREATE INDEX IF NOT EXISTS idx_fighter_stats_fighter_id ON fighters.fighter_stats(fighter_id);

-- Add unique constraints
ALTER TABLE fighters.fighter_stats ADD CONSTRAINT unique_fighter_stat UNIQUE (fighter_id);
ALTER TABLE fighters.fighters ADD CONSTRAINT pf_fighters_name_debut_timestamp_key UNIQUE (name, debut_timestamp);

-- Insert divisions
INSERT INTO fighters.divisions (id, name, value)
VALUES
    (12, 'Women''s Featherweight Division', 'womens_featherweight'),
    (11, 'Women''s Bantamweight Division', 'womens_bantamweight'),
    (10, 'Women''s Flyweight Division', 'womens_flyweight'),
    (9, 'Women''s Strawweight Division', 'womens_strawweight'),
    (8, 'Heavyweight Division', 'heavyweight'),
    (7, 'Light Heavyweight Division', 'lightheavyweight'),
    (6, 'Middleweight Division', 'middleweight'),
    (5, 'Welterweight Division', 'welterweight'),
    (4, 'Lightweight Division', 'lightweight'),
    (3, 'Featherweight Division', 'featherweight'),
    (2, 'Bantamweight Division', 'bantamweight'),
    (1, 'Flyweight Division', 'flyweight'),
    (0, 'No Division', 'no_division')
ON CONFLICT (id) DO NOTHING;