BEGIN;

CREATE TEXT SEARCH DICTIONARY simple_no_stop (
    TEMPLATE = pg_catalog.simple,
    stopwords = 'empty' );

CREATE TEXT SEARCH CONFIGURATION simple_no_stop (
    PARSER = pg_catalog."default" );

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR asciiword WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR word WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR numword WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR email WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR url WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR host WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR sfloat WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR version WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR hword_numpart WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR hword_part WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR hword_asciipart WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR numhword WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR asciihword WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR hword WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR url_path WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR file WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR "float" WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR "int" WITH simple_no_stop;

ALTER TEXT SEARCH CONFIGURATION simple_no_stop
    ADD MAPPING FOR uint WITH simple_no_stop;

CREATE OR REPLACE FUNCTION get_full_name(VARCHAR, VARCHAR) RETURNS VARCHAR AS $$
    SELECT COALESCE($1, ''::varchar) || ' '::varchar || COALESCE($2, ''::varchar)
$$ LANGUAGE sql IMMUTABLE;

CREATE OR REPLACE FUNCTION get_trigram_tsvector(TEXT) RETURNS tsvector AS $$
    SELECT
    to_tsvector(
        'simple_no_stop'::regconfig,
        array_to_string(
            show_trgm(
                lower($1)
            ), ' '::text
        )
    )
$$ LANGUAGE sql IMMUTABLE;

CREATE OR REPLACE FUNCTION search_for_trigram_fts(TEXT) RETURNS SETOF RECORD LANGUAGE SQL STABLE AS $$
    SELECT id, group_id, leechers, seeds, infohash, name,
        ts_rank(get_trigram_tsvector(
            name),
            to_tsquery(array_to_string(show_trgm(lower($1)), '|'::text))) AS rank
    FROM torrents
    WHERE
        get_trigram_tsvector(name) @@
        to_tsquery(array_to_string(show_trgm(lower($1)), '|'::text))

        UNION

    SELECT id, group_id, leechers, seeds, infohash, name,
        ts_rank(get_trigram_tsvector(name),
            to_tsquery(array_to_string(show_trgm(lower($1)), '|'::text))) AS rank
    FROM torrents
    WHERE
        get_trigram_tsvector(name) @@
        to_tsquery(array_to_string(show_trgm(lower($1)), '|'::text))

$$;

CREATE OR REPLACE FUNCTION search_for_name_fts(TEXT) RETURNS SETOF RECORD LANGUAGE SQL STABLE AS $$
    SELECT id, group_id, leechers, seeds, infohash, name,
        ts_rank(
            to_tsvector('simple_no_stop'::regconfig, name),
            to_tsquery(array_to_string(regexp_split_to_array($1, E'\\s+'), '|'))
        )
    FROM torrents
        WHERE
        to_tsvector(name) @@
        to_tsquery(array_to_string(regexp_split_to_array($1, E'\\s+'), '|'))
$$;

CREATE OR REPLACE FUNCTION
    search_for_name_trgm(TEXT)
    RETURNS SETOF RECORD
    LANGUAGE SQL STABLE AS
$$
    SELECT id, group_id, leechers, seeds, infohash, name,
        ts_rank(get_trigram_tsvector(name),
            to_tsquery(array_to_string(show_trgm(lower($1)), '|'::text))
        ) AS rank
    FROM torrents
    WHERE
        get_trigram_tsvector(name) @@
        to_tsquery(array_to_string(show_trgm(lower($1)), '|'::text))
$$;

CREATE OR REPLACE FUNCTION search_for_name(TEXT) RETURNS SETOF RECORD LANGUAGE SQL STABLE AS $$
    SELECT id, group_id, leechers, seeds, infohash, name, sources, summed_rank +
        similarity(COALESCE(name, ''), $1) +
        (
            CASE WHEN soundex($1) = soundex(name) THEN 1 ELSE 0 END +
            CASE WHEN metaphone($1, 10) = metaphone(name, 10) THEN 1 ELSE 0 END
         ) / 2 AS rank
   FROM (

        SELECT id, group_id, leechers, seeds, infohash, name, SUM(rank) AS summed_rank, ARRAY_ACCUM(source) AS sources FROM (
            SELECT id, group_id, leechers, seeds, infohash, name, rank, 'trigram_fts'::text AS source
            FROM search_for_trigram_fts($1) AS a(id iNTEGER,group_id iNTEGER,leechers iNTEGER,seeds iNTEGER, infohash VARCHAR, name VARCHAR, rank REAL)
                UNION
            SELECT id, group_id, leechers, seeds, infohash, name, rank, 'name_fts'::text AS source
            FROM search_for_name_fts($1) AS a(id iNTEGER,group_id iNTEGER,leechers iNTEGER,seeds iNTEGER, infohash VARCHAR, name VARCHAR,  rank REAL)
                UNION
            SELECT id, group_id, leechers, seeds, infohash, name, rank, 'name_trgm'::text AS source
            FROM search_for_name_trgm($1) AS a(id iNTEGER,group_id iNTEGER,leechers iNTEGER,seeds iNTEGER, infohash VARCHAR, name VARCHAR, rank REAL)
        ) a
        GROUP BY id, name, group_id, leechers, seeds, infohash
    ) b
    ORDER BY rank DESC
$$;



CREATE INDEX name_trgm_ix ON torrents USING GIST (get_trigram_tsvector(name));

CREATE INDEX name_fts_ix ON torrents USING gist (to_tsvector('simple_no_stop'::regconfig,
                name));

COMMIT;