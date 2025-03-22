package MySql

import (
	"filestorage/model"
	"fmt"
	"log"
)

func FetchUploadedFiles(username string, limit, offset int) ([]model.FileMetadata, int, error) {
	var files []model.FileMetadata
	var totalRecords int

	fmt.Println("offset========", offset)

	// Query to fetch totalrecords and metadata from db
	query := "SELECT file_name, original_filename, uploaded_by, file_size, uploaded_timestamp, (SELECT COUNT(*) FROM file_metadata WHERE uploaded_by = ?) AS total_count FROM file_metadata WHERE uploaded_by = ? ORDER BY uploaded_timestamp DESC LIMIT ? OFFSET ?"
	log.Println("Executing SQL Query:", query, "Params:", username, limit, offset)

	rows, err := Db.Query(query, username, username, limit, offset)
	if err != nil {
		log.Println("Error fetching paginated files:", err)
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var file model.FileMetadata
		if err := rows.Scan(&file.FileName, &file.OriginalFileName, &file.UploadedBy, &file.Size, &file.UploadTime, &totalRecords); err != nil {
			log.Println("Error scanning file row:", err)
			return nil, 0, err
		}
		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error reading file rows:", err)
		return nil, 0, err
	}

	return files, totalRecords, nil

}
