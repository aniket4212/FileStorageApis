package MySql

import (
	"filestorage/model"
	"log"
)

func FetchUserDetailsFromDB(userName string) (model.UserDetailsFromDB, error) {
	var UserDetailsFromMySQL model.UserDetailsFromDB

	fetchUserDetailsQuery := "SELECT username, password, storage_quota, used_storage FROM users WHERE username = ?;"

	log.Println("Query to retrieve user details from DB::", fetchUserDetailsQuery)
	log.Println("Username being queried:", userName)

	err := Db.QueryRow(fetchUserDetailsQuery, userName).Scan(&UserDetailsFromMySQL.UserName, &UserDetailsFromMySQL.Password, &UserDetailsFromMySQL.StorageQuota, &UserDetailsFromMySQL.UsedStorage)
	if err != nil {
		log.Println("Error fetching user details:", err)
		return model.UserDetailsFromDB{}, err
	}

	log.Println("User details fetched successfully:", UserDetailsFromMySQL)
	return UserDetailsFromMySQL, nil
}
