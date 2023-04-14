package plugin

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PluginOptionType int

const (
	PluginOptionTypeNone   PluginOptionType = 0
	PluginOptionTypeString PluginOptionType = 1
	PluginOptionTypeBool   PluginOptionType = 2
	PluginOptionTypeInt    PluginOptionType = 3
	PluginOptionTypeUint   PluginOptionType = 4
)

type PluginOption struct {
	Name         string
	Type         PluginOptionType
	CustomEnvVar string
	Description  string
	DefaultValue interface{}
	Dest         interface{}
}

func (p *PluginOption) AddToFlagSet(fs *flag.FlagSet, flagPrefix string) error {
	flagName := fmt.Sprintf("%s%s", flagPrefix, p.Name)
	switch p.Type {
	case PluginOptionTypeString:
		fs.StringVar(p.Dest.(*string), flagName, p.DefaultValue.(string), p.Description)
	case PluginOptionTypeBool:
		fs.BoolVar(p.Dest.(*bool), flagName, p.DefaultValue.(bool), p.Description)
	case PluginOptionTypeInt:
		fs.IntVar(p.Dest.(*int), flagName, p.DefaultValue.(int), p.Description)
	case PluginOptionTypeUint:
		fs.UintVar(p.Dest.(*uint), flagName, p.DefaultValue.(uint), p.Description)
	default:
		return fmt.Errorf("unknown plugin option type %d for option %s", p.Type, p.Name)
	}
	return nil
}

func (p *PluginOption) ProcessEnvVars(envPrefix string) error {
	envVars := []string{
		// Automatically generate env var from specified prefix and option name
		strings.ToUpper(
			strings.ReplaceAll(
				fmt.Sprintf(
					"%s%s",
					envPrefix,
					p.Name,
				),
				"-",
				"_",
			),
		),
	}
	// Also check any custom env var specified
	if p.CustomEnvVar != "" {
		envVars = append(envVars, p.CustomEnvVar)
	}
	for _, envVar := range envVars {
		if value, ok := os.LookupEnv(envVar); ok {
			switch p.Type {
			case PluginOptionTypeString:
				*(p.Dest.(*string)) = value
			case PluginOptionTypeBool:
				value, err := strconv.ParseBool(value)
				if err != nil {
					return fmt.Errorf("error processing env vars: %s", err)
				}
				*(p.Dest.(*bool)) = value
			case PluginOptionTypeInt:
				// We limit to 32-bit to not get inconsistent behavior on 32-bit platforms
				value, err := strconv.ParseInt(value, 10, 32)
				if err != nil {
					return fmt.Errorf("error processing env vars: %s", err)
				}
				*(p.Dest.(*int)) = int(value)
			case PluginOptionTypeUint:
				// We limit to 32-bit to not get inconsistent behavior on 32-bit platforms
				value, err := strconv.ParseUint(value, 10, 32)
				if err != nil {
					return fmt.Errorf("error processing env vars: %s", err)
				}
				*(p.Dest.(*uint)) = uint(value)
			default:
				return fmt.Errorf("unknown plugin option type %d for option %s", p.Type, p.Name)
			}
		}
	}
	return nil
}

func (p *PluginOption) ProcessConfig(pluginData map[interface{}]interface{}) error {
	if optionData, ok := pluginData[p.Name]; ok {
		switch p.Type {
		case PluginOptionTypeString:
			switch value := optionData.(type) {
			case string:
				*(p.Dest.(*string)) = value
			default:
				return fmt.Errorf("invalid value for option '%s': expected string and got %T", p.Name, optionData)
			}
		case PluginOptionTypeBool:
			switch value := optionData.(type) {
			case bool:
				*(p.Dest.(*bool)) = value
			default:
				return fmt.Errorf("invalid value for option '%s': expected bool and got %T", p.Name, optionData)
			}
		case PluginOptionTypeInt:
			switch value := optionData.(type) {
			case int:
				*(p.Dest.(*int)) = int(value)
			default:
				return fmt.Errorf("invalid value for option '%s': expected int and got %T", p.Name, optionData)
			}
		case PluginOptionTypeUint:
			switch value := optionData.(type) {
			case int:
				*(p.Dest.(*uint)) = uint(value)
			default:
				return fmt.Errorf("invalid value for option '%s': expected uint and got %T", p.Name, optionData)
			}
		default:
			return fmt.Errorf("unknown plugin option type %d for option %s", p.Type, p.Name)

		}
	}
	return nil
}