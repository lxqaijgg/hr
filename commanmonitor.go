package main
import (
	"os"
	"fmt"
	"time"
	"bytes"
	"errors"
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

var ddurl string="https://oapi.dingtalk.com/robot/send?access_token=3498c49e3b7406a1a0eb1f39e800a491e727581b22fe3b6154cd92cb9bae67b0"

//var outurl []string=[]string{"https://www.facebook.com","https://twitter.com","https://aws.amazon.com"}
var ca chan int=make(chan int,3)
var flag bool=true
func climbwallcheck(url string,cb chan int) {
	resp,err:=http.Get(url)
	if err!=nil{
		cb <- 50
		return
	}
	if resp.StatusCode != http.StatusOK{ 
		cb <- 50
	}else{
		cb <- 200
		chkerr(url+" 返回状态吗: "+strconv.Itoa(resp.StatusCode)+" StatusOK码: "+strconv.Itoa(http.StatusOK),errors.New(""))
	}
}

func chkerr(content string,err error){
        var cst, _ =time.LoadLocation("Asia/Shanghai")
        if  err != nil {
                openfileobj,_:=os.OpenFile("commanmonitor.log",os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
                defer openfileobj.Close()
                openfileobj.Write([]byte(strings.Join([]string{time.Now().In(cst).Format("2016-01-02-15:04:05"),">>",content,"--->",fmt.Sprintf(">>%s<<",err),"\n"},"")))        }
}

func push_alert_dd(content string){
        sent_jsonbt,_:=json.Marshal(ddjg{Msgtype:"text",Text:map[string]string{"content":content}})
        prereq,_:=http.NewRequest("POST",ddurl,bytes.NewReader(sent_jsonbt))
        prereq.Header.Set("Content-Type","application/json;charset=UTF-8")
        httpcli:=http.Client{}
        resp,_:=httpcli.Do(prereq)
        respbodybts,_:=ioutil.ReadAll(resp.Body)
        defer resp.Body.Close()
        chkerr(string(respbodybts),errors.New(""))
}

type ddjg struct{
        Msgtype         string                  `json:"msgtype"`
        Text            map[string]string       `json:"text"`
}

func main(){
	for{
		changetime:
		ta:=time.Now().Format("2006-01-02 15:04:05")
		checkfalse:
		suma:=0
		climbwallcheck("https://www.facebook.com",ca)
		climbwallcheck("https://twitter.com",ca)
		climbwallcheck("https://aws.amazon.com",ca)
		for v:=0;v<cap(ca);v++ {
			suma+= <- ca
		}
		if suma<200 {
			flag=false
//			fmt.Println("发生告警\n级别: warning\n时间: "+ta+"\n详情: 翻墙网络可能不通")
			push_alert_dd("发生告警\n级别: warning\n时间: "+ta+"\n详情: 翻墙网络可能不通")
			goto checktrue
		}
		checktrue:		
		for{
			if !flag{
				time.Sleep(10*time.Hour)
				goto checkfalse
			}else{
				time.Sleep(30*time.Second)
				goto changetime
			}
		}
	}

}
