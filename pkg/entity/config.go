package entity

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// The name of the device which you will be watching for.
	// This app was built for yushakobo's Quick Paint device
	// However it will likely work with others.
	TargetDeviceName string `yaml:"target_device_name"`
	// A list of Macros to watch for... it should follow the
	// following format:
	//
	// macros:
	//   - name: xxxx
	//     on_press:
	//       entrypoint: "echo"
	//       args: ["hello", "world"]
	//     on_release:
	//       entrypoint: "echo"
	//       args: ["hello", "world"]
	//     on_hold:
	//       entrypoint: "echo"
	//       args: ["hello", "world"]
	//
	// given that a given macro can have one or all of the
	// handlers configured.
	CustomMacro  []Macro       `yaml:"macros"`
	CodeMacroMap map[int]Macro `yaml:"-"`
}

func generateDefaultConfigFile(name string) (*Config, error) {
	config := &Config{
		TargetDeviceName: "yushakobo Quick Paint Consumer Control",
		CustomMacro: []Macro{
			{
				Name: "KEY_MACRO1",
				OnPress: &Command{
					Type: "bash",
					Runner: &BashCommand{
						Entrypoint: "echo",
						Args:       []string{"hello", "world"},
					},
				},
			},
		},
	}
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}
	f, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Write(yamlData)
	return config, err
}

func openConfigFile(name string) (*Config, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	config := &Config{}
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)
	return config, err
}

func GetConfig(name string) (*Config, error) {
	log.Println("Loading config from:", name)
	spew.Config.DisableMethods = true

	var config *Config
	var err error

	dir := filepath.Dir(name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	if _, err := os.Stat(name); os.IsNotExist(err) {
		config, err = generateDefaultConfigFile(name)
	} else {
		config, err = openConfigFile(name)
	}

	return config, err
}

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type configAlias Config
	aux := configAlias(*c)
	if err := unmarshal(&aux); err != nil {
		return err
	}
	*c = Config(aux)
	codeMap := make(map[string]int)
	for i, val := 1, 656; val <= 685; i++ {
		codeMap["KEY_MACRO"+strconv.Itoa(i)] = val
		val++
	}
	c.CodeMacroMap = make(map[int]Macro)
	for _, cmd := range c.CustomMacro {
		code, ok := codeMap[cmd.Name]
		if !ok {
			continue
		}
		c.CodeMacroMap[code] = cmd
	}
	return nil
}

func (c *Config) UdevDeviceName() string {
	return "\"" + c.TargetDeviceName + "\""
}

type Macro struct {
	// The name of this must be one of KEY_MACRO{1-30}
	// if not it will be ignored.
	Name       string           `yaml:"name"`
	OnPress    *Command         `yaml:"on_press,omitempty"`
	OnRelease  *Command         `yaml:"on_release,omitempty"`
	OnHold     *Command         `yaml:"on_hold,omitempty"`
	HandlerMap map[int]*Command `yaml:"-"`
}

func (m *Macro) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tmpMacro Macro
	aux := tmpMacro(*m)
	if err := unmarshal(&aux); err != nil {
		return err
	}
	*m = Macro(aux)
	codeMap := make(map[string]int)
	for i, val := 1, 656; val <= 685; i++ {
		codeMap["KEY_MACRO"+strconv.Itoa(i)] = val
		val++
	}
	m.HandlerMap = make(map[int]*Command)
	if m.OnPress != nil {
		m.HandlerMap[1] = m.OnPress
	}
	if m.OnRelease != nil {
		m.HandlerMap[0] = m.OnRelease
	}
	if m.OnHold != nil {
		m.HandlerMap[2] = m.OnHold
	}
	return nil
}
