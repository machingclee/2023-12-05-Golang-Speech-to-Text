package azure

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	ffmpeg "github.com/u2takey/ffmpeg-go"
)

func createNewDir(dir string) error {
	_, err := os.Stat(dir)
	if err != nil {
		fmt.Printf("%s does not exist\n", dir)
		err := os.Mkdir(dir, os.ModePerm)
		if err != nil {
			fmt.Println("folder cannot be created")
		}
	}
	return nil
}

var audioFormatRegex = regexp.MustCompile(`\.mp3|\.m4a`)

func DownloadAndTranscode(fileUrl string) (tmpFilePath string, outFilePath string, err error) {
	split := strings.Split(fileUrl, "/")
	if len(split) == 0 {
		err = errors.New("invalid file url")
		return
	}
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	newtmpDir := path.Join(cwd, "tmp")
	newOutDir := path.Join(cwd, "out")

	err = createNewDir(newtmpDir)
	if err != nil {
		return
	}
	err = createNewDir(newOutDir)
	if err != nil {
		return
	}

	filename := split[len(split)-1]
	outFilename := audioFormatRegex.ReplaceAllString(filename, ".wav")
	tmpFilePath = path.Join(newtmpDir, filename)
	outFilePath = path.Join(newOutDir, outFilename)

	out, _ := os.Create(tmpFilePath)
	defer out.Close()
	resp, _ := http.Get(fileUrl)
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}

	err = ffmpeg.
		Input(tmpFilePath).
		Output(outFilePath, ffmpeg.KwArgs{
			"acodec": "pcm_s16le",
			"ac":     1,
			"ar":     16000,
		}).
		OverWriteOutput().ErrorToStdOut().Run()

	if err != nil {
		return
	}

	return
}
