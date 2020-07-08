package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
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

	if err != nil {
		return result, err
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		if reg.MatchString(text) {
			var user User

			if err = user.UnmarshalJSON([]byte(text)); err != nil {
				return nil, err
			}

			if reg.MatchString(user.Email) {
				domain := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
				num := result[domain]
				result[domain] = num + 1
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}

	return result, nil
}
