package consumer

// Route 在这里写路由
func (s *Scheduler) Route() (err error) {
	_ = s.Sub("sample-topic", s.BizCtx.DealMsg, 3, 10)
	_ = s.Sub("sample-topic1", s.BizCtx.DealMsg2, 1, 10)
	return err
}
