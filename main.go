package main

import (
	// "encoding/xml"
	"fmt"
	"golang.org/x/text/width"
	"io"
	"log"
	"net/http"
	// "net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/beevik/etree"
)

const (
	/*
		BaseDir is where you want the resulting articles to go. The data will be in a structure like this (subdir names
		based on the consts defined):

			BaseDir/
				content/
					posts/
						Some-post-title/
							index.md
							ImagesDir/
								post-image-1.jpg
								post-image-2.png
						A-different-post-title/
							index.md
							ImagesDir/
								post-image-for-this-other-arcicle.jpg

	*/
	BaseDir    = "/home/yamadatt/git/hugoplate/content/english/blog"
	ContentDir = "content"
	PostsDir   = "posts"
	ImagesDir  = "images"

	// The path to the exported XML file containing all posts
	WordPressXMLFile = "./WordPress.2023-11-06.xml"

	// The path to the images export dir
	LocalImageDir = "/home/yamadatt/git/wordpress_to_hugo/"
)

func main() {

	var ConvertMD = false
	var firstImage string

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(WordPressXMLFile); err != nil {
		panic(err)
	}

	for _, item := range doc.FindElements("//item") {
		postType := item.SelectElement("wp:post_type")
		if postType.Text() == "attachment" {
			continue
		}

		title := strings.Replace(item.SelectElement("title").Text(), ":", "：", 1)
		content := item.SelectElement("content:encoded").Text()
		dateStr := item.SelectElement("pubDate").Text()
		tags := extractTags(item)
		postid := item.SelectElement("post_id").Text()
		// println("postid", postid)

		// postDir := filepath.Join(BaseDir, ContentDir, PostsDir, formatyyyymmdd(dateStr)+strings.ReplaceAll(title, " ", "-"))
		postDir := filepath.Join(BaseDir, ContentDir, PostsDir, formatyyyymmdd(dateStr)) + width.Fold.String(title)
		_ = os.MkdirAll(filepath.Join(postDir, ImagesDir), os.ModePerm)

		// Copy images
		// for i, imgURL := range extractImageURLs(content) {

		// 	// 画像ファイルパスを便利に使いたいため、Parseする。
		// 	parsedImgURL, _ := url.Parse(imgURL)
		// 	pathParts := strings.Split(parsedImgURL.Path, "/")

		// 	imgName := pathParts[len(pathParts)-1]
		// 	destPath := filepath.Join(postDir, ImagesDir, imgName)

		// 	img_file_name := DownloadImage(parsedImgURL.String(), destPath)

		// 	//HTMLに記述されているイメージパスを書き換える
		// 	content = strings.Replace(content, parsedImgURL.String(), filepath.Join(ImagesDir, img_file_name), -1)

		// 	// アイキャッチに設定するファイルのサイズがゼロだとhugoでエラーになるのでゼロバイトファイルをはじく
		// 	filesize, _ := FileSizeCheck(destPath)

		// 	if i == 1 && img_file_name != "" && filesize != 0 {
		// 		firstImage = filepath.Join(ImagesDir, img_file_name)
		// 	}

		// }

		frontMatter := fmt.Sprintf("---\ntitle: %s\nimage: \"%s\"\ndate: %s\nslug: %s\ndraft: false\ntags: [%s]\nsummary: \ncategory: \"\"\ntype: Post\n---\n", title, firstImage, formatDate(dateStr), postid, strings.Join(tags, ", "))

		//ConverMDがtrueかfalseでmarkdownに変換を決める
		var contentString string

		if ConvertMD == true {
			converter := md.NewConverter("", true, nil)
			markdown, err := converter.ConvertString(content)
			if err != nil {
				panic(err)
			}

			contentString = frontMatter + string(markdown)

		} else {
			contentString = frontMatter + string(content)
		}

		err := os.WriteFile(filepath.Join(postDir, "index.md"), []byte(contentString), 0644)
		if err != nil {
			panic(err)
		}

	}
}

func FileSizeCheck(filepath string) (int, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	size64 := info.Size()

	var size int
	if int64(int(size64)) == size64 {
		size = int(size64)
	}

	return size, nil
}

func DownloadImage(imgurl string, downloadpath string) string {

	response, e := http.Get(imgurl)
	if e != nil {
		log.Fatal(e)
	}

	_, filename := path.Split(imgurl)

	defer response.Body.Close()
	//open a file for writing

	file, err := os.Create(downloadpath)
	if err != nil {
		log.Fatal(err)
	}
	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	fmt.Printf("Download Success !  %s \n", downloadpath)
	return filename
}

func extractTags(item *etree.Element) []string {
	var tags []string

	for _, category := range item.SelectElements("category") {
		fmt.Println(category.Text())
		if domain := category.SelectAttrValue("domain", ""); domain != "post_tag" {
			tags = append(tags, category.Text())
		}
	}
	return tags
}

func extractImageURLs(content string) []string {

	stringReader := strings.NewReader(string(content))

	var urls []string

	doc, err := goquery.NewDocumentFromReader(stringReader)
	if err != nil {
		fmt.Print("html parse failed")
	}
	doc.Find("img").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("src")
		if url[:10] != "data:image" {
			urls = append(urls, url)
		}
	})
	return urls
}

func copyFile(src, dest string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	err = os.WriteFile(dest, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

func formatDate(date string) string {
	t, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		panic(err)
	}
	return t.Format(time.RFC3339)
}

func formatyyyymmdd(date string) string {
	t, err := time.Parse(time.RFC1123Z, date)
	if err != nil {
		panic(err)
	}
	return t.Format("20060102-")
}

func scrapeAmeblo(url string) {

