package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecuteStage(in In, done In, stage Stage) Out {
	stageIn := make(Bi)
	go func() {
		defer close(stageIn)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				stageIn <- val
			}
		}
	}()
	stageOut := stage(stageIn)
	return stageOut
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	var out Out

	for _, stage := range stages {
		out = ExecuteStage(in, done, stage)
		in = out
	}

	return out
}
