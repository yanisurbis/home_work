package main //nolint:golint,stylecheck

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	reg, err := regexp.Compile("\\." + domain)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {

		var user User
		if err = json.Unmarshal([]byte(scanner.Text()), &user); err != nil {
			return nil, err
		}
		matched := reg.MatchString(user.Email)

		if matched {
			domain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[domain]
			result[domain] = num + 1
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return result, nil
}

func main() {
	data := `{"Id":1,"Name":"Howard Mendoza","Username":"0Oliver","Email":"aliquid_qui_ea@Browsedrive.gov","Phone":"6-866-899-36-79","Password":"InAQJvsq","Address":"Blackbird Place 25"}
{"Id":2,"Name":"Jesse Vasquez","Username":"qRichardson","Email":"mLynch@broWsecat.com","Phone":"9-373-949-64-00","Password":"SiZLeNSGn","Address":"Fulton Hill 80"}
{"Id":3,"Name":"Clarence Olson","Username":"RachelAdams","Email":"RoseSmith@Browsecat.com","Phone":"988-48-97","Password":"71kuz3gA5w","Address":"Monterey Park 39"}
{"Id":4,"Name":"Gregory Reid","Username":"tButler","Email":"5Moore@Teklist.net","Phone":"520-04-16","Password":"r639qLNu","Address":"Sunfield Park 20"}
{"Id":5,"Name":"Janice Rose","Username":"KeithHart","Email":"nulla@Linktype.com","Phone":"146-91-01","Password":"acSBF5","Address":"Russell Trail 61"}`

	GetDomainStat(bytes.NewBufferString(data), "com")
}
