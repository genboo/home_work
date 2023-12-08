package hw06pipelineexecution

import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if len(stages) == 0 {
		out := make(Bi)
		close(out)
		return out
	}
	result := sync.Map{}
	order := make([]interface{}, 0)
	var wg sync.WaitGroup
	output := func(init interface{}, v In) {
		var r interface{}
		select {
		case r = <-pipeline(v, stages...):
			result.Store(init, r)
		case <-done:
		}
		wg.Done()
	}
	values := make(Bi)
	for v := range in {
		wg.Add(1)
		order = append(order, v)
		go output(v, values)
		values <- v
	}
	close(values)
	out := make(Bi)
	// ожидаем завершения всех пайплайнов, чтобы записать результаты в канал и закрыть его
	go func() {
		wg.Wait()
		for _, v := range order {
			val, ok := result.Load(v)
			if ok {
				out <- val
			}
		}
		close(out)
	}()
	return out
}

func pipeline(in In, stages ...Stage) Out {
	var out Out
	for _, stage := range stages {
		if out == nil {
			out = stage(in)
		} else {
			out = stage(out)
		}
	}
	return out
}
