package config

type Config struct {
	Mysql      Mysql      `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	SystemInfo SystemInfo `mapstructure:"system-info" json:"systemInfo" yaml:"system-info"`
	Mixin      Mixin      `mapstructure:"mixin" json:"mixin" yaml:"mixin"`
	DApp       DApp       `mapstructure:"dapp" json:"dapp" yaml:"dapp"`
	Zap        Zap        `mapstructure:"zap" json:"zap" yaml:"zap"`
	Deribit    Deribit    `mapstructure:"deribit" json:"deribit" yaml:"deribit"`
	Redis      Redis      `mapstructure:"redis" json:"redis" yaml:"redis"`
}

type SystemInfo struct {
	ServerPort string `mapstructure:"server-port" json:"server" yaml:"server-port"`
	Domain     string `mapstructure:"domain" json:"domain" yaml:"domain"`
	Api        string `mapstructure:"api" json:"api" yaml:"api"`
	System     string `mapstructure:"system" json:"system" yaml:"system"`
}

type Mysql struct {
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Path         string `mapstructure:"path" json:"path" yaml:"path"`
	Dbname       string `mapstructure:"db-name" json:"dbname" yaml:"db-name"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`
	LogLevel     int    `mapstructure:"log-level" json:"logLevel" yaml:"log-level"`
	CACert       string `mapstructure:"ca-cert"`
	ClientCert   string `mapstructure:"client-cert"`
	ClientKey    string `mapstructure:"client-key"`
}

type Mixin struct {
	Pin        string   `mapstructure:"pin" json:"pin" yaml:"pin"`
	ClientId   string   `mapstructure:"client-id" json:"client_id" yaml:"client-id"`
	SessionId  string   `mapstructure:"session-id" json:"session_id" yaml:"session-id"`
	PinToken   string   `mapstructure:"pin-token" json:"pin_token" yaml:"pin-token"`
	PrivateKey string   `mapstructure:"private-key" json:"private_key" yaml:"private-key"`
	Receivers  []string `mapstructure:"receivers" json:"receivers" yaml:"receivers"`
}

type Zap struct {
	Level         string `mapstructure:"level" json:"level" yaml:"level"`
	Filename      string `mapstructure:"filename" json:"filename" yaml:"filename"`
	AdminFilename string `mapstructure:"admin-filename" json:"admin-filename" yaml:"admin-filename"`
	MaxSize       int    `mapstructure:"max-size" json:"max-size" yaml:"max-size"`
	MaxAge        int    `mapstructure:"max-age" json:"max-age" yaml:"max-age"`
	MaxBackups    int    `mapstructure:"max-backups" json:"max-backups" yaml:"max-backups"`
	ConsoleOut    bool   `mapstructure:"console-out" json:"console-out" yaml:"console-out"`
}

type Deribit struct {
	ClientId               string `mapstructure:"client_id" json:"client_id" yaml:"client_id"`
	Secret                 string `mapstructure:"secret" json:"secret" yaml:"secret"`
	SyncMarketInterval     int64  `mapstructure:"sync_market_price_interval" json:"sync_market_price_interval" yaml:"sync_market_price_interval"`
	SyncIndexPriceInterval int64  `mapstructure:"sync_index_price_interval" json:"sync_index_price_interval" yaml:"sync_index_price_interval"`
}

type Redis struct {
	Address  string `mapstructure:"address" json:"address" yaml:"address"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Db       int    `mapstructure:"db" json:"db" yaml:"db"`
}

type DApp struct {
	AppID               string   `yaml:"app_id" json:"app_id" mapstructure:"app_id"`
	SessionID           string   `yaml:"session_id" json:"session_id" mapstructure:"session_id"`
	Secret              string   `yaml:"secret" json:"secret" mapstructure:"secret"`
	Pin                 string   `yaml:"pin" json:"pin" mapstructure:"pin"`
	PinToken            string   `yaml:"pin_token" json:"pin_token" mapstructure:"pin_token"`
	PrivateKey          string   `yaml:"private_key" json:"private_key" mapstructure:"private_key"`
	Receivers           []string `yaml:"receivers" json:"receivers" mapstructure:"receivers"`
	Master              bool     `yaml:"master" json:"master" mapstructure:"master"`
	Assets              Assets   `yaml:"assets" json:"assets" mapstructure:"assets"`
	Threshold           int64    `yaml:"threshold"  json:"threshold" mapstructure:"threshold"`
	GroupConversationId string   `yaml:"group_conversation_id"  json:"group_conversation_id" mapstructure:"group_conversation_id"`
}

type Assets struct {
	BTC  string `yaml:"btc" json:"btc" mapstructure:"btc"`
	XIN  string `yaml:"xin" json:"xin" mapstructure:"xin"`
	USDT string `yaml:"usdt" json:"usdt" mapstructure:"usdt"`
	USDC string `yaml:"usdc" json:"usdc" mapstructure:"usdc"`
	PUSD string `yaml:"pusd" json:"pusd" mapstructure:"pusd"`
	ETH  string `yaml:"eth" json:"eth" mapstructure:"eth"`
	CNB  string `yaml:"cnb" json:"cnb" mapstructure:"cnb"`
}
