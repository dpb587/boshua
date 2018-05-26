package blobstore

type BlobstoreConfig struct {
	AWS AWSBlobstoreConfig `yaml:"aws"`
}

type AWSBlobstoreConfig struct {
	Host      string `yaml:"host"`
	Bucket    string `yaml:"bucket"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
}
