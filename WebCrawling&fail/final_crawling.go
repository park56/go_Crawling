// 최종 실습 예제

// http - https 변화로 인한 웹 클롤링 실패
// 대상 사이트 : 루리웹(ruluweb.com)

package main

import (
	"bufio"
	_ "bufio"
	"fmt"
	"net/http"
	"os"
	_ "os"
	"strings"
	"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// https://github.com/yhat/scrape  -> 스크랩, 크롤링 패키지, 사용하기 어려움
// http://go-colly.org/docs/  -> 스크앱,크롤링 라이브러리, goquery 기반 강력하고 쉬운 패키지(가장 많이 사용)
// https://github.com/PuerkitoBio/goquery -> 쉬운 HTML Pasing 지원

// 스크랩 대상 url
const (
	//	urlRoot = "https://www.fmkorea.com/"
	urlRoot = "https://bbs.ruliweb.com/"
)

// 첫번째 방문(메인페이지) 대상으로 원하는 url 파싱 후 반환하는 함수
/*func parseMainNodes(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n.Parent.Parent, "class") == "a1 sub"
	}
	return false
}*/

func parseMainNodes(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n.Parent.Parent, "class") == "list dot"
	}
	return false
}

// 에러체크
func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// 동기화를 위한 작업 그룹 생성
var wg sync.WaitGroup

//	Url대상이 되는 페이지(서브페이지) 대상으로 원하는 내용을 파싱 후 반환
func ScrapeContents(url string, fn string) {
	// 작업 종료 알림
	defer wg.Done()

	// Get 방식 요청
	resp, err := http.Get(url)
	errCheck(err)

	// 요청 Body 닫기
	defer resp.Body.Close()

	// 응답 데이터(Html)
	root, err := html.Parse(resp.Body)
	errCheck(err)

	// Response 데이터(html)의 원하는 부분 파싱
	matchNode := func(n *html.Node) bool {
		return n.DataAtom == atom.A && scrape.Attr(n, "class") == "deco"
	}

	// 파일 스크림 생성(열기) -> 파일명, 옵션, 권한
	file, err := os.OpenFile(`C:\Users\DGSW\go_study\src\StudyGo\WebCrawling&fail\scrape\`+fn+".txt", os.O_CREATE|os.O_RDWR, os.FileMode(777))
	errCheck(err)

	// 메소드 종료 시 파일 닫기
	defer file.Close()

	// 쓰기 버퍼 생성
	w := bufio.NewWriter(file)

	// matchNode 메소드를 사용해서 원하는 노드 순회(Iterator) 하며 출력
	for _, g := range scrape.FindAll(root, matchNode) {
		// Url 및 해당 데이터 출력
		fmt.Println("result : ", scrape.Text(g))
		// 파싱 데이터 -> 버퍼에 기록
		w.WriteString(scrape.Text(g) + "\r\n")
	}
	w.Flush()
}

func main() {

	// 메인 페이지 Get 방식 요청
	resp, err := http.Get(urlRoot) // response : 응답, request : 요청
	errCheck(err)

	// 요청 Body 닫기
	defer resp.Body.Close()

	// 응답 데이터(HTML)
	root, err := html.Parse(resp.Body)
	errCheck(err)
	//body, err := ioutil.ReadAll(res.Body) // res.body의 값 담기

	//ParseMainNodes 메소드를 크롤링(스크래핑) 대상 URL 추출
	urlList := scrape.FindAll(root, parseMainNodes)

	for _, link := range urlList {
		// 대상 url 1차 출략
		//fmt.Println("Check Main Link : ", link, idx)
		//fmt.Println("TargetUrl : ", scrape.Attr(link, "href"))
		fileName := strings.Replace(scrape.Attr(link, "href"), "https://bbs.ruliweb.com/family/", "", 1)
		//fmt.Println(strings.Replace("golang golang golang golang", "ng", "nd", 3))    // old를 new로 n개만큼 replace
		//fmt.Println("fileName : ", fileName)

		// 작업 대기열 추가
		wg.Add(1) // Done 개수와 일치
		// 고루틴 시작 -> 작업 대기열 개수와 같아야 함
		go ScrapeContents(scrape.Attr(link, "href"), fileName)
	}
	// 모든 작업 종료까지 대기
	wg.Wait()
}
