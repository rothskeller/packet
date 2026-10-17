= PackItForms 4.0.17 Release Notes

PackItForms 4.0.17 introduces these changes:

- **KISS Connections and the Kenwood TH-D75:**  There is a new "KISS"
  connection method, for TNCs that support only KISS mode, such as the one
  built into the Kenwood TH-D75.  With a KISS TNC, PackItForms manages the
  AX.25 connection to the BBS itself.  Choose "KISS" as the Connection Method
  on the Incident Settings page, and then choose "Kenwood TH-D75" (or "Generic
  KISS TNC" for other KISS TNCs) as the TNC Type.  On the TH-D75, put the
  built-in TNC in KISS mode at 1200 baud (the display shows "KISS12"), set its
  KISS interface (Menu 983) to USB or Bluetooth to match how the radio is
  connected, and tune its data band (Menu 506) to the BBS frequency.
  PackItForms waits patiently for a slow BBS to answer before resending a
  packet, and adjusts that wait to how quickly the BBS actually answers.

- **Link Monitor:**  The new "Link Monitor" item on the Connect menu (shortcut
  M) opens a panel on the right side of the incident window, beside the message
  list, that shows the exchange with the BBS as it happens.  For KISS
  connections it shows each AX.25 packet sent and received, in the style of a
  packet monitor, with the text it carries, and also (dimmed, marked "other")
  the packets between other stations that the radio hears on the frequency,
  which the "Other traffic" checkbox hides; for other connections it shows the
  text exchanged with the BBS.  It keeps the recent history of each incident's
  connections, so it can be opened after a connection problem to see what
  happened.  Drag its left edge to resize it.

- **Leaving the TNC in Good Order:**  When a connection fails, is aborted, or
  is cut short because PackItForms is closing, PackItForms now puts the TNC
  back in order before letting go of it.  A command-mode TNC is brought back
  to command mode, disconnected, and has its settings restored, and a KISS
  connection sends the BBS a final disconnect.  The TH-D75's TNC is also
  reset, which discards any packets it was still holding.  If a session was
  cut off without that cleanup (because PackItForms was killed, say), the
  next connection does it first.

- **Serial Port Selection:**  Clicking the Serial Port field on the Incident
  Settings page now lists every serial port detected on the computer, even
  when the field already has a value.  Previously, a saved port hid the
  others.

- **Native Apple Silicon Build:**  The Mac version of PackItForms is now a
  single "universal" program that runs natively on both Apple Silicon and
  Intel Macs.  Previously it ran on Apple Silicon only through Rosetta.

Under the covers, PackItForms 4.0.17 is built with Go 1.27 (now the minimum
version needed to build it from source), and all of the libraries it uses have
been updated to their current versions.  Existing incidents and their settings
are unaffected.

= PackItForms 4.0 Release Notes

PackItForms 4.0 introduces these changes:

- **PDF Rendering:**  When viewing an existing message, or after submitting a
  new or edited message to Outpost, the message is now rendered in PDF format
  and displayed as a PDF.  This ensures exact compatibility with the paper
  forms, as well as proper pagination and footers.

- **Automatic Form Updates:**  Once each day, PackItForms will check with the
  county ARES/RACES website to see if any forms have changed.  If so, it will
  automatically download and install them.  System administrator privileges
  are not required for these forms updates.  Manual installs will still be
  needed when there are changes to the PackItForms software.

- **New Hospital Forms:**  All of the medical forms have been rewritten: Allied
  Health Status, Hospital Bed Availability (HAvBed), and Hospital Status
  (formerly Medical Facility Status).  Also, they are now included in the
  standard distribution.  There is no longer a separate installer for the
  Hospital Net.  All operators now get all county forms, including the hospital
  forms.

- **New ICS-213 Form:**  The ICS-213 form has been rewritten to match all of the
  other forms in styling and behavior.

- **Updates to Veoci Forms:**  Minor changes have been made to the Veoci forms
  that were rolled out in August.  In particular:
  - Field numbers have been added to most fields, and to the individual items
    in some checkbox and radio button fields.  This will make make voice
    transmission of the forms more efficient.
  - The blue "-- choose one --" label has been removed from the Jurisdiction
    drop-down fields on the PDFs, so that blank messages can be printed without
    that label.

- **Requires Current Platform:**  PackItForms 4.0 is not tested or supported on
  Windows versions prior to Windows 10; on 32-bit processors; or with
  non-current browsers or PDF readers.

Under the covers, PackItForms 4.0 has been completely rewritten.  It is now
written in Go and built on the same Go packet library as the weekly packet
practice server and the packet shell.  As a result, it is much smaller than its
predecessor (22MB vs. 207MB) and faster as well.

Source code and maintenance documentation for PackItForms 4.0 is at
github.com/rothskeller/packet, and is freely licensed (BSD).  The forms and
their definitions are at github.com/scc-ares-races/forms and are copyrighted by
SCCo ARES/RACES.

If you encounter any problems with PackItForms 4.0, please report them on the
packet@scc-ares-races.groups.io mailing list.  Please include the appropriate
dated log file from C:\PackItForms\Log.
