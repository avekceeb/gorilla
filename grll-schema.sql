-- status dictionary
DROP TABLE IF EXISTS status CASCADE;
CREATE TABLE "status" (
	"id"   SMALLINT     NOT NULL,
	"name" VARCHAR(63)  NOT NULL,
	CONSTRAINT "pk_status" PRIMARY KEY ("id")
);
INSERT INTO status (id,name) VALUES (1,'passed');
INSERT INTO status (id,name) VALUES (2,'failed');
INSERT INTO status (id,name) VALUES (3,'skipped');


-- test table
DROP TABLE IF EXISTS test CASCADE;
CREATE TABLE "test" (
	"id"   SERIAL NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	CONSTRAINT "pk_test" PRIMARY KEY ("id"),
	CONSTRAINT "un_test" UNIQUE ("name")
);

-- test results
-- TODO: ??? multiple results for same run + same test
DROP TABLE IF EXISTS result CASCADE;
CREATE TABLE "result" (
	"testrun" INT NOT NULL,
	"test"    INT NOT NULL,
	"status"  INT NOT NULL,
	"message" VARCHAR(1023) DEFAULT NULL
	--"DURATION" INTERVAL DAY (2) TO SECOND (6),
);

DROP TABLE IF EXISTS measurement CASCADE;
CREATE TABLE "measurement" (
	"testrun" INT NOT NULL,
	"test"    INT NOT NULL,
	"val"     INT NOT NULL,
	"unit"    VARCHAR(1023) DEFAULT NULL
	--"DURATION" INTERVAL DAY (2) TO SECOND (6),
);

-- test run
DROP TABLE IF EXISTS testrun CASCADE;
-- TODO: unique
CREATE TABLE "testrun" (
	"id"   SERIAL NOT NULL,
	"name" VARCHAR(255),
	"ts"   TIMESTAMPTZ,
	"link" VARCHAR(1023) DEFAULT NULL,
	CONSTRAINT "pk_testrun" PRIMARY KEY ("id")
);

-- tags for run
DROP TABLE IF EXISTS tag CASCADE;
CREATE TABLE "tag" (
	"id" SERIAL NOT NULL,
	"name" VARCHAR(255) NOT NULL,
	CONSTRAINT "pk_tag" PRIMARY KEY ("id"),
	CONSTRAINT "un_tag" UNIQUE ("name")
);

-- stick arbitrary tags to testrun
DROP TABLE IF EXISTS testrun_tag CASCADE;
CREATE TABLE "testrun_tag" (
	"testrun" INT NOT NULL,
	"tag"     INT NOT NULL
);

-----------------------------------------------
CREATE OR REPLACE FUNCTION new_testrun(n VARCHAR, t TIMESTAMPTZ, l VARCHAR)
	RETURNS INTEGER AS
$BODY$
DECLARE v_id INTEGER;
BEGIN
	INSERT INTO testrun ("name", "ts", "link")
	    VALUES (n, t, l)
	    RETURNING id INTO v_id;
	RETURN v_id;
END;
$BODY$
  LANGUAGE plpgsql VOLATILE;

-----------------------------------------------
-- black woodoo magic
-- return id of existing or insert new
CREATE OR REPLACE FUNCTION new_tag(n VARCHAR)
	RETURNS INTEGER AS
$$
DECLARE ni INTEGER;
BEGIN
WITH s AS (
    SELECT id
    FROM tag
    WHERE name = n
), i AS (
    INSERT INTO tag ("name")
    SELECT n
    WHERE NOT EXISTS (SELECT 1 FROM s)
    RETURNING id
)
SELECT id
FROM i
UNION ALL
SELECT id
FROM s INTO ni;
	RETURN ni;
END;
$$
  LANGUAGE plpgsql VOLATILE;

-----------------------------------------------

CREATE OR REPLACE FUNCTION new_test(n varchar)
	RETURNS INTEGER AS
$$
DECLARE test_id INTEGER;
BEGIN
WITH s AS (
    SELECT id
    FROM "test"
    WHERE name = n
), i AS (
    INSERT INTO "test" ("name")
    SELECT n
    WHERE NOT EXISTS (SELECT 1 FROM s)
    RETURNING id
)
SELECT id
    FROM i
    UNION ALL
    SELECT id
        FROM s INTO test_id;
RETURN test_id;
END;
$$
  LANGUAGE plpgsql VOLATILE;


CREATE OR REPLACE FUNCTION new_result(t VARCHAR, r VARCHAR, m VARCHAR, tr INTEGER) RETURNS void
AS $$
DECLARE test_id   INTEGER;
DECLARE status_id INTEGER;
BEGIN
	SELECT new_test(t) INTO test_id;
	SELECT id FROM status WHERE name = r INTO status_id;
	INSERT INTO result ("testrun", "test", "status", "message")
		VALUES (tr, test_id, status_id, m);
END;
$$
  LANGUAGE plpgsql VOLATILE;


CREATE OR REPLACE FUNCTION new_measurement(t VARCHAR, v REAL, u VARCHAR, tr INTEGER)
    RETURNS void
AS $$
DECLARE test_id   INTEGER;
BEGIN
	SELECT new_test(t) INTO test_id;
	INSERT INTO measurement ("testrun", "test", "val", "unit")
		VALUES (tr, test_id, v, u);
END;
$$
  LANGUAGE plpgsql VOLATILE;


-----------------------------------------------

CREATE OR REPLACE FUNCTION set_tag(r INTEGER, t VARCHAR) RETURNS void
AS $$
BEGIN
	INSERT INTO testrun_tag ("testrun", "tag")
		SELECT r, new_tag(t);
END;
$$
  LANGUAGE plpgsql VOLATILE;

-----------------------------------------------
CREATE OR REPLACE FUNCTION get_testruns_by_tag(VARCHAR) RETURNS SETOF testrun AS
$$
	SELECT DISTINCT testrun.id, testrun.name, testrun.ts, testrun.link
		FROM testrun
		INNER JOIN testrun_tag ON testrun_tag.testrun = testrun.id
		INNER JOIN tag ON tag.id = testrun_tag.tag
		WHERE tag.name = $1 ;
$$
  language 'sql';
-----------------------------------------------

