package migadu

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

//func (c *Client) GetCookie(ctx context.Context) http.Cookie

type Dashboard map[string][]*UsageStats

type UsageStats struct {
	Date      time.Time
	Domain    string
	Receiving uint64
	Sending   uint64
	Storage   uint64
}

const (
	AdminHost = "https://admin.migadu.com"
)

func (c *Client) GetDomains(ctx context.Context) ([]string, error) {
	req, err := c.GetV1ReqBuilder().
		SetMethod(http.MethodGet).
		AddPath(MailboxesPath).
		Build()
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get Cookies status code %d", resp.StatusCode)
	}
	c.Cookies = resp.Cookies()

	req, err = NewReqBuilder().
		SetMethod(http.MethodGet).
		SetHost(AdminHost).
		AddPath(DomainsPath).
		SetBasicAuth(c.Email, c.APIKey).
		Build()
	if err != nil {
		return nil, err
	}
	fmt.Println(req)
	resp, err = c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code %d", resp.StatusCode)
	}

	file, err := os.Create("/Users/xavier/Developer/goProject/migadu-go/test.html")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	w, err := io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(w)

	//doc, err := goquery.NewDocumentFromReader(resp.Body)
	//if err != nil {
	//	fmt.Print("Error loading HTML document:", err)
	//	return nil, err
	//}
	//
	//selection := doc.Find("html > body > main > div:nth-child(2)")
	//fmt.Println(selection.Text())
	//selection := doc.Find("html > body > main > div:nth-child(2) > div > div > table > tbody > tr")
	//if selection.Length() == 0 {
	//	fmt.Println("Element not found")
	//	return nil, fmt.Errorf("element not found")
	//}
	/////html/body/main/div[2]/div/div/table/tbody/tr/td[1]/strong/a
	//var r []string
	//selection.Each(func(i int, s *goquery.Selection) {
	//	// For each item found, get the band and title
	//	idStr := s.Find("td:nth-child(1)").Text()
	//	fmt.Printf("idStr: %s\n", idStr)
	//	r = append(r, idStr)
	//})
	return []string{}, nil
}
