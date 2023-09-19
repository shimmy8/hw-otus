package hw10programoptimization

import (
	"bufio"
	"fmt"
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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	reader := bufio.NewReader(r)
	i := 0
	for {
		line, readErr := reader.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			return result, err
		}

		var user User
		if err := easyjson.Unmarshal(line, &user); err != nil {
			return result, err
		}
		result[i] = user
		i++

		if readErr == io.EOF {
			return result, nil
		}
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	domainRegexp, err := regexp.Compile(".+@(.+\\." + domain + ")$")
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		found := domainRegexp.FindStringSubmatch(user.Email)
		if len(found) > 0 {
			fullDomain := strings.ToLower(found[1])
			result[fullDomain]++
		}
	}
	return result, nil
}
