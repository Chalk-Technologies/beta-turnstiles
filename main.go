package main

import (
	"beta-turnstiles/internal/checkCode"
	"beta-turnstiles/internal/config"
	"beta-turnstiles/internal/relay"
	"flag"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"log"
	"strconv"
	"strings"
)

func main() {
	// flags needed for accounting stuff
	// get wise api key
	// get desired month
	var headless = flag.Bool("headless", false, "run in background without UI")

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = relay.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer relay.CleanUp()
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
		container.NewHBox(
			header,
			widget.NewButton("settings", func() { editConfig(w) }),
		),
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

func editConfig(w fyne.Window) {
	var modal *widget.PopUp
	demoMode := binding.NewBool()
	singleMode := binding.NewBool()
	dirOut := binding.NewBool()
	highMode := binding.NewBool()
	apiKey := binding.NewString()
	relayPin := config.GlobalConfig.RelayPin

	demoMode.Set(config.GlobalConfig.DemoMode)
	singleMode.Set(config.GlobalConfig.SingleMode)
	dirOut.Set(config.GlobalConfig.DirectionOut)
	highMode.Set(config.GlobalConfig.HighMode)
	//relayPin.Set(config.GlobalConfig.RelayPin)
	if config.GlobalConfig.ApiKey != nil {
		apiKey.Set(*config.GlobalConfig.ApiKey)
	}
	pinSelectEntry := widget.NewSelectEntry([]string{"17", "27", "22", "4", "18", "23", "24", "25", "5", "6"})
	pinSelectEntry.OnChanged = func(s string) {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Println(err.Error())
		} else {
			relayPin = i
		}
	}
	submit := func() {
		dm, err := demoMode.Get()
		if err != nil {
			log.Println(err.Error())
			return
		}
		sm, err := singleMode.Get()
		if err != nil {
			log.Println(err.Error())
			return
		}

		do, err := dirOut.Get()
		if err != nil {
			log.Println(err.Error())
			return
		}

		ak, err := apiKey.Get()
		if err != nil {
			log.Println(err.Error())
			return
		}

		hm, err := highMode.Get()
		if err != nil {
			log.Println(err.Error())
			return
		}

		newConfig := config.Config{
			DemoMode:     dm,
			SingleMode:   sm,
			DirectionOut: do,
			ApiKey:       &ak,
			RelayPin:     relayPin,
			HighMode:     hm,
		}
		err = config.StoreConfig(newConfig)
		if err != nil {
			log.Println(err.Error())
		} else {
			modal.Hide()
		}

	}
	modal = widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel("Settings"),
			widget.NewCheckWithData("Demo mode", demoMode),
			widget.NewCheckWithData("Single mode", singleMode),
			widget.NewCheckWithData("Direction out", dirOut),
			widget.NewCheckWithData("io pin high mode", highMode),
			pinSelectEntry,
			widget.NewEntryWithData(apiKey),
			container.NewHBox(
				widget.NewButton("cancel", func() { modal.Hide() }),
				widget.NewButton("submit", submit),
			),
		),
		w.Canvas(),
	)
	modal.Show()
	return
}
