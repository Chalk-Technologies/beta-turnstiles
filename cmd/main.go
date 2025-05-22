package main

import (
	"beta-turnstiles/internal/checkCode"
	"beta-turnstiles/internal/config"
	"beta-turnstiles/internal/relay"
	"flag"
	"fmt"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"log"
	"strings"
)

func main() {
	// flags needed for accounting stuff
	// get wise api key
	// get desired month
	var headless = flag.Bool("headless", false, "run in background without UI")
	log.Printf("got headless flag %v", headless)
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = relay.Init()
	if err != nil {
		log.Fatal(err)
	}
	// flags needed for everything
	flag.Parse()
	if *headless {
		// idk yet TODO
	} else {
		runUI()
	}

	return
}

var scannedText = binding.NewString()
var errorText = binding.NewString()
var lastCheckinInfo = binding.NewString()

func runUI() {
	// open window to handle checkins
	a := app.New()
	w := a.NewWindow("BETA Turnstiles App")

	// initiate content
	header := widget.NewLabel("Scan a tag")
	scannedTextField := widget.NewEntryWithData(scannedText)
	scannedTextField.OnSubmitted = scanText
	errorTextField := widget.NewLabelWithData(errorText)
	errorContainer := container.NewVBox(errorTextField)
	successTextField := widget.NewLabelWithData(lastCheckinInfo)
	w.SetContent(container.NewVBox(
		header,
		scannedTextField,
		errorContainer,
		successTextField,
	))
	w.ShowAndRun()
	// focus text field
	w.Canvas().Focus(scannedTextField)
}

var permittedPrefixes = []string{"PA_", "SU_", "CL_", "SE_", "SL_", "EV_"} //  "BCK_"} // "GC_"}

func scanText(text string) {
	errorText.Set("")       // ignoring error
	lastCheckinInfo.Set("") // ignoring error
	// check if this is a good checkin string
	isValid := false
	for _, prefix := range permittedPrefixes {
		if strings.HasPrefix(text, prefix) {
			isValid = true
			break
		}
	}
	if isValid {
		//if !config.GlobalConfig.SingleMode {
		//	// check code
		//	err := checkCode.CheckCode(text)
		//	if err != nil {
		//		errorText.Set(err.Error())
		//		scannedText.Set("")
		//		return
		//	}
		//}

		// check in
		result, err := checkCode.ConsumeCode(text)
		if err != nil {
			errorText.Set(err.Error())
		} else {
			relay.TriggerRelay()
			lastCheckinInfo.Set(result)
		}
	} else {
		errorText.Set(fmt.Sprintf("%s is not permitted", text))
	}
	err := scannedText.Set("")
	if err != nil {
		log.Println("Failed to reset text", err.Error())
	}
}

// todo call setup api
func testSetup() {

}
