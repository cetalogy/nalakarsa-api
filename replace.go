//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	err := filepath.WalkDir("internal/service", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".go") {
			return nil
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		original := string(content)
		modified := strings.ReplaceAll(original, ".Profile.", ".")
		
		// Some edge cases might just be `.Profile` at the end
		// but wait, we only want to replace access to profile fields.
		// `.Profile.` covers `m.User.Profile.FullName` -> `m.User.FullName`
		// Wait, what about `n.Actor.Profile` where it is passed as a value?
		// Usually it's `.Profile.`
		
		// Let's also check for `, Profile:` or `.Profile` assignment
		
		if modified != original {
			err = os.WriteFile(path, []byte(modified), 0644)
			if err != nil {
				return err
			}
			fmt.Printf("Updated %s\n", path)
		}
		return nil
	})
	
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
