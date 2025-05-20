package config

// System 系统配置
type System struct {
	Env           string `mapstructure:"env" json:"env" yaml:"env"`                                  // 环境
	Port          string `mapstructure:"port" json:"port" yaml:"port"`                               // 端口
	DbType        string `mapstructure:"db_type" json:"db_type" yaml:"db_type"`                      // 数据库类型
	UseRedis      bool   `mapstructure:"use_redis" json:"use_redis" yaml:"use_redis"`                // 使用redis
	UseMultipoint bool   `mapstructure:"use_multipoint" json:"use_multipoint" yaml:"use_multipoint"` // 多点登录拦截
	OssType       string `mapstructure:"oss_type" json:"oss_type" yaml:"oss_type"`                   // 存储类型
	UseHttps      bool   `mapstructure:"use_https" json:"use_https" yaml:"use_https"`                // 使用https
	JwtSecret     string `mapstructure:"jwt_secret" json:"jwt_secret" yaml:"jwt_secret"`             // jwt密钥
	JwtExpire     int    `mapstructure:"jwt_expire" json:"jwt_expire" yaml:"jwt_expire"`             // jwt过期时间
}

// Log 日志配置
type Log struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"`                            // 级别
	Format       string `mapstructure:"format" json:"format" yaml:"format"`                         // 输出格式
	Prefix       string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`                         // 日志前缀
	Director     string `mapstructure:"director" json:"director"  yaml:"director"`                  // 日志文件夹
	ShowLine     bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`                // 显示行
	LogInConsole bool   `mapstructure:"log_in_console" json:"log_in_console" yaml:"log_in_console"` // 输出控制台
	MaxSize      int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`                   // 日志文件大小
	MaxAge       int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`                      // 日志文件保存天数
	MaxBackups   int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`          // 日志文件保存数量
}

// Mysql 数据库配置
type Mysql struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`                               // 服务器地址
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               // 端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                         // 高级配置
	DbName       string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`                      // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 数据库密码
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"` // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`                   // 是否开启Gorm全局日志
	LogZap       bool   `mapstructure:"log_zap" json:"log_zap" yaml:"log_zap"`                      // 是否通过zap写入日志文件
}

// Redis 配置
type Redis struct {
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                      // redis的哪个数据库
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`                // 服务器地址:端口
	Password string `mapstructure:"password" json:"password" yaml:"password"`    // 密码
	PoolSize int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"` // 连接池大小
}

// CORS 跨域配置
type CORS struct {
	Mode             string   `mapstructure:"mode" json:"mode" yaml:"mode"`                                        // 模式
	Whitelist        []string `mapstructure:"whitelist" json:"whitelist" yaml:"whitelist"`                         // 白名单
	AllowCredentials bool     `mapstructure:"allow_credentials" json:"allow_credentials" yaml:"allow_credentials"` // 允许携带cookie
	AllowMethods     string   `mapstructure:"allow_methods" json:"allow_methods" yaml:"allow_methods"`             // 允许方法
	AllowHeaders     string   `mapstructure:"allow_headers" json:"allow_headers" yaml:"allow_headers"`             // 允许头部
	ExposeHeaders    string   `mapstructure:"expose_headers" json:"expose_headers" yaml:"expose_headers"`          // 暴露头部
	MaxAge           int      `mapstructure:"max_age" json:"max_age" yaml:"max_age"`                               // 预检请求有效期
}

// Local 本地上传配置
type Local struct {
	Path      string `mapstructure:"path" json:"path" yaml:"path"`                   // 本地文件访问路径
	StorePath string `mapstructure:"store_path" json:"store_path" yaml:"store_path"` // 本地文件存储路径
}

// Qiniu 七牛云配置
type Qiniu struct {
	Zone          string `mapstructure:"zone" json:"zone" yaml:"zone"`                                  // 存储区域
	Bucket        string `mapstructure:"bucket" json:"bucket" yaml:"bucket"`                            // 空间名称
	ImgPath       string `mapstructure:"img_path" json:"img_path" yaml:"img_path"`                      // CDN加速域名
	UseHTTPS      bool   `mapstructure:"use_https" json:"use_https" yaml:"use_https"`                   // 是否使用https
	AccessKey     string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`                // 秘钥AK
	SecretKey     string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`                // 秘钥SK
	UseCdnDomains bool   `mapstructure:"use_cdn_domains" json:"use_cdn_domains" yaml:"use_cdn_domains"` // 上传是否使用CDN上传加速
}

// Upload 上传配置
type Upload struct {
	Local Local `mapstructure:"local" json:"local" yaml:"local"`
	Qiniu Qiniu `mapstructure:"qiniu" json:"qiniu" yaml:"qiniu"`
}

// Configuration 总配置结构
type Configuration struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Log    Log    `mapstructure:"log" json:"log" yaml:"log"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	CORS   CORS   `mapstructure:"cors" json:"cors" yaml:"cors"`
	Upload Upload `mapstructure:"upload" json:"upload" yaml:"upload"`
}
