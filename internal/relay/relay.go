package relay

import (
	"beta-turnstiles/internal/config"
	"errors"
	"log"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/conn/v3/pin/pinreg"
	"periph.io/x/host/v3"
	"time"
)

var pin gpio.PinIO

func Init() error {
	log.Println("Initializing pin connections")
	_, err := host.Init()
	if err != nil {
		return err
	}
	log.Printf("Initializing relay pin %v\n", config.GlobalConfig.RelayPin)
	p := gpioreg.ByName(config.GlobalConfig.RelayPin)
	if p == nil {
		allPins := gpioreg.All()
		for _, i := range allPins {
			log.Printf("%v is connected: %v\n", i.Name(), pinreg.IsConnected(i))
		}
		return errors.New("failed to find relay pin")
	}
	pin = p
	if config.GlobalConfig.HighMode {
		log.Println("Setting pin to high")
		err = pin.Out(gpio.High)
	} else {
		log.Println("Setting pin to low")
		err = pin.Out(gpio.Low)
	}
	return err
}

func CleanUp() error {
	log.Println("Cleaning up pin connections")
	if config.GlobalConfig.HighMode {
		return pin.Out(gpio.High)
	} else {
		return pin.Out(gpio.Low)
	}
}
func TriggerRelay() {
	log.Println("Toggling pin state")
	// will error out if pin was not initialized
	if config.GlobalConfig.HighMode {
		err := pin.Out(gpio.Low)
		if err != nil {
			log.Printf("Failed to toggle pin state %v\n", err)
			return
		}
		time.Sleep(time.Duration(config.GlobalConfig.SignalDurationMS) * time.Millisecond)
		err = pin.Out(gpio.High)
		if err != nil {
			log.Printf("Failed to toggle pin state %v\n", err)
		}
	} else {
		err := pin.Out(gpio.High)
		if err != nil {
			log.Printf("Failed to toggle pin state %v\n", err)
			return
		}
		time.Sleep(time.Duration(config.GlobalConfig.SignalDurationMS) * time.Millisecond)
		err = pin.Out(gpio.Low)
		if err != nil {
			log.Printf("Failed to toggle pin state %v\n", err)
		}
	}
}
