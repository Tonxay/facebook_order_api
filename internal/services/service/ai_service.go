package service

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AIService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewAIService(ctx context.Context, apiKey string) (*AIService, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-1.5-flash")

	// --- ກຳນົດບຸກຄະລິກ (Persona) ເປັນພາສາລາວ ---
	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{
			genai.Text("ເຈົ້າແມ່ນ 'ແອດມິນກ້າມປູ' (Admin Kampoo) ຜູ້ຊ່ວຍຂາຍເຄື່ອງອອນລາຍທີ່ໜ້າຮັກ. " +
				"ພາສາຫຼັກ: ໃຫ້ຕອບເປັນພາສາລາວ (Lao Language) ເທົ່ານັ້ນ. " +
				"ບຸກຄະລິກ: ສຸພາບ, ເປັນກັນເອງ, ຍິ້ມແຍ້ມ. " +
				"ຄຳສັບທີ່ຕ້ອງໃຊ້: ໃຊ້ຄຳລົງທ້າຍວ່າ 'ເຈົ້າ' ຫຼື 'ໂດຍ' ສະເໝີເພື່ອຄວາມສຸພາບ. " +
				"ຕົວຢ່າງການຕອບ: 'ສະບາຍດີເຈົ້າ! ມີຫຍັງໃຫ້ແອດມິນຊ່ວຍເຫຼືອບໍ່ເຈົ້າ? ✨ 🦀'"),
		},
	}

	return &AIService{client: client, model: model}, nil
}

func (s *AIService) Responses(ctx context.Context, history []*genai.Content, userMsg string) (string, []*genai.Content, error) {
	session := s.model.StartChat()
	session.History = history

	resp, err := session.SendMessage(ctx, genai.Text(userMsg))
	if err != nil {
		return "", nil, err
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "ຂໍອະໄພເຈົ້າ, ລະບົບຂັດຂ້ອງໜ້ອຍໜຶ່ງ ລົບກວນລູກຄ້າພິມໃໝ່ອີກຄັ້ງໄດ້ບໍ່ເຈົ້າ? 🦀", session.History, nil
	}

	finalText := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	return finalText, session.History, nil
}
