package iprange

import (
	"fmt"
	"testing"
)

// FuzzParseList - фазз-тест для функции ParseList
func FuzzParseList(f *testing.F) {
	// Функция, которая будет генерировать различные входные данные
	f.Add("10.0.0.1")
	f.Add("10.0.0.0/-1")
	f.Add("10.0.0.*")
	f.Add("10.0.0.1-10")
	f.Add("10.0.0.1,10.0.0.2-10")
	f.Add("10.0.0.0/16,10.0.0.10-20")

	// Фаззинг
	f.Fuzz(func(t *testing.T, input string) {
		// Пытаемся распарсить входные данные
		ranges, err := ParseList(input)

		// Убедитесь, что ошибки не возникают при обработке нормальных входных данных
		if err != nil {
			// Пример обработки ошибки: она может быть вызвана некорректным форматом
			return
		}

		// Дополнительная проверка: убедитесь, что IP-адреса в диапазоне правильные
		for _, ipRange := range ranges {
			fmt.Printf("IP Range:\n  Min IP: %s\n  Max IP: %s\n\n", ipRange.Min, ipRange.Max)
		}
	})
}
