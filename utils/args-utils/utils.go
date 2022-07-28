package args_utils

type ArgsBuilder struct {
	args             []string
	addOrReplaceArgs []map[string]string
}

func NewArgsBuilder(args []string) *ArgsBuilder {
	return &ArgsBuilder{
		args:             args,
		addOrReplaceArgs: []map[string]string{},
	}
}

func (u *ArgsBuilder) ReplaceOrAddArgForFlags(name string, flags []string, newParam string) *ArgsBuilder {
	flagsToReplaceGroup := map[string]string{}
	for _, flag := range flags {
		flagsToReplaceGroup[flag] = newParam
	}
	u.addOrReplaceArgs = append(u.addOrReplaceArgs, flagsToReplaceGroup)
	return u
}

func (u *ArgsBuilder) Build() *Args {

	for _, flagGroup := range u.addOrReplaceArgs {
		var replaced bool
		for s, s2 := range flagGroup {
			for i := range u.args {
				if u.args[i] == s {
					u.args[i+1] = s2
					replaced = true
				}
			}

		}
		if !replaced {
			for s, s2 := range flagGroup {
				u.args = append(u.args, s, s2)
				break
			}
		}
	}

	return NewArgs(u.args)

}

type Args struct {
	args []string
}

func NewArgs(args []string) *Args {
	return &Args{args: args}
}

func (u *Args) GetValueFromFlag(flag string) string {

	for i, _ := range u.args {
		if u.args[i] == flag {
			return u.args[i+1]
		}
	}
	return ""
}

func (u *Args) GetValueFromFlags(flags []string) string {

	for i := range flags {
		if val := u.GetValueFromFlag(flags[i]); val != "" {
			return val
		}
	}

	return ""
}

func (u *Args) GetArgs() []string {

	return u.args
}

func SplitArgs(args []string, availableCommands []string) (cttArgs []string, commandArgs []string) {

	if len(args) <= 1 {
		return args, nil
	}

	for i, val := range args {
		for _, cmd := range availableCommands {
			if cmd == val {
				return args[:i], args[i:]
			}
		}
	}

	return args, nil
}

func SplitToCommand(args []string) (cmd string, commandArgs []string) {

	switch len(args) {
	case 0:
	case 1:
		cmd = args[0]
	default:
		cmd = args[0]
		commandArgs = args[1:]
	}
	return

}
