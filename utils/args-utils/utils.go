package args_utils

import "golang.org/x/exp/slices"

type ArgsBuilder struct {
	args             []string
	addOrReplaceArgs map[string]string
}

func NewArgsBuilder(args []string) *ArgsBuilder {
	return &ArgsBuilder{
		args:             args,
		addOrReplaceArgs: map[string]string{},
	}
}

func (u *ArgsBuilder) ReplaceOrAddArgForFlag(flag string, newParam string) *ArgsBuilder {
	u.addOrReplaceArgs[flag] = newParam
	return u
}

func (u *ArgsBuilder) ReplaceOrAddArgForFlags(flag []string, newParam string) *ArgsBuilder {
	for i := range flag {
		u.ReplaceOrAddArgForFlag(flag[i], newParam)
	}
	return u
}

func (u *ArgsBuilder) Build() *Args {

	for s, s2 := range u.addOrReplaceArgs {
		var replaced bool
		for i := range u.args {
			if u.args[i] == s {
				u.args[i+1] = s2
				replaced = true
			}
		}
		if !replaced {
			u.args = append(u.args, s, s2)
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

	var splitIndex int

	if len(args) <= 1 {
		cttArgs = args
		return
	}

	for i, val := range args {
		if slices.Contains(availableCommands, val) {
			splitIndex = i
		}
	}

	if splitIndex == 0 {
		cttArgs = args
		return
	}

	if splitIndex != 0 {
		cttArgs = args[:splitIndex]
		commandArgs = args[splitIndex:]
	}

	return

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
