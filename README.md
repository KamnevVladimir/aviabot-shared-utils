# Shared Utils Library

Утилиты и провайдеры для всех микросервисов.

## Назначение

- HTTP клиенты
- Валидация данных
- Генераторы ID
- Провайдеры времени
- Конфигурация

## Содержимое

- `http/` - HTTP клиенты
- `validation/` - валидация
- `providers/` - ID генераторы, время
- `config/` - конфигурация

## Использование

```go
import "github.com/KamnevVladimir/aviabot-shared-utils/providers"

// Генератор ID
generator := providers.NewIDGenerator()
id := generator.Generate()

// Валидация
validator := validation.NewValidator()
err := validator.Validate(data)
```

## Версионирование

- v1.0.1 - текущая версия
- Используется во всех микросервисах
