package boshreleasedpb

type releaseConfig struct {
  Name_ string `yaml:"name"`
  FinalName_ string `yaml:"final_name"`
  Blobstore releaseConfigBlobstore `yaml:"blobstore"`
}

func (c *releaseConfig) Merge(a releaseConfig) {
  if a.Name_ != "" {
    c.Name_ = a.Name_
  }

  if a.FinalName_ != "" {
    c.FinalName_ = a.FinalName_
  }

  c.Blobstore.Merge(a.Blobstore)
}

type releaseConfigBlobstore struct {
  Provider string `yaml:"provider"`
  Options releaseConfigBlobstoreS3 `yaml:"options"`
}

func (c *releaseConfigBlobstore) Merge(a releaseConfigBlobstore) {
  if a.Provider != "" {
    c.Provider = a.Provider
  }

  c.Options.Merge(a.Options)
}

type releaseConfigBlobstoreS3 struct {
  BucketName string `yaml:"bucket_name"`
  Host string `yaml:"host"`
  AccessKeyID string `yaml:"access_key_id"`
  SecretAccessKey string `yaml:"secret_access_key"`
}

func (c *releaseConfigBlobstoreS3) Merge(a releaseConfigBlobstoreS3) {
  if a.BucketName != "" {
    c.BucketName = a.BucketName
  }

  if a.Host != "" {
    c.Host = a.Host
  }

  if a.AccessKeyID != "" {
    c.AccessKeyID = a.AccessKeyID
  }

  if a.SecretAccessKey != "" {
    c.SecretAccessKey = a.SecretAccessKey
  }
}
