package dev

import (
	"fmt"
	"os"
	"path"
	"plugin"
	"reflect"

	"gopkg.in/yaml.v3"
)

var (
	ErrReadingBinding     = "error reading binding file"
	ErrYamlUnmarshal      = "error unmarshalling YAML"
	ErrPluginExec         = "error running plugin"
	ErrPluginLookup       = "failed to lookup symbol"
	ErrPluginMethodLookup = "method not found in the plugin"
	ErrPluginErr          = "plugin method execution failed"
)

type PluginBindings struct {
	SecretFunction string `yaml:"secret_function"`
	SecretSymbol   string `yaml:"secret_symbol"`
}

func RunPlugin(pluginPath string, filepath string) error {

	base_dir := os.Getenv("BASE_DIR")
	secretFilePath := path.Join(base_dir, "uploads/", os.Getenv("DEV_FILE"))

	fileContent, err := os.ReadFile(secretFilePath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("%s: %v\n", ErrReadingBinding, err)
	}

	var bindings PluginBindings

	err = yaml.Unmarshal(fileContent, &bindings)
	if err != nil {

		fmt.Println(err)
		return fmt.Errorf(" %s: %v\n", ErrYamlUnmarshal, err)
	}

	p, err := plugin.Open(pluginPath)

	if err != nil {
		return fmt.Errorf("%s: %v", ErrPluginExec, err)
	}

	symbol, err := p.Lookup(bindings.SecretSymbol)
	if err != nil {
		return fmt.Errorf(" %v", err)
	}

	v := reflect.ValueOf(symbol).MethodByName(bindings.SecretFunction)
	if !v.IsValid() {

		fmt.Println(err)
		return fmt.Errorf("%s", ErrPluginLookup)
	}

	inputs := []reflect.Value{reflect.ValueOf(filepath)}
	results := v.Call(inputs)

	if len(results) == 1 && !results[0].IsNil() {
		err, _ = results[0].Interface().(error)

		fmt.Println(err)
		return fmt.Errorf("%v", err)
	}

	return nil
}
