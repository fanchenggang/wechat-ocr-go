//go:build cgo
// +build cgo

package main

import (
	"fmt"
)

func main() {
	wechatOCR, err := NewWechatOCR("wcocr.dll")
	if err != nil {
		fmt.Printf("Failed to create WechatOCR instance: %v\n", err)
		return
	}
	//项目自带的
	for i := 0; i < 100; i++ {
		result := OcrDefault(wechatOCR, "test.png")
		fmt.Println(result)
	}

	wechatOCR.stop_ocr.Call()

	//使用自己电脑上的wechat
	//ocrExe := "C:\\Users\\Administrator\\AppData\\Roaming\\Tencent\\WeChat\\XPlugin\\Plugins\\WeChatOCR\\7079\\extracted\\WeChatOCR.exe"
	//wechatDir := "D:\\SOFTWARE\\Tencent\\WeChat\\[3.9.11.25]"
	//result2 := OcrCustom(ocrExe, wechatDir, "test.png")
	//fmt.Println(result2)
}
