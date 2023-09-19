package hw10programoptimization

import (
	"bufio"
	"io"
	"regexp"
	"strings"

	easyjson "github.com/mailru/easyjson"
)

type User struct {
	ID       int    `json:"ID"`       //nolint:tagliatelle
	Name     string `json:"Name"`     //nolint:tagliatelle
	Username string `json:"Username"` //nolint:tagliatelle
	Email    string `json:"Email"`    //nolint:tagliatelle
	Phone    string `json:"Phone"`    //nolint:tagliatelle
	Password string `json:"Password"` //nolint:tagliatelle
	Address  string `json:"Address"`  //nolint:tagliatelle
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)

	var user User

	domainRegexp, err := regexp.Compile(".+@(.+\\." + domain + ")$")
	if err != nil {
		return result, err
	}

	for scanner.Scan() {
		line := scanner.Bytes()
		if err := easyjson.Unmarshal(line, &user); err != nil {
			return result, err
		}

		submatch := domainRegexp.FindStringSubmatch(user.Email)
		if len(submatch) > 0 {
			fullDomain := strings.ToLower(submatch[1])
			result[fullDomain]++
		}
	}

	return result, nil
}
