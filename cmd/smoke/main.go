package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	tigeropen "tigeropen/src"
)

// 轻量自测工具：读取 env 文件构造客户端，顺序拉取资产与持仓做健康检查。
func main() {
	envPath := "text.env"
	if len(os.Args) > 1 {
		envPath = os.Args[1]
	}

	env, err := loadEnvFile(envPath)
	if err != nil {
		log.Fatalf("load env file: %v", err)
	}

	client, err := buildClient(env)
	if err != nil {
		log.Fatalf("init client: %v", err)
	}

	ctx := context.Background()

	assets, err := client.GetAssets(ctx, tigeropen.AssetsRequest{Segment: true, MarketValue: true})
	if err != nil {
		log.Printf("assets failed: %v", err)
	} else {
		fmt.Printf("assets code=%d items=%d\n", assets.Response.Code, len(assets.Assets.Items))
	}

	positions, err := client.GetPositions(ctx, tigeropen.PositionsRequest{})
	if err != nil {
		log.Printf("positions failed: %v", err)
	} else {
		fmt.Printf("positions code=%d items=%d\n", positions.Response.Code, len(positions.Positions.Items))
	}
}

// buildClient 根据 env 中的账号信息初始化 SDK 客户端，并自动处理区域网关。
func buildClient(env map[string]string) (*tigeropen.Client, error) {
	cfg := tigeropen.Config{
		TigerID:    env["tiger_id"],
		Account:    env["account"],
		SecretKey:  env["secret_key"],
		PrivateKey: env["private_key_pk1"],
		Lang:       "zh_CN",
	}

	if cfg.PrivateKey == "" {
		cfg.PrivateKey = env["private_key_pk8"]
	}

	// 部分 license (如 TBNZ/TBSG) 需要访问区域网关。
	if envURL := env["server_url"]; envURL != "" {
		cfg.ServerURL = envURL
	} else if strings.EqualFold(env["license"], "TBNZ") || strings.EqualFold(env["license"], "TBSG") {
		cfg.ServerURL = "https://openapi.tigerfintech.com/hkg/gateway"
	}

	return tigeropen.NewClient(cfg)
}

// loadEnvFile 解析简单的 key=value 文本，忽略空行和注释。
func loadEnvFile(path string) (map[string]string, error) {
	out := make(map[string]string)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		out[key] = val
	}
	return out, nil
}
