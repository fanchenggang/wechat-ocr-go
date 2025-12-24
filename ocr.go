package main

/*
#cgo LDFLAGS: -L. -lwcocr
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"sync"
	"syscall"
	"unsafe"
)

const ocrExe = "path\\WeChatOCR\\WeChatOCR.exe"
const wechatDir = "path"

// WechatOCR 定义了 DLL 中的函数和回调接口
type WechatOCR struct {
	dll         *syscall.LazyDLL
	wechat_ocr  *syscall.LazyProc
	stop_ocr    *syscall.LazyProc
	callbackPtr uintptr
	mu          sync.Mutex
	resultChan  chan string
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
	stop_ocr := dll.NewProc("stop_ocr")
	if stop_ocr == nil {
		return nil, fmt.Errorf("failed to find function 'stop_ocr' in DLL: %s", dllPath)
	}

	// 创建单例回调函数
	ocr := &WechatOCR{
		dll:        dll,
		wechat_ocr: wechat_ocr,
		stop_ocr:   stop_ocr,
		resultChan: make(chan string, 1),
	}

	// 初始化回调函数（只创建一次）
	ocr.callbackPtr = syscall.NewCallback(func(result *C.char) uintptr {
		ocr.resultChan <- C.GoString(result)
		return 0
	})

	return ocr, nil
}

// CallWechatOCR 调用 DLL 中的 wechat_ocr 函数
func (w *WechatOCR) CallWechatOCR(ocrExe, wechatDir, imgFn string) error {
	ocrExeWStr, err := syscall.UTF16PtrFromString(ocrExe)
	if err != nil {
		return fmt.Errorf("failed to convert ocrExe to UTF16: %v", err)
	}

	wechatDirWStr, err := syscall.UTF16PtrFromString(wechatDir)
	if err != nil {
		return fmt.Errorf("failed to convert wechatDir to UTF16: %v", err)
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	ret, _, err := w.wechat_ocr.Call(
		uintptr(unsafe.Pointer(ocrExeWStr)),
		uintptr(unsafe.Pointer(wechatDirWStr)),
		uintptr(unsafe.Pointer(syscall.StringBytePtr(imgFn))),
		w.callbackPtr, // 使用缓存的回调函数指针
	)
	if ret == 0 {
		return err
	}
	return nil
}

func OcrCustom(wechatOCR *WechatOCR, ocrExe, wechatDir, imgFn string) *Result {
	return ocr(wechatOCR, ocrExe, wechatDir, imgFn)
}

func OcrDefault(wechatOCR *WechatOCR, imgFn string) *Result {
	return ocr(wechatOCR, ocrExe, wechatDir, imgFn)
}

func ocr(wechatOCR *WechatOCR, ocrExe, wechatDir, imgFn string) *Result {
	err := wechatOCR.CallWechatOCR(ocrExe, wechatDir, imgFn)
	if err != nil {
		fmt.Printf("Failed to call wechat_ocr: %v\n", err)
		return nil
	}

	// 等待回调结果
	resp := <-wechatOCR.resultChan
	r := &Result{}
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
