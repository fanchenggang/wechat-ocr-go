package main

/*
#cgo LDFLAGS: -L. -lwcocr
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"syscall"
	"unsafe"
)

const ocrExe = "path\\WeChatOCR\\WeChatOCR.exe"
const wechatDir = "path"

// WechatOCR 定义了 DLL 中的函数和回调接口
type WechatOCR struct {
	dll        *syscall.LazyDLL
	wechat_ocr *syscall.LazyProc
}

// SetResCallback 定义了回调函数的接口
type SetResCallback func(result string)

// NewWechatOCR 创建一个新的 WechatOCR 实例
func NewWechatOCR(dllPath string) (*WechatOCR, error) {
	dll := syscall.NewLazyDLL(dllPath)
	if dll == nil {
		return nil, fmt.Errorf("failed to load DLL: %s", dllPath)
	}

	wechat_ocr := dll.NewProc("wechat_ocr")
	if wechat_ocr == nil {
		return nil, fmt.Errorf("failed to find function 'wechat_ocr' in DLL: %s", dllPath)
	}

	return &WechatOCR{
		dll:        dll,
		wechat_ocr: wechat_ocr,
	}, nil
}

// CallWechatOCR 调用 DLL 中的 wechat_ocr 函数
func (w *WechatOCR) CallWechatOCR(ocrExe, wechatDir, imgFn string, callback SetResCallback) error {
	ocrExeWStr, err := syscall.UTF16PtrFromString(ocrExe)
	if err != nil {
		return fmt.Errorf("failed to convert ocrExe to UTF16: %v", err)
	}

	wechatDirWStr, err := syscall.UTF16PtrFromString(wechatDir)
	if err != nil {
		return fmt.Errorf("failed to convert wechatDir to UTF16: %v", err)
	}

	callbackPtr := syscall.NewCallback(func(result *C.char) uintptr {
		callback(C.GoString(result))
		return 0
	})

	ret, _, err := w.wechat_ocr.Call(
		uintptr(unsafe.Pointer(ocrExeWStr)),
		uintptr(unsafe.Pointer(wechatDirWStr)),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(imgFn))),
		uintptr(callbackPtr),
	)
	if ret == 0 {
		return err
	}
	return nil
}
func OcrCustom(wechatOCR *WechatOCR, ocrExe, wechatDir, imgFn string) *Result {
	//ocrExe := "C:\\Users\\Administrator\\AppData\\Roaming\\Tencent\\WeChat\\XPlugin\\Plugins\\WeChatOCR\\7079\\extracted\\WeChatOCR.exe"
	//wechatDir := "D:\\SOFTWARE\\Tencent\\WeChat\\[3.9.11.25]"
	return ocr(wechatOCR, ocrExe, wechatDir, imgFn)
}
func OcrDefault(wechatOCR *WechatOCR, imgFn string) *Result {
	return ocr(wechatOCR, ocrExe, wechatDir, imgFn)
}
func ocr(wechatOCR *WechatOCR, ocrExe, wechatDir, imgFn string) *Result {

	result := make(chan string, 1)
	defer close(result)

	err := wechatOCR.CallWechatOCR(ocrExe, wechatDir, imgFn, func(res string) {
		result <- res
	})
	if err != nil {
		fmt.Printf("Failed to call wechat_ocr: %v\n", err)
		return nil
	}
	r := &Result{}
	resp, _ := <-result
	json.Unmarshal([]byte(resp), &r)
	return r

}

type Result struct {
	Errcode     int         `json:"errcode"`
	OcrResponse []OcrResult `json:"ocr_response"`
}
type OcrResult struct {
	Text string `json:"text"`
}
