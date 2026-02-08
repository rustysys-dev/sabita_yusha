package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type SystemdStatus string

const (
	SystemdStatusRunning  SystemdStatus = "Running"
	SystemdStatusFailed   SystemdStatus = "Failed"
	SystemdStatusNotFound SystemdStatus = "NotFound"
	SystemdStatusSkipped  SystemdStatus = "Skipped"
	SystemdStatusUpdated  SystemdStatus = "Updated"
)

const (
	systemdSystemDir = "/etc/systemd/system"
	installBinPath   = "/usr/local/bin/sabita_yusha"
	serviceFileName  = "sabita_yusha.service"
	serviceTemplate  = `[Unit]
Description=Sabita Yusha
After=multi-user.target

[Service]
ExecStart=%s -systemd
Restart=on-failure
User=%s
SupplementaryGroups=input

[Install]
WantedBy=multi-user.target
`
)

// EnsureSystemdService checks if the app is registered as a system-level
// systemd service and prompts the user to register it if not. Registration
// requires sudo privileges. All errors are non-fatal; they are logged as
// warnings and the application continues.
func EnsureSystemdService(skipped, update bool) (SystemdStatus, error) {
	if !hasSystemd() {
		log.Println("systemd not detected, skipping service registration check")
		return SystemdStatusNotFound, nil
	}

	if skipped {
		return SystemdStatusSkipped, nil
	}

	servicePath := filepath.Join(systemdSystemDir, serviceFileName)

	if !isServiceInstalled(servicePath) {
		if !promptUserForRegistration() {
			log.Println("skipping systemd service registration")
			return SystemdStatusSkipped, nil
		}

		if err := installBinary(); err != nil {
			log.Println("WARNING: failed to install binary:", err)
			return SystemdStatusFailed, err
		}

		content, err := generateServiceFile()
		if err != nil {
			log.Println("WARNING: failed to generate service file:", err)
			return SystemdStatusFailed, err
		}

		if err := registerService(servicePath, content); err != nil {
			log.Println("WARNING: failed to register systemd service:", err)
			return SystemdStatusFailed, err
		}

		if err := startService(serviceFileName); err != nil {
			log.Println("WARNING: failed to start systemd service:", err)
			return SystemdStatusFailed, err
		}

		log.Println("systemd service registered and enabled successfully")
		return SystemdStatusRunning, nil
	}

	if update {
		if err := stopService(serviceFileName); err != nil {
			log.Println("WARNING: failed to stop systemd service:", err)
			return SystemdStatusFailed, err
		}

		if err := installBinary(); err != nil {
			log.Println("WARNING: failed to install binary:", err)
			return SystemdStatusFailed, err
		}

		if err := startService(serviceFileName); err != nil {
			log.Println("WARNING: failed to start systemd service:", err)
			return SystemdStatusFailed, err
		}

		log.Println("systemd service updated and restarted successfully")
		return SystemdStatusUpdated, nil
	}

	log.Println("systemd service already registered")
	return SystemdStatusRunning, nil
}

func hasSystemd() bool {
	_, err := exec.LookPath("systemctl")
	return err == nil
}

func isServiceInstalled(servicePath string) bool {
	_, err := os.Stat(servicePath)
	return err == nil
}

func promptUserForRegistration() bool {
	fmt.Println("sabita_yusha is not registered as a systemd service.")
	fmt.Println("this action will copy this executable to", installBinPath)
	fmt.Println("it will then register and enable it as a systemd service.")
	fmt.Print("(requires sudo) [y/N]: ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Println("WARNING: unable to read user input:", err)
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

func generateServiceFile() ([]byte, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to determine current user: %w", err)
	}

	fileContentsStr := fmt.Sprintf(serviceTemplate, installBinPath, currentUser.Username)
	return []byte(fileContentsStr), nil
}

func installBinary() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to determine executable path: %w", err)
	}

	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("unable to resolve executable symlinks: %w", err)
	}

	if out, err := exec.Command("sudo", "cp", execPath, installBinPath).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to copy binary to %s: %w, output: %s", installBinPath, err, string(out))
	}

	return nil
}

func registerService(servicePath string, content []byte) error {
	cmd := exec.Command("sudo", "tee", servicePath)
	cmd.Stdin = strings.NewReader(string(content))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to write service file %s: %w, output: %s", servicePath, err, string(out))
	}

	if out, err := exec.Command("sudo", "systemctl", "daemon-reload").CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl daemon-reload failed: %w, output: %s", err, string(out))
	}

	if out, err := exec.Command("sudo", "systemctl", "enable", serviceFileName).CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl enable failed: %w, output: %s", err, string(out))
	}

	return nil
}

func stopService(service string) error {
	if out, err := exec.Command("sudo", "systemctl", "stop", service).CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl stop failed: %w, output: %s", err, string(out))
	}
	return nil
}

func startService(service string) error {
	if out, err := exec.Command("sudo", "systemctl", "start", service).CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl start failed: %w, output: %s", err, string(out))
	}
	return nil
}
