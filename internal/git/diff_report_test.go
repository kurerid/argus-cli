package git

import (
	"argus-agent/internal/llm"
	"argus-agent/internal/report"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMockAnalysisRendersReport(t *testing.T) {
	service, err := report.NewService()
	if err != nil {
		t.Fatal(err)
	}
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(previousDir) })

	if err := writeDiffReport(service, "main", "release", "repo", &llm.Response{Content: mockLLMAnalysis}); err != nil {
		t.Fatal(err)
	}
	content, err := os.ReadFile(filepath.Join("diff-report.html"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"main", "release", "Показаны примерные данные", "Статистика недоступна",
		"Пользователь может указать комментарий", "Доля заказов, требующих уточнения",
		"В запрос создания заказа", "совместимость старых клиентов", "ограничена ли длина",
		"описание продукта не предусматривает", "Проверить создание заказа",
		"example-frontend", "example-api", "Проверки рекомендованы и не выполнялись Argus",
	} {
		if !strings.Contains(string(content), want) {
			t.Errorf("report does not contain %q", want)
		}
	}
}
