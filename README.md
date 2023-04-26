# cbr_exchange_rate

API предоставляет данные по переводу различных валют в рубли за последние 90 дней.

## Запуск сервиса
Выполните команду для клонирования репозитория:
```bash
git clone git@github.com:zd4r/cbr_exchange_rate.git
```
Перейдите в папку проекта:
```bash
cd cbr_exchange_rate
```
И выполните команду для запуска сервиса:
```bash
make compose-build-up
```
## Пример использования
Получение данных о котировках валют (минимальные, максимальные и средние значения) за последние 90 дней (для форматирования ответа необходима утилита [jq](https://github.com/stedolan/jq):
```bash
curl http://localhost:8080/v1/dynamic_quotes | jq .
```

Пример ответа:
```bash
{
  "quotes": [
    {
      "currency": "Австралийский доллар",
      "min_quote": {
        "value": 48.9735,
        "date": "26.01.2023"
      },
      "max_quote": {
        "value": 55.2845,
        "date": "15.04.2023"
      },
      "avg_quote": {
        "value": 51.6745
      }
    },
    {
      "currency": "Азербайджанский манат",
      "min_quote": {
        "value": 40.5631,
        "date": "26.01.2023"
      },
      "max_quote": {
        "value": 48.4699,
        "date": "08.04.2023"
      },
      "avg_quote": {
        "value": 44.7684
      }
    },
    
    ...
    
    {
      "currency": "Японская иена",
      "min_quote": {
        "value": 0.529,
        "date": "26.01.2023"
      },
      "max_quote": {
        "value": 0.6255,
        "date": "08.04.2023"
      },
      "avg_quote": {
        "value": 0.5725
      }
    }
  ]
}
```
`currency` - название валюты

`min_quote.value` - значение минимального курса валюты

`min_quote.date` - дата минимального значения курса валюты

`max_quote.value` - значение максимального курса валюты

`max_quote.date` - дата максимального значения курса валюты
    
`avg_quote` - среднее значение курса рубля
