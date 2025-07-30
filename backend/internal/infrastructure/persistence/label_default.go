package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
)

func getDefaultLabels() []label.Label {
	owner := auth.NewSystemOwner()
	return []label.Label{
		label.NewLabel(uuid.MustParse("387b2bea-8f73-46e2-8223-be8e0db6ec48"), owner, "Music",
			"#FF81C784", time.Now(), time.Now()),
		label.NewLabel(uuid.MustParse("825bd63a-1572-4a94-bd2e-ea7168941345"), owner, "Internet",
			"#FF64B5F6", time.Now(), time.Now()),
		label.NewLabel(uuid.MustParse("ccb4e552-82d3-46ba-ab61-60c792988976"), owner, "Mobile",
			"#FFFFD54F", time.Now(), time.Now()),
		label.NewLabel(uuid.MustParse("49d713ba-ac07-4e6d-896e-5e9965aa68c0"), owner, "Utilities",
			"#FF9575CD", time.Now(), time.Now()),
		label.NewLabel(uuid.MustParse("4b3144a9-a98b-4b6e-bf70-5e0f69dffc6f"), owner, "Streaming",
			"#FFF06292", time.Now(), time.Now()),
		label.NewLabel(uuid.MustParse("e7145075-8250-45fc-99dc-3e2172e073de"), owner, "Gaming",
			"#FF4DB6AC", time.Now(), time.Now()),
		label.NewLabel(uuid.MustParse("4a3aaa69-edf7-4175-af25-eb93a9cd2816"), owner, "Software",
			"#FF7986CB", time.Now(), time.Now()),
	}

}
