PRAGMA foreign_keys=OFF;
BEGIN;
CREATE TABLE new_session (
    id                integer  PRIMARY KEY,
    callsign          text     NOT NULL,
    name              text     NOT NULL,
    prefix            text     NOT NULL,
    start             datetime NOT NULL,
    end               datetime NOT NULL,
    excludefromweek   boolean  NOT NULL,
    reporttotext      text     NOT NULL,
    reporttohtml      text     NOT NULL,
    tobbses           text     NOT NULL,
    downbbses         text     NOT NULL,
    messagetypes      text     NOT NULL,
    modified          boolean  NOT NULL,
    running           boolean  NOT NULL,
    imported          boolean  NOT NULL,
    report            text     NOT NULL
);
INSERT INTO new_session SELECT id, callsign, name, prefix, start, end, excludefromweek, reportto, '', tobbses, downbbses, messagetypes, modified, running, imported, report FROM session;
DROP TABLE session;
ALTER TABLE new_session RENAME TO session;
CREATE UNIQUE INDEX session_call_end_idx ON session (callsign, end);
CREATE INDEX session_end_idx ON session (end);
CREATE INDEX session_running_idx ON session (running) WHERE running;
COMMIT;