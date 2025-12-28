package tigeropen

import (
	"bytes"
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"
)

const (
	defaultServerURL = "https://openapi.tigerfintech.com/gateway"
	defaultCharset   = "UTF-8"
	defaultSignType  = "RSA"
	defaultVersion   = "2.0"
	defaultUserAgent = "openapi-go-sdk-0.1.0"
)

// Config carries connection settings.
type Config struct {
	TigerID        string
	Account        string
	SecretKey      string
	PrivateKey     string // RSA key, with or without PEM markers.
	TigerPublicKey string
	ServerURL      string
	DeviceID       string
	NotifyURL      string
	Charset        string
	SignType       string
	Version        string
	Lang           string
	Token          string
	Timeout        time.Duration
	HTTPClient     *http.Client
}

// Client executes signed OpenAPI requests.
type Client struct {
	cfg        Config
	privateKey *rsa.PrivateKey
	httpClient *http.Client
	userAgent  string
}

// NewClient builds a Client from Config.
func NewClient(cfg Config) (*Client, error) {
	if cfg.TigerID == "" {
		return nil, errors.New("tiger_id is required")
	}
	if cfg.PrivateKey == "" {
		return nil, errors.New("private key is required")
	}

	priv, err := parsePrivateKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}

	if cfg.ServerURL == "" {
		cfg.ServerURL = defaultServerURL
	}
	if cfg.Charset == "" {
		cfg.Charset = defaultCharset
	}
	if cfg.SignType == "" {
		cfg.SignType = defaultSignType
	}
	if cfg.Version == "" {
		cfg.Version = defaultVersion
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 15 * time.Second
	}
	if cfg.Lang == "" {
		cfg.Lang = "en_US"
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: cfg.Timeout}
	}

	userAgent := defaultUserAgent

	return &Client{
		cfg:        cfg,
		privateKey: priv,
		httpClient: httpClient,
		userAgent:  userAgent,
	}, nil
}

// GetAssets queries account assets.
func (c *Client) GetAssets(ctx context.Context, req AssetsRequest) (*AssetsResult, error) {
	biz := req.toBiz(c.cfg)
	resp, err := c.call(ctx, "assets", biz)
	if err != nil {
		return nil, err
	}
	var wrapper assetsWrapper
	if len(resp.Data) > 0 {
		if err := json.Unmarshal(resp.Data, &wrapper); err != nil {
			return nil, fmt.Errorf("decode assets data: %w", err)
		}
		var payload AssetsData
		if err := payload.attachRawFrom(wrapper); err != nil {
			return nil, err
		}
		return &AssetsResult{Response: resp, Assets: payload}, nil
	}
	return &AssetsResult{Response: resp, Assets: AssetsData{}}, nil
}

// GetPositions queries current positions.
func (c *Client) GetPositions(ctx context.Context, req PositionsRequest) (*PositionsResult, error) {
	biz := req.toBiz(c.cfg)
	resp, err := c.call(ctx, "positions", biz)
	if err != nil {
		return nil, err
	}
	var wrapper positionsWrapper
	if len(resp.Data) > 0 {
		if err := json.Unmarshal(resp.Data, &wrapper); err != nil {
			return nil, fmt.Errorf("decode positions data: %w", err)
		}
		var payload PositionsData
		if err := payload.attachRawFrom(wrapper); err != nil {
			return nil, err
		}
		return &PositionsResult{Response: resp, Positions: payload}, nil
	}
	return &PositionsResult{Response: resp, Positions: PositionsData{}}, nil
}

// PlaceOrder submits an order and returns the global order id.
func (c *Client) PlaceOrder(ctx context.Context, order Order) (*OrderResult, error) {
	biz := order.toBiz(c.cfg)
	resp, err := c.call(ctx, "place_order", biz)
	if err != nil {
		return nil, err
	}
	var payload OrderIDData
	if len(resp.Data) > 0 {
		if err := json.Unmarshal(resp.Data, &payload); err != nil {
			return nil, fmt.Errorf("decode order response: %w", err)
		}
		payload.normalize()
	}
	result := &OrderResult{Response: resp, Order: payload}
	if payload.Code != "" && payload.Code != "0" {
		return result, fmt.Errorf("order rejected code=%s msg=%s", payload.Code, payload.Message)
	}
	return result, nil
}

