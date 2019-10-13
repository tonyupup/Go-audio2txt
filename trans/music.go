package trans

/*
#cgo LDFLAGS: -Wl,-rpath="${SRCDIR}/lib"
#cgo LDFLAGS: -L./lib -lvad
int split(const char *,const char *);
#include "stdlib.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"unsafe"
)

type TransTask struct {
	ID             int32
	src, dest, Out string
}

//NewTransTask create new TransTask obj.
func NewTransTask(id int32, src string) (*TransTask, error) {
	//check file is exists.
	if fi, err := os.Stat(src); err != nil {
		if os.IsNotExist(err) {
			return nil, os.ErrExist
		}
		log.Println("Name: ", fi.Name(), "\tSize: ", fi.Size())
		log.Fatalln(err.Error())
	}
	_, fn := filepath.Split(src)
	dn := filepath.Join(os.TempDir(), fn)
	os.Mkdir(dn, os.ModeDir|os.ModePerm)
	return &TransTask{id, src, dn, filepath.Join(dn, "out")}, nil
}

//Audio2PCM16k convent audio to pcm and rate is 16k.
func (t *TransTask) Audio2PCM16k() (string, error) {
	outfile := filepath.Join(t.dest, "x.pcm")
	if _, err := os.Stat(outfile); err != nil {
		if os.IsNotExist(err) {
			x := exec.Command("ffmpeg", "-i", t.src, "-threads", strconv.Itoa(runtime.NumCPU()), "-f", "s16le", "-ar", "16000", "-ac", "1", "-acodec", "pcm_s16le", outfile)
			o, e := x.CombinedOutput()
			result := string(o)
			if e != nil {
				panic(e)
			}
			return result, nil
		}
	}
	return "ok", nil
}

//Split split pcm to short part
func (t *TransTask) Split() (fils []string, e error) {
	src := C.CString(filepath.Join(t.dest, "x.pcm"))
	dest := C.CString(t.Out)
	defer C.free(unsafe.Pointer(src))
	defer C.free(unsafe.Pointer(dest))
	os.Mkdir(t.Out, os.ModeDir|os.ModePerm)
	if x := int(C.split(src, dest)); x != 0 {
		return nil, fmt.Errorf("Some Error happend in C.split ,code :%d", x)
	}
	if rd, err := ioutil.ReadDir(t.Out); err != nil {
		return nil, err
	} else {
		for _, file := range rd {
			fils = append(fils, file.Name())
		}
	}
	return
}

//Free delete temp files
func (t *TransTask) Free() error {
	return os.RemoveAll(t.dest)
}
