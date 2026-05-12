# Weekly Packet Practice Message Analysis

Analysis of a weekly packet practice message happens in two phases.  First,
Analyze() is called on the message, which produces an Analysis structure with
information about all of the message metadata and problems.  This gets returned
to the caller, which will use it to send response messages.  After the response
messages are successfully sent, the caller will call the Commit method of the
Analysis structure to write the analysis to the database.  As a result, if we
are unable to send the response messages, we'll come back during the next wppsvr
pass, re-analyze the message, and retry sending the responses.

Message analysis starts by parsing the message and determining whether it should
be analyzed at all.  We don't analyze unparseable messages, bounce messages,
delivery and read receipts, or messages that have already been analyzed.  Then
the code runs a set of analysis checks.

The analysis checks are modular.  Each check is identified by a problem code,
which is a single word in PascalCase.  After analysis is complete, the database
contains the list of problem codes of the problems that were found in the
message.

Each analysis check is implemented with a Problem structure, which contains the
following fields:

* `Code` is a problem code that identifies the problem.  It is a single word in
  PascalCase. After analysis is complete, the database contains the list of
  problem codes of the problems that were found in the message.
* `ifnot` is a list of other Problems that preclude checking for this problem.
  For a given problem `p`, all checks in `p.ifnot` are run before `p` is, and if
  any of them find a problem, `p` is not run at all.
* `detect` is the function that detects whether the message has the problem.
* `Variables` is a map.  The keys are the names of variables that might be used
  in the response text for the problem.  The values are functions that return
  the values of those variables for a given analysis.

The set of known analysis checks is stored in the `Problems` variable, which is
a map from problem code to Problem structure.  New analysis checks are
registered by adding entries for them to this map in a `func init()`.
(Corresponding entries must be added to the `problems` map in `config.yaml` as
well, so that wppsvr knows what to do when a problem of this type is found.)

The message analysis code runs all of the registered analysis checks, in an
order that satisfies the `ifnot` constraints.  (The order is otherwise random.)

## Response Text

When a problem is found with a message, we often want to notify the message
sender of the problem.  The wppsvr code finds the entry for the problem in the
`problems` map in `config.yaml`.  If that entry has a `response` key, we will
send a problem response to the message sender, containing a description of the
problem.  (If the message has multiple problems with `response` keys, they get a
single problem response containing the descriptions of all such problems.)

The `response` key maps to a string containing is the description of the
problem.  That string may include `{VARIABLE}` references, which are places
where variables should be interpolated into the description.  The variable names
must either be defined in the `Variables` map in the Problem structure for the
problem, or they must be one of the well-known global variables:

* `AMSGTYPE`: name of the message type, preceded by the "a" or "an"
* `FROMBBS`: name of the BBS from which the message was sent
* `FROMCALLSIGN`: call sign of the sender of the message
* `MSGDATE`: date of the message, in the form "2006-01-02 at 15:04"
* `MSGTYPE`: name of the message type
* `SESSIONBBSES`: list of names of BBSes to which messages should be sent
* `SESSIONDATE`: date of the end of the practice session, e.g. "January 2"
* `SESSIONNAME`: name of the practice session, e.g. "SVECS Net"
* `TOBBS`: name of the BBS to which the message was sent
* `TOCALLSIGN`: name (call sign) of the mailbox to which the message was sent

## How to Add a New Analysis Check

To add a new analysis check, follow these steps.

1. Determine a code for the new problem type.  It must be a single word, and
   for consistency, should be in PascalCase.  Make sure the code you've chosen
   isn't already in use.
2. Add an entry for the new problem type to the `problems` map in
   `testdata/config.yaml`.  See the comments in `config.yaml` for details.  In
   particular, write the response text so that you know what variables you want
   to interpolate.
3. Determine which other problem checks should preclude yours, i.e., if that
   other problem is found, yours shouldn't be checked.
4. Create a new Go source file in this directory with a name relevant to your
   new check.  (You could add code to an existing source file instead, if it
   fits topically.)
5. In that source file, define a `ProbCodeName` variable (where CodeName is the
   code you chose in step 1), as a pointer to a Problem structure.  Fill in the
   `Code` and `detect` fields of the structure.  If you identified any
   precluding checks in step 3, list them in the `ifnot` field of the structure.
   If you used any variables in your response text in step 2, other than the
   well-known global variables listed above, define those variables in the
   `Variables` field of the structure.
6. Create a `func init()` in your source file (or add to the existing one if
   there already is one).  In it, add an entry to the global `Problems` map,
   mapping your code name to your Problem structure.
7. Add one or more `codename.yaml` files in appropriate place(s) in the
   `testdata` tree, containing test cases for the new analysis check.  See the
   existing test cases for examples.
8. Run the package tests and make sure they pass.

Once the new analysis check is added and tested, commit the code and put it into
production.  **At the same time,** copy the entry for the new problem type from
the `testdata/config.yaml` file to the production `config.yaml` file.
