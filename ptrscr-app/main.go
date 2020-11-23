package main

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"io/ioutil"
	"os"

	"github.com/MakeNowJust/hotkey"
	"github.com/atotto/clipboard"
	"github.com/bay0/ptrscr/ptrscr-app/gist"
	"github.com/bay0/ptrscr/ptrscr-app/icon"
	"github.com/bay0/ptrscr/ptrscr-app/logging"
	"github.com/bay0/ptrscr/ptrscr-app/registry"
	"github.com/bay0/ptrscr/ptrscr-app/utils"
	"github.com/disintegration/imaging"
	"github.com/getlantern/systray"
	"github.com/google/go-github/github"
	"github.com/kbinani/screenshot"
	log "github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"golang.org/x/oauth2"
	"gopkg.in/toast.v1"
)

type Image struct {
	Bounds image.Rectangle
	Data   *image.RGBA
}

func init() {
	logging.Init()
}

func main() {
	onExit := func() {
		log.Println("Closing app")
	}

	systray.Run(onReady, onExit)
}

func onReady() {
	//check registry for accessToken
	accessToken := registry.GetStringFromLocalMachine(`Software\PTRSCR`, "token")

	if accessToken == "" {
		log.Println("Token doesnt exist")
		dialog.Message("%s", "Please add your gist token").Title("PTRSCR").Info()
		filename, err := dialog.File().Filter("Token file", "txt").Title("PTRSCR").Load()
		if err != nil {
			log.Println(err)
		}

		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(content))
		registry.CreateRegistryKey(`PTRSCR`, "token", string(content))
		accessToken = registry.GetStringFromLocalMachine(`Software\PTRSCR`, "token")
	}

	//github
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	imgBytes, err := utils.GetImageBytesFromURL("https://raw.githubusercontent.com/bay0/ptrscr/master/ptrscr-app/icon/icon.png")
	if err != nil {
		log.Fatal(err)
	}
	// Decode image back from the byte slice.
	watermark, err := imaging.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		log.Fatal(err)
	}
	watermark = imaging.Resize(watermark, 50, 50, imaging.Lanczos)
	watermarkPadding := 5

	isAppEnabled := true

	//systray here
	systray.SetTemplateIcon(icon.Data, icon.Data)
	systray.SetTitle("PTRSCR")
	systray.SetTooltip("PTRSCR")
	mQuitOrig := systray.AddMenuItem("Quit", "Quit the whole app")
	go func() {
		<-mQuitOrig.ClickedCh
		log.Println("Requesting quit")
		systray.Quit()
		log.Println("Finished quitting")
	}()

	go func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("PTRSCR")
		systray.SetTooltip("PTRSCR")

		listenerEnabled := systray.AddMenuItemCheckbox("Enabled", "Check Me", true)

		for {
			select {
			case <-listenerEnabled.ClickedCh:
				if listenerEnabled.Checked() {
					listenerEnabled.Uncheck()
					listenerEnabled.SetTitle("Disabled")
					isAppEnabled = false
				} else {
					listenerEnabled.Check()
					listenerEnabled.SetTitle("Enabled")
					isAppEnabled = true
				}
			}
		}
	}()

	hkey := hotkey.New()
	hkey.Register(hotkey.None, hotkey.F10, func() {
		if isAppEnabled {
			n := screenshot.NumActiveDisplays()

			var images []Image

			for i := 0; i < n; i++ {
				bounds := screenshot.GetDisplayBounds(i)

				img, err := screenshot.CaptureRect(bounds)
				if err != nil {
					panic(err)
				}

				images = append(images, Image{
					Bounds: bounds,
					Data:   img,
				})
			}

			maxWidth := 0
			maxHeight := 0

			for _, img := range images {
				maxWidth += img.Bounds.Dx()

				if maxHeight < img.Bounds.Dy() {
					maxHeight = img.Bounds.Dy()
				}
			}

			dst := imaging.New(maxWidth, maxHeight, color.NRGBA{255, 255, 255, 255})

			for _, img := range images {
				dst = imaging.Paste(dst, img.Data, image.Pt(img.Bounds.Min.X, img.Bounds.Min.Y))
				//watermark
				dst = imaging.Overlay(dst, watermark, image.Pt(img.Bounds.Min.X+watermarkPadding, img.Bounds.Min.Y+watermarkPadding), 0.8)
			}

			var buf bytes.Buffer
			err = imaging.Encode(&buf, dst, imaging.PNG)
			if err != nil {
				log.Fatal(err)
			}
			pngBytes := buf.Bytes()

			res, err := gist.Create(client, pngBytes, utils.BuildFileName())
			if err == nil {
				finalURL := "https://ptrscr.dev/?id=" + res.GetID()

				log.Printf("Uploaded: %s", finalURL)

				clipboard.WriteAll(finalURL)
				notification := toast.Notification{
					AppID:   "ptrscr-app-snap",
					Title:   "Upload successful",
					Message: "Open in browser",
					Actions: []toast.Action{
						{"protocol", "Open", finalURL},
					},
				}
				err := notification.Push()
				if err != nil {
					log.Fatalln(err)
				}
			} else {
				log.Info(err)
				os.Exit(0)
			}
		}
	})
}
