package main

import (
	"os"
	"strings"

	vdfparser "github.com/Jleagle/steam-go/steamvdf"
	vdfgenerator "github.com/mdouchement/vdf"
)

type ErrVdfNotFound []string

func (err ErrVdfNotFound) Error() string {
	return "key: " + strings.Join(err, " - ") + " not found"
}

func applyLaunchOptions(value string, path string, overrite bool) error {
	kv, err := vdfparser.ReadFile(path)
	if err != nil {
		return err
	}

	vdf := kv.ToMapOuter()
	apps, ok := vdf["UserLocalConfigStore"].(map[string]any)["Software"].(map[string]any)["Valve"].(map[string]any)["Steam"].(map[string]any)["apps"].(map[string]any)
	keypath := []string{"Software", "Valve", "Steam", "apps"}
	if !ok {
		return ErrVdfNotFound(keypath)
	}

	for _, child := range apps {
		childMap := child.(map[string]any)
		_, ok := childMap["LaunchOptions"]
		if overrite || !ok {
			childMap["LaunchOptions"] = value
		}
	}

	err = os.Rename(path, path+".bak")
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return vdfgenerator.GenerateIO(f, vdf)
}
