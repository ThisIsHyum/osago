# OSAGO

Go-клиент для [OpenScheduleApi](github.com/thisishyum/openscheduleapi), сгенерированный на основе OpenAPI-спецификации.

Соответствует версии API v0.5.0

## Использование
```golang
client := osago.NewClient("http://localhost:3505", 10*time.Second)
schedule, err := client.GetScheduleForToday(ctx, groupID, client.ScheduleFilter{})
if err != nil {
    // handle error
}
```
