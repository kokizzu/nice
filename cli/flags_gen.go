// Code generated by generate_flags.py; DO NOT EDIT.

package cli

// bool

func BoolVar(register Register, p *bool, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return register.RegisterFlag(newFlag(newBoolValue(p), opts))
}

func Bool(register Register, name string, options ...FlagOptionApplyer) *bool {
	p := new(bool)
	_ = BoolVar(register, p, name, options...)
	return p
}

// string

func StringVar(register Register, p *string, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return register.RegisterFlag(newFlag(newStringValue(p), opts))
}

func String(register Register, name string, options ...FlagOptionApplyer) *string {
	p := new(string)
	_ = StringVar(register, p, name, options...)
	return p
}

// int

func IntVar(register Register, p *int, name string, options ...FlagOptionApplyer) error {
	var opts FlagOptions
	opts.applyName(name)
	opts.applyFlagOptions(options)

	return register.RegisterFlag(newFlag(newIntValue(p), opts))
}

func Int(register Register, name string, options ...FlagOptionApplyer) *int {
	p := new(int)
	_ = IntVar(register, p, name, options...)
	return p
}
