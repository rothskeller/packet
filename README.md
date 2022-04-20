# Packet Radio Software

This repository contains software related to packet radio:

* Package `pktmsg` is a packet message decoder and encoder.  It understands
  RFC-4155 and RFC-5322 email encoding, SCCo-standard subject line encoding,
  PackItForms and similar form encodings, and Outpost-specific feature
  encodings.
* Package `jnos` is a library for communicating with JNOS BBS servers.  It
  includes transport adapters for RF (via a serial connection to a KPC 3 Plus
  TNC) and telnet.  It also includes a rudimentary JNOS BBS simulator for
  testing purposes.
* Package `wppsvr` is the program that receives, responds to, and reports on
  SCCo weekly packet practice messages.

## Legal Text

This software was written by Steve Roth, KC6RSC.

Copyright © 2021–2022 by Steven Roth <steve@rothskeller.net>

Permission to use, copy, modify, and/or distribute this software for any purpose
with or without fee is hereby granted.

DISCLAIMER: THIS SOFTWARE IS WITHOUT WARRANTY.
