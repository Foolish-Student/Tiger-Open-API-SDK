package tigeropen

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// APIResponse 映射老虎返回的响应包结构。
type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Sign    string          `json:"sign"`
}

// Success 表示返回码是否为 0。
func (r APIResponse) Success() bool {
	return r.Code == 0
}

// NormalizeData 将字符串载荷转为原始 JSON，便于解码。
func (r *APIResponse) NormalizeData() {
	if len(r.Data) == 0 {
		return
	}
	trimmed := bytes.TrimSpace(r.Data)
	if len(trimmed) > 0 && trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(r.Data, &s); err == nil {
			r.Data = []byte(s)
		}
	}
}

type AssetsRequest struct {
	Account      string
	SubAccounts  []string
	Segment      bool
	MarketValue  bool
	BaseCurrency string
	Consolidated *bool
	SecretKey    string
	Language     string
}

func (r AssetsRequest) toBiz(cfg Config) map[string]interface{} {
	account := r.Account
	if account == "" {
		account = cfg.Account
	}
	secret := r.SecretKey
	if secret == "" {
		secret = cfg.SecretKey
	}
	lang := r.Language
	if lang == "" {
		lang = cfg.Lang
	}

	biz := map[string]interface{}{}
	if account != "" {
		biz["account"] = account
	}
	if secret != "" {
		biz["secret_key"] = secret
	}
	if r.Segment {
		biz["segment"] = r.Segment
	}
	if r.MarketValue {
		biz["market_value"] = r.MarketValue
	}
	if len(r.SubAccounts) > 0 {
		biz["sub_accounts"] = r.SubAccounts
	}
	if r.BaseCurrency != "" {
		biz["base_currency"] = r.BaseCurrency
	}
	if r.Consolidated != nil {
		biz["consolidated"] = *r.Consolidated
	}
	if lang != "" {
		biz["lang"] = lang
	}
	return biz
}

type PositionsRequest struct {
	Account        string
	Symbol         string
	SecType        string
	Currency       string
	Market         string
	SubAccounts    []string
	Expiry         string
	Strike         *float64
	PutCall        string
	AssetQuoteType string
	SecretKey      string
	Language       string
}

func (r PositionsRequest) toBiz(cfg Config) map[string]interface{} {
	account := r.Account
	if account == "" {
		account = cfg.Account
	}
	secret := r.SecretKey
	if secret == "" {
		secret = cfg.SecretKey
	}
	lang := r.Language
	if lang == "" {
		lang = cfg.Lang
	}

	biz := map[string]interface{}{}
	if account != "" {
		biz["account"] = account
	}
	if secret != "" {
		biz["secret_key"] = secret
	}
	if r.Symbol != "" {
		biz["symbol"] = r.Symbol
	}
	if r.SecType != "" {
		biz["sec_type"] = r.SecType
	}
	if r.Currency != "" {
		biz["currency"] = r.Currency
	}
	if r.Market != "" {
		biz["market"] = r.Market
	}
	if len(r.SubAccounts) > 0 {
		biz["sub_accounts"] = r.SubAccounts
	}
	if r.Expiry != "" {
		biz["expiry"] = r.Expiry
	}
	if r.Strike != nil {
		biz["strike"] = *r.Strike
	}
	if r.PutCall != "" {
		biz["right"] = r.PutCall
	}
	if r.AssetQuoteType != "" {
		biz["asset_quote_type"] = r.AssetQuoteType
	}
	if lang != "" {
		biz["lang"] = lang
	}
	return biz
}

type CancelOrderRequest struct {
	Account   string
	ID        *int64
	OrderID   *int64
	SecretKey string
	Language  string
}

