package slack

import "gitlab.com/distributed_lab/logan/v3/errors"

type response struct {
	payload    interface{}
	nextCursor string
}

func (c *client) paginationWrapper(wrapperFunc func() (response, error), priority int) (*response, error) {
	item, err := addFunctionInPQueue(
		c.pqueues.BotPQueue,
		wrapperFunc,
		[]any{},
		priority,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add function in pqueue")
	}

	err = item.Response.Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to get response from api")
	}

	response, ok := item.Response.Value.(response)
	if !ok {
		return nil, errors.New("failed to convert response")
	}

	return &response, err
}
