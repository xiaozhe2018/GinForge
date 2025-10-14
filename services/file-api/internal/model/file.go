package model

import (
	"time"

	"goweb/pkg/model"
)

// FileRecord 文件记录
type FileRecord struct {
	model.BaseModel
	OriginalName  string    `json:"original_name" gorm:"type:varchar(255);not null;comment:原始文件名"`
	FileName      string    `json:"file_name" gorm:"type:varchar(255);not null;index;comment:存储文件名"`
	RelativePath  string    `json:"relative_path" gorm:"type:varchar(500);not null;uniqueIndex;comment:相对路径"`
	FileSize      int64     `json:"file_size" gorm:"type:bigint;not null;comment:文件大小(字节)"`
	MimeType      string    `json:"mime_type" gorm:"type:varchar(100);comment:MIME类型"`
	FileType      string    `json:"file_type" gorm:"type:varchar(50);index;comment:文件类型(image/video/document/other)"`
	Hash          string    `json:"hash" gorm:"type:varchar(64);index;comment:文件哈希(MD5)"`
	StorageType   string    `json:"storage_type" gorm:"type:varchar(50);default:local;comment:存储类型(local/oss/s3)"`
	URL           string    `json:"url" gorm:"type:varchar(500);comment:访问URL"`
	UploadedBy    uint      `json:"uploaded_by" gorm:"type:int unsigned;index;comment:上传用户ID"`
	UserType      string    `json:"user_type" gorm:"type:varchar(50);comment:用户类型(admin/user/merchant)"`
	UploadIP      string    `json:"upload_ip" gorm:"type:varchar(50);comment:上传IP"`
	UploadTime    time.Time `json:"upload_time" gorm:"type:datetime;index;comment:上传时间"`
	DownloadCount int       `json:"download_count" gorm:"type:int;default:0;comment:下载次数"`
	Status        int       `json:"status" gorm:"type:tinyint;default:1;index;comment:状态(1:正常 2:已删除)"`
	Description   string    `json:"description" gorm:"type:varchar(500);comment:文件描述"`
	Tags          string    `json:"tags" gorm:"type:varchar(255);comment:文件标签(逗号分隔)"`
}

// TableName 表名
func (FileRecord) TableName() string {
	return "file_records"
}

// FileUploadLog 文件上传日志
type FileUploadLog struct {
	model.BaseModel
	FileName     string    `json:"file_name" gorm:"type:varchar(255);not null;comment:文件名"`
	FileSize     int64     `json:"file_size" gorm:"type:bigint;comment:文件大小"`
	UploadedBy   uint      `json:"uploaded_by" gorm:"type:int unsigned;index;comment:上传用户ID"`
	UploadIP     string    `json:"upload_ip" gorm:"type:varchar(50);comment:上传IP"`
	UploadTime   time.Time `json:"upload_time" gorm:"type:datetime;index;comment:上传时间"`
	Success      bool      `json:"success" gorm:"type:tinyint(1);comment:是否成功"`
	ErrorMessage string    `json:"error_message" gorm:"type:text;comment:错误信息"`
	Duration     int64     `json:"duration" gorm:"type:bigint;comment:上传耗时(毫秒)"`
	UploadType   string    `json:"upload_type" gorm:"type:varchar(20);comment:上传类型(simple/chunked)"`
}

// TableName 表名
func (FileUploadLog) TableName() string {
	return "file_upload_logs"
}

// FileDownloadLog 文件下载日志
type FileDownloadLog struct {
	model.BaseModel
	FileID       uint      `json:"file_id" gorm:"type:int unsigned;index;comment:文件ID"`
	FileName     string    `json:"file_name" gorm:"type:varchar(255);comment:文件名"`
	DownloadedBy uint      `json:"downloaded_by" gorm:"type:int unsigned;index;comment:下载用户ID"`
	DownloadIP   string    `json:"download_ip" gorm:"type:varchar(50);comment:下载IP"`
	DownloadTime time.Time `json:"download_time" gorm:"type:datetime;index;comment:下载时间"`
	Success      bool      `json:"success" gorm:"type:tinyint(1);comment:是否成功"`
}

// TableName 表名
func (FileDownloadLog) TableName() string {
	return "file_download_logs"
}
