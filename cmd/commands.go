package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(stateCmd)
	rootCmd.AddCommand(statesCmd)
	rootCmd.AddCommand(onCmd)
	rootCmd.AddCommand(offCmd)
	rootCmd.AddCommand(toggleCmd)
	rootCmd.AddCommand(sceneCmd)
	rootCmd.AddCommand(scriptCmd)
	rootCmd.AddCommand(automationCmd)
	rootCmd.AddCommand(climateCmd)
	rootCmd.AddCommand(cameraCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(callCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(daemonCmd)
}

func getEntities() ([]byte, error) {
	ef, err := entitiesFile()
	if err == nil {
		if isDaemonAlive() {
			if _, statErr := os.Stat(ef); statErr == nil {
				data, readErr := os.ReadFile(ef)
				if readErr == nil {
					return data, nil
				}
			}
		}
	}
	return doAPIRequest("GET", "/api/states", nil)
}

func getDomain(entity string) string {
	parts := strings.SplitN(entity, ".", 2)
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

var stateCmd = &cobra.Command{
	Use:     "state <entity_id>",
	Aliases: []string{"get"},
	Short:   "Get entity state",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entity := args[0]
		resp, err := doAPIRequest("GET", "/api/states/"+entity, nil)
		if err != nil {
			return err
		}
		var data struct {
			State string `json:"state"`
		}
		if err := json.Unmarshal(resp, &data); err != nil {
			return err
		}
		if data.State == "" {
			fmt.Println("unknown")
		} else {
			fmt.Println(data.State)
		}
		return nil
	},
}

var statesCmd = &cobra.Command{
	Use:   "states <entity_id>",
	Short: "Get full entity state with attributes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entity := args[0]
		resp, err := doAPIRequest("GET", "/api/states/"+entity, nil)
		if err != nil {
			return err
		}
		var out bytes.Buffer
		json.Indent(&out, resp, "", "  ")
		fmt.Println(out.String())
		return nil
	},
}

var onCmd = &cobra.Command{
	Use:     "on <entity_id> [brightness]",
	Aliases: []string{"turn_on"},
	Short:   "Turn on (optional brightness 0-255)",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		entity := args[0]
		domain := getDomain(entity)
		payload := map[string]interface{}{"entity_id": entity}
		if len(args) == 2 {
			brightness, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid brightness: %v", err)
			}
			payload["brightness"] = brightness
		}
		_, err := doAPIRequest("POST", "/api/services/"+domain+"/turn_on", payload)
		if err != nil {
			return err
		}
		fmt.Printf("✓ %s turned on\n", entity)
		return nil
	},
}

var offCmd = &cobra.Command{
	Use:     "off <entity_id>",
	Aliases: []string{"turn_off"},
	Short:   "Turn off",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entity := args[0]
		domain := getDomain(entity)
		_, err := doAPIRequest("POST", "/api/services/"+domain+"/turn_off", map[string]interface{}{
			"entity_id": entity,
		})
		if err != nil {
			return err
		}
		fmt.Printf("✓ %s turned off\n", entity)
		return nil
	},
}

var toggleCmd = &cobra.Command{
	Use:   "toggle <entity_id>",
	Short: "Toggle on/off",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		entity := args[0]
		domain := getDomain(entity)
		_, err := doAPIRequest("POST", "/api/services/"+domain+"/toggle", map[string]interface{}{
			"entity_id": entity,
		})
		if err != nil {
			return err
		}
		fmt.Printf("✓ %s toggled\n", entity)
		return nil
	},
}

var sceneCmd = &cobra.Command{
	Use:   "scene <name>",
	Short: "Activate scene",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scene := args[0]
		if !strings.HasPrefix(scene, "scene.") {
			scene = "scene." + scene
		}
		_, err := doAPIRequest("POST", "/api/services/scene/turn_on", map[string]interface{}{
			"entity_id": scene,
		})
		if err != nil {
			return err
		}
		fmt.Printf("✓ Scene %s activated\n", scene)
		return nil
	},
}

var scriptCmd = &cobra.Command{
	Use:   "script <name>",
	Short: "Run script",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		script := args[0]
		if !strings.HasPrefix(script, "script.") {
			script = "script." + script
		}
		_, err := doAPIRequest("POST", "/api/services/script/turn_on", map[string]interface{}{
			"entity_id": script,
		})
		if err != nil {
			return err
		}
		fmt.Printf("✓ Script %s executed\n", script)
		return nil
	},
}

