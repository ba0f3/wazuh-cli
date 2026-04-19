package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/ba0f3/wazuh-cli/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage wazuh-cli configuration file",
	Long:  `Get, set, or list values in the configuration file (~/.config/wazuh/config.json).`,
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get a value from the config file",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfigFileOnly()
		if err != nil {
			if os.IsNotExist(err) {
				cfg = &config.Config{}
			} else {
				return err
			}
		}

		m := configToMap(cfg)

		if len(args) == 0 {
			// List all
			return listConfig(m)
		}

		key := args[0]
		val, ok := m[key]
		if !ok {
			return fmt.Errorf("config key %q not found", key)
		}

		if val == "" {
			fmt.Fprintf(os.Stderr, "Config key %q is not set\n", key)
		} else {
			fmt.Println(val)
		}
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a value in the config file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfigFileOnly()
		if err != nil {
			// If file doesn't exist, start with a fresh one
			if os.IsNotExist(err) {
				cfg = &config.Config{
					Timeout: config.DefaultTimeout,
					Output:  config.DefaultOutput,
					Pretty:  true,
				}
			} else {
				return err
			}
		}

		key := args[0]
		val := args[1]

		if err := setConfigValue(cfg, key, val); err != nil {
			return err
		}

		if err := cfg.Save(configPath); err != nil {
			return err
		}

		displayVal := val
		if key == "password" {
			displayVal = "********"
		}
		fmt.Fprintf(os.Stderr, "✓ Set %s = %s\n", key, displayVal)
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all values in the config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfigFileOnly()
		if err != nil {
			if os.IsNotExist(err) {
				cfg = &config.Config{}
			} else {
				return err
			}
		}
		return listConfig(configToMap(cfg))
	},
}

var configDeleteCmd = &cobra.Command{
	Use:   "delete <key>",
	Short: "Delete a value from the config file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfigFileOnly()
		if err != nil {
			return err
		}

		key := args[0]
		if err := deleteConfigValue(cfg, key); err != nil {
			return err
		}

		if err := cfg.Save(configPath); err != nil {
			return err
		}

		fmt.Fprintf(os.Stderr, "✓ Deleted %s\n", key)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configGetCmd, configSetCmd, configListCmd, configDeleteCmd)
}

// loadConfigFileOnly loads ONLY the config file, ignoring flags/env.
func loadConfigFileOnly() (*config.Config, error) {
	path := configPath
	if path == "" {
		path = config.DefaultConfigPath()
	}

	cfg := &config.Config{}
	// We use a small hack here: Load with no overrides to get the file values.
	// But Load also merges env/dotenv. To get ONLY the file, we need a custom loader
	// or just use the Save/Load logic from the config package if it were public.
	// Since I wrote internal/config/config.go, I know I can just use json.Unmarshal.

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config JSON: %w", err)
	}

	return cfg, nil
}

func configToMap(cfg *config.Config) map[string]string {
	m := make(map[string]string)
	m["url"] = cfg.URL
	m["user"] = cfg.User
	m["password"] = cfg.Password
	m["insecure"] = strconv.FormatBool(cfg.Insecure)
	m["ca_cert"] = cfg.CACert
	m["client_cert"] = cfg.ClientCert
	m["client_key"] = cfg.ClientKey
	m["timeout"] = strconv.Itoa(cfg.Timeout)
	m["output"] = cfg.Output
	m["pretty"] = strconv.FormatBool(cfg.Pretty)
	m["debug"] = strconv.FormatBool(cfg.Debug)
	m["quiet"] = strconv.FormatBool(cfg.Quiet)
	return m
}

func setConfigValue(cfg *config.Config, key, val string) error {
	switch strings.ToLower(key) {
	case "url":
		cfg.URL = val
	case "user":
		cfg.User = val
	case "password":
		cfg.Password = val
	case "insecure":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Insecure = b
	case "ca_cert":
		cfg.CACert = val
	case "client_cert":
		cfg.ClientCert = val
	case "client_key":
		cfg.ClientKey = val
	case "timeout":
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		cfg.Timeout = i
	case "output":
		cfg.Output = val
	case "pretty":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Pretty = b
	case "debug":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Debug = b
	case "quiet":
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		cfg.Quiet = b
	default:
		return fmt.Errorf("unknown config key %q", key)
	}
	return nil
}

func deleteConfigValue(cfg *config.Config, key string) error {
	switch strings.ToLower(key) {
	case "url":
		cfg.URL = ""
	case "user":
		cfg.User = ""
	case "password":
		cfg.Password = ""
	case "insecure":
		cfg.Insecure = false
	case "ca_cert":
		cfg.CACert = ""
	case "client_cert":
		cfg.ClientCert = ""
	case "client_key":
		cfg.ClientKey = ""
	case "timeout":
		cfg.Timeout = 0
	case "output":
		cfg.Output = ""
	case "pretty":
		cfg.Pretty = false
	case "debug":
		cfg.Debug = false
	case "quiet":
		cfg.Quiet = false
	default:
		return fmt.Errorf("unknown config key %q", key)
	}
	return nil
}

func listConfig(m map[string]string) error {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		if m[k] != "" {
			val := m[k]
			if k == "password" {
				val = "********"
			}
			fmt.Printf("%s = %s\n", k, val)
		}
	}
	return nil
}