func (r CancelOrderRequest) toBiz(cfg Config) map[string]interface{} {
	account := r.Account
	if account == "" {
		account = cfg.Account
	}
	secret := r.SecretKey
	if secret == "" {
		secret = cfg.SecretKey
	}
	lang := r.Language
	if lang == "" {
		lang = cfg.Lang
	}

	biz := map[string]interface{}{}
	if account != "" {
		biz["account"] = account
	}
	if secret != "" {
		biz["secret_key"] = secret
	}
	if r.OrderID != nil {
		biz["order_id"] = *r.OrderID
	}
	if r.ID != nil {
		biz["id"] = *r.ID
	}
	if lang != "" {
		biz["lang"] = lang
	}
	return biz
}

type OrdersRequest struct {
	Account       string
	SecretKey     string
	Symbol        string
	SecType       string
	Market        string
	Status        string
	StartTime     *int64
	EndTime       *int64
	Limit         *int
	NextPageToken string
	SegType       string
	Language      string
}

func (r OrdersRequest) toBiz(cfg Config) map[string]interface{} {
	account := r.Account
	if account == "" {
		account = cfg.Account
	}
	secret := r.SecretKey
	if secret == "" {
		secret = cfg.SecretKey
	}
	lang := r.Language
	if lang == "" {
		lang = cfg.Lang
	}

	biz := map[string]interface{}{}
	if account != "" {
		biz["account"] = account
	}
	if secret != "" {
		biz["secret_key"] = secret
	}
	if r.Symbol != "" {
		biz["symbol"] = r.Symbol
	}
	if r.SecType != "" {
		biz["sec_type"] = r.SecType
	}
	if r.Market != "" {
		biz["market"] = r.Market
	}
	if r.Status != "" {
		biz["status"] = r.Status
	}
	if r.StartTime != nil {
		biz["start_time"] = *r.StartTime
	}
	if r.EndTime != nil {
		biz["end_time"] = *r.EndTime
	}
	if r.Limit != nil {
		biz["limit"] = *r.Limit
	}
	if r.NextPageToken != "" {
		biz["next_page_token"] = r.NextPageToken
	}
	if r.SegType != "" {
		biz["seg_type"] = r.SegType
	}
	if lang != "" {
		biz["lang"] = lang
	}
	return biz
}

type Contract struct {
	Symbol      string
	Currency    string
	SecType     string
	Exchange    string
	LocalSymbol string
	Expiry      string
	Strike      *float64
	PutCall     string
	Multiplier  string
}

func (c Contract) toBiz() map[string]interface{} {
	biz := map[string]interface{}{}
	if c.Symbol != "" {
		biz["symbol"] = c.Symbol
	}
	if c.Currency != "" {
		biz["currency"] = c.Currency
	}
	if c.SecType != "" {
		biz["sec_type"] = c.SecType
	}
	if c.Exchange != "" {
		biz["exchange"] = c.Exchange
	}
	if c.LocalSymbol != "" {
		biz["local_symbol"] = c.LocalSymbol
	}
	if c.Expiry != "" {
		biz["expiry"] = c.Expiry
	}
	if c.Strike != nil {
		biz["strike"] = *c.Strike
	}
	if c.PutCall != "" {
		biz["right"] = c.PutCall
	}
	if c.Multiplier != "" {
		biz["multiplier"] = c.Multiplier
	}
	return biz
}

type ContractLeg struct {
	Contract
	Ratio  int
	Action string
}

func (c ContractLeg) toBiz() map[string]interface{} {
	biz := c.Contract.toBiz()
	if c.Ratio != 0 {
		biz["ratio"] = c.Ratio
	}
	if c.Action != "" {
		biz["action"] = c.Action
	}
	return biz
}

type Order struct {
	Account            string
	SecretKey          string
	Contract           Contract
	Action             string
	OrderType          string
	Quantity           float64
	QuantityScale      *int
	LimitPrice         *float64
	AuxPrice           *float64
	TrailStopPrice     *float64
	TrailingPercent    *float64
	PercentOffset      *float64
	TimeInForce        string
	OutsideRTH         *bool
	AdjustLimit        *bool
	UserMark           string
	ExpireTime         *int64
	ComboType          string
	ContractLegs       []ContractLeg
	TotalCashAmount    *float64
	TradingSessionType string
	OrderID            *int64
	ID                 *int64
	Language           string
}

func (o Order) toBiz(cfg Config) map[string]interface{} {
	account := o.Account
	if account == "" {
		account = cfg.Account
	}
	secret := o.SecretKey
	if secret == "" {
		secret = cfg.SecretKey
	}
	lang := o.Language
	if lang == "" {
		lang = cfg.Lang
	}

	biz := map[string]interface{}{}
	if account != "" {
		biz["account"] = account
	}
	if secret != "" {
		biz["secret_key"] = secret
	}
	for k, v := range o.Contract.toBiz() {
		biz[k] = v
	}
	if len(o.ContractLegs) > 0 {
		legs := make([]interface{}, 0, len(o.ContractLegs))
		for _, leg := range o.ContractLegs {
			legs = append(legs, leg.toBiz())
		}
		biz["contract_legs"] = legs
		if _, exists := biz["sec_type"]; !exists {
			biz["sec_type"] = "MLEG"
		}
	}
	if o.ID != nil {
		biz["id"] = *o.ID
	}
	if o.OrderID != nil {
		biz["order_id"] = *o.OrderID
	}
	if o.OrderType != "" {
		biz["order_type"] = o.OrderType
	}
	if o.Action != "" {
		biz["action"] = o.Action
	}
	biz["total_quantity"] = o.Quantity
	if o.QuantityScale != nil {
		biz["total_quantity_scale"] = *o.QuantityScale
	}
	if o.LimitPrice != nil {
		biz["limit_price"] = *o.LimitPrice
	}
	if o.AuxPrice != nil {
		biz["aux_price"] = *o.AuxPrice
	}
	if o.TrailStopPrice != nil {
		biz["trail_stop_price"] = *o.TrailStopPrice
	}
	if o.TrailingPercent != nil {
		biz["trailing_percent"] = *o.TrailingPercent
	}
	if o.PercentOffset != nil {
		biz["percent_offset"] = *o.PercentOffset
	}
	if o.TimeInForce != "" {
		biz["time_in_force"] = o.TimeInForce
	}
	if o.OutsideRTH != nil {
		biz["outside_rth"] = *o.OutsideRTH
	}
	if o.AdjustLimit != nil {
		biz["adjust_limit"] = *o.AdjustLimit
	}
	if o.UserMark != "" {
		biz["user_mark"] = o.UserMark
	}
	if o.ExpireTime != nil {
		biz["expire_time"] = *o.ExpireTime
	}
	if o.ComboType != "" {
		biz["combo_type"] = o.ComboType
	}
	if o.TotalCashAmount != nil {
		biz["cash_amount"] = *o.TotalCashAmount
	}
	if o.TradingSessionType != "" {
		biz["trading_session_type"] = o.TradingSessionType
	}
	if lang != "" {
		biz["lang"] = lang
	}
	return biz
}

type AssetsData struct {
	Items     []AssetItem
	IsSuccess bool
}

type AssetItem struct {
	Account            string          `json:"account"`
	Currency           string          `json:"currency,omitempty"`
	NetLiquidation     float64         `json:"netLiquidation,omitempty"`
	EquityWithLoan     float64         `json:"equityWithLoan,omitempty"`
	AvailableFunds     float64         `json:"availableFunds,omitempty"`
	BuyingPower        float64         `json:"buyingPower,omitempty"`
	Cash               float64         `json:"cash,omitempty"`
	GrossPositionValue float64         `json:"grossPositionValue,omitempty"`
	UnrealizedPnL      float64         `json:"unrealizedPnL,omitempty"`
	RealizedPnL        float64         `json:"realizedPnL,omitempty"`
	MaintMarginReq     float64         `json:"maintMarginReq,omitempty"`
	InitMarginReq      float64         `json:"initMarginReq,omitempty"`
	UpdateTime         int64           `json:"updateTime,omitempty"`
	Raw                json.RawMessage `json:"-"`
}

