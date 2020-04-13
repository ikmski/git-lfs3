package usecase

// LockService is ...
type LockService interface {
	Lock(req *LockRequest) (*LockResult, error)
	Unlock(req *UnlockRequest) (*LockResult, error)
	List(req *LockListRequest) (*LockListResult, error)
	Verify(req *LockVerifyRequest) (*LockVerifyResult, error)
}

type lockService struct {
	LockRepository LockRepository
}

func NewLockService(lockRepo LockRepository) LockService {
	return &lockService{
		LockRepository: lockRepo,
	}
}

func (s *lockService) Lock(req *LockRequest) (*LockResult, error) {

	return nil, nil
}

func (s *lockService) Unlock(req *UnlockRequest) (*LockResult, error) {

	return nil, nil
}

func (s *lockService) List(req *LockListRequest) (*LockListResult, error) {

	return nil, nil
}

func (s *lockService) Verify(req *LockVerifyRequest) (*LockVerifyResult, error) {

	return nil, nil
}
