package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
)

func getDefaultLabels() []label.Label {
	return []label.Label{
		label.NewLabelWithoutValidation(uuid.MustParse("387b2bea-8f73-46e2-8223-be8e0db6ec48"), "Music", true,
			"#81C784", time.Now(), time.Now()),
		label.NewLabelWithoutValidation(uuid.MustParse("825bd63a-1572-4a94-bd2e-ea7168941345"), "Internet", true,
			"#64B5F6", time.Now(), time.Now()),
		label.NewLabelWithoutValidation(uuid.MustParse("ccb4e552-82d3-46ba-ab61-60c792988976"), "Mobile", true,
			"#FFD54F", time.Now(), time.Now()),
		label.NewLabelWithoutValidation(uuid.MustParse("49d713ba-ac07-4e6d-896e-5e9965aa68c0"), "Utilities", true,
			"#9575CD", time.Now(), time.Now()),
		label.NewLabelWithoutValidation(uuid.MustParse("4b3144a9-a98b-4b6e-bf70-5e0f69dffc6f"), "Streaming", true,
			"#F06292", time.Now(), time.Now()),
		label.NewLabelWithoutValidation(uuid.MustParse("e7145075-8250-45fc-99dc-3e2172e073de"), "Gaming", true,
			"#4DB6AC", time.Now(), time.Now()),
		label.NewLabelWithoutValidation(uuid.MustParse("4a3aaa69-edf7-4175-af25-eb93a9cd2816"), "Software", true,
			"#7986CB", time.Now(), time.Now()),
	}

}
