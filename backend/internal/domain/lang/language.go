package lang

import (
	"context"
	"sort"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

const (
	ContextKey = "X-Lang"
)

type Info interface {
	Preferences() []Preference
	MostPreferred() Preference
}

type info struct {
	preferences []Preference
}

func (i info) Preferences() []Preference {
	return i.preferences
}

func (i info) MostPreferred() Preference {
	return i.preferences[0]
}

type Preference interface {
	Locale() string
	Weight() float32
	PreferredCurrency() currency.Unit
}

type preference struct {
	tag    language.Tag
	weight float32
}

func (p preference) Locale() string {
	return p.tag.String()
}

func (p preference) Weight() float32 {
	return p.weight
}

func (p preference) PreferredCurrency() currency.Unit {
	reg, conf := p.tag.Region()
	if conf >= language.Exact {
		c, ok := currency.FromRegion(reg)
		if ok {
			return c
		}
		return currency.USD
	}
	return currency.USD
}

func newPreference(tag language.Tag, weight float32) Preference {
	return &preference{
		tag:    tag,
		weight: weight,
	}
}

func GetDefault() Info {
	return &info{
		preferences: []Preference{
			newPreference(language.English, 1.0),
		},
	}
}

func ParseAcceptLanguage(acceptLanguage string) (Info, error) {
	tags, q, err := language.ParseAcceptLanguage(acceptLanguage)
	if err != nil {
		return nil, err
	}

	prefs := make([]Preference, len(tags))
	for i, tag := range tags {
		prefs[i] = newPreference(tag, q[i])
	}

	sort.Slice(prefs, func(i, j int) bool {
		return prefs[i].Weight() > prefs[j].Weight()
	})

	return &info{
		preferences: prefs,
	}, nil
}

func FromContext(ctx context.Context) Info {
	if ctx == nil {
		return GetDefault()
	}

	value := ctx.Value(ContextKey)
	if value == nil {
		return GetDefault()
	}

	d, ok := value.(Info)
	if !ok {
		return GetDefault()
	}

	return d
}
