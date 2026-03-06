package bootstrap

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const defaultAppEnv = "development"

func LoadEnvByEnvironment() (string, []string, error) {
	appEnv := strings.TrimSpace(os.Getenv("APP_ENV"))
	if appEnv == "" {
		appEnv = defaultAppEnv
	}

	rootDir, err := findProjectRoot()
	if err != nil {
		// In production containers we usually run only the compiled binary
		// (without repository files such as go.mod/.env). In this case,
		// keep running with env vars already provided by the platform.
		return appEnv, nil, nil
	}

	alreadyDefined := map[string]bool{}
	for _, entry := range os.Environ() {
		key := strings.SplitN(entry, "=", 2)[0]
		alreadyDefined[key] = true
	}

	candidates := []string{
		filepath.Join(rootDir, ".env"),
		filepath.Join(rootDir, fmt.Sprintf(".env.%s", appEnv)),
	}

	loaded := make([]string, 0, len(candidates))
	for _, candidate := range candidates {
		ok, loadErr := loadEnvFile(candidate, alreadyDefined)
		if loadErr != nil {
			return "", nil, loadErr
		}
		if ok {
			loaded = append(loaded, candidate)
		}
	}

	return appEnv, loaded, nil
}

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("go.mod not found from working dir %s", wd)
		}
		dir = parent
	}
}

func loadEnvFile(path string, alreadyDefined map[string]bool) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		key, value, ok, parseErr := parseEnvLine(scanner.Text())
		if parseErr != nil {
			return false, fmt.Errorf("%s:%d: %w", path, lineNo, parseErr)
		}
		if !ok {
			continue
		}
		if alreadyDefined[key] {
			continue
		}
		if err := os.Setenv(key, value); err != nil {
			return false, err
		}
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func parseEnvLine(raw string) (key, value string, ok bool, err error) {
	line := strings.TrimSpace(raw)
	if line == "" || strings.HasPrefix(line, "#") {
		return "", "", false, nil
	}

	line = strings.TrimPrefix(line, "\ufeff")
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "export ") {
		line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
	}

	parts := strings.SplitN(line, "=", 2)
	if len(parts) != 2 {
		return "", "", false, fmt.Errorf("invalid env declaration")
	}

	key = strings.TrimSpace(parts[0])
	if key == "" {
		return "", "", false, fmt.Errorf("empty env key")
	}

	value = strings.TrimSpace(parts[1])
	if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) >= 2 {
		value = value[1 : len(value)-1]
	} else if strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'") && len(value) >= 2 {
		value = value[1 : len(value)-1]
	}

	return key, value, true, nil
}
