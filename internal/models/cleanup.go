package models

import (
	bolt "go.etcd.io/bbolt"
)

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func DailyDataCleanup_TESTMODE(db *bolt.DB) {
	//go DeleteALLEndpoint(db)
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func DailyDataCleanup(db *bolt.DB) {
	ClearupDeleteendpointsFromOwnr(db)
	go ClearLogs(db)
	go ClearUserVerificationData(db)
}

// -----------------------------------------------------------------
//
// -----------------------------------------------------------------
func ClearupDeleteendpointsFromOwnr(db *bolt.DB) {

	userModel := &UserModel{DB: db}
	endPointModel := &EndPointModel{DB: db}

	for _, u := range userModel.List() {
		newOwnList := make([]string, 0)
		for _, epid := range u.OwnedEndPoints {

			_, err := endPointModel.Get(epid)
			if err == nil {
				newOwnList = append(newOwnList, epid)
			}
		}

		u.OwnedEndPoints = newOwnList

		userModel.Save(u, false)
	}

}

// ------------------------------------------------------
//
// ------------------------------------------------------
func ClearLogs(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(getLogTableName())
		tx.CreateBucketIfNotExists(getLogTableName())
		return nil
	})

	// clear end point log entrie

	epm := &EndPointModel{DB: db}

	eps := epm.List()

	for _, ep := range eps {
		ep.EndPointCallLog = make([]EndPointCallLog, 0)
		epm.Save(ep, "")
	}

}

// ------------------------------------------------------
//
// ------------------------------------------------------
func ClearUserVerificationData(db *bolt.DB) {

	userModel := &UserModel{DB: db}

	db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(userModel.GetVerificationTableName())
		tx.CreateBucketIfNotExists(userModel.GetVerificationTableName())

		tx.DeleteBucket(userModel.GetPasswordResetTableName())
		tx.CreateBucketIfNotExists(userModel.GetPasswordResetTableName())
		return nil
	})

}

// ------------------------------------------------------
//
// ------------------------------------------------------
func DeleteALLEndpoint(db *bolt.DB) {

	epm := &EndPointModel{DB: db}

	for _, e := range epm.List() {
		epm.Delete(e.ID)
	}
}
