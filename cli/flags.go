package cli

type Flag struct {
	Value     Value
	Short     string
	Long      string
	Usage     Usager
	Necessary Necessary

	// NOTE(SuperPaintman):
	//     The first version had "Aliases" for flags. It's quite handy to have
	//     (e.g. --dry and --dry-run) but at the same time makes API a bit
	//     confusing because of duplication logic.
	//
	//     We can override flags on aliases collision or remove the alias from the
	//     original list but it makes API quite unpredictable for developers.
	//
	//     I decided to remove aliases. It's not so commonly used feature and
	//     developers can easely make a workaround if they need it.
}

func newFlag(value Value, opts FlagOptions) Flag {
	return Flag{
		Value:     value,
		Short:     opts.Short,
		Long:      opts.Long,
		Usage:     opts.Usage,
		Necessary: opts.Necessary,
	}
}

func (f *Flag) Type() string {
	if t, ok := f.Value.(Typer); ok {
		return t.Type()
	}

	return ""
}

func (f *Flag) Required() bool {
	return f.Necessary == Required
}

func (f *Flag) String() string {
	v := "Flag("
	v += f.Type()

	if f.Short != "" {
		v += ","
		v += "-" + f.Short
	}

	if f.Long != "" {
		if v == "" {
			v += ","
		} else {
			v += "/"
		}

		v += "--" + f.Long
	}
	v += ")"

	return v
}

//go:generate python ./generate_flags.py
