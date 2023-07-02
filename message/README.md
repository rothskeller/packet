# How to Update PackItForms

The canonical source of information about the various form types is the source
code for PackItForms itself — particularly the HTML form files.  For this
reason, copies of those files are maintained is this repository for easy
reference during development.  When a new version of PackItForms comes, the new
HTML form files should be copied into this repository, using the instructions
below.

## Step 1:  Update PackItForms Source Code

If you haven't checked out the PackItForms source code before:

    > cd ..       # (i.e., the directory above "packet")
    > git clone https://github.com/jmkristian/pack-it-forms
    > cd packet   # (i.e., return to the "packet" directory)

If you have it already checked out and you need to update it:

    > cd ../pack-it-forms
    > git pull
    > cd ../packet

## Step 2:  Update the HTML Files in the Packet Tree

The packet source code tree contains copies of the PackItForms form HTML files
from various versions of PackItForms.  Some of these may have changed and need
to be updated.  Also, there may be new versions of some forms that need to be
copied.  To do these things:

    > go run ./cmd/get-pifo-html

(Note:  get-pifo-html assumes that the PackItForms source code is at
`../pack-it-forms` and that the output files should go in form-tagged
subdirectories of `./message`.  Both of these assumptions are correct if you're
following these instructions.  But if you need to do something else, you can
override those paths on the get-pifo-html command line.)

If the new PackItForms version has a new type of form in it, this command will
fail with a message like:

    Tag vSCCo.99 contains an unknown form form-doughnut-order.html.

If that happens, you'll need to edit `./cmd/get-pifo-html/main.go` and add that
form to the table at the top of the file.  If it's a form we don't care about,
it can be mapped to an empty string.  Otherwise, it needs to be mapped to an
identifier resembling the form name.  The identifier must start with a lowercase
letter and contain only lowercase letters and digits.  Once that change is made,
run get-pifo-html again.

The HTML files are stored in `./message/«formtype»/*.html`.  Note that any
"include" directives in the original HTML file have been followed; the resulting
HTML files are self-contained.

By default, get-pifo-html only gets HTML files from PackItForms versions that
are new since the last time it was run.  If for some reason you need to
regenerate them all, remove the `./message/tags-read` file, remove all existing
`*.html` files, and re-run get-pifo-html.

## Step 3:  Update the Form Definitions

Once the new HTML form files have been retrieved, compare them with the previous
versions to determine what has changed, and update the form-specific source code
to match.  The details will vary case-by-case.

## Step 4:  Commit Changes

After appropriate testing, commit all changes to the repository.
