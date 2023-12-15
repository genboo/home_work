package hw06pipelineexecution

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
	out := in
	for _, stage := range stages {
		out = executeStage(out, done, stage)
	}
	return out
}

func executeStage(in, done In, stage Stage) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
		}()
		val := stage(in)
		for {
			select {
			case v, ok := <-val:
				if !ok {
					return
				}
				out <- v
			case <-done:
				return
			}
		}
	}()
	return out
}
