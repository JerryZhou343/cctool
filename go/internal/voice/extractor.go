package voice

import (
	"github.com/JerryZhou343/cctool/go/internal/status"
	"os"
	"os/exec"
	"path/filepath"
)

type Extractor struct {
	sampleRate string
	ffmpage    string
}

func NewExtractor(sampleRate, ffmpeg string) *Extractor {
	ret := &Extractor{
		sampleRate: sampleRate,
		ffmpage:    ffmpeg,
	}
	if ffmpeg != "" {
		ret.AddEnv(ffmpeg)
	}
	return ret
}

func (e *Extractor) ExtractAudio(src string, dst string) (err error) {
	if err = e.Valid(); err != nil {
		return err
	}
	cmd := exec.Command("ffmpeg", "-i", src, "-ar", e.sampleRate, dst)
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (e *Extractor) Valid() (err error) {
	ts := exec.Command("ffmpeg", "-version")
	//ts.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	if err = ts.Run(); err != nil {
		return status.ErrFFmpegeCheckFailed
	}
	return nil
}

func (e *Extractor) AddEnv(path string) (err error) {
	var (
		absPath string
	)
	absPath, err = filepath.Abs(path)
	if err != nil {
		return status.ErrFileNotExits
	}
	_, err = os.Stat(absPath)
	if err != nil {
		return
	}
	if os.IsNotExist(err) {
		err = status.ErrFileNotExits
		return
	}
	os.Setenv("PATH", absPath)
	return
}
