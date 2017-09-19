package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"runtime"
	"time"
)

var verboseMode bool

// TaskBo ... 任务定义参数集
type TaskBo struct {
	Filename         string      // 文件名
	UserRegexp       string      // 用户搜索匹配模型
	SplitChars       []string    // 切割策略
	SplitResultIndex int         // 切割结果索引值
	BlockIndex       int         // 块索引
	StartLineIndex   int         // 当前任务起始行号
	LinesCount       int         // 当前任务行数
	DataCh           chan string // 数据通道
	DoneCh           chan int    // 同步通道
}

// 应用程序入口...
func main() {
	filenamePtr := flag.String("f", "", "指定一个文件作为统计对象")
	userRegexpPtr := flag.String("r", "", "指定一个正则表达式用于执行逐行搜索")
	splitPolicyPtr := flag.String("s", "", "指定切割策略，示例：-s [:]3 表示以冒号切割搜索到的字符并取索引为3的元素")
	verboseModePtr := flag.Bool("v", false, "指定使用详细模式")
	flag.Parse()

	if *filenamePtr == "" {
		log.Fatal("[错误]必须指定一个有效的文件作为统计对象")
		os.Exit(1)
	}
	if *userRegexpPtr == "" {
		log.Fatal("[错误]必须指定一个有效的正则表达式作为逐行搜索模型")
		os.Exit(2)
	}
	if *splitPolicyPtr != "" {

	}

	verboseMode = *verboseModePtr
	if verboseMode {
		log.Printf("[调试]确认文件的总行数")
	}
	lineCount, err := GetLineCount(*filenamePtr)
	if err != nil {
		log.Fatal("[错误]确认文件总行数失败，任务被迫中止")
		os.Exit(3)
	}
	if verboseMode {
		log.Printf("[调试]文件的总行数：%+v", lineCount)
	}
	if lineCount < 0 {
		log.Printf("[提示]没有有效的行数可以用于进一步进行搜索匹配")
	}

	err = Search(*filenamePtr, *userRegexpPtr, lineCount)
	if err != nil {
		log.Fatal("[错误]搜索匹配任务失败")
		os.Exit(4)
	}
	os.Exit(0)
}

// GetLineCount ... 获得文件的总行数
func GetLineCount(filename string) (lineCount int, err error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("[错误]打开文件失败，错误描述：", err)
		return 0, err
	}
	defer file.Close()

	buf := make([]byte, 64*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			log.Fatal("[错误]读取文件失败，错误描述：", err)
			return count, err
		}
	}
}

// Search ... 执行搜索
func Search(filename string, userRegexp string, lineCount int) (err error) {
	startTime := time.Now()
	blockCount := runtime.NumCPU()
	blockSize := lineCount
	if lineCount < 10000 {
		blockCount = 1
		blockSize++ // 只作为一块数据处理的时候，需要收尾
		log.Printf("[提示]小文件只需要启用1个CPU执行任务")
	} else {
		blockSize = lineCount / blockCount
		blockCount++ // 之所以要再加1，是因为最后可能有“零头”，零头不给分配独立的Go routine
		log.Printf("[提示]大文件同时启用%+v个CPU执行任务，每个CPU承担%+v行的解析任务", blockCount-1, blockSize)
	}
	log.Printf("[提示]文件名：%+v", filename)
	log.Printf("[提示]当前文件大小：%+v", GetFriendlyFileSizeDesc(filename))
	log.Printf("[提示]匹配模式：%+v", userRegexp)

	doneCh := make(chan int, blockCount)   // 带缓存功能的同步通道
	dataCh := make(chan string, blockSize) // 带缓存功能的数据通道
	encountered := make(map[string]bool)   // 去重map
	for i := 0; i < blockCount; i++ {
		taskBo := TaskBo{
			Filename:       filename,
			UserRegexp:     userRegexp,
			BlockIndex:     i,
			StartLineIndex: i * blockSize,
			LinesCount:     blockSize,
			DataCh:         dataCh,
			DoneCh:         doneCh}
		go Task(taskBo)
	}

	doneCount := 0
	for {
		select {
		case data := <-dataCh:
			encountered[data] = true
		case <-doneCh:
			doneCount++
		}
		if doneCount >= blockCount {
			break
		}
	}

	endTime := time.Now()
	log.Printf("[提示]分析完成，文件共%+v行，找到有效数据%+v个，任务耗时：%+v", lineCount+1, len(encountered), endTime.Sub(startTime))
	return nil
}

// Task ... 具体任务
func Task(taskBo TaskBo) (err error) {
	if verboseMode {
		log.Printf("[调试]正在解析，块索引：%+v，起始行：%+v，结束行：%+v",
			taskBo.BlockIndex, taskBo.StartLineIndex, taskBo.StartLineIndex+taskBo.LinesCount)
	}
	file, err := os.Open(taskBo.Filename)
	if err != nil {
		log.Fatal("[错误]打开文件失败，错误描述：", err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0

	re := regexp.MustCompile(taskBo.UserRegexp)
	for scanner.Scan() {
		lineText := scanner.Text()
		count++
		if count >= taskBo.StartLineIndex {
			if count < taskBo.LinesCount {
				if len(lineText) == 0 {
					continue
				}
				elements := re.FindAllString(lineText, -1)
				for v := range elements {
					// 将每个找到的结果都推入数据通道
					// log.Println(elements[v])
					taskBo.DataCh <- elements[v]
				}
			} else {
				break
			}
		}

		err := scanner.Err()
		switch {
		case err == io.EOF:
			break

		case err != nil:
			log.Fatal("[错误]读取文件失败，错误描述：", err)
			break
		}
	}
	taskBo.DoneCh <- 1
	return nil
}

// GetFriendlyFileSizeDesc ... 获得一个友好的文件大小描述信息
func GetFriendlyFileSizeDesc(filename string) string {
	f, _ := os.Stat(filename)
	size := f.Size()
	desc := ""
	if size < 1024*1024 {
		desc = fmt.Sprintf("%-6.3fKB", float32(size)/1024)
	} else if size < 1024*1024*1024 {
		desc = fmt.Sprintf("%-7.3fMB", float32(size)/1024/1024)
	} else {
		desc = fmt.Sprintf("%-7.3fGB", float32(size)/1024/1024/1024)
	}
	return desc
}
