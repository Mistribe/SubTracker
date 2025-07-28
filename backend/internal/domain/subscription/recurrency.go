package subscription

type RecurrencyType string

const (
	OneTime    RecurrencyType = "one_time"
	Monthly    RecurrencyType = "monthly"
	Quarterly  RecurrencyType = "quarterly"
	HalfYearly RecurrencyType = "half_yearly"
	Yearly     RecurrencyType = "yearly"
	Custom     RecurrencyType = "custom"
)