// CancelOrder cancels an order by global id or account order id.
func (c *Client) CancelOrder(ctx context.Context, req CancelOrderRequest) (*OrderResult, error) {
	biz := req.toBiz(c.cfg)
	resp, err := c.call(ctx, "cancel_order", biz)
	if err != nil {
		return nil, err
	}
	var payload OrderIDData
	if len(resp.Data) > 0 {
		if err := json.Unmarshal(resp.Data, &payload); err != nil {
			return nil, fmt.Errorf("decode cancel response: %w", err)
		}
		payload.normalize()
	}
	result := &OrderResult{Response: resp, Order: payload}
	if payload.Code != "" && payload.Code != "0" {
		return result, fmt.Errorf("cancel rejected code=%s msg=%s", payload.Code, payload.Message)
	}
	return result, nil
}

func (c *Client) call(ctx context.Context, method string, biz map[string]interface{}) (APIResponse, error) {
	if biz == nil {
		biz = map[string]interface{}{}
	}

	bizContent, err := marshalBizContent(biz)
	if err != nil {
		return APIResponse{}, fmt.Errorf("marshal biz_content: %w", err)
	}

	params := map[string]interface{}{
		"method":      method,
		"version":     c.cfg.Version,
		"biz_content": bizContent,
		"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		"tiger_id":    c.cfg.TigerID,
		"charset":     c.cfg.Charset,
		"sign_type":   c.cfg.SignType,
	}
	if c.cfg.DeviceID != "" {
		params["device_id"] = c.cfg.DeviceID
	}
	if c.cfg.NotifyURL != "" {
		params["notify_url"] = c.cfg.NotifyURL
	}

	signContent, err := buildSignContent(params)
	if err != nil {
		return APIResponse{}, fmt.Errorf("build sign content: %w", err)
	}
	signature, err := signSHA1WithRSA(c.privateKey, []byte(signContent))
	if err != nil {
		return APIResponse{}, fmt.Errorf("sign content: %w", err)
	}
	params["sign"] = signature

	body, err := marshalRequestBody(params)
	if err != nil {
		return APIResponse{}, fmt.Errorf("encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.ServerURL, bytes.NewReader(body))
	if err != nil {
		return APIResponse{}, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json;charset="+c.cfg.Charset)
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "Keep-Alive")
	req.Header.Set("User-Agent", c.userAgent)
	if c.cfg.Token != "" {
		req.Header.Set("Authorization", c.cfg.Token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return APIResponse{}, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return APIResponse{}, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return APIResponse{}, fmt.Errorf("read response: %w", err)
	}

	var result APIResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return APIResponse{}, fmt.Errorf("decode response: %w", err)
	}
	result.NormalizeData()

	return result, nil
}

func parsePrivateKey(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		block = pemBlockFromBase64(key)
	}
	if block == nil {
		return nil, errors.New("unable to parse private key")
	}
	if block.Type == "PRIVATE KEY" {
		pkcs8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		if rsaKey, ok := pkcs8.(*rsa.PrivateKey); ok {
			return rsaKey, nil
		}
		return nil, errors.New("private key is not RSA")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func pemBlockFromBase64(key string) *pem.Block {
	clean := strings.TrimSpace(key)
	if clean == "" {
		return nil
	}
	var builder strings.Builder
	builder.WriteString("-----BEGIN RSA PRIVATE KEY-----\n")
	for len(clean) > 64 {
		builder.WriteString(clean[:64])
		builder.WriteString("\n")
		clean = clean[64:]
	}
	builder.WriteString(clean)
	builder.WriteString("\n-----END RSA PRIVATE KEY-----")
	block, _ := pem.Decode([]byte(builder.String()))
	return block
}

func marshalRequestBody(params map[string]interface{}) ([]byte, error) {
	encoded, err := marshalWithSpaces(params)
	if err != nil {
		return nil, err
	}
	return []byte(encoded), nil
}

// marshalBizContent matches the python sdk: sort keys and compact separators.
func marshalBizContent(biz map[string]interface{}) (string, error) {
	return marshalDeterministic(biz)
}
