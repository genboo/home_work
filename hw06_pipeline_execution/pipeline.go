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
	output := func(init interface{}) {
		r := <-pipeline(init, done, stages...)
		if r != nil {
			result.Store(init, r)
		}
		wg.Done()
	}
	for v := range in {
		wg.Add(1)
		order = append(order, v)
		go output(v)
	}
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

func pipeline(r interface{}, done In, stages ...Stage) Out {
	out := make(Bi)
	go func() {
		for _, stage := range stages {
			stageChn := make(Bi)
			go func() {
				stageChn <- r
				close(stageChn)
			}()
			select {
			case r = <-stage(stageChn):
				continue
			case <-done:
				close(out)
				return
			}
		}
		out <- r
		close(out)
	}()
	return out
}
