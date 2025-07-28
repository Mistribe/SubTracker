package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
)

func getDefaultLabels() []label.label {
	return []label.label{
		label.NewLabel(uuid.MustParse("387b2bea-8f73-46e2-8223-be8e0db6ec48"), nil, "Music", true,
			"#FF81C784", time.Now(), time.Now(), false),
		label.NewLabel(uuid.MustParse("825bd63a-1572-4a94-bd2e-ea7168941345"), nil, "Internet", true,
			"#FF64B5F6", time.Now(), time.Now(), false),
		label.NewLabel(uuid.MustParse("ccb4e552-82d3-46ba-ab61-60c792988976"), nil, "Mobile", true,
			"#FFFFD54F", time.Now(), time.Now(), false),
		label.NewLabel(uuid.MustParse("49d713ba-ac07-4e6d-896e-5e9965aa68c0"), nil, "Utilities", true,
			"#FF9575CD", time.Now(), time.Now(), false),
		label.NewLabel(uuid.MustParse("4b3144a9-a98b-4b6e-bf70-5e0f69dffc6f"), nil, "Streaming", true,
			"#FFF06292", time.Now(), time.Now(), false),
		label.NewLabel(uuid.MustParse("e7145075-8250-45fc-99dc-3e2172e073de"), nil, "Gaming", true,
			"#FF4DB6AC", time.Now(), time.Now(), false),
		label.NewLabel(uuid.MustParse("4a3aaa69-edf7-4175-af25-eb93a9cd2816"), nil, "Software", true,
			"#FF7986CB", time.Now(), time.Now(), false),
	}

}
