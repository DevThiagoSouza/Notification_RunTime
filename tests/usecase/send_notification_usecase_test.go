package usecase_test

import (
	"context"
	"testing"

	"apiRest/internal/entity"
	"apiRest/internal/usecase"
)

type mockRepo struct {
	saved []*entity.Notification
}

func (m *mockRepo) Save(ctx context.Context, n *entity.Notification) error {
	m.saved = append(m.saved, n)
	return nil
}

type mockPublisher struct {
	published []*entity.Notification
}

func (m *mockPublisher) Publish(ctx context.Context, n *entity.Notification) error {
	m.published = append(m.published, n)
	return nil
}

func TestSendNotificationUseCaseExecute(t *testing.T) {
	repo := &mockRepo{}
	pub := &mockPublisher{}
	uc := &usecase.SendNotificationUseCase{Repo: repo, Publisher: pub}

	n := &entity.Notification{
		ID:      "1",
		UserID:  "user-1",
		Content: "novo comentario",
		Type:    "comment",
	}

	if err := uc.Execute(context.Background(), n); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if len(repo.saved) != 1 {
		t.Fatalf("expected 1 saved notification, got %d", len(repo.saved))
	}
	if len(pub.published) != 1 {
		t.Fatalf("expected 1 published notification, got %d", len(pub.published))
	}
	if repo.saved[0].Status != entity.StatusPending {
		t.Fatalf("expected status %s, got %s", entity.StatusPending, repo.saved[0].Status)
	}
}
