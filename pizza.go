package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/gvalkov/golang-evdev"
)

const targetWord = "pizza"

func main() {
	layout, err := detectKeyboardLayout()
	if err != nil {
		log.Printf("Warning: Could not detect keyboard layout, defaulting to 'us'. Error: %v", err)
	}
	log.Printf("Detected keyboard layout: %s", layout)

	keyMap := generateKeyMap(layout)

	devices, err := findKeyboards()
	if err != nil {
		log.Fatalf("Error finding keyboards: %v", err)
	}
	if len(devices) == 0 {
		log.Fatalf("No keyboards found. Make sure the program is run with sufficient privileges (e.g., as root or in the 'input' group).")
	}

	keyChan := make(chan evdev.InputEvent, 10)
	var wg sync.WaitGroup

	for _, dev := range devices {
		wg.Add(1)
		go listenToDevice(dev, keyChan, &wg)
	}

	log.Println("Monitoring keyboard inputs...")
	processKeyEvents(keyChan, keyMap)

	wg.Wait()
}

func listenToDevice(device *evdev.InputDevice, keyChan chan<- evdev.InputEvent, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Printf("Listening on: %s (%s)", device.Name, device.Phys)
	for {
		event, err := device.ReadOne()
		if err != nil {
			log.Printf("Error reading from device %s: %v. Stopping listener.", device.Name, err)
			return
		}

		if event.Type == evdev.EV_KEY && event.Value == 1 { // Key press
			keyChan <- *event
		}
	}
}

func processKeyEvents(keyChan <-chan evdev.InputEvent, keyMap map[uint16]rune) {
	var buffer []rune

	for event := range keyChan {
		code := event.Code

		if code == evdev.KEY_ENTER || code == evdev.KEY_SPACE || code == evdev.KEY_ESC {
			buffer = nil
			continue
		}

		char, ok := keyMap[code]
		if !ok {
			continue
		}

		buffer = append(buffer, char)

		if len(buffer) > len(targetWord) {
			buffer = buffer[1:]
		}

		if string(buffer) == targetWord {
			log.Println("Target word 'pizza' detected. Locking the system.")
			if err := lockSystem(); err != nil {
				log.Printf("Error locking the system: %v", err)
			}
			buffer = nil
		}
	}
}

func findKeyboards() ([]*evdev.InputDevice, error) {
	var keyboards []*evdev.InputDevice
	devices, err := evdev.ListInputDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to list input devices: %w", err)
	}

	for _, dev := range devices {
		isKeyboard := false
		// Brute-force approach: Iterate over all device capabilities.
		// This avoids issues with the map key data types.
		for _, caps := range dev.Capabilities {
			for _, key := range caps {
				if key.Code == evdev.KEY_A {
					isKeyboard = true
					break
				}
			}
			if isKeyboard {
				break
			}
		}

		if isKeyboard {
			keyboards = append(keyboards, dev)
		}
	}
	return keyboards, nil
}

func lockSystem() error {
	cmd := exec.Command("loginctl", "lock-session")
	if err := cmd.Run(); err == nil {
		log.Println("System locked via loginctl.")
		return nil
	}

	log.Println("loginctl failed, trying fallback with 'xset'...")
	cmd = exec.Command("xset", "s", "activate")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("neither 'loginctl' nor 'xset' could be executed: %w", err)
	}
	log.Println("System locked via xset.")
	return nil
}

func detectKeyboardLayout() (string, error) {
	cmd := exec.Command("localectl", "status")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "us", fmt.Errorf("could not execute 'localectl status': %w", err)
	}

	scanner := bufio.NewScanner(&out)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "X11 Layout") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				layout := strings.TrimSpace(parts[1])
				return strings.Split(layout, ",")[0], nil
			}
		}
	}
	return "us", fmt.Errorf("X11 Layout not found in 'localectl' output")
}

func generateKeyMap(layout string) map[uint16]rune {
	m := map[uint16]rune{
		evdev.KEY_A: 'a', evdev.KEY_B: 'b', evdev.KEY_C: 'c', evdev.KEY_D: 'd',
		evdev.KEY_E: 'e', evdev.KEY_F: 'f', evdev.KEY_G: 'g', evdev.KEY_H: 'h',
		evdev.KEY_I: 'i', evdev.KEY_J: 'j', evdev.KEY_K: 'k', evdev.KEY_L: 'l',
		evdev.KEY_M: 'm', evdev.KEY_N: 'n', evdev.KEY_O: 'o', evdev.KEY_P: 'p',
		evdev.KEY_Q: 'q', evdev.KEY_R: 'r', evdev.KEY_S: 's', evdev.KEY_T: 't',
		evdev.KEY_U: 'u', evdev.KEY_V: 'v', evdev.KEY_W: 'w', evdev.KEY_X: 'x',
	}

	if layout == "de" || layout == "qwertz" {
		m[evdev.KEY_Y] = 'z'
		m[evdev.KEY_Z] = 'y'
	} else {
		m[evdev.KEY_Y] = 'y'
		m[evdev.KEY_Z] = 'z'
	}
	return m
}