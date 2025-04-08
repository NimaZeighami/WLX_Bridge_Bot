package configs

import (
    "fmt"
    "gopkg.in/yaml.v3"
    "os"
)

// LoadDynamicConfig reads and parses the YAML file into a nested map of strings.
func LoadDynamicConfig(filePath string) (map[string]map[string]map[string]string, error) {
    file, err := os.Open(filePath)
    if err != nil {
        return nil, fmt.Errorf("error opening config file: %v", err)
    }
    defer file.Close()

    var config map[string]map[string]map[string]string
    decoder := yaml.NewDecoder(file)
    if err := decoder.Decode(&config); err != nil {
        return nil, fmt.Errorf("error decoding YAML file: %v", err)
    }

    return config, nil
}
