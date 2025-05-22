package relay

import (
	"beta-turnstiles/internal/config"
	"github.com/stianeikeland/go-rpio"
	"time"
)

var pin *rpio.Pin

func Init() error {
	err := rpio.Open()
	if err != nil {
		return err
	}
	p := rpio.Pin(config.GlobalConfig.RelayPin)
	pin = &p
	if config.GlobalConfig.HighMode {
		pin.High()
	} else {
		pin.Low()
	}
	return nil
}

func CleanUp() error {
	if config.GlobalConfig.HighMode {
		pin.High()
	} else {
		pin.Low()
	}
	return rpio.Close()
}
func TriggerRelay() {
	// will error out if pin was not initialized
	pin.Toggle()
	time.Sleep(100 * time.Millisecond)
	pin.Toggle()
}
