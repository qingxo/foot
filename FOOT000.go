package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"tesou.io/platform/foot-parent/foot-api/common/base"
	launch2 "tesou.io/platform/foot-parent/foot-core/launch"
	"tesou.io/platform/foot-parent/foot-core/module/leisu/constants"
	"tesou.io/platform/foot-parent/foot-core/module/leisu/service"
	"tesou.io/platform/foot-parent/foot-core/module/leisu/utils"
	"tesou.io/platform/foot-parent/foot-spider/launch"
	"time"
)

func init() {

}

func main() {
	var input string
	if len(os.Args) > 1 {
		input = strings.ToLower(os.Args[1])
	} else {
		input = ""
	}

	switch input {
	case "init":
		launch2.GenTable()
		launch2.TruncateTable()
	case "spider":
		launch.Spider(4)
	case "analy":
		launch2.Analy()
	case "limit":
		pubLimitService := new(service.PubLimitService)
		publimit := pubLimitService.GetPublimit()
		bytes, _ := json.Marshal(publimit)
		fmt.Println("发布限制信息为:" + string(bytes))
	case "price":
		priceService := new(service.PriceService)
		price := priceService.GetPrice()
		bytes, _ := json.Marshal(price)
		fmt.Println("收费信息为:" + string(bytes))
	case "matchpool":
		//测试从雷速获取可发布的比赛池
		readCloser := utils.Get(constants.MATCH_LIST_URL)
		reader := bufio.NewReader(readCloser)
		for {
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				break;
			} else if err != nil {
				fmt.Println(err)
				break;
			} else {
				fmt.Println(string(line))
			}
		}
		//尝试获取比赛列表
		poolService := new(service.MatchPoolService)
		list := poolService.GetMatchList()
		for _, e := range list {
			bytes, _ := json.Marshal(e)
			fmt.Println(string(bytes))
		}
	case "pub":
		pubService := new(service.PubService)
		pubService.PubBJDC()
	case "auto":
		for {
			base.Log.Info("--------程序开始运行--------")
			//1.安装数据库
			//2.配置好数据库连接,打包程序发布
			//3.程序执行流程,周期定制定为一天三次
			//3.1 FS000Application 爬取数据
			launch.Spider(4)
			//3.2 FC002AnalyApplication 分析得出推荐列表
			launch2.Analy()
			//3.3 FW001PubApplication 执行发布到雷速
			pubService := new(service.PubService)
			pubService.PubBJDC()
			base.Log.Info("--------程序周期结束--------")
			time.Sleep(time.Duration(pubService.CycleTime()) * time.Minute)
		}
	default:
		fmt.Println("usage: init|spider|analy|pub|auto")
	}

}