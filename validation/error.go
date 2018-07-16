package validation

// Error for valid.
type Error struct {
	Msg string
	Field string
}

// Return the Error's Msg.
func (e *Error) String() string {
	if e == nil {
		return ""
	}
	return e.Msg
}

// Implement Error interface.
func (e *Error) Error() string { return e.String() }

// error message template.
type MsgTmpl map[string]string

// error message template map.
var MsgTmplMap  = map[string]MsgTmpl{
	"en": MsgTmplEn,
}

// English error message template.
var MsgTmplEn = MsgTmpl{
	"Required":     "Can not be empty",
	"Min":          "Minimum is %d",
	"Max":          "Maximum is %d",
	"Range":        "Range is %d to %d",
	"MinSize":      "Minimum size is %d",
	"MaxSize":      "Maximum size is %d",
	"Length":       "Required length is %d",
	"Alpha":        "Must be valid alpha characters",
	"Numeric":      "Must be valid numeric characters",
	"AlphaNumeric": "Must be valid alpha or numeric characters",
	"Match":        "Must match %s",
	"NoMatch":      "Must not match %s",
	"AlphaDash":    "Must be valid alpha or numeric or dash(-_) characters",
	"Email":        "Must be a valid email address",
	"IP":           "Must be a valid ip address",
	"Mobile":       "Must be valid mobile number",
	"Tel":          "Must be valid telephone number",
	"Phone":        "Must be valid telephone or mobile phone number",
}

// Set error message template.
func (c MsgTmpl) SetDefaultMessage(msg map[string]string) {
	if len(msg) == 0 {
		return
	}

	for name, tmpl := range msg {
		c[name] = tmpl
	}
}