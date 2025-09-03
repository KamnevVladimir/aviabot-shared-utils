# Aviabot Shared Utils

Утилиты и вспомогательные компоненты для микросервисов Aviabot.

## Содержимое

- `providers` - Провайдеры (IDGenerator implementations)
- `config` - Общие конфигурации
- `http` - HTTP клиенты
- `validation` - Валидаторы

## Использование

```go
import (
    "github.com/KamnevVladimir/aviabot-shared-utils/providers"
    "github.com/KamnevVladimir/aviabot-shared-utils/config"
    "github.com/KamnevVladimir/aviabot-shared-utils/http"
    "github.com/KamnevVladimir/aviabot-shared-utils/validation"
)
```

## Зависимости

- `github.com/KamnevVladimir/aviabot-shared-core` - Базовые интерфейсы

## Установка

```bash
go get github.com/KamnevVladimir/aviabot-shared-utils@v1.0.0
```
