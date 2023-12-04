package azure

import (
	"fmt"
	"testing"
)

func TestTranscode(t *testing.T) {
	fileUrl := "https://audioasr.blob.core.windows.net/dev-voices/voice_james.lee@wonderbricks.com_day_2023-11-18_hktime_18-59-58.m4a"
	tmpPath, outPath, err := DownloadAndTranscode(fileUrl)
	fmt.Println(tmpPath)
	fmt.Println(outPath)
	if err != nil {
		fmt.Println("[Error]:", err)
	}
}
