// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/profiles/seccomp"
	"github.com/docker/engine-api/types"
	"github.com/opencontainers/runc/libcontainer/configs"
)

var operators = map[configs.Operator]types.Operator{
	configs.NotEqualTo:           types.OpNotEqual,
	configs.LessThan:             types.OpLessThan,
	configs.LessThanOrEqualTo:    types.OpLessEqual,
	configs.EqualTo:              types.OpEqualTo,
	configs.GreaterThanOrEqualTo: types.OpGreaterEqual,
	configs.GreaterThan:          types.OpGreaterThan,
	configs.MaskEqualTo:          types.OpMaskedEqual,
}

var actions = map[configs.Action]types.Action{
	configs.Kill:  types.ActKill,
	configs.Errno: types.ActErrno,
	configs.Trap:  types.ActTrap,
	configs.Allow: types.ActAllow,
	configs.Trace: types.ActTrace,
}

var archs = map[string]types.Arch{
	"x86":         types.ArchX86,
	"amd64":       types.ArchX86_64,
	"x32":         types.ArchX32,
	"arm":         types.ArchARM,
	"arm64":       types.ArchAARCH64,
	"mips":        types.ArchMIPS,
	"mips64":      types.ArchMIPS64,
	"mips64n32":   types.ArchMIPS64N32,
	"mipsel":      types.ArchMIPSEL,
	"mipsel64":    types.ArchMIPSEL64,
	"mipsel64n32": types.ArchMIPSEL64N32,
}

// convertOperator converts a Seccomp comparison operator.
func convertOperator(in configs.Operator) (types.Operator, error) {
	if op, ok := operators[in]; ok {
		return op, nil
	}
	return "", fmt.Errorf("operator %s is not a valid operator for seccomp", in)
}

// convertAction converts a Seccomp rule match action.
func convertAction(in configs.Action) (types.Action, error) {
	if act, ok := actions[in]; ok {
		return act, nil
	}
	return "", fmt.Errorf("action %s is not a valid action for seccomp", in)
}

// convertArch converts a Seccomp comparison arch.
func convertArch(in string) (types.Arch, error) {
	if arch, ok := archs[in]; ok {
		return arch, nil
	}
	return "", fmt.Errorf("arch %s is not a valid arch for seccomp", in)
}

func translateSeccomp(config *configs.Seccomp) (newConfig *types.Seccomp, err error) {
	if config == nil {
		return nil, nil
	}

	if config.DefaultAction == 0 && len(config.Syscalls) == 0 {
		return nil, nil
	}

	newConfig = new(types.Seccomp)
	newConfig.Syscalls = []*types.Syscall{}

	// if config.Architectures == 0 then libseccomp will figure out the architecture to use
	if len(config.Architectures) > 0 {
		newConfig.Architectures = []types.Arch{}
		for _, arch := range config.Architectures {
			newArch, err := convertArch(arch)
			if err != nil {
				return nil, err
			}
			newConfig.Architectures = append(newConfig.Architectures, newArch)
		}
	}

	// Convert default action
	newConfig.DefaultAction, err = convertAction(config.DefaultAction)
	if err != nil {
		return nil, err
	}

	// Loop through all syscall blocks and convert them to docker-engine format
	for _, call := range config.Syscalls {
		newAction, err := convertAction(call.Action)
		if err != nil {
			return nil, err
		}

		newCall := types.Syscall{
			Name:   call.Name,
			Action: newAction,
			Args:   []*types.Arg{},
		}

		// Loop through all the arguments of the syscall and convert them
		for _, arg := range call.Args {
			newOp, err := convertOperator(arg.Op)
			if err != nil {
				return nil, err
			}

			newArg := types.Arg{
				Index:    arg.Index,
				Value:    arg.Value,
				ValueTwo: arg.ValueTwo,
				Op:       newOp,
			}

			newCall.Args = append(newCall.Args, &newArg)
		}

		newConfig.Syscalls = append(newConfig.Syscalls, &newCall)
	}

	return newConfig, nil
}

// saves the default seccomp profile as a json file so people can use it as a
// base for their own custom profiles
func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f := filepath.Join(wd, "default.json")

	// get the default profile
	p := seccomp.GetDefaultProfile()
	s, err := translateSeccomp(p)
	if err != nil {
		panic(err)
	}

	// write the default profile to the file
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(f, b, 0644); err != nil {
		panic(err)
	}
}
