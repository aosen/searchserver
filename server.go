package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"searchserver/routers"
	"strconv"

	"github.com/aosen/kernel"
	"github.com/aosen/search"
	"github.com/aosen/search/indexer"
	"github.com/aosen/search/pipeline"
	"github.com/aosen/search/ranker"
	"github.com/aosen/search/segmenter"
	_ "github.com/go-sql-driver/mysql"
	"github.com/larspensjo/config"
	"gopkg.in/redis.v3"
)

//所有请求都会经过的开始处理方法
func initHandler(w http.ResponseWriter, r *http.Request, g kernel.G) {
	log.Printf("%s %s ip: %s start handler", r.Method, r.URL, r.RemoteAddr)
}

//所有请求都会经过的结束处理方法
func endHandler(w http.ResponseWriter, r *http.Request, g kernel.G) {
	log.Printf("%s %s ip: %s end handler", r.Method, r.URL, r.RemoteAddr)
}

//如果没有对应处理方法，则调用此方法
func defaultHandler(w http.ResponseWriter, r *http.Request, g kernel.G) {
}

var configFile = flag.String("configfile", "conf/app.ini", "General configuration file")
var settings = make(map[string]interface{})

func loadConf() map[string]interface{} {
	flag.Parse()
	//set config file std
	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		log.Fatalf("Fail to find", *configFile, err)
	}
	if cfg.HasSection("topic") {
		section, err := cfg.SectionOptions("topic")
		if err == nil {
			for _, v := range section {
				options, err := cfg.String("topic", v)
				if err == nil {
					settings[v] = options
				} else {
					log.Fatal(err)
				}
			}
		} else {
			log.Fatal(err)
		}
	}
	return settings
}

var db *sql.DB

func mountDB() *sql.DB {
	dbinfo, err := kernel.GetSetting(settings, "DBINFO")
	checkError(err)
	db, _ = sql.Open("mysql", dbinfo)
	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(30)
	checkError(db.Ping())
	return db
}

var cc *redis.Client

func mountRedis() *redis.Client {
	addr, err := kernel.GetSetting(settings, "REDISADDR")
	checkError(err)
	passwd, err := kernel.GetSetting(settings, "REDISPASSWORD")
	checkError(err)
	redisdb, err := kernel.GetSetting(settings, "REDISDB")
	checkError(err)
	db, err := strconv.ParseInt(redisdb, 10, 64)
	checkError(err)
	cc = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: passwd, // no password set
		DB:       db,     // use default DB
	})

	_, e := cc.Ping().Result()
	checkError(e)
	return cc
}

//全局变量
var g kernel.G

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//加载配置文件
	loadConf()
	debug, e := kernel.GetSetting(settings, "DEBUG")
	checkError(e)
	//挂载数据库
	mountDB()
	//挂载redis
	mountRedis()
	dict, e1 := kernel.GetSetting(settings, "DICT")
	checkError(e1)
	stop, e2 := kernel.GetSetting(settings, "STOP")
	checkError(e2)
	//初始化分词, 采用引擎自带分词器
	seg := segmenter.InitChinaCut(dict)
	//初始化搜索引擎
	var searcher search.Engine
	//获取mongodb集合数
	indextmp, e3 := kernel.GetSetting(settings, "INDEXSTORENUM")
	checkError(e3)
	indexstorenum, err := strconv.Atoi(indextmp)
	checkError(err)
	//获取mongodb库名
	mongodbname, e4 := kernel.GetSetting(settings, "MONGODBNAME")
	checkError(e4)
	//获取mongodb连接url
	mongodburl, e5 := kernel.GetSetting(settings, "MONGODBURL")
	checkError(e5)
	//获取mongodb集合前缀
	collectionprefix, e5 := kernel.GetSetting(settings, "COLLECTIONPREFIX")
	checkError(e5)
	searcher.Init(search.EngineInitOptions{
		//分词器采用引擎自带的分词器
		Segmenter:     seg,
		StopTokenFile: stop,
		UsePersistentStorage: func() bool {
			if debug == "True" {
				return false
			} else {
				return true
			}
		}(),
		IndexerInitOptions: &search.IndexerInitOptions{
			IndexType: search.LocationsIndex,
			BM25Parameters: &search.BM25Parameters{
				K1: 2.0,
				B:  0.75,
			},
		},
		//pipline采用引擎自带的mongo pipline
		SearchPipline: pipeline.InitMongo(
			mongodbname,
			indexstorenum,
			mongodburl,
			collectionprefix),
		//索引器接口实现，采用自带的wukong索引器
		CreateIndexer: func() search.SearchIndexer {
			return indexer.NewWuKongIndexer()
		},
		//排序器生成方法
		CreateRanker: func() search.SearchRanker {
			return ranker.NewWuKongRanker()
		},
	})
	g = kernel.G{
		//可以处理的http方法字典
		Ml: map[string]string{
			"GET":     "Get",
			"POST":    "Post",
			"OPTIONS": "Options",
			"HEAD":    "Head",
			"PUT":     "Put",
			"DELETE":  "Delete",
			"CONNECT": "Connect",
		},
		Init:           initHandler,
		DefaultHandler: defaultHandler,
		End:            endHandler,
		DB:             db,
		CC:             cc,
		Settings:       settings,
		DIY: map[string]interface{}{
			"seg":      seg,
			"searcher": searcher,
		},
	}
}

func main() {
	port, err := kernel.GetSetting(settings, "PORT")
	checkError(err)
	host, err := kernel.GetSetting(settings, "HOST")
	checkError(err)
	log.Printf("server run on %s:%s", host, port)
	http.Handle("/", routers.Register(&g))
	err = http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
