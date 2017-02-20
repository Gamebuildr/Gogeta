package sourcesystem

import "testing"

func TestRepoLimitReturnsFalseWhenRepoSmall(t *testing.T) {
	repo := SourceRepository{}
	limit := repo.SizeLimitsReached(2500)

	if limit {
		t.Errorf("Expected %v, got %v", false, limit)
	}
}

func TestRepoLimitReturnsTrueWhenRepoTooLarge(t *testing.T) {
	repo := SourceRepository{}
	limitOne := repo.SizeLimitsReached(3000)
	limitTwo := repo.SizeLimitsReached(3500)

	if !limitOne {
		t.Errorf("Expected %v, got %v", true, limitOne)
	}
	if !limitTwo {
		t.Errorf("Expected %v, got %v", true, limitTwo)
	}
}
