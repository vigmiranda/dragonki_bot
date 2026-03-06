package telegram

import (
	"bytes"
	"encoding/base64"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vigmiranda/dragonki/internal/application/view"
)

func TestToTelegramResponse_AssetImage(t *testing.T) {
	resp := view.GameResponse{
		Text:    "ok",
		MenuKey: "MAIN_MENU",
		Image: &view.Image{
			Type: view.ImageAsset,
			Path: "assets/a.png",
		},
	}

	got := ToTelegramResponse(resp, true, "test")

	if got.Message != "ok" || !got.ShowMenu || got.CameFrom != "test" {
		t.Fatalf("unexpected response metadata: %+v", got)
	}
	if _, ok := got.Img.(tgbotapi.FilePath); !ok {
		t.Fatalf("expected FilePath image, got %T", got.Img)
	}
	if got.Menu == nil {
		t.Fatalf("expected menu mapped")
	}
}

func TestToTelegramResponse_Base64Image(t *testing.T) {
	raw := []byte("dragonki")
	resp := view.GameResponse{
		Text:    "ok",
		MenuKey: "MAIN_MENU",
		Image: &view.Image{
			Type:   view.ImageBase64,
			Name:   "pix.png",
			Base64: base64.StdEncoding.EncodeToString(raw),
		},
	}

	got := ToTelegramResponse(resp, false, "pix")

	fileBytes, ok := got.Img.(tgbotapi.FileBytes)
	if !ok {
		t.Fatalf("expected FileBytes image, got %T", got.Img)
	}
	if fileBytes.Name != "pix.png" {
		t.Fatalf("expected name pix.png, got %s", fileBytes.Name)
	}
	if !bytes.Equal(fileBytes.Bytes, raw) {
		t.Fatalf("unexpected decoded bytes: %q", string(fileBytes.Bytes))
	}
}
