package action

import (
	"github.com/mylxsw/adanos-alert/internal/repository"
)

type HttpAction struct {
	Histories []repository.Trigger
}

func NewHttpAction() *HttpAction {
	return &HttpAction{Histories: make([]repository.Trigger, 0)}
}

func (h *HttpAction) Handle(trigger repository.Trigger, grp repository.MessageGroup) error {
	h.Histories = append(h.Histories, trigger)
	return nil
}
