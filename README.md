# Packet Radio Software

This repository contains software related to packet radio:

* Package `cmd` contains some utility commands related to the other packages in
  this repository.
* Package `envelope` understands how to extract packet message content out of
  RFC-4155 and RFC-4322 email encodings, and how to save and restore them from
  local files.
* Package `incident` knows how to store a set of related messages for an
  incident in a directory, and generate ICS-309 logs for them.
* Package `jnos` is a library for communicating with JNOS BBS servers.  It
  includes transport adapters for RF (via a serial connection to a KPC 3 Plus
  TNC) and telnet.  It also includes a rudimentary JNOS BBS simulator for
  testing purposes.
* Package `message` is a registry for message types, allowing messages of
  various types to be decoded, created, and manipulated.  It also contains the
  base implementation that underlies all message types.
* Package `xscmsg` has subpackages with message type definitions for all of the
  public Santa Clara County standard message types.  The main `xscmsg` has a
  function to register them all.

## Legal Text

This software was written by Steve Roth, KC6RSC.

Copyright © 2021–2023 by Steven Roth <steve@rothskeller.net>

See LICENSE.txt for license details.
