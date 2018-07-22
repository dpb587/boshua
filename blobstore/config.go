package blobstore

type BlobstoreConfig struct {
	S3 S3BlobstoreConfig `yaml:"s3"`
}

type S3BlobstoreConfig struct {
	Host      string `yaml:"host"`
	Bucket    string `yaml:"bucket"`
	Prefix    string `yaml:"prefix"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}
