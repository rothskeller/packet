package xscmsg

type fieldContainer struct {
	fields []FormField
}

func (c fieldContainer) FieldCount() int { return len(c.fields) }
func (c fieldContainer) FieldByIndex(idx int) FormField {
	if idx < 0 || idx >= len(c.fields) {
		return nil
	}
	return c.fields[idx]
}
func (c fieldContainer) FieldByTag(tag string) FormField {
	for _, f := range c.fields {
		if f.Tag() == tag {
			return f
		}
	}
	return nil
}
func (c fieldContainer) FieldByKey(key FieldKey) FormField {
	for _, f := range c.fields {
		if f.Key() == key {
			return f
		}
	}
	return nil
}
func (c fieldContainer) FieldValue(tag string) Value {
	if f := c.FieldByTag(tag); f != nil {
		return f.Value()
	}
	return ""
}
func (c fieldContainer) KeyValue(key FieldKey) Value {
	if f := c.FieldByKey(key); f != nil {
		return f.Value()
	}
	return ""
}
func (c fieldContainer) Calculate() {
	for _, f := range c.fields {
		f.Calculate()
	}
}
func (c fieldContainer) Validate() (problems []string) {
	for _, f := range c.fields {
		if problem := f.Validate(); problem != "" {
			problems = append(problems, problem)
		}
	}
	return problems
}
func (c fieldContainer) ValidatePIFO() (problems []string) {
	for _, f := range c.fields {
		if problem := f.ValidatePIFO(); problem != "" {
			problems = append(problems, problem)
		}
	}
	return problems
}
