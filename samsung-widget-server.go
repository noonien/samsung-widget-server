package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
)

var widgetDir = flag.String("widgetDir", "Package",
	"The path under which the widgets are found")

type compression struct {
	Size int64  `xml:"size,attr"`
	Type string `xml:"type,attr"`
}

type widget struct {
	Id          string      `xml:"id,attr"`
	Title       string      `xml:"title"`
	Compression compression `xml:"compression"`
	Description string      `xml:"description"`
	Download    string      `xml:"download"`
}

type response struct {
	XMLName xml.Name `xml:"rsp"`
	Stat    string   `xml:"stat,attr"`
	List    []widget `xml:"list>widget"`
}

func WidgetListing(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	log.Printf("Widget listing requested from %s", ip)

	wgFiles := make(map[string]os.FileInfo)

	files, err := ioutil.ReadDir(*widgetDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		nameComp := strings.Split(file.Name(), "_")
		name := strings.Join(nameComp[:len(nameComp)-3], "_")

		ofile, ok := wgFiles[name]
		if !ok || file.ModTime().After(ofile.ModTime()) {
			wgFiles[name] = file
		}
	}

	log.Printf("Available widgets are: ")
	var widgets []widget
	for name, file := range wgFiles {
		log.Printf(" >> %s", name)
		widgets = append(widgets, widget{
			Id:          name,
			Title:       name,
			Compression: compression{Size: file.Size(), Type: "zip"},
			Description: "",
			Download: fmt.Sprintf("http://%s/w/%s",
				r.Host, file.Name()),
		})
	}

	w.Header().Set("Content-Type", "text/xml")
	io.WriteString(w, xml.Header)
	xmlEnc := xml.NewEncoder(w)
	xmlEnc.Encode(response{Stat: "ok", List: widgets})
}

type WidgetHandler struct {
	fileHandler http.Handler
}

func (ws *WidgetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	widgetFile := path.Base(r.URL.Path)

	log.Printf("Serving widget file %s to %s", widgetFile, ip)

	r.URL.Path = "/" + widgetFile
	ws.fileHandler.ServeHTTP(w, r)
	log.Printf("Done")
}

func GetWidgetHandler() *WidgetHandler {
	return &WidgetHandler{http.FileServer(http.Dir(*widgetDir))}
}

func main() {
	flag.Parse()
	log.Print("Starting Samsung Widget Server")
	log.Printf("Serving widgets from %s", *widgetDir)

	http.HandleFunc("/widgetlist.xml", WidgetListing)
	http.Handle("/w/", GetWidgetHandler())
	log.Fatal(http.ListenAndServe(":80", nil))
}
