package install

import (
    "fmt"
    "os/exec"
)

// InstallPackages installs a list of packages using apt-get
func InstallPackages(packages []string) error {
    for _, pkg := range packages {
        cmd := exec.Command("sudo", "apt", "install", "-y", pkg)
        output, err := cmd.CombinedOutput()
        if err != nil {
            return fmt.Errorf("failed to install package %s: %v\nOutput: %s", pkg, err, output)
        }
        fmt.Printf("Successfully installed package: %s\n", pkg)
    }
    return nil
}