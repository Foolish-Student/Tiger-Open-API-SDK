# Tiger OpenAPI Go SDK (minimal)

受 `openapi-python-sdk` 启发的轻量 Go 版本，只覆盖了基础交易接口：

- 获取资产 `GetAssets`
- 获取持仓 `GetPositions`
- 下单 `PlaceOrder`
- 撤单 `CancelOrder`

签名、`biz_content` 组装规则与 Python SDK 保持一致（RSA+SHA1，按参数排序拼接后签名）。

## 安装

```bash
go get ./src
```

或在其他项目的 `go.mod` 引入本仓库：

```bash
require tigeropen v0.0.0
replace tigeropen => ../Tiger-Open-API-SDK
```

## 快速开始

```go
package main

import (
	"context"
	"fmt"

	"tigeropen/src"
)

func main() {
	cfg := src.Config{
		TigerID:    "your-tiger-id",
		Account:    "your-account",
		SecretKey:  "your-secret-key",
		PrivateKey: `-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQ...
-----END RSA PRIVATE KEY-----`,
		// 可选：ServerURL, DeviceID, NotifyURL, Token, Lang
	}

	client, err := src.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// 资产
	assets, err := client.GetAssets(ctx, src.AssetsRequest{Segment: true, MarketValue: true})
	if err != nil {
		panic(err)
	}
	fmt.Printf("assets code=%d items=%d\n", assets.Response.Code, len(assets.Assets.Items))

	// 持仓
	positions, err := client.GetPositions(ctx, src.PositionsRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("positions %d\n", len(positions.Positions.Items))

	// 下单（示例：美股限价买入）
	limit := 1.23
	order, err := client.PlaceOrder(ctx, src.Order{
		Action:     "BUY",
		OrderType:  "LMT",
		Quantity:   10,
		LimitPrice: &limit,
		TimeInForce:"DAY",
		Contract: src.Contract{
			Symbol:  "AAPL",
			SecType: "STK",
			Currency:"USD",
		},
	})
	if err != nil {
		// err 会包含 code/message，仍可从 order.Order 查看返回数据
		fmt.Println("place order failed:", err)
	}
	fmt.Println("order id:", order.Order.ID)

	// 撤单
	if order.Order.ID != 0 {
		_, _ = client.CancelOrder(ctx, tigeropen.CancelOrderRequest{ID: &order.Order.ID})
	}
}
```

### 字段对照

- `Order`、`Contract`、`CancelOrderRequest` 的字段名与 Python SDK 中的 `PlaceModifyOrderParams`/`CancelOrderParams` 一致。
- `Language`/`Lang` 默认 `en_US`，需要其它语言自行覆盖。
- `AssetsData.Items[i].Raw`、`PositionsData.Items[i].Raw` 保留了原始响应片段，方便自行解析更多字段。

## 注意事项

- 仅实现基础接口，行情推送、组合单等高级能力暂未覆盖。
- 请求签名使用 RSA+SHA1，与官方文档一致，确保私钥与虎 ID、账号配置正确。
- 服务端返回的 `code`（Envelope 或 data 内的 code）不为 0 时会返回 error，但同时会附带已解析的响应数据，便于调试。
