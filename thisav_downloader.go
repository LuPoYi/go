package main

import (
  "fmt"
  "log"

  "github.com/PuerkitoBio/goquery"
  "gopkg.in/cheggaaa/pb.v1"
  "io/ioutil"
  "net/http"
  "regexp"
  "io"
  "os"
)

const targetURL  = "http://www.thisav.com/videos?o=tf"

func StartCrawler() {
  doc, err := goquery.NewDocument(targetURL) 
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("#content .video_box").Each(func(i int, s *goquery.Selection) {
    video_url := ""
    video_page_url, _ := s.Find("a").Attr("href")
    img := s.Find("img")
    title, _ := img.Attr("alt")
    full_title := title + ".flv"
    video_url, err = getThisavVideoUrl(video_page_url)
    downloadVideoFromLink(video_url, full_title)
  })
}

func getThisavVideoUrl(video_page_url string) (video_url string, err error) {
  response, err := http.Get(video_page_url)
  if err != nil {
    return
  } else {
    defer response.Body.Close()
    contents, err := ioutil.ReadAll(response.Body)
    var video_url = regexp.MustCompile(`so.addVariable\('file','(.*\.flv)'\);`)
    if video_url.Match(contents) {
      return video_url.FindStringSubmatch(string(contents))[1], err
    }
  }
  return
}

func downloadVideoFromLink(downloadLink string, filename string) {
  res, err := http.Get(downloadLink)
  if err != nil {
    log.Fatal(err)
  }
  defer res.Body.Close()

  if _, err := os.Stat(filename); os.IsNotExist(err) {
    out, err := os.Create(filename)

    if err != nil {
      log.Fatal(err)
    }
    defer out.Close()

    fmt.Println("Downloading ", filename, downloadLink)

    filesize := int(res.ContentLength)
    bar := pb.New(filesize).SetUnits(pb.U_BYTES)
    bar.Start()
    rd := bar.NewProxyReader(res.Body)
    io.Copy(out, rd)

    fmt.Println("Download OK ", filename)
  } else {
    fmt.Println(filename, " already exist.")
  }
}

func main() {
  fmt.Printf("Start!\n")
  StartCrawler()
  fmt.Printf("End!\n")
}
