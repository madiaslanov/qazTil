package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/madiaslanov/qazTil/internal/domain"
)

func TestProgressListFillsEveryCategory(t *testing.T) {
	studied := time.Date(2026, 9, 22, 15, 0, 0, 0, time.UTC)
	categories := &memCategories{items: []domain.Category{
		{ID: 1, NameKK: "Сәлемдесу", NameRU: "Приветствия"},
		{ID: 2, NameKK: "Отбасы", NameRU: "Семья"},
	}}
	stats := newMemProgress()
	if err := stats.Record(context.Background(), 1, true, studied); err != nil {
		t.Fatal(err)
	}
	if err := stats.Record(context.Background(), 1, false, studied); err != nil {
		t.Fatal(err)
	}

	list, err := NewProgressService(categories, stats).List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 {
		t.Fatalf("len = %d", len(list))
	}
	if list[0].CategoryNameRU != "Приветствия" || list[0].TotalAnswers != 2 || list[0].CorrectAnswers != 1 {
		t.Fatalf("first = %+v", list[0])
	}
	if list[0].LastStudiedAt == nil || !list[0].LastStudiedAt.Equal(studied) {
		t.Fatalf("last studied = %v", list[0].LastStudiedAt)
	}
	if list[1].CategoryNameRU != "Семья" || list[1].TotalAnswers != 0 || list[1].CorrectAnswers != 0 || list[1].LastStudiedAt != nil {
		t.Fatalf("second = %+v", list[1])
	}
}
