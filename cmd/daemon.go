package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var daemonInterval int

func daemonDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".ha-cli")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func pidFile() (string, error) {
	dir, err := daemonDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "daemon.pid"), nil
}

func entitiesFile() (string, error) {
	dir, err := daemonDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "entities.json"), nil
}

func logFile() (string, error) {
	dir, err := daemonDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "daemon.log"), nil
}

func writePID(pid int) error {
	pf, err := pidFile()
	if err != nil {
		return err
	}
	return os.WriteFile(pf, []byte(strconv.Itoa(pid)), 0644)
}

func readPID() (int, error) {
	pf, err := pidFile()
	if err != nil {
		return 0, err
	}
	data, err := os.ReadFile(pf)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}

func isDaemonAlive() bool {
	pid, err := readPID()
	if err != nil {
		return false
	}
	return isProcessAlive(pid)
}

func runForeground(interval int) {
	lf, err := logFile()
	if err != nil {
		fmt.Fprintln(os.Stderr, "daemon: cannot get log path:", err)
		os.Exit(1)
	}

	f, err := os.OpenFile(lf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "daemon: cannot open log:", err)
		os.Exit(1)
	}
	defer f.Close()

	logger := log.New(f, "", log.LstdFlags)

	if err := writePID(os.Getpid()); err != nil {
		logger.Println("ERROR write PID:", err)
		os.Exit(1)
	}

	pf, _ := pidFile()
	defer os.Remove(pf)

	logger.Printf("daemon started pid=%d interval=%ds", os.Getpid(), interval)

	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	defer ticker.Stop()

	fetch := func() {
		data, err := doAPIRequest("GET", "/api/states", nil)
		if err != nil {
			logger.Println("ERROR fetch:", err)
			return
		}
		ef, err := entitiesFile()
		if err != nil {
			logger.Println("ERROR entities path:", err)
			return
		}
		tmp := ef + ".tmp"
		if err := os.WriteFile(tmp, data, 0644); err != nil {
			logger.Println("ERROR write tmp:", err)
			return
		}
		if err := os.Rename(tmp, ef); err != nil {
			logger.Println("ERROR rename:", err)
			return
		}
		logger.Printf("entities updated bytes=%d", len(data))
	}

	fetch()

	for range ticker.C {
		fetch()
	}
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Manage background entity sync daemon",
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start daemon in background",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isDaemonAlive() {
			pid, _ := readPID()
			fmt.Printf("daemon already running pid=%d\n", pid)
			return nil
		}

		interval := daemonInterval
		if !cmd.Flags().Changed("interval") && appConfig.Interval > 0 {
			interval = appConfig.Interval
		}

		exe, err := os.Executable()
		if err != nil {
			return fmt.Errorf("cannot get executable path: %w", err)
		}

		lf, err := logFile()
		if err != nil {
			return err
		}
		logOut, err := os.OpenFile(lf, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("cannot open log file: %w", err)
		}

		child := exec.Command(exe, "daemon", "_run",
			fmt.Sprintf("--interval=%d", interval),
		)
		child.Env = os.Environ()
		child.Stdout = logOut
		child.Stderr = logOut
		child.Stdin = nil

		setSysProcAttr(child)

		if err := child.Start(); err != nil {
			logOut.Close()
			return fmt.Errorf("failed to start daemon: %w", err)
		}
		logOut.Close()

		fmt.Printf("✓ daemon started pid=%d\n", child.Process.Pid)
		return nil
	},
}

var daemonRunCmd = &cobra.Command{
	Use:    "_run",
	Hidden: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		runForeground(daemonInterval)
		return nil
	},
}

var daemonStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the running daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !isDaemonAlive() {
			fmt.Println("daemon not running")
			return nil
		}
		pid, err := readPID()
		if err != nil {
			return fmt.Errorf("cannot read PID: %w", err)
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("cannot find process: %w", err)
		}
		if err := proc.Kill(); err != nil {
			return fmt.Errorf("cannot kill process: %w", err)
		}
		pf, _ := pidFile()
		os.Remove(pf)
		fmt.Printf("✓ daemon stopped pid=%d\n", pid)
		return nil
	},
}

var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show daemon status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isDaemonAlive() {
			pid, _ := readPID()
			fmt.Printf("daemon running pid=%d\n", pid)
		} else {
			fmt.Println("daemon not running")
		}
		return nil
	},
}

func init() {
	daemonCmd.PersistentFlags().IntVar(&daemonInterval, "interval", 300, "sync interval in seconds")
	daemonCmd.AddCommand(daemonStartCmd)
	daemonCmd.AddCommand(daemonRunCmd)
	daemonCmd.AddCommand(daemonStopCmd)
	daemonCmd.AddCommand(daemonStatusCmd)
}
