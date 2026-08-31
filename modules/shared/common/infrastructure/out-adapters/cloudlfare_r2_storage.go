package commonoutadapters

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"
	commonports "github.com/PurpleSavage/monekai-server/modules/shared/common/application/ports"
	"github.com/PurpleSavage/monekai-server/modules/shared/common/config"
	globalerrors "github.com/PurpleSavage/monekai-server/modules/shared/common/infrastructure/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type CloudFlareAdapterService struct {
	client    *s3.Client
	presigner *s3.PresignClient
}

func NewCloudFlareAdapterService() commonports.StoragePort {

	cfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion("auto"),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.Envs.R2AccessKeyID,
				config.Envs.R2SecretAccessKey,
				"",
			),
		),
		func(o *awsconfig.LoadOptions) error {
			o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
			return nil
		},
	)

	if err != nil {
		panic(
			globalerrors.NewAppError(
				500,
				"R2 Configuration Error",
				"Could not initialize Cloudflare R2 config",
				err,
			),
		)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {

		o.BaseEndpoint = aws.String(
			fmt.Sprintf(
				"https://%s.r2.cloudflarestorage.com",
				config.Envs.R2AccountID,
			),
		)
	})

	presigner := s3.NewPresignClient(client)

	return &CloudFlareAdapterService{
		client:    client,
		presigner: presigner,
	}
}

func (c *CloudFlareAdapterService) UploadFile(
	key string,
	body io.Reader,
	size int64,
	contentType string,
) (string, error) {

	_, err := c.client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket:      aws.String(config.Envs.R2BucketName),
			Key:         aws.String(key),
			Body:        body,
			ContentType: aws.String(contentType),
			ContentLength: aws.Int64(size),
		},
	)

	if err != nil {
		log.Printf("R2 PUT OBJECT ERROR: %+v", err)
		return "", globalerrors.NewAppError(
			500,
			"R2 Upload Error",
			"Could not upload file to Cloudflare R2",
			err,
		)
	}

	publicURL := fmt.Sprintf(
		"%s/%s",
		config.Envs.R2PublicURL,
		key,
	)

	return publicURL, nil
}

func (c *CloudFlareAdapterService) GeneratePResignedURL(
	ctx context.Context,
	key string,
	contentType string,
) (string, error) {
	result, err := c.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(config.Envs.R2BucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = 15 * time.Minute
	})

	if err != nil {
		log.Printf("R2 PRESIGN ERROR: %+v", err)
		return "", globalerrors.NewAppError(500, "R2 Presign Error", "Could not generate presigned URL", err)
	}

	return result.URL, nil
}