package main

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/shirou/gopsutil/v4/process"
)

var ErrProcessNotFound = fmt.Errorf("process not found")

func killProcess(name string) ([]string, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range processes {
		n, err := p.Name()
		if err != nil {
			return nil, err
		}
		if n == name {
			args, err := p.CmdlineSlice()
			if err != nil {
				return nil, err
			}

			return args, p.Kill()
		}
	}
	return nil, ErrProcessNotFound
}

func ptr[T any](value T) *T {
	return &value
}

func main() {
	defaultUserId := getDefaultUserId()
	_ = kong.Parse(&CLI,
		kong.Name("SteamGlobalLaunchOptions"),
		kong.Description("A CLI tool to apply launch options for all Steam games at once."),
		kong.UsageOnError(),
		kong.ConfigureHelp(
			kong.HelpOptions{
				Compact: true,
				Summary: true,
			},
		),
		kong.Vars{"default_user_id": defaultUserId},
	)
	if CLI.UserId == "" {
		log.Fatalln("Default user id have not been found. You must provide user id.")
	}

	if CLI.RestoreSteam == nil {
		CLI.RestoreSteam = ptr(true)
	}
	args, err := killProcess("steam")
	if errors.Is(err, ErrProcessNotFound) {
		if CLI.RestoreSteam == nil {
			*CLI.RestoreSteam = false
		}
	} else if err != nil {
		log.Fatalln(err)
	}

	err = applyLaunchOptions(
		CLI.Value,
		filepath.Join(parsePath(SteamUserdata), CLI.UserId, ConfigFilePath),
		CLI.Overrite,
	)
	if err != nil {
		log.Fatalln(err)
	}

	if *CLI.RestoreSteam {
		fmt.Println(args)
		if len(args) <= 1 {
			_ = exec.Command("steam").Start()
		} else {
			_ = exec.Command("steam", args[1:]...).Start()
		}
	}
}
