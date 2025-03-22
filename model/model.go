package model

type Configurations struct {
	Prefix string
	Server struct {
		Host string
		Port string
	}
	MysqlConf              MysqlConfig
	SecretKey              string
	TokenValidityInSeconds int
	BaseDir                string
	DefaultStorageQuota    int
}

type MysqlConfig struct {
	Username     string
	Password     string
	Net          string
	Address      string
	DatabaseName string
}

type UserReqBody struct {
	Username string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserDetailsFromDB struct {
	UserName     string `json:"userName"`
	Password     string `json:"password"`
	StorageQuota int64  `json:"storage_quota"`
	UsedStorage  int64  `json:"used_storage"`
}

type FileMetadata struct {
	FileName         string `json:"file_name"`
	OriginalFileName string `json:"original_filename"`
	UploadedBy       string `json:"uploaded_by"`
	Size             int64  `json:"file_size"`
	UploadTime       string `json:"uploaded_timestamp"`
}

type FileResponse struct {
	FileName   string `json:"file_name"`
	UploadedBy string `json:"uploaded_by"`
	Size       string `json:"file_size"`
	UploadTime string `json:"uploaded_timestamp"`
}
