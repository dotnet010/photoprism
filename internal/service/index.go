package service

import (
	"sync"

	"github.com/photoprism/photoprism/internal/photoprism"
)

var onceIndex sync.Once

func initIndex() {
	services.Index = photoprism.NewIndex(Config(), Classify(), NsfwDetector(), Convert())
}

func Index() *photoprism.Index {
	onceIndex.Do(initIndex)

	return services.Index
}
