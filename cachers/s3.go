package cachers

import (
	"bytes"
	"io/ioutil"

	"github.com/ssalevan/cachepix/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func NewS3Cacher(conf *config.S3CacherConfig) *S3Cacher {
	return &S3Cacher{
		conf: conf,
	}
}

type S3Cacher struct {
	conf *config.S3CacherConfig

	awsSession *session.Session
	s3Client   *s3.S3
	s3Uploader *s3manager.Uploader
}

func (s *S3Cacher) Init() error {
	var err error

	awsCredentials := credentials.NewStaticCredentials(
		s.conf.AccessKeyId, s.conf.SecretAccessKey, "")
	awsConfig := aws.NewConfig().WithRegion(s.conf.Region).WithCredentials(awsCredentials)
	s.awsSession, err = session.NewSessionWithOptions(session.Options{
		Config: *awsConfig,
	})
	if err != nil {
		return err
	}

	s.s3Client = s3.New(s.awsSession)
	s.s3Uploader = s3manager.NewUploader(s.awsSession)

	_, err = s.s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s.conf.Bucket),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				return nil
			default:
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (s *S3Cacher) Get(url string) (bool, []byte, error) {
	var contents []byte
	result, err := s.s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.conf.Bucket),
		Key:    aws.String(url),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return false, contents, nil
			default:
				return false, contents, err
			}
		} else {
			return false, contents, err
		}
	}

	contents, err = ioutil.ReadAll(result.Body)
	return true, contents, err
}

func (s *S3Cacher) Name() string {
	return "s3"
}

func (s *S3Cacher) Set(url string, contents []byte) error {
	result, err := s.s3Client.ListObjects(&s3.ListObjectsInput{
		Bucket:  aws.String(s.conf.Bucket),
		Prefix:  aws.String(url),
		MaxKeys: aws.Int64(1),
	})
	if err != nil {
		return err
	}
	if len(result.Contents) == 0 {
		_, err = s.s3Uploader.Upload(&s3manager.UploadInput{
			Bucket: aws.String(s.conf.Bucket),
			Key:    aws.String(url),
			Body:   bytes.NewReader(contents),
		})
	}
	return err
}
