package azure

import (
	"errors"
	"fmt"
	"time"

	"github.com/Microsoft/cognitive-services-speech-sdk-go/audio"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/common"
	"github.com/Microsoft/cognitive-services-speech-sdk-go/speech"
)

type result struct {
	msg   string
	ended bool
}

func RecognizeFromWavFile(filepath string) (results []string, err error) {
	subscription := "d66de8b96b0c499185743e15381ce23a"
	region := "eastus"

	audioConfig, err := audio.NewAudioConfigFromWavFileInput(filepath)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return nil, err
	}
	defer audioConfig.Close()
	config, err := speech.NewSpeechConfigFromSubscription(subscription, region)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return nil, err
	}
	defer config.Close()
	config.SetSpeechRecognitionLanguage("zh-HK")
	config.EnableDictation()
	config.SetOutputFormat(common.Detailed)
	speechRecognizer, err := speech.NewSpeechRecognizerFromConfig(config, audioConfig)
	if err != nil {
		fmt.Println("Got an error: ", err)
		return nil, err
	}
	defer speechRecognizer.Close()

	resultChannel := make(chan result)
	speechRecognizer.SessionStarted(func(event speech.SessionEventArgs) {
		defer event.Close()
		fmt.Println("Session Started (ID=", event.SessionID, ")")
	})
	speechRecognizer.Recognized(func(event speech.SpeechRecognitionEventArgs) {
		fmt.Println("RECOGNIZED", event.Result.Text)
		resultChannel <- result{msg: event.Result.Text, ended: false}
	})
	speechRecognizer.Canceled(func(event speech.SpeechRecognitionCanceledEventArgs) {
		fmt.Println("CANCELED")
		resultChannel <- result{msg: "", ended: true}
		speechRecognizer.StopContinuousRecognitionAsync()
	})
	speechRecognizer.SessionStopped(func(event speech.SessionEventArgs) {
		defer event.Close()
		speechRecognizer.StopContinuousRecognitionAsync()
		fmt.Println("Session Stopped (ID=", event.SessionID, ")")
	})
	speechRecognizer.StartContinuousRecognitionAsync()

RangeLoop:
	for {
		select {
		case outcome := <-resultChannel:
			{
				msg := outcome.msg
				ended := outcome.ended
				if msg != "" {
					results = append(results, msg)
				}
				if ended {
					break RangeLoop
				}
			}
		case <-time.After(5 * time.Second):
			fmt.Println("Timed out")
			err = errors.New("Time out, exceeding 5 seconds")
			break
		}
	}
	return
}
