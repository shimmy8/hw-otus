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
		for n := range in {
			select {
			case <-done:
				return
			default:
				stageIn <- n
			}
		}
	}()
	stageOut := stage(stageIn)
	return stageOut
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Out)

	for _, stage := range stages {
		out = ExecuteStage(in, done, stage)
		in = out
	}

	return out
}
