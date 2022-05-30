package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"zadatak/config"

	"github.com/tidwall/gjson"
)

const URL_Login = "https://instagram.com/accounts/login/ajax/"
const URL_Base = "https://instagram.com"

//Request struct for JSON request to get username from frontend
type Request struct {
	Username string `json:"username"`
}

//Response struct for JSON response from backend
type Data struct {
	Engagement string `json:"engagement"`
}

func main() {
	fmt.Println("Starting server...")

	http.HandleFunc("/calculate", engagement)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		panic(err)
	}

}

//handler for calculating engagement and parsing username
func engagement(w http.ResponseWriter, r *http.Request) {



	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/json; charset=utf-8")

	body := Request{}
	json.NewDecoder(r.Body).Decode(&body)

	engagement := calculateEngagement(body.Username)


	EngagementResponse := Data{
		Engagement: engagement,
	}
	json.NewEncoder(w).Encode(EngagementResponse)




}

//send request to instagram for username with csrf and sessionid and return calculated engagement
func calculateEngagement(username string) string {

	sessionCookie, sessionCsrf := login()
	// fmt.Println(sessionCookie, sessionCsrf)
	client := &http.Client{}
	println("Get Account Information...")
	req, _ := http.NewRequest("GET", "https://www.instagram.com/"+username, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36")
	req.Header.Set("cookie", sessionCookie)
	req.Header.Set("x-csrftoken", sessionCsrf)
	req.Header.Set("referer", URL_Base+"/")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body))

	window := between(string(body), `window._sharedData = `, `;</script>`)

	followers := gjson.Get(window, "entry_data.ProfilePage.0.graphql.user.edge_followed_by.count")
	fmt.Println(followers)
	edges := gjson.Get(window, "entry_data.ProfilePage.0.graphql.user.edge_owner_to_timeline_media.edges").Array()
	comments := float64(0)
	likes := float64(0)
	for _, edge := range edges {
		comment := gjson.Get(edge.String(), "node.edge_media_to_comment")
		comments += comment.Get("count").Float()
		like := gjson.Get(edge.String(), "node.edge_liked_by")
		likes += like.Get("count").Float()
	}

	engagement := (likes + comments) / followers.Float()

	engagementString := strconv.FormatFloat(engagement, 'f', 2, 64)

	return engagementString

}

//login to instagram and return cookie and csrf token for session
func login() (string, string) {
	username := config.Config("INSTA_USERNAME")
	password := config.Config("INSTA_PASSWORD")

	var mid string
	var csrftoken string

	println("Initiate...")
	req, _ := http.NewRequest("GET", URL_Base, nil)
	client := &http.Client{}
	req.Header.Set("cookie", "ig_cb=1")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "csrftoken" {
			csrftoken = cookies[i].Value
		}
		if cookies[i].Name == "mid" {
			mid = cookies[i].Value
		}
	}

	data := url.Values{}
	data.Set("username", username)
	data.Add("password", password)
	req, _ = http.NewRequest("POST", URL_Login, bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value") // This makes it work
	req.Header.Set("x-csrftoken", csrftoken)
	req.Header.Set("cookie", fmt.Sprintf("csrftoken=%s; mid=%s;", csrftoken, mid))
	req.Header.Set("referer", URL_Base)
	req.Header.Set("user-agent", "Instagram 10.26.0 (iPhone9,1; iOS 10_2_1; en_US; en-US; scale=2.00; 750x1334; 1605844) AppleWebKit/420+")

	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	cookies = resp.Cookies()

	var sessionCsrf string
	var sessionCookie string

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == "csrftoken" && cookies[i].Value != "" {
			sessionCsrf = cookies[i].Value
		}

		sessionCookie += fmt.Sprintf("%s=%s; ", cookies[i].Name, cookies[i].Value)
	}

	return sessionCookie, sessionCsrf
}

//helper function for getting string between two strings
func between(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}
