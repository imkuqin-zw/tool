package validation

func NewValidation(conf ...string) *Validation {
	v := &Validation{Lang: DefMsgLang}
	if len(conf) > 0 {
		v.Lang = conf[0]
	}
	return v
}
