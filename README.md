# Packet Radio Software

This repository contains software related to packet radio:

* Package `cmd` contains some utility commands related to the other packages in
  this repository.
* Package `envelope` understands how to extract packet message content out of
  RFC-4155 and RFC-4322 email encodings, and how to save and restore them from
  local files.
* Package `jnos` is a library for communicating with JNOS BBS servers.  It
  includes transport adapters for RF (via a serial connection to a KPC 3 Plus
  TNC) and telnet.  It also includes a rudimentary JNOS BBS simulator for
  testing purposes.
* Package `message` is a registry for message types, allowing messages of
  various types to be decoded, created, and manipulated.  It has subpackages
  with message type definitions for all of the public Santa Clara County
  standard message types.
* Package `wppsvr` is the program that receives, responds to, and reports on
  SCCo weekly packet practice messages.

## Legal Text

This software was written by Steve Roth, KC6RSC.

Copyright © 2021–2023 by Steven Roth <steve@rothskeller.net>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted.

DISCLAIMER: THIS SOFTWARE IS WITHOUT WARRANTY.
