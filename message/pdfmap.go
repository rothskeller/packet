package message

// PDFMapper is the interface honored by a Field.PDFMap value.
type PDFMapper interface {
	RenderPDF(*Field) []PDFField
}

// PDFField describes the rendering of a single PDF field, as returned by the
// PDFMap interface in a Field structure.
type PDFField struct {
	// Name is the name of the field in the PDF.
	Name string
	// Value is the value to be placed in the PDF field.
	Value string
	// Size is the font size to be used for the PDF field; it can be zero
	// if supplied by the PDF template or the Type.PDFSize value.
	Size float64
}

// NoPDFField is a zero implementation of PDFMapper that returns no mappings.
type NoPDFField struct{}

func (NoPDFField) RenderPDF(*Field) []PDFField { return nil }

// PDFName is a simple implementation of PDFMapper.  It maps the field value,
// unedited, to the specified PDF field name, with the default font size.
type PDFName string

func (n PDFName) RenderPDF(f *Field) []PDFField {
	value := *f.Value
	if f.Choices != nil {
		value = f.Choices.ToHuman(value)
	}
	return []PDFField{{Name: string(n), Value: value}}
}

// PDFNameMap is an implementation of PDFMapper.  The first element of the slice
// is the PDF field name.  The remaining elements are (PIFO value, PDF value)
// pairs.  If the current value of the field matches any of the PIFO values in
// the slice, it is mapped to the corresponding PDF value.  Otherwise it is
// rendered unchanged.  All values are rendered with the default font size.
type PDFNameMap []string

func (m PDFNameMap) RenderPDF(f *Field) []PDFField {
	value := *f.Value
	var mapped bool
	for i := 1; i < len(m)-1; i += 2 {
		if value == m[i] {
			value, mapped = m[i+1], true
			break
		}
	}
	if !mapped && f.Choices != nil {
		value = f.Choices.ToHuman(value)
	}
	return []PDFField{{Name: m[0], Value: value}}
}

// PDFMapFunc converts a PDF mapping function into an interface that satisfies
// PDFMapper.
type PDFMapFunc func(*Field) []PDFField

func (fn PDFMapFunc) RenderPDF(f *Field) []PDFField {
	return fn(f)
}
