package state

type DefaultState struct{}

func (s *DefaultState) Enter()         {}
func (s *DefaultState) ExecuteBefore() {}
func (s *DefaultState) Execute()       {}
func (s *DefaultState) ExecuteAfter()  {}
func (s *DefaultState) End()           {}
