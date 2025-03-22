package MySql

import (
	"filestorage/model"
	"fmt"
	"log"
)

func StoreFileMetadata(metadata model.FileMetadata) error {
	query := `INSERT INTO file_metadata (file_name, original_filename, uploaded_by, file_size, uploaded_timestamp) 
              VALUES (?, ?, ?, ?, ?)`

	_, err := Db.Exec(query, metadata.FileName, metadata.OriginalFileName, metadata.UploadedBy, metadata.Size, metadata.UploadTime)
	if err != nil {
		log.Println("Error inserting file metadata:", err)
		return err
	}

	log.Println("File metadata stored successfully for:", metadata.FileName)
	return nil
}

func UpdateUserStorage(username string, fileSize int64) error {
	query := `UPDATE users SET used_storage = used_storage + ? WHERE username = ? AND (storage_quota - used_storage) >= ?`

	result, err := Db.Exec(query, fileSize, username, fileSize)
	if err != nil {
		log.Println("Error updating user storage:", err)
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		log.Printf("Insufficient storage for user %s. Upload failed.\n", username)
		return fmt.Errorf("insufficient storage available")
	}
	log.Printf("Updated used storage for user %s. Increased by %d bytes\n", username, fileSize)
	return nil
}
