package family

func memberUniqueComparer(a Member, b Member) bool {
	return a.Id() == b.Id()
}

func memberComparer(a Member, b Member) bool {
	return a.Equal(b)
}
