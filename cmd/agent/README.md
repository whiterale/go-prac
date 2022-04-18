# cmd/agent

## Задание для трека «Сервис сбора метрик и алертинга»

Разработайте агент по сбору рантайм-метрик и их последующей отправке на сервер по протоколу HTTP. Разработку нужно вести с использованием шаблона.

### Агент должен собирать метрики двух типов

* gauge, тип float64

* counter, тип int64

В качестве источника метрик используйте пакет runtime.

### Нужно собирать следующие метрики

* "Alloc", тип: gauge

* "BuckHashSys", тип: gauge

* "Frees", тип: gauge

* "GCCPUFraction", тип: gauge

* "GCSys", тип: gauge

* "HeapAlloc", тип: gauge

* "HeapIdle", тип: gauge

* "HeapInuse", тип: gauge

* "HeapObjects", тип: gauge

* "HeapReleased", тип: gauge

* "HeapSys", тип: gauge

* "LastGC", тип: gauge

* "Lookups", тип: gauge

* "MCacheInuse", тип: gauge

* "MCacheSys", тип: gauge

* "MSpanInuse", тип: gauge

* "MSpanSys", тип: gauge

* "Mallocs", тип: gauge

* "NextGC", тип: gauge

* "NumForcedGC", тип: gauge

* "NumGC", тип: gauge

* "OtherSys", тип: gauge

* "PauseTotalNs", тип: gauge

* "StackInuse", тип: gauge

* "StackSys", тип: gauge

* "Sys", тип: gauge

* "TotalAlloc", тип: gauge

### К метрикам пакета runtime добавьте другие


* "PollCount", тип: counter — счётчик, увеличивающийся на 1 при каждом обновлении метрики из пакета runtime (на каждый pollInterval — см. ниже).

* "RandomValue", тип: gauge — обновляемое рандомное значение.

Репортинг

* pollInterval — 2 секунды

* reportInterval — 10 секунд

Метрики нужно отправлять по протоколу HTTP, методом POST:
по умолчанию на адрес: 127.0.0.1, порт: 8080;
в формате: http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>;
Content-Type: text/plain.

Агент должен штатно завершаться по сигналам: 
 
*syscall.SIGTERM, 

*syscall.SIGINT, 

*syscall.SIGQUIT.
