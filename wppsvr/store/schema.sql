-- Database schema for the packet-checkins application.

-- The message table stores all received messages.
CREATE TABLE message (
    id           text     PRIMARY KEY,
    hash         text     NOT NULL UNIQUE,
    deliverytime datetime NOT NULL,
    message      text     NOT NULL,
    session      integer  NOT NULL REFERENCES session,
    fromaddress  text     NOT NULL,
    fromcallsign text     NOT NULL,
    frombbs      text     NOT NULL,
    tobbs        text     NOT NULL,
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
    responseto text     NOT NULL REFERENCES message,
    sendto     text     NOT NULL,
    subject    text     NOT NULL,
    body       text     NOT NULL,
    sendtime   datetime NOT NULL,
    sendercall text     NOT NULL,
    senderbbs  text     NOT NULL
);

-- The retrieval table contains the time of the most recent successful retrieval
-- from each combination of call sign and BBS.
CREATE TABLE retrieval (
    callsign text     NOT NULL,
    bbs      text     NOT NULL,
    time     datetime NOT NULL,
    PRIMARY KEY (callsign, bbs)
) WITHOUT ROWID;

-- The session table describes all sessions.
CREATE TABLE session (
    id                     integer  PRIMARY KEY,
    callsign               text     NOT NULL,
    name                   text     NOT NULL,
    prefix                 text     NOT NULL,
    start                  datetime NOT NULL,
    end                    datetime NOT NULL,
    generateweeksummary    boolean  NOT NULL,
    excludefromweeksummary boolean  NOT NULL,
    reportto               text     NOT NULL,
    tobbses                text     NOT NULL,
    downbbses              text     NOT NULL,
    retrievefrombbses      text     NOT NULL,
    retrieveat             text     NOT NULL,
    messagetypes           text     NOT NULL,
    modified               boolean  NOT NULL,
    running                boolean  NOT NULL,
    report                 text     NOT NULL
);
CREATE UNIQUE INDEX session_call_end_idx ON session (callsign, end);
CREATE INDEX session_end_idx ON session (end);
CREATE INDEX session_running_idx ON session (running) WHERE running;