var automationCmd = &cobra.Command{
	Use:     "automation <name>",
	Aliases: []string{"trigger"},
	Short:   "Trigger automation",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		auto := args[0]
		if !strings.HasPrefix(auto, "automation.") {
			auto = "automation." + auto
		}
		_, err := doAPIRequest("POST", "/api/services/automation/trigger", map[string]interface{}{
			"entity_id": auto,
		})
		if err != nil {
			return err
		}
		fmt.Printf("✓ Automation %s triggered\n", auto)
		return nil
	},
}

var climateCmd = &cobra.Command{
	Use:     "climate <entity> <temp>",
	Aliases: []string{"temp"},
	Short:   "Set temperature",
	Args:    cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		entity := args[0]
		temp, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			return fmt.Errorf("invalid temperature: %v", err)
		}
		_, err = doAPIRequest("POST", "/api/services/climate/set_temperature", map[string]interface{}{
			"entity_id":   entity,
			"temperature": temp,
		})
		if err != nil {
			return err
		}
		fmt.Printf("✓ %s set to %v°\n", entity, temp)
		return nil
	},
}

var cameraCmd = &cobra.Command{
	Use:     "camera <entity_id|list> [output_file]",
	Aliases: []string{"snapshot", "snap"},
	Short:   "Capture image from camera entity or list cameras",
	Args:    cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] == "list" {
			resp, err := getEntities()
			if err != nil {
				return err
			}

			var states []struct {
				EntityID string `json:"entity_id"`
				State    string `json:"state"`
			}
			if err := json.Unmarshal(resp, &states); err != nil {
				return err
			}

			for _, s := range states {
				if strings.HasPrefix(s.EntityID, "camera.") {
					fmt.Printf("%s: %s\n", s.EntityID, s.State)
				}
			}
			return nil
		}

		entity := args[0]
		useFrigate := appConfig.FrigateURL != ""

		var img []byte
		var err error

		if useFrigate {
			cameraName := strings.TrimPrefix(entity, "camera.")
			img, err = doFrigateCameraRequest(cameraName)
		} else {
			if !strings.HasPrefix(entity, "camera.") {
				entity = "camera." + entity
			}
			img, err = doCameraRequest(entity)
		}

		if err != nil {
			return err
		}

		var outFile string
		if len(args) == 2 {
			outFile = args[1]
		} else {
			baseName := strings.TrimPrefix(entity, "camera.")
			if useFrigate {
				outFile = "fg." + baseName + ".jpg"
			} else {
				outFile = "ha." + baseName + ".jpg"
			}
		}

		if err := os.WriteFile(outFile, img, 0644); err != nil {
			return fmt.Errorf("failed to write image: %v", err)
		}

		fmt.Printf("✓ Image saved to %s (%d bytes)\n", outFile, len(img))
		return nil
	},
}

var listCmd = &cobra.Command{
	Use:   "list [domain]",
	Short: "List entities (lights, switches, all)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := "all"
		if len(args) > 0 {
			filter = args[0]
		}
		if filter != "all" && strings.HasSuffix(filter, "s") {
			filter = filter[:len(filter)-1]
		}

		resp, err := getEntities()
		if err != nil {
			return err
		}

		var states []struct {
			EntityID string `json:"entity_id"`
		}
		if err := json.Unmarshal(resp, &states); err != nil {
			return err
		}

		var result []string
		for _, s := range states {
			if filter == "all" || strings.HasPrefix(s.EntityID, filter+".") {
				result = append(result, s.EntityID)
			}
		}

		sort.Strings(result)
		for _, e := range result {
			fmt.Println(e)
		}
		return nil
	},
}

var searchCmd = &cobra.Command{
	Use:   "search <pattern>",
	Short: "Search entities by name",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := strings.ToLower(args[0])

		resp, err := getEntities()
		if err != nil {
			return err
		}

		var states []struct {
			EntityID string `json:"entity_id"`
			State    string `json:"state"`
		}
		if err := json.Unmarshal(resp, &states); err != nil {
			return err
		}

		for _, s := range states {
			if strings.Contains(strings.ToLower(s.EntityID), pattern) {
				fmt.Printf("%s: %s\n", s.EntityID, s.State)
			}
		}
		return nil
	},
}

