package config

import (
	"github.com/spf13/viper"
)

var DefaultS3CacherConfig = &S3CacherConfig{
	Bucket:          "REPLACEME",
	Region:          "us-west-1",
	AccessKeyId:     "REPLACEME",
	SecretAccessKey: "REPLACEME",
}

type S3CacherConfig struct {
	ConfigBase

	Bucket string
	Region string

	AccessKeyId     string
	SecretAccessKey string
}

func (p *S3CacherConfig) ConfigureViper() {
	p.setConfig("s3_cacher.bucket",
		DefaultS3CacherConfig.Bucket)
	p.setConfig("s3_cacher.region",
		DefaultS3CacherConfig.Region)
	p.setConfig("s3_cacher.access_key_id",
		DefaultS3CacherConfig.AccessKeyId)
	p.setConfig("s3_cacher.secret_access_key",
		DefaultS3CacherConfig.SecretAccessKey)
}

func (p *S3CacherConfig) ReadConfig() {
	p.Bucket = viper.GetString("s3_cacher.bucket")
	p.Region = viper.GetString("s3_cacher.region")
	p.AccessKeyId = viper.GetString("s3_cacher.access_key_id")
	p.SecretAccessKey = viper.GetString("s3_cacher.secret_access_key")
}
