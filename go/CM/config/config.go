package config

import (
    "fmt"
    "os"
)

// WriteConfig writes the configuration content to the specified file
func WriteConfig(filePath, content string) error {
    confFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
        return fmt.Errorf("Failed to open file: %v", err)
    }
    defer confFile.Close()

    _, err = confFile.WriteString(content)
    if err != nil {
        return fmt.Errorf("Failed to write to file: %v", err)
    }

    fmt.Printf("Configuration written to %s\n", filePath)
    return nil
}