type PositionsData struct {
	Items     []Position
	IsSuccess bool
}

type Position struct {
	Account       string          `json:"account"`
	Symbol        string          `json:"symbol"`
	SecType       string          `json:"sec_type"`
	Currency      string          `json:"currency"`
	Market        string          `json:"market,omitempty"`
	Position      float64         `json:"position"`
	AverageCost   float64         `json:"avgCost,omitempty"`
	MarketPrice   float64         `json:"marketPrice,omitempty"`
	MarketValue   float64         `json:"marketValue,omitempty"`
	UnrealizedPnL float64         `json:"unrealizedPnL,omitempty"`
	RealizedPnL   float64         `json:"realizedPnL,omitempty"`
	UpdateTime    int64           `json:"updateTime,omitempty"`
	Raw           json.RawMessage `json:"-"`
}

// IntOrString allows numeric codes encoded as numbers or quoted strings.
type IntOrString int

func (v *IntOrString) UnmarshalJSON(data []byte) error {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 || bytes.Equal(trimmed, []byte("null")) {
		*v = 0
		return nil
	}
	if trimmed[0] == '"' {
		var s string
		if err := json.Unmarshal(trimmed, &s); err != nil {
			return err
		}
		if s == "" {
			*v = 0
			return nil
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*v = IntOrString(i)
		return nil
	}
	var i int
	if err := json.Unmarshal(trimmed, &i); err != nil {
		return err
	}
	*v = IntOrString(i)
	return nil
}

type OrderIDData struct {
	OrderID    int64           `json:"orderId,omitempty"`
	AltOrderID int64           `json:"order_id,omitempty"`
	ID         int64           `json:"id,omitempty"`
	SubIDs     []int64         `json:"subIds,omitempty"`
	Orders     json.RawMessage `json:"orders,omitempty"`
	Code       IntOrString     `json:"code,omitempty"`
	Message    string          `json:"message,omitempty"`
}

func (o *OrderIDData) normalize() {
	if o.OrderID == 0 {
		o.OrderID = o.AltOrderID
	}
}

type AssetsResult struct {
	Response APIResponse
	Assets   AssetsData
}

type PositionsResult struct {
	Response  APIResponse
	Positions PositionsData
}

type OrderResult struct {
	Response APIResponse
	Order    OrderIDData
}

type OrdersData struct {
	Items         []json.RawMessage
	NextPageToken string
	IsSuccess     bool
}

type OrdersResult struct {
	Response APIResponse
	Orders   OrdersData
}

func (a *AssetsData) attachRawFrom(wrapper assetsWrapper) error {
	a.IsSuccess = wrapper.IsSuccess
	for _, raw := range wrapper.Items {
		var item AssetItem
		if err := json.Unmarshal(raw, &item); err != nil {
			return fmt.Errorf("decode asset item: %w", err)
		}
		item.Raw = raw
		a.Items = append(a.Items, item)
	}
	return nil
}

func (p *PositionsData) attachRawFrom(wrapper positionsWrapper) error {
	p.IsSuccess = wrapper.IsSuccess
	for _, raw := range wrapper.Items {
		var item Position
		if err := json.Unmarshal(raw, &item); err != nil {
			return fmt.Errorf("decode position item: %w", err)
		}
		item.Raw = raw
		p.Items = append(p.Items, item)
	}
	return nil
}

type assetsWrapper struct {
	Items     []json.RawMessage `json:"items"`
	IsSuccess bool              `json:"is_success"`
}

type positionsWrapper struct {
	Items     []json.RawMessage `json:"items"`
	IsSuccess bool              `json:"is_success"`
}

type ordersWrapper struct {
	Items         []json.RawMessage `json:"items"`
	NextPageToken string            `json:"nextPageToken"`
	IsSuccess     bool              `json:"is_success"`
}
