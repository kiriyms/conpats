package pipeline

type stageFunc[I any, O any] func(I) O

type stage[I any, O any] struct {
	fn      stageFunc[I, O]
	workers int
}

type Pipeline[I any] struct {
	stages []any
}

func New[I any]() *Pipeline[I] {
	return &Pipeline[I]{
		stages: make([]any, 0),
	}
}

func AddStage[I, O any](p *Pipeline[I], fn func(I) O, workers int) *Pipeline[O] {
	if workers <= 0 {
		workers = 1
	}

	newStages := append(make([]any, 0), p.stages...)
	newStages = append(newStages, stage[I, O]{fn: fn, workers: workers})

	return &Pipeline[O]{
		stages: newStages,
	}
}


