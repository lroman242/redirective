package main

import (
	"encoding/json"
	"fmt"
	"github.com/lroman242/redirective/tracer"
	"github.com/raff/godet"
	"io/ioutil"
	"log"
	"net/url"
)

func main() {
	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		fmt.Println("cannot connect to Chrome instance:", err)
		return
	}
	chr := tracer.NewChromeTracer(remote)

	targetUrl, _ := url.Parse("https://www.google.com.ua")
	//targetUrl, _ := url.Parse("http://compre.vc/v2/331179766c6")
	//targetUrl, _ := url.Parse("http://trk.indoleads.com/5a435e1bb6920d3422293737")
	//targetUrl, _ := url.Parse("http://trk.indoleads.com/5ad61882010d9")
	redirects, err := chr.GetTrace(targetUrl)
	if err != nil {
		log.Fatalln(err)
	}

	b, err := json.Marshal(redirects)
	if err != nil {
		fmt.Println(err)
		return
	}
	ioutil.WriteFile("./result.json", b, 0644)
	//
	//var mainFrameId string
	//
	//// connect to Chrome instance
	//remote, err := godet.Connect("localhost:9222", false)
	//if err != nil {
	//	fmt.Println("cannot connect to Chrome instance:", err)
	//	return
	//}
	//
	//// disconnect when done
	//defer remote.Close()
	//
	//remote.EnableRequestInterception(true)
	//
	//// get list of open tabs
	//tabs, _ := remote.TabList("")
	//fmt.Println(tabs)
	//
	//// install some callbacks
	//remote.CallbackEvent(godet.EventClosed, func(params godet.Params) {
	//	fmt.Println("RemoteDebugger connection terminated.")
	//})
	//
	//i := 0
	//remote.CallbackEvent("Network.requestWillBeSent", func(params godet.Params) {
	//	if _, ok := params["redirectResponse"]; ok && params["type"] == "Document" {
	//		fmt.Printf("Document URL: %s\n", params["documentURL"])
	//
	//		b, err := json.Marshal(params)
	//		if err != nil {
	//			fmt.Println(err)
	//			return
	//		}
	//		ioutil.WriteFile(fmt.Sprintf("./%d-%s.json", i, mainFrameId), b, 0644)
	//		i++
	//	}
	//})
	//
	////remote.CallbackEvent("Network.requestIntercepted", func(params godet.Params) {
	////	fmt.Println("Network.requestIntercepted")
	////	fmt.Printf("%v+\n", params)
	////	//fmt.Println("requestWillBeSent",
	////	//	params["type"],
	////	//	params["documentURL"],
	////	//	params["request"].(map[string]interface{})["url"])
	////})
	////
	////remote.CallbackEvent("Network.responseReceived", func(params godet.Params) {
	////fmt.Println("Network.requestIntercepted")
	////fmt.Printf("%v+\n", params)
	////fmt.Println("responseReceived",
	////	params["type"],
	////	params["response"].(map[string]interface{})["url"])
	////})
	//
	////remote.CallbackEvent("Log.entryAdded", func(params godet.Params) {
	////	entry := params["entry"].(map[string]interface{})
	////	fmt.Println("LOG", entry["type"], entry["level"], entry["text"])
	////})
	//
	//// block loading of most images
	////_ = remote.SetBlockedURLs("*.jpg", "*.png", "*.gif")
	//
	//// create new tab
	//tab, _ := remote.NewTab("https://www.google.com")
	//defer func() {
	//	remote.CloseTab(tab)
	//}()
	//
	//// enable event processing
	//remote.RuntimeEvents(true)
	//remote.NetworkEvents(true)
	//remote.PageEvents(true)
	//remote.DOMEvents(true)
	//remote.LogEvents(true)
	//
	//// navigate in existing tab
	//_ = remote.ActivateTab(tabs[0])
	//
	////remote.StartPreciseCoverage(true, true)
	//
	//// re-enable events when changing active tab
	//remote.AllEvents(true) // enable all events
	//
	////remote.Navigate("https://www.google.com")
	////mainFrameId, _ = remote.Navigate("http://compre.vc/v2/331179766c6")
	////remote.Navigate("http://trk.indoleads.com/5a435e1bb6920d3422293737")
	//
	//mainFrameId, _ = remote.Navigate("http://trk.indoleads.com/5ad61882010d9")
	//fmt.Printf("main frame id - %s\n", mainFrameId)
	//
	//// evaluate Javascript expression in existing context
	////res, _ := remote.EvaluateWrap(`
	////        console.log("hello from godet!")
	////        return 42;
	////    `)
	////fmt.Println(res)
	//time.Sleep(5 * time.Second)
	//
	//// take a screenshot
	//_ = remote.SaveScreenshot(fmt.Sprintf("%s.png", mainFrameId), 0644, 100, true)
	//
	////time.Sleep(time.Second)
	//
	//// or save page as PDF
	////_ = remote.SavePDF("page.pdf", 0644, godet.PortraitMode(), godet.Scale(0.5), godet.Dimensions(6.0, 2.0))
	//
	//// if err := remote.SetInputFiles(0, []string{"hello.txt"}); err != nil {
	////     fmt.Println("setInputFiles", err)
	//// }
	//
	////time.Sleep(5 * time.Second)
	//
	////remote.StopPreciseCoverage()
	//
	////r, err := remote.GetPreciseCoverage(true)
	////if err != nil {
	////	fmt.Println("error profiling", err)
	////} else {
	////	fmt.Println(r)
	////}
	//
	////// Allow downloads
	////_ = remote.SetDownloadBehavior(godet.AllowDownload, "/tmp/")
	////_, _ = remote.Navigate("http://httpbin.org/response-headers?Content-Type=text/plain;%20charset=UTF-8&Content-Disposition=attachment;%20filename%3d%22test.jnlp%22")
	////
	////time.Sleep(time.Second)
	////
	////// Block downloads
	////_ = remote.SetDownloadBehavior(godet.DenyDownload, "")
	////_, _ = remote.Navigate("http://httpbin.org/response-headers?Content-Type=text/plain;%20charset=UTF-8&Content-Disposition=attachment;%20filename%3d%22test.jnlp%22")

}
