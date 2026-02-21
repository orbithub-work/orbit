package services

import "context"

func (s *ScanService) GetImportProgress() ImportProgress {
	s.mu.Lock()
	res := s.progress
	res.QueueSize = len(s.scanQueue)
	s.mu.Unlock()

	count, _ := s.assetService.assets.GetPendingCount(context.Background())
	res.PendingCount = count
	return res
}

func (s *ScanService) updateDir(index int, f func(p *ImportDirectoryProgress)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.progress.Directories) {
		return
	}
	p := s.progress.Directories[index]
	f(&p)
	s.progress.Directories[index] = p
}
