package types

import (
	"github.com/golang/protobuf/proto"
)

// CalculateTaskDependencyGraph combines the static workflow definition with the dynamic invocation to determine
// the current state of the workflow.
func CalculateTaskDependencyGraph(wf *Workflow, invoc *WorkflowInvocation) map[string]*TaskStatus {
	target := map[string]*TaskStatus{}

	addStaticTasks(wf, invoc, target)

	// Add dynamic tasks overriding similarly named static tasks in case of conflicts
	addDynamicTasks(invoc, target)

	return target
}

func addStaticTasks(wf *Workflow, invoc *WorkflowInvocation, target map[string]*TaskStatus) {
	for id, task := range wf.Spec.Tasks {
		taskStatus := invoc.Status.Tasks[id]
		target[id] = &TaskStatus{
			Task:       task,
			Invocation: taskStatus,
		}
	}
}

func addDynamicTasks(invoc *WorkflowInvocation, target map[string]*TaskStatus) {
	mapping := map[string]string{}
	for id, task := range invoc.Status.DynamicTasks {
		taskStatus := invoc.Status.Tasks[id]
		target[id] = &TaskStatus{
			Task:       task,
			Invocation: taskStatus,
		}

		for depId, depParams := range task.Requires {
			if depParams.Type == TaskDependencyParameters_DYNAMIC_OUTPUT {
				mapping[depId] = id
			}
		}
	}

	// Reroute dependencies to also depend on the outputted task of dynamic tasks.
	for taskId, task := range target {
		for depId, depParams := range task.Requires {
			if outputTask, ok := mapping[depId]; ok && depParams.Type != TaskDependencyParameters_DYNAMIC_OUTPUT {
				cloned := proto.Clone(task.Task) // TODO maybe clone all the things
				t, ok := cloned.(*Task)
				if !ok {
					panic("Invalid clone")
				}
				t.Requires[outputTask] = &TaskDependencyParameters{}
				target[taskId].Task = t
			}
		}
	}
}

type TaskStatus struct {
	*Task
	// Invocation is nil if not yet invoked
	Invocation *TaskInvocation
}

// Tasks gets both static as well as dynamic tasks of a workflow invocation.
func Tasks(wf *Workflow, wfi *WorkflowInvocation) map[string]*Task {
	tasks := map[string]*Task{}
	for id, task := range wf.Spec.Tasks {
		tasks[id] = task
	}
	for id, task := range wfi.Status.DynamicTasks {
		tasks[id] = task
	}
	return tasks
}
