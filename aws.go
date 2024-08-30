package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"log/slog"
	"time"
)

func listUsers(ctx context.Context, svc *iam.Client) ([]types.User, error) {
	slog.Debug("listUser")
	var users []types.User
	var marker *string

	for {
		input := &iam.ListUsersInput{
			MaxItems: aws.Int32(500),
			Marker:   marker,
		}
		resp, err := svc.ListUsers(ctx, input)
		if err != nil {
			return users, err
		}
		users = append(users, resp.Users...)
		marker = resp.Marker

		slog.Debug("truncated", "is", resp.IsTruncated)
		if !resp.IsTruncated {
			break
		}
	}

	return users, nil
}

func listAccessKeys(ctx context.Context, svc *iam.Client, username string) ([]types.AccessKeyMetadata, error) {
	slog.Debug("listAccessKeys")
	var accessKeys []types.AccessKeyMetadata
	var marker *string

	for {
		input := &iam.ListAccessKeysInput{
			UserName: aws.String(username),
			MaxItems: aws.Int32(500),
			Marker:   marker,
		}
		resp, err := svc.ListAccessKeys(ctx, input)
		if err != nil {
			return accessKeys, err
		}
		accessKeys = append(accessKeys, resp.AccessKeyMetadata...)
		marker = resp.Marker

		if !resp.IsTruncated {
			break
		}
	}

	return accessKeys, nil
}

type AccessKeyInfo struct {
	AccessKeyId string
	UserName    string
}

func GetExpiredAccessKeys(ctx context.Context, expireSec int) ([]AccessKeyInfo, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	svc := iam.NewFromConfig(cfg)

	users, err := listUsers(ctx, svc)
	if err != nil {
		return nil, err
	}

	var accessKeys []types.AccessKeyMetadata
	for _, user := range users {
		ret, err := listAccessKeys(ctx, svc, *user.UserName)
		if err != nil {
			return nil, err
		}

		accessKeys = append(accessKeys, ret...)
	}

	// elapsed := time.Now().Add(time.Second * time.Duration(expireSec))
	threshold := time.Second * time.Duration(expireSec)

	var ret []AccessKeyInfo
	for _, key := range accessKeys {
		elapsed := time.Since(*key.CreateDate)
		if elapsed > threshold {
			ret = append(ret, AccessKeyInfo{
				AccessKeyId: *key.AccessKeyId,
				UserName:    *key.UserName,
			})
		}
	}

	return ret, nil
}
