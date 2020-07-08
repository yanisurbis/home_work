package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
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
		if err = user.UnmarshalJSON([]byte(scanner.Text())); err != nil {
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
