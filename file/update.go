package file

// Update all entries in items database
func Update() error {
	r := NewItemRepository()

	err := r.AppendOrRemove()

	if err != nil {
		return err
	}

	return r.UpdateAllFingerprints()
}
