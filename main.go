package main

import "fmt"

func main() {
	//项目自带的
	result := OcrDefault("test.png")
	fmt.Println(result)

	//使用自己电脑上的wechat
	ocrExe := "C:\\Users\\Administrator\\AppData\\Roaming\\Tencent\\WeChat\\XPlugin\\Plugins\\WeChatOCR\\7079\\extracted\\WeChatOCR.exe"
	wechatDir := "D:\\SOFTWARE\\Tencent\\WeChat\\[3.9.11.25]"
	result2 := OcrCustom(ocrExe, wechatDir, "test.png")
	fmt.Println(result2)
}
