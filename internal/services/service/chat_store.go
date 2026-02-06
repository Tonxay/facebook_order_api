package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/redis/go-redis/v9"
)

type ChatStore struct {
	redis *redis.Client
	ttl   time.Duration
}

// โครงสร้างสำหรับ Serialize ลง Redis
type SerializableMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChatStore(client *redis.Client) *ChatStore {
	return &ChatStore{
		redis: client,
		ttl:   48 * time.Hour, // เก็บไว้ 2 วันตามโจทย์
	}
}

// ดึงประวัติการคุย
func (s *ChatStore) GetHistory(ctx context.Context, senderID string) ([]*genai.Content, error) {
	val, err := s.redis.Get(ctx, "chat:"+senderID).Result()
	if err == redis.Nil {
		return []*genai.Content{}, nil
	} else if err != nil {
		return nil, err
	}

	var tempMsgs []SerializableMsg
	json.Unmarshal([]byte(val), &tempMsgs)

	var history []*genai.Content
	for _, m := range tempMsgs {
		history = append(history, &genai.Content{
			Role:  m.Role,
			Parts: []genai.Part{genai.Text(m.Content)},
		})
	}
	return history, nil
}

// บันทึกประวัติการคุย
func (s *ChatStore) SaveHistory(ctx context.Context, senderID string, history []*genai.Content) error {
	var toSave []SerializableMsg
	// เก็บแค่ 10-14 ข้อความล่าสุดพอ เพื่อประหยัด Token และ RAM
	start := 0
	if len(history) > 14 {
		start = len(history) - 14
	}

	for _, h := range history[start:] {
		for _, part := range h.Parts {
			if text, ok := part.(genai.Text); ok {
				toSave = append(toSave, SerializableMsg{
					Role:    h.Role,
					Content: string(text),
				})
			}
		}
	}

	jsonData, _ := json.Marshal(toSave)
	return s.redis.Set(ctx, "chat:"+senderID, jsonData, s.ttl).Err()
}
