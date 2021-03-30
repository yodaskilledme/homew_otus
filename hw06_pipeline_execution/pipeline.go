package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := in

	for _, stage := range stages {
		dataCh := make(Bi)

		go func(dataCh Bi, out Out) {
			defer close(dataCh)

			for {
				select {
				case <-done:
					return
				case v, ok := <-out:
					if !ok {
						return
					}
					dataCh <- v
				}
			}
		}(dataCh, out)

		out = stage(dataCh)
	}

	return out
}
