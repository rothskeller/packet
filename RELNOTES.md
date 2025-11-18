= PackItForms 4.0 Release Notes

PackItForms 4.0 introduces these changes:

- **PDF Rendering:** When viewing an existing message, or after submitting a
  new or edited message to Outpost, the message is now rendered in PDF format
  and displayed as a PDF. This ensures exact compatibility with the paper
  forms, as well as proper pagination and footers. For messages whose forms
  don't contain routing fields (primarily, the Allied Health Facility Status
  form), the PDF rendering will start with an SCCo-standard Radio Routing Slip.

- **Automatic Form Updates:** Once each day, PackItForms will check with the
  county ARES/RACES website to see if any forms have changed. If so, it will
  automatically download and install them. System administrator privileges
  are not required for these forms updates. Manual installs will still be
  needed when there are changes to the PackItForms software.

- **Merged Hospital Forms:** There is no longer a separate installer for the
  Hospital Net. All operators now get all county forms, including the hospital
  forms.

- **New Hospital Bed Availability Form:** The hospital bed availability form
  (HAvBed) has been substantially changed. PackItForms can still receive and
  print the old form, but will only create messages with the new one.

- **Updates to Veoci Forms:** Minor changes have been made to the Veoci forms
  that were rolled out in August. In particular:
  - Field numbers have been added to most fields, and to the individual items
    in some checkbox and radio button fields. This will make make voice
    transmission of the forms more efficient.
  - The blue "-- choose one --" label has been removed from the Jurisdiction
    drop-down fields on the PDFs, so that blank messages can be printed without
    that label.

- **Requires Current Platform:** PackItForms 4.0 is not tested or supported on
  Windows versions prior to Windows 10; on 32-bit processors; or with
  non-current browsers or PDF readers.

Under the covers, PackItForms 4.0 has been completely rewritten. It is now
written in Go and built on the same Go packet library as the weekly packet
practice server and the unofficial alternative packet shell. As a result, it
is much smaller than its predecessor (22MB vs. 207MB) and faster as well.

Source code and maintenance documentation for PackItForms 4.0 is at
github.com/rothskeller/packet, and is freely licensed (BSD). The forms and
their definitions are at github.com/scc-ares-races/forms and are copyrighted by
SCCo ARES/RACES.

If you encounter any problems with PackItForms 4.0, please report them on the
packet@scc-ares-races.groups.io mailing list. Please include the appropriate
dated log file from C:\PackItForms\Log.
