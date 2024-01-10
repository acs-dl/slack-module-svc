package slack

import (
	"container/heap"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/acs-dl/slack-module-svc/internal/pqueue"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

func addFunctionInPQueue(pq *pqueue.PriorityQueue, function any, functionArgs []any, priority int) (*pqueue.QueueItem, error) {
	queueItem := &pqueue.QueueItem{
		Id:       getFunctionSignature(function, functionArgs),
		Func:     function,
		Args:     functionArgs,
		Priority: priority,
	}
	heap.Push(pq, queueItem)

	item, err := pq.WaitUntilInvoked(queueItem.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to wait until invoked")
	}

	err = pq.RemoveById(queueItem.Id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to remove by id")
	}

	return item, nil
}

func getFunctionName(function interface{}) string {
	splitName := strings.Split(runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name(), ".")
	return splitName[len(splitName)-1]
}

func getFunctionSignature(function interface{}, args []interface{}) string {
	signatureParts := []string{getFunctionName(function), "("}
	signatureParts = append(signatureParts)

	for _, arg := range args {
		signatureParts = append(signatureParts, fmt.Sprintf("%v", arg))
	}

	signatureParts = append(signatureParts, ")")

	return strings.Join(signatureParts, " ")
}
