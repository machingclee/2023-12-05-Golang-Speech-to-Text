package azure

import (
	"fmt"
	"log"
	"testing"
)

func TestConfigizeFromWavFile(t *testing.T) {
	file := "/workspaces/2023-11-18-wb-golang-azure-speech/transcode/out/voice_james.lee@wonderbricks.com_day_2023-11-18_hktime_18-59-58.wav"
	transcriptionResult, err := RecognizeFromWavFile(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[Transcription Result]", transcriptionResult)
}
