package MySql

import "fmt"

func RegisterUserIfNotExists(username, password string, DefaultStorageQuota int) error {
	var exists bool

	err := Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	_, err = Db.Exec("INSERT INTO users (username, password, storage_quota) VALUES (?, ?, ?)", username, password, DefaultStorageQuota)
	if err != nil {
		return err
	}

	fmt.Println("User inserted successfully")
	return nil
}
