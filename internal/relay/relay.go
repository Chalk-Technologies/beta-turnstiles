package relay

import (
	"beta-turnstiles/internal/config"
	"github.com/stianeikeland/go-rpio"
	"log"
	"time"
)

var pin *rpio.Pin

func Init() error {
	log.Println("Initializing pin connections")
	err := rpio.Open()
	if err != nil {
		return err
	}
	log.Printf("Initializing relay pin %v\n", config.GlobalConfig.RelayPin)
	p := rpio.Pin(config.GlobalConfig.RelayPin)
	pin = &p
	if config.GlobalConfig.HighMode {
		log.Println("Setting pin to high")
		pin.High()
	} else {
		log.Println("Setting pin to low")
		pin.Low()
	}
	return nil
}

func CleanUp() error {
	log.Println("Cleaning up pin connections")
	if config.GlobalConfig.HighMode {
		pin.High()
	} else {
		pin.Low()
	}
	return rpio.Close()
}
func TriggerRelay() {
	log.Println("Toggling pin state")
	// will error out if pin was not initialized
	pin.Toggle()
	time.Sleep(100 * time.Millisecond)
	pin.Toggle()
}
