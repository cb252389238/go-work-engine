# 这是一个工作流任务编排库|This is a workflow task orchestration library.

### 安装|install
```go
go get github.com/cb252389238/go-work-engine
```

### 使用说明|direction
```go
package main

import (
	"fmt"
	"github.com/cb252389238/go-work-engine"
	"log"
)

type task1 struct {
	name string
}

func (t *task1) GetTaskName() string {
	return t.name
}
func (t *task1) Run(val any) (any, error) {
	fmt.Println("task1")
	return val, nil
}

type task2 struct {
	name string
}

func (t *task2) GetTaskName() string {
	return t.name
}
func (t *task2) Run(val any) (any, error) {
	fmt.Println("task2")
	return val, nil
}

type task3 struct {
	name string
}

func (t *task3) GetTaskName() string {
	return t.name
}
func (t *task3) Run(val any) (any, error) {
	fmt.Println("task3")
	return val, nil
}

type task4 struct {
	name string
}

func (t *task4) GetTaskName() string {
	return t.name
}
func (t *task4) Run(val any) (any, error) {
	fmt.Println("task4")
	return val, nil
}

func main() {
	//注册任务
	//任务必实现WorkTask接口
	t1 := &task1{name: "task1"}
	t2 := &task2{name: "task2"}
	t3 := &task3{name: "task3"}
	w, err := work.Register(t1, t2, t3)
	if err != nil {
		log.Fatal(err)
	}
	t4 := &task4{name: "task4"}
	//临时添加新的任务
	w.AddTask(t4)
	//创建执行计划  planId必须是唯一的，否则计划任务会覆盖
	p := w.Plan("plan1")
	//往计划内添加任务,返回成功节点和失败节点，可以分别在两个分支节点上分别做后续操作
	okNode, failNode := p.Append(t1.GetTaskName())
	{
		//失败节点
		failNode.Append(t2.GetTaskName())
	}
	{
		//成功节点
		okNode, _ = okNode.Append(t3.GetTaskName())
		{
			okNode.Append(t4.GetTaskName())
		}
	}
	err = w.Do("plan1")
	if err != nil {
		log.Fatal(err)
	}
}


```