package behavior

import "delivery/domain"

type Command interface {
	Execute()
	Undo()
	Name() string
}

type ChangeStatusCommand struct {
	delivery *domain.Delivery
	next     domain.DeliveryStatus
	prev     domain.DeliveryStatus
}

func NewChangeStatusCommand(d *domain.Delivery, next domain.DeliveryStatus) *ChangeStatusCommand {
	return &ChangeStatusCommand{delivery: d, next: next}
}

func (c *ChangeStatusCommand) Execute() {
	c.prev = c.delivery.Status
	c.delivery.Status = c.next
}
func (c *ChangeStatusCommand) Undo() { c.delivery.Status = c.prev }
func (c *ChangeStatusCommand) Name() string { return "ChangeStatus→" + string(c.next) }

type History struct{ commands []Command }

func (h *History) Run(c Command) { c.Execute(); h.commands = append(h.commands, c) }
func (h *History) UndoLast() Command {
	if len(h.commands) == 0 {
		return nil
	}
	last := h.commands[len(h.commands)-1]
	h.commands = h.commands[:len(h.commands)-1]
	last.Undo()
	return last
}