var callCmd = &cobra.Command{
	Use:   "call <domain> <svc> [json]",
	Short: "Call any service",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		service := args[1]
		data := "{}"
		if len(args) == 3 {
			data = args[2]
		}

		resp, err := doAPIRequest("POST", "/api/services/"+domain+"/"+service, data)
		if err != nil {
			return err
		}
		var out bytes.Buffer
		json.Indent(&out, resp, "", "  ")
		if out.Len() > 0 && out.String() != "null\n" && out.String() != "[]\n" {
			fmt.Println(out.String())
		}
		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get HA instance info",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := doAPIRequest("GET", "/api/", nil)
		if err != nil {
			return err
		}
		var out bytes.Buffer
		json.Indent(&out, resp, "", "  ")
		fmt.Println(out.String())
		return nil
	},
}

func getConfigPaths() (string, string, error) {
	var configDir string
	if envDir := os.Getenv("HA_CONFIG_DIR"); envDir != "" {
		configDir = envDir
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", "", fmt.Errorf("could not get home directory: %v", err)
		}
		configDir = filepath.Join(home, ".ha-cli")
	}
	return configDir, filepath.Join(configDir, "config.json"), nil
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Setup configuration interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter Home Assistant URL (e.g. http://192.168.1.100:8123): ")
		url, _ := reader.ReadString('\n')
		url = strings.TrimSpace(url)

		fmt.Print("Enter Home Assistant Long-Lived Access Token: ")
		token, _ := reader.ReadString('\n')
		token = strings.TrimSpace(token)

		fmt.Print("Enter Frigate URL (optional, e.g. http://192.168.1.100:5000): ")
		frigateURL, _ := reader.ReadString('\n')
		frigateURL = strings.TrimSpace(frigateURL)

		if url == "" || token == "" {
			return fmt.Errorf("URL and Token cannot be empty")
		}

		fmt.Print("Enter daemon sync interval in seconds (default: 300): ")
		intervalStr, _ := reader.ReadString('\n')
		intervalStr = strings.TrimSpace(intervalStr)
		interval := 300
		if intervalStr != "" {
			parsed, err := strconv.Atoi(intervalStr)
			if err != nil || parsed <= 0 {
				return fmt.Errorf("invalid interval: must be a positive integer")
			}
			interval = parsed
		}

		config := Config{
			URL:        url,
			Token:      token,
			Interval:   interval,
			FrigateURL: frigateURL,
		}

		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return err
		}

		configDir, configPath, err := getConfigPaths()
		if err != nil {
			return err
		}

		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory %s: %v", configDir, err)
		}

		err = os.WriteFile(configPath, data, 0600)
		if err != nil {
			return fmt.Errorf("failed to write config to %s: %v", configPath, err)
		}

		fmt.Printf("✓ Configuration saved to %s\n", configPath)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a specific configuration key",
	RunE: func(cmd *cobra.Command, args []string) error {
		validKeys := []string{"url", "token", "interval", "frigate_url"}

		if len(args) == 0 {
			return fmt.Errorf("missing config key. Available keys: %s", strings.Join(validKeys, ", "))
		}

		key := strings.ToLower(args[0])

		isValid := false
		for _, k := range validKeys {
			if key == k {
				isValid = true
				break
			}
		}

		if !isValid {
			return fmt.Errorf("invalid config key '%s'. Available keys: %s", key, strings.Join(validKeys, ", "))
		}

		if len(args) < 2 {
			return fmt.Errorf("missing value for key '%s'", key)
		}

		value := args[1]

		configDir, configPath, err := getConfigPaths()
		if err != nil {
			return err
		}

		var config Config
		if data, err := os.ReadFile(configPath); err == nil {
			json.Unmarshal(data, &config)
		}

		switch key {
		case "url":
			config.URL = value
		case "token":
			config.Token = value
		case "interval":
			val, err := strconv.Atoi(value)
			if err != nil || val <= 0 {
				return fmt.Errorf("invalid interval: must be a positive integer")
			}
			config.Interval = val
		case "frigate_url":
			config.FrigateURL = value
		}

		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %v", err)
		}

		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(configPath, data, 0600)
		if err != nil {
			return fmt.Errorf("failed to write config: %v", err)
		}

		fmt.Printf("✓ Configuration '%s' updated\n", key)
		return nil
	},
}
