package main

import (
	"fmt"
)

const (
	ConfigBaseDir = iota
	ConfigServerPort
	ConfigMonkeyPatcher // Deprecated (This is the best we can do when something should not be used)
	ConfigDuckTyper
	ConfigAreWeHuman
	ConfigOrAreWeDancers
	ConfigFrobnitz
	ConfigDroids
)

// Create a new Config
func New(key int, value interface{}) *Config {
	return &Config{
		next_config: nil,
		key:         key,
		value:       value,
	}
}

type Config struct {
	// Link to next item
	next_config *Config
	// My key
	key int
	// My value
	value interface{}
}

// IsNext indicates if we are at the end of the list
// Usually, there is no reason for userland code to call this
func (c *Config) IsNext() bool {
	return c.next_config != nil
}

// Next returns the next item
// Usually, there is no reason for userland code to call this
func (c *Config) Next() *Config {
	return c.next_config
}

// Add a new config item to the list, and return it as the new head of the list.
func (c *Config) Add(key int, value interface{}) *Config {
	return &Config{
		next_config: c,
		key:         key,
		value:       value,
	}
}

// There are other methods, like GetAll, that I did not do here.

// Get finds and returns the first matchng value in the list.
// If no value is found, it returns the default value.
// The 'ok' flag indicates whether the value was found (true)
// or whether it was not found and the default was used (false).
func (c *Config) Get(key int, default_value interface{}) (interface{}, bool) {
	if c.key == key {
		// If it's a match, return this
		return c.value, true
	} else if !c.IsNext() {
		// If there is no next value, return default
		return default_value, false
	}
	// In all other cases, send it to the next config item
	return c.next_config.Get(key, default_value)
}

func main() {
	println("Initializing a basic config")
	config := New(ConfigFrobnitz, "hello").
		Add(ConfigAreWeHuman, false).
		Add(ConfigOrAreWeDancers, true)

	// Get an item back out
	human, ok := config.Get(ConfigAreWeHuman, true)
	fmt.Printf("Are we human? %v And was that the default? %v\n", human.(bool), ok)

	// Get a string
	frob, _ := config.Get(ConfigFrobnitz, "goodbye")
	fmt.Printf("%s world\n", frob.(string))

	// And this is what a miss looks like
	if oops, ok := config.Get(ConfigDroids, 123); !ok {
		fmt.Printf("These are not the droids you are looking for. %#v\n", oops)
	}

	// Later on, something can override a config and add new params.
	fmt.Println("Adding some new droids")
	config = config.Add(ConfigDroids, 211)
	if droids, ok := config.Get(ConfigDroids, 123); ok {
		fmt.Printf("These ARE the droids you are looking for. %#v\n", droids)
	}
}
