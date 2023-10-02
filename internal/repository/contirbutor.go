package repository

func (r *Repository) GetContributorByUser(userID uint) (uint, error) {
	var contributor uint

	sql := `SELECT User_ID FROM Contributor WHERE User_ID = ?`
	if err := r.db.Raw(sql, userID).First(&contributor).Error; err != nil {
		return contributor, err
	}

	return contributor, nil
}
