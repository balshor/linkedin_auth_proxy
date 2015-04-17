package providers

import (
        "bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
        "net/url"

	"github.com/bitly/go-simplejson"
	"github.com/bitly/google_auth_proxy/api"
)

type LinkedInProvider struct {
        *ProviderData
}

func NewLinkedInProvider(p *ProviderData) *LinkedInProvider {
        p.ProviderName = "LinkedIn"
	if p.LoginUrl.String() == "" {
	        p.LoginUrl = &url.URL{Scheme: "https",
		        Host: "www.linkedin.com",
			Path: "/uas/oauth2/authorization"}
	}
	if p.RedeemUrl.String() == "" {
	        p.RedeemUrl = &url.URL{Scheme: "https",
		        Host: "www.linkedin.com",
			Path: "/uas/oauth2/accessToken"}
	}
	if p.ProfileUrl.String() == "" {
	        p.ProfileUrl = &url.URL{Scheme: "https",
		        Host: "www.linkedin.com",
			Path: "/v1/people/~/email-address"}
        }
	if p.Scope == "" {
	        p.Scope = "r_emailaddress r_basicprofile"
	}
	return &LinkedInProvider{ProviderData: p}
}

func (p *LinkedInProvider) GetEmailAddress(unused_auth_response *simplejson.Json,
        access_token string) (string, error) {
        if access_token == "" {
                return "", errors.New("missing access token")
        }
        params := url.Values{}
        req, err := http.NewRequest("GET", p.ProfileUrl.String()+"?format=json", bytes.NewBufferString(params.Encode()))
        if err != nil {
           	return "", err
        }
        req.Header.Set("Accept", "application/json")
	req.Header.Set("x-li-format", "json")
        req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	json, err := api.Request(req)
	if err != nil {
	        log.Printf("failed making request %s", err)
		return "", err
	}

	email, err := json.String()
	if err != nil {
	        log.Printf("failed making request %s", err)
		return "", err
	}
        return email, nil
}
