package service

import (
	"encoding/json"
)

type CachedShortUrl struct {
	Url            string `redis:"url"`
	FraudDetection bool   `redis:"fraud_detection"`
	Disabled       bool   `redis:"disabled"`
	NoReferrer     bool   `redis:"no_referrer"`
	UtmSource      string `redis:"utm_source"`
	UtmMedium      string `redis:"utm_medium"`
	UtmCampaign    string `redis:"utm_campaign"`
	UtmTerm        string `redis:"utm_term"`
	UtmContent     string `redis:"utm_content"`
}

func (m *CachedShortUrl) MarshalBinary() ([]byte, error) {
	return json.Marshal(m)
}

func (m *CachedShortUrl) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}

func (m *CachedShortUrl) GeneratedUrl() string {
	var utm string = ""
	if m.UtmSource != "" {
		if utm == "" {
			utm += "?utm_source=" + m.UtmSource
		} else {
			utm += "&utm_source=" + m.UtmSource
		}
	}

	if m.UtmMedium != "" {
		if utm == "" {
			utm += "?utm_medium=" + m.UtmMedium
		} else {
			utm += "&utm_medium=" + m.UtmMedium
		}
	}

	if m.UtmCampaign != "" {
		if utm == "" {
			utm += "?utm_campaign=" + m.UtmCampaign
		} else {
			utm += "&utm_campaign=" + m.UtmCampaign
		}
	}

	if m.UtmTerm != "" {
		if utm == "" {
			utm += "?utm_term=" + m.UtmTerm
		} else {
			utm += "&utm_term=" + m.UtmTerm
		}
	}

	if m.UtmContent != "" {
		if utm == "" {
			utm += "?utm_content=" + m.UtmContent
		} else {
			utm += "&utm_content=" + m.UtmContent
		}
	}

	return m.Url + utm
}
