package work

import (
	"errors"
	"fmt"
	"sync"
)

type WorkTask interface {
	Run(val any) (any, error) //执行任务
	GetTaskName() string      //获取任务名称  任务名要在引擎内唯一
}

type WorkEngine struct {
	tasks map[string]WorkTask
	plans map[string]*Plan
	l     sync.RWMutex
}

type Plan struct {
	taskName        string //任务名
	successNextPlan *Plan
	failNextPlan    *Plan
}

// 注册任务
func Register(task ...WorkTask) (*WorkEngine, error) {
	w := &WorkEngine{}
	w.tasks = make(map[string]WorkTask)
	w.plans = make(map[string]*Plan)
	for _, t := range task {
		if _, ok := w.tasks[t.GetTaskName()]; ok {
			continue
		}
		w.tasks[t.GetTaskName()] = t
	}
	return w, nil
}

// 判断任务是否已添加
func (w *WorkEngine) IsExist(task WorkTask) bool {
	_, ok := w.tasks[task.GetTaskName()]
	return ok
}

// 添加任务
func (w *WorkEngine) AddTask(task WorkTask) {
	w.tasks[task.GetTaskName()] = task
}

// 删除任务
func (w *WorkEngine) DelTask(task WorkTask) {
	delete(w.tasks, task.GetTaskName())
}

// 创建新工作计划
func (w *WorkEngine) Plan(planId string) *Plan {
	w.plans[planId] = &Plan{}
	return w.plans[planId]
}

// 插入任务
func (p *Plan) Append(taskName string) (*Plan, *Plan) {
	p.taskName = taskName
	p.successNextPlan = &Plan{}
	p.failNextPlan = &Plan{}
	return p.successNextPlan, p.failNextPlan
}

// 执行计划
func (w *WorkEngine) Do(planId string) error {
	w.l.Lock()
	defer w.l.Unlock()
	plan, ok := w.plans[planId]
	if !ok {
		return errors.New("not found plan")
	}
	var val any
	for plan != nil {
		if plan.taskName == "" {
			break
		}
		task, ok := w.tasks[plan.taskName]
		if !ok {
			return fmt.Errorf("not found task:%s", plan.taskName)
		}
		runVal, err := task.Run(val)
		if err != nil {
			if plan.failNextPlan != nil {
				plan = plan.failNextPlan
			} else {
				return fmt.Errorf("task:%s run error:%s", plan.taskName, err)
			}
		} else {
			val = runVal
			if plan.successNextPlan != nil {
				plan = plan.successNextPlan
			} else {
				return nil
			}
		}
	}
	return nil
}
