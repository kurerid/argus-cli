package git

import (
	"argus-agent/internal/llm"
	"argus-agent/internal/report"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

// mockLLMAnalysis is example data, not the result of analyzing the supplied diff.
const mockLLMAnalysis = `{
  "summary": "Пример анализа изменений оформления заказа.",
  "product": {
    "scenarios": [{"description": "Пользователь может указать комментарий к заказу.", "files": [{"repository": "example-frontend", "path": "src/checkout/OrderForm.vue"}]}],
    "effects": [{"description": "Комментарии могут сократить число уточнений по заказу.", "files": [{"repository": "example-frontend", "path": "src/checkout/OrderForm.vue"}], "metrics": ["Доля заказов, требующих уточнения"]}]
  },
  "development": {
    "changes": [{"description": "В запрос создания заказа добавлено поле комментария.", "files": [{"repository": "example-api", "path": "internal/orders/handler.go"}]}],
    "risks": [{"description": "Нужно проверить совместимость старых клиентов с новым полем.", "files": [{"repository": "example-api", "path": "internal/orders/handler.go"}]}],
    "unclear": [{"description": "Из diff неясно, ограничена ли длина комментария.", "files": [{"repository": "example-api", "path": "internal/orders/handler.go"}]}],
    "contradictions": [{"description": "В примере описание продукта не предусматривает комментарии к заказам.", "files": [{"repository": "example-frontend", "path": "src/checkout/OrderForm.vue"}]}]
  },
  "qa": [{"description": "Проверить создание заказа с комментарием и без него.", "files": [{"repository": "example-frontend", "path": "src/checkout/OrderForm.vue"}, {"repository": "example-api", "path": "internal/orders/handler.go"}]}]
}`

func compareDiff(
	ai llm.Provider,
	reportService reportGenerator,
	pathToRepo, left, right, aiModel string,
	temperature float64,
	maxOutputTokens int64,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	diff, err := getDiff(ctx, pathToRepo, left, right)
	if err != nil {
		return err
	}
	if diff == "" {
		return NO_DIFFERENCES_ERROR
	}

	req := llm.Request{
		Model:           aiModel,
		Message:         newDiffAnalysisPropmt(diff),
		Temperature:     temperature,
		MaxOutputTokens: maxOutputTokens,
	}

	resp, err := getDiffAnalysis(ctx, ai, req)
	if err != nil {
		return err
	}

	err = writeDiffReport(
		reportService,
		left,
		right,
		filepath.Base(pathToRepo),
		resp,
	)
	if err != nil {
		return err
	}

	return nil
}

func getDiff(
	ctx context.Context,
	pathToRepo, left, right string,
) (string, error) {
	cmd := exec.CommandContext(
		ctx,
		"git",
		"diff",
		fmt.Sprintf("%s...%s", left, right),
	)
	cmd.Dir = pathToRepo
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get branches diff: %w", err)
	}
	return string(output), nil
}

func newDiffAnalysisPropmt(diff string) string {
	return fmt.Sprintf("my message to make html-report %s", diff)
}

func getDiffAnalysis(
	ctx context.Context,
	ai llm.Provider,
	req llm.Request,
) (*llm.Response, error) {
	// resp, err := ai.Generate(ctx, req)
	// if err != nil {
	// 	err = fmt.Errorf("failed to get diff html-report: %w", err)
	// 	return nil, err
	// }
	// if resp.FinishReason != string(responses.ResponseStatusCompleted) {
	// 	return nil, fmt.Errorf(
	// 		"diff analysis did not complete: %s",
	// 		resp.FinishReason,
	// 	)
	// }
	return &llm.Response{Content: mockLLMAnalysis, FinishReason: "completed"}, nil
}

func writeDiffReport(
	reportService reportGenerator,
	left, right, repoName string,
	resp *llm.Response,
) error {
	analysis, err := report.ParseLLMAnalysis(resp.Content)
	if err != nil {
		return err
	}

	data := report.Report{
		Product:   repoName,
		CreatedAt: time.Now(),
		Branches:  report.Branches{Base: left, Target: right},
		Repo:      report.Repository{Name: repoName},
		Completeness: report.Completeness{Reasons: []string{
			"Показаны примерные данные: анализ diff моделью пока не выполняется.",
		}},
		Analysis: analysis,
	}
	return reportService.Generate(data)
}
