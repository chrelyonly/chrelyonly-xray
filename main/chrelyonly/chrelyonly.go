package chrelyonly

import (
	"fmt"
	"github.com/xtls/xray-core/infra/conf"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

var MyGlobalConfig *conf.Config

// MyEndFunc 启动结束方法
func MyEndFunc() {
	for _, inboundConfig := range MyGlobalConfig.InboundConfigs {
		fmt.Println("协议:", inboundConfig.Protocol, ",监听地址", inboundConfig.ListenOn, ",端口:", inboundConfig.PortList.Range[0].From)
	}
	if len(MyGlobalConfig.InboundConfigs) > 0 {
		for _, inboundConfig := range MyGlobalConfig.InboundConfigs {
			if inboundConfig.Protocol == "http" {
				fmt.Println("win终端快乐设置: set http_proxy=", inboundConfig.ListenOn, ":", inboundConfig.PortList.Range[0].From, " && ", "set https_proxy=", inboundConfig.ListenOn, ":", inboundConfig.PortList.Range[0].From)
				fmt.Println("linux终端快乐设置: export http_proxy=", inboundConfig.ListenOn, ":", inboundConfig.PortList.Range[0].From, " && ", "export https_proxy=", inboundConfig.ListenOn, ":", inboundConfig.PortList.Range[0].From)
			}
		}
		fmt.Println("将快乐星球设置到上述地址既可,0.0.0.0可替换为任何可访问的ip")
	}
}

// MyInit 开始初始化接口
func MyInit() {
	// 检查文件是否存在
	if !fileExists("geosite.dat") {
		//下载geosite
		err := downloadFile("https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/geosite.dat", "geosite.dat")
		if err != nil {
			fmt.Printf("下载失败(请手动下载文件至程序运行目录): %v\n", err)
		} else {
			//fmt.Println("下载成功!")
		}
	}
	// 检查文件是否存在
	if !fileExists("geoip.dat") {
		//下载geosite
		err := downloadFile("https://cdn.jsdelivr.net/gh/Loyalsoldier/v2ray-rules-dat@release/geoip.dat", "geoip.dat")
		if err != nil {
			fmt.Printf("下载失败(请手动下载文件至程序运行目录): %v\n", err)
		} else {
			//fmt.Println("下载成功!")
		}
	}
}
func downloadFile(urlStr string, filename string) error {
	// 定义代理URL
	//proxyStr := "http://127.0.0.1:20809" // 替换为实际的代理URL
	//proxyURL, err := url.Parse(proxyStr)
	//if err != nil {
	//	fmt.Printf("解析代理URL出错: %v\n", err)
	//}
	// 自定义http.Client，使用系统代理
	client := &http.Client{
		Transport: &http.Transport{
			//Proxy: http.ProxyURL(proxyURL),
		},
	}
	//开始下载
	fmt.Println("开始下载区域数据库(若无法下载请手动下载文件至运行目录):", urlStr)
	// 发送HTTP请求
	resp, err := client.Get(urlStr)
	if err != nil {
		return fmt.Errorf("无法发送请求: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败: %s", resp.Status)
	}

	// 获取文件大小（Content-Length）
	contentLength := resp.ContentLength

	// 获取程序执行目录
	execDir, err := os.Executable()
	if err != nil {
		return fmt.Errorf("无法获取执行目录: %v", err)
	}
	execDir = filepath.Dir(execDir)

	// 创建文件
	filePath := filepath.Join(execDir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer file.Close()

	// 创建带进度显示的Reader
	progressReader := &ProgressReader{
		Reader:   resp.Body,
		Total:    contentLength,
		Filename: filename,
	}

	// 将响应内容写入文件
	_, err = io.Copy(file, progressReader)
	if err != nil {
		return fmt.Errorf("无法保存文件: %v", err)
	}

	fmt.Printf("\n文件已下载到: %s\n", filePath)
	return nil
}

// fileExists 检查文件是否存在
func fileExists(file string) bool {
	// 获取程序执行目录
	execDir, err := os.Executable()
	// 创建文件
	filePath := filepath.Join(execDir, file)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// ProgressReader wraps an io.Reader and reports the progress of the read.
type ProgressReader struct {
	Reader   io.Reader
	Total    int64
	Progress int64
	Filename string
}

func (p *ProgressReader) Read(b []byte) (int, error) {
	n, err := p.Reader.Read(b)
	if n > 0 {
		p.Progress += int64(n)
		p.printProgress()
	}
	return n, err
}

func (p *ProgressReader) printProgress() {
	percent := float64(p.Progress) / float64(p.Total) * 100
	fmt.Printf("\r下载中 [%s]: %.2f%% (%d/%d bytes)", p.Filename, percent, p.Progress, p.Total)
}
