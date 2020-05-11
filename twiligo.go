package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type message struct {
	to      string
	from    string
	message string
}

type config struct {
	accountSid string
	authToken  string
	url        string
}

func main() {
	c := &config{}
	var ok bool
	if c.accountSid, ok = os.LookupEnv("TWILIO_ACCOUNT_SID"); !ok {
		log.Fatalf("Need to set %v environment variable\n", "TWILIO_ACCOUNT_SID")
	}
	if c.authToken, ok = os.LookupEnv("TWILIO_AUTH_TOKEN"); !ok {
		log.Fatalf("Need to set %v environment variable\n", "TWILIO_AUTH_TOKEN")
	}
	c.url = fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", c.accountSid)

	m := &message{
		message: getProverb(),
	}
	if m.to, ok = os.LookupEnv("TO"); !ok {
		log.Fatalf("Need to set %v environment variable\n", "TO")
	}
	if m.from, ok = os.LookupEnv("FROM"); !ok {
		log.Fatalf("Need to set %v environment variable\n", "FROM")
	}

	if err := m.sendSms(c); err != nil {
		log.Fatal(err)
	}

}

func getProverb() string {
	proverbs := []string{
		"Don't communicate by sharing memory, share memory by communicating.",
		"Concurrency is not parallelism.",
		"Channels orchestrate; mutexes serialize.",
		"The bigger the interface, the weaker the abstraction.",
		"Make the zero value useful.",
		"interface{} says nothing.",
		"Gofmt's style is no one's favorite, yet gofmt is everyone's favorite.",
		"A little copying is better than a little dependency.",
		"Syscall must always be guarded with build tags.",
		"Cgo must always be guarded with build tags.",
		"Cgo is not Go.",
		"With the unsafe package there are no guarantees.",
		"Clear is better than clever.",
		"Reflection is never clear.",
		"Errors are values.",
		"Don't just check errors, handle them gracefully.",
		"Design the architecture, name the components, document the details.",
		"Documentation is for users.",
		"Don't panic.",
	}

	rand.Seed(time.Now().Unix())
	return proverbs[rand.Intn(len(proverbs))]
}

func (m *message) sendSms(c *config) error {
	msgData := url.Values{}
	msgData.Set("To", m.to)
	msgData.Set("From", m.from)
	msgData.Set("Body", m.message)

	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", c.url, &msgDataReader)
	req.SetBasicAuth(c.accountSid, c.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("POST to response url: %s", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("POST returned a non 200 error code: %v", res.StatusCode)
	}

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("decode response data: %s", err)
	}

	fmt.Println(data["sid"])
	return nil
}
