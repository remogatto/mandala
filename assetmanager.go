package gorgasm

import (
	"git.tideland.biz/goas/loop"
	"io"
	"unsafe"
)

var (
	AssetsPath string = "android"
)

type LoadAssetResponse struct {
	Buffer io.Reader
	Error  error
}

type LoadAssetRequest struct {
	Filename string
	Response chan LoadAssetResponse
}

// Run runs the nativeEventsLoop.
// The loop handles native input events.
func assetLoopFunc(activity chan unsafe.Pointer, request chan interface{}) loop.LoopFunc {
	var act unsafe.Pointer
	return func(l loop.Loop) error {
		for {
			select {
			case act = <-activity:
			case untypedRequest := <-request:
				switch request := untypedRequest.(type) {
				case LoadAssetRequest:
					file, err := loadAsset(act, request.Filename)
					request.Response <- LoadAssetResponse{file, err}
				}
			}
		}
	}
}

// func LoadAsset(filename string) <-chan io.Reader {
// 	command := LoadAssetCommand{
// 		Filename: filename,
// 		Buffer:   make(chan io.Reader),
// 	}
// 	Assets <- command
// 	return command.Buffer
// }
