package project

func init() {
	content["/cmd/conf/app.toml"] = appTomlConf()
	content["/cmd/conf/db.toml"] = dbTomlConf()
	content["/cmd/conf/redis.toml"] = redisConf()
	content["/cmd/main.go"] = mainTemplate()
}

func appTomlConf() string {
	return `[other]
listen_addr = ":8000"
service_name = "{{.PackageName}}"
trace_key = "Trace-ID"
repository_request_timeout = 5
prometheus_listen_addr = ":9090"`
}

func dbTomlConf() string {
	return `addr = "root:123123@tcp(127.0.0.1:3306)/xxxx?charset=utf8&parseTime=True&loc=Local"
max_open_conns = 16
max_idle_conns = 8
conn_max_life_time = 300

[cache]
#缓存时间
expires = 300
#地址
addr = "127.0.0.1:6379"
#密码
password = ""
#redis 库
db = 0
#重试次数, 默认不重试
max_retries = 0
#连接池大小
pool_size = 32
#读取超时时间 3秒
read_timeout = 3
#写入超时时间 3秒
write_timeout = 3
#连接空闲时间 300秒
idle_timeout = 300
#检测死连接,并清理 默认60秒
idle_check_frequency = 60
#连接最长时间，300秒
max_conn_age = 300
#如果连接池已满 等待可用连接的时间默认 8秒
pool_timeout = 8`
}

func redisConf() string {
	return `#地址
addr = "127.0.0.1:6379"
#密码
password = ""
#redis 库
db = 0
#重试次数, 默认不重试
max_retries = 0
#连接池大小
pool_size = 32
#读取超时时间 3秒
read_timeout = 3
#写入超时时间 3秒
write_timeout = 3
#连接空闲时间 300秒
idle_timeout = 300
#检测死连接,并清理 默认60秒
idle_check_frequency = 60
#连接最长时间，300秒
max_conn_age = 300
#如果连接池已满 等待可用连接的时间默认 8秒
pool_timeout = 8`
}

func mainTemplate() string {
	return `package main

	import (
		"time"
		_ "github.com/jinzhu/gorm/dialects/mysql"
		"github.com/8treenet/freedom"
		_ "{{.PackagePath}}/controllers"
		_ "{{.PackagePath}}/repositorys"
		"{{.PackagePath}}/models/config"
		"github.com/8treenet/gcache"
		"github.com/go-redis/redis"
		"github.com/jinzhu/gorm"
		"github.com/kataras/iris"
		"github.com/sirupsen/logrus"
		
	)
	
	func main() {
		install()
	
		//http2 h2c 服务
		//h2caddrRunner := freedom.CreateH2CRunner(config.Get().App.Other["listen_addr"].(string))
		addrRunner := iris.Addr(config.Get().App.Other["listen_addr"].(string))
		freedom.Run(addrRunner, config.Get().App)
	}
	
	func install() {
		//installLogrus()
		//installDatabase()
		//installRedis()
	}
	
	
	func installDatabase() {
		freedom.InstallGorm(func() (db *gorm.DB, cache gcache.Plugin) {
			conf := config.Get().DB
			var e error
			db, e = gorm.Open("mysql", conf.Addr)
			if e != nil {
				freedom.Logger().Fatal(e.Error())
			}
	
			db.DB().SetMaxIdleConns(conf.MaxIdleConns)
			db.DB().SetMaxOpenConns(conf.MaxOpenConns)
			db.DB().SetConnMaxLifetime(time.Duration(conf.ConnMaxLifeTime) * time.Second)
			
			/*
				启用缓存中间件
				cfg := config.Get().DB.Cache
				ropt := gcache.RedisOption{
					Addr:               cfg.Addr,
					Password:           cfg.Password,
					DB:                 cfg.DB,
					MaxRetries:         cfg.MaxRetries,
					PoolSize:           cfg.PoolSize,
					ReadTimeout:        time.Duration(cfg.ReadTimeout) * time.Second,
					WriteTimeout:       time.Duration(cfg.WriteTimeout) * time.Second,
					IdleTimeout:        time.Duration(cfg.IdleTimeout) * time.Second,
					IdleCheckFrequency: time.Duration(cfg.IdleCheckFrequency) * time.Second,
					MaxConnAge:         time.Duration(cfg.MaxConnAge) * time.Second,
					PoolTimeout:        time.Duration(cfg.PoolTimeout) * time.Second,
				}
				opt := gcache.DefaultOption{}
				opt.Expires = cfg.Expires      //缓存时间，默认60秒。范围 30-900
				opt.Level = gcache.LevelSearch //缓存级别，默认LevelSearch。LevelDisable:关闭缓存，LevelModel:模型缓存， LevelSearch:查询缓存
				//缓存中间件 注入到Gorm
				cache = gcache.AttachDB(db, &opt, &ropt)
			*/
			return
		})
	}
	
	func installRedis() {
		freedom.InstallRedis(func() (client *redis.Client) {
			cfg := config.Get().Redis
			opt := &redis.Options{
				Addr:               cfg.Addr,
				Password:           cfg.Password,
				DB:                 cfg.DB,
				MaxRetries:         cfg.MaxRetries,
				PoolSize:           cfg.PoolSize,
				ReadTimeout:        time.Duration(cfg.ReadTimeout) * time.Second,
				WriteTimeout:       time.Duration(cfg.WriteTimeout) * time.Second,
				IdleTimeout:        time.Duration(cfg.IdleTimeout) * time.Second,
				IdleCheckFrequency: time.Duration(cfg.IdleCheckFrequency) * time.Second,
				MaxConnAge:         time.Duration(cfg.MaxConnAge) * time.Second,
				PoolTimeout:        time.Duration(cfg.PoolTimeout) * time.Second,
			}
			client = redis.NewClient(opt)
			if e := client.Ping().Err(); e != nil {
				freedom.Logger().Fatal(e.Error())
			}
			return
		})
	}
	
	func installLogrus() {
		logrus.SetLevel(logrus.InfoLevel)
		logrus.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
		freedom.Logger().Install(logrus.StandardLogger())
	}
	`
}
