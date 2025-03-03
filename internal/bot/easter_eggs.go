package bot

import (
	"bufio"
	"math/rand"
	"os"
	"strings"
	"time"
)

func CheckEasterEgg(message string) string {
	keywords, err := loadLines("internal/bot/input.txt")
	if err != nil {
		return ""
	}

	responses, err := loadLines("internal/bot/output.txt")
	if err != nil || len(responses) == 0 {
		return ""
	}

	loweredMessage := strings.ToLower(strings.ReplaceAll(message, " ", ""))

	for _, keyword := range keywords {
		if strings.Contains(loweredMessage, strings.ToLower(keyword)) {
			return getRandomResponse(responses)
		}
	}

	return ""
}

func loadLines(filename string) ([]string, error) {
	var lines []string

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	return lines, scanner.Err()
}

func getRandomResponse(responses []string) string {
	rand.Seed(time.Now().UnixNano())
	return responses[rand.Intn(len(responses))]
}
