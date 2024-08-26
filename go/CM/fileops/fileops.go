package fileops

import (
    "fmt"
    "io"
    "os"
    "os/exec"
)

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("failed to open source file: %v", err)
    }
    defer sourceFile.Close()

    destFile, err := os.Create(dst)
    if err != nil {
        return fmt.Errorf("failed to create destination file: %v", err)
    }
    defer destFile.Close()

    _, err = io.Copy(destFile, sourceFile)
    if err != nil {
        return fmt.Errorf("failed to copy file from %s to %s: %v", src, dst, err)
    }

    fmt.Printf("Successfully copied file from %s to %s\n", src, dst)
    return nil
}

// CopyDirectory copies a directory from src to dst
func CopyDirectory(src, dst string) error {
    cmd := exec.Command("cp", "-r", src, dst)
    output, err := cmd.CombinedOutput()
    if err != nil {
        return fmt.Errorf("failed to copy directory from %s to %s: %v\nOutput: %s", src, dst, err, output)
    }
    fmt.Printf("Successfully copied directory from %s to %s\n", src, dst)
    return nil
}

// CreateDirectories creates multiple directories from a given list
func CreateDirectories(dirs []string) error {
    for _, dir := range dirs {
        err := os.MkdirAll(dir, 0755)
        if err != nil {
            return fmt.Errorf("failed to create directory %s: %v", dir, err)
        }
        fmt.Printf("Directory %s created successfully\n", dir)
    }
    return nil
    }
// runDebootstrap executes the debootstrap command
func RunDebootstrap() error {
    cmd := exec.Command("sudo", "debootstrap", "jammy", "/srv/nfs/jammy")

    // Attach stdout and stderr to the command's output streams
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // Start the command
    if err := cmd.Start(); err != nil {
        return fmt.Errorf("failed to start debootstrap: %v", err)
    }

    // Wait for the command to complete
    if err := cmd.Wait(); err != nil {
        return fmt.Errorf("debootstrap failed: %v", err)
    }

    fmt.Println("Debootstrap completed successfully!")
    return nil
}
