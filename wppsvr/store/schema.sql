-- Database schema for the packet-checkins application.

-- The message table stores all received messages.
CREATE TABLE message (
    id           text     PRIMARY KEY,
    hash         text     NOT NULL UNIQUE,
    deliverytime datetime NOT NULL,
    message      text     NOT NULL,
    session      integer  NOT NULL REFERENCES session ON DELETE CASCADE,
    fromaddress  text     NOT NULL,
    fromcallsign text     NOT NULL,
    frombbs      text     NOT NULL,
    tobbs        text     NOT NULL,
    jurisdiction text     NOT NULL,
    messagetype  text     NOT NULL,
    subject      text     NOT NULL,
    problems     text     NOT NULL,
    actions      integer  NOT NULL
);
CREATE INDEX message_session_idx ON message (session);

-- The msgnum table keeps track of which local message numbers have been used
-- for each prefix.
CREATE TABLE msgnum (
    prefix text    PRIMARY KEY,
    num    integer NOT NULL
) WITHOUT ROWID;

-- The response table stores all outgoing responses to incoming messages.
CREATE TABLE response (
	id         text     PRIMARY KEY,
    responseto text     NOT NULL REFERENCES message ON DELETE CASCADE,
    sendto     text     NOT NULL,
    subject    text     NOT NULL,
    body       text     NOT NULL,
    sendtime   datetime NOT NULL,
    sendercall text     NOT NULL,
    senderbbs  text     NOT NULL
);

-- The retrieval table contains a row for each scheduled retrieval for each
-- session, describing the retrieval parameters and the last time that retrieval
-- was successfully completed.
CREATE TABLE retrieval (
    session           integer  REFERENCES session ON DELETE CASCADE,
    interval          text     NOT NULL,
    bbs               text     NOT NULL,
    mailbox           text     NOT NULL,
    dontkillmessages  boolean  NOT NULL,
    dontsendresponses boolean  NOT NULL,
    lastrun           datetime NOT NULL
);
CREATE INDEX retrieval_session_idx ON retrieval (session);

-- The session table describes all sessions.
CREATE TABLE session (
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
CREATE UNIQUE INDEX session_call_end_idx ON session (callsign, end);
CREATE INDEX session_end_idx ON session (end);
CREATE INDEX session_running_idx ON session (running) WHERE running;
