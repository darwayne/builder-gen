package generator

import (
	"flag"
	"strings"
)

func annotateFlagInfo(comment string, data *Data) {
	f := flag.NewFlagSet(magicString+" flags", flag.ContinueOnError)

	f.BoolVar(&data.Globals,
		"with-globals", false,
		"set this flag if you want to generate global functions as well")
	f.BoolVar(&data.UseOptionalBools,
		"with-optional-bools", false,
		"set this flag if you want bool opts to be optional")
	f.BoolVar(&data.NoBuilder, "no-builder", false,
		"set this flag if you want to exclude creating the builder object")
	f.StringVar(&data.Prefix, "prefix", "",
		"if set this will be the prefix of your global functions. Note: with-globals option required")
	f.StringVar(&data.Suffix, "suffix", "",
		"if set this will be the suffix of your global functions. Note: with-globals option required")
	str := strings.Split(comment, magicString+" ")

	if len(str) == 1 {
		return
	}

	var args []string
	for _, arg := range strings.Split(str[1], " ") {
		arg = strings.TrimSpace(arg)
		if arg == "" {
			continue
		}
		args = append(args, arg)
	}

	_ = f.Parse(args)
}
