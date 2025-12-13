package auth

import (
	"context"

	"github.com/nedpals/supabase-go"
)

// Client 封装 Supabase 客户端（中文说明：用于与 Supabase Auth API 交互）
type Client struct {
	client *supabase.Client
}

// NewClient 创建 Supabase 客户端
func NewClient(url, key string) *Client {
	if url == "" || key == "" {
		return nil
	}
	return &Client{
		client: supabase.CreateClient(url, key),
	}
}

// SignUp 注册新用户
func (c *Client) SignUp(email, password string) (*supabase.User, error) {
	user, err := c.client.Auth.SignUp(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// SignIn 登录用户
func (c *Client) SignIn(email, password string) (*supabase.AuthenticatedDetails, error) {
	user, err := c.client.Auth.SignIn(context.Background(), supabase.UserCredentials{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}
