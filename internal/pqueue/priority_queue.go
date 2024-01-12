package pqueue

import (
	"context"
	"errors"
	"sync"
	"time"

	"gitlab.com/distributed_lab/logan/v3"
)

type PriorityQueueInterface interface {
	WaitUntilInvoked(id string) (*QueueItem, error)
	ProcessQueue(requestLimit int64, timeLimit time.Duration, stop chan struct{})
}

type PQueues struct {
	UserPQueue *PriorityQueue
	BotPQueue  *PriorityQueue
}

func NewPQueues(log *logan.Entry) PQueues {
	return PQueues{
		UserPQueue: NewPriorityQueue(log).(*PriorityQueue),
		BotPQueue:  NewPriorityQueue(log).(*PriorityQueue),
	}
}

type PriorityQueue struct {
	queueArray []*QueueItem
	queueMap   sync.Map
	log        *logan.Entry
}

func NewPriorityQueue(log *logan.Entry) PriorityQueueInterface {
	return &PriorityQueue{
		queueArray: make([]*QueueItem, 0),
		log:        log,
	}
}

func (pq *PriorityQueue) Len() int { return len(pq.queueArray) }

func (pq *PriorityQueue) Less(i, j int) bool {
	return pq.queueArray[i].Priority > pq.queueArray[j].Priority
}

func (pq *PriorityQueue) Swap(i, j int) {
	pq.queueArray[i], pq.queueArray[j] = pq.queueArray[j], pq.queueArray[i]
	pq.queueArray[i].index = i
	pq.queueArray[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*QueueItem)

	pqMapItem, exists := pq.queueMap.Load(item.Id)
	if exists {
		pqItem, ok := pqMapItem.(*QueueItem)
		if !ok {
			panic("failed to convert map element")
		}
		pqItem.Amount++
		return
	}

	n := len(pq.queueArray)
	item.index = n
	item.invoked = PROCESSING
	item.Amount++
	pq.queueArray = append(pq.queueArray, item)
	pq.queueMap.Store(item.Id, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.queueArray
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	pq.queueArray = old[0 : n-1]
	pq.queueMap.Delete(item.Id)

	return item
}

func (pq *PriorityQueue) RemoveById(id string) error {
	item, err := pq.getElement(id)
	if err != nil {
		return err
	}

	if item.Amount > 1 {
		item.Amount--
		return nil
	}

	pq.queueArray = append(pq.queueArray[:item.index], pq.queueArray[item.index+1:]...)
	pq.queueMap.Delete(item.Id)

	pq.FixIndexesInPQueue()
	return nil
}

func (pq *PriorityQueue) FixIndexesInPQueue() {
	for i, queueItem := range pq.queueArray {
		if queueItem.index != i {
			queueItem.index = i
		}
	}
}

func (pq *PriorityQueue) getElement(id string) (*QueueItem, error) {
	pqMapItem, exists := pq.queueMap.Load(id)
	if !exists {
		return nil, errors.New("element not found")
	}

	pqItem, ok := pqMapItem.(*QueueItem)
	if !ok {
		return nil, errors.New("failed to convert map element")
	}

	return pqItem, nil
}

func (pq *PriorityQueue) WaitUntilInvoked(id string) (*QueueItem, error) {
	pq.log.Infof("waiting until invoked for `%s`", id)

	item, err := pq.getElement(id)
	if err != nil {
		return nil, err
	}

	item.waitInvoked()

	return item, nil
}

func (pq *PriorityQueue) ProcessQueue(requestLimit int64, timeLimit time.Duration, stop chan struct{}) {
	ticker := time.NewTicker(timeLimit / time.Duration(requestLimit))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pq.processNextItem()
		case <-stop:
			return
		}
	}
}

func (pq *PriorityQueue) processNextItem() {
	for i := 0; i < pq.Len(); i++ {
		item := pq.queueArray[i]
		if item == nil {
			continue
		}

		if item.invoked == INVOKED {
			continue
		}

		item.callFunction()
		return
	}
}

func PQueuesInstance(ctx context.Context) *PQueues {
	return ctx.Value("pqueues").(*PQueues)
}

func CtxPQueues(entry interface{}, ctx context.Context) context.Context {
	return context.WithValue(ctx, "pqueues", entry)
}
