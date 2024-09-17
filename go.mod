module github.com/rothskeller/packet

go 1.21

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/phpdave11/gofpdf v1.4.2
	github.com/phpdave11/gofpdi v1.0.13
	github.com/rothskeller/pdf v1.3.0
	go.bug.st/serial v1.6.0
	golang.org/x/exp v0.0.0-20240808152545-0cdaa3abc0fa
	golang.org/x/net v0.12.0
)

require github.com/pkg/errors v0.9.1 // indirect

require (
	github.com/creack/goselect v0.1.2 // indirect
	github.com/go-pdf/fpdf v0.9.0
	golang.org/x/sys v0.10.0 // indirect
)

replace github.com/rothskeller/pdf => ../../pdf
