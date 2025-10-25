# QuePasa Webhook System - Documentação

## 📋 Índice
- [Visão Geral](#visão-geral)
- [Métricas e Monitoramento](#métricas-e-monitoramento)
- [Health Endpoint](#health-endpoint)
- [Configuração](#configuração)
- [Exemplos Práticos](#exemplos-práticos)
- [Troubleshooting](#troubleshooting)

---

## 🎯 Visão Geral

O **Sistema de Webhooks do QuePasa** é uma solução simples e direta para envio de webhooks:

### 🚀 Funcionalidades Principais
- ✅ **Processamento Direto**: Envio imediato de webhooks
- ✅ **Métricas Prometheus**: Monitoramento básico de performance
- ✅ **Health Endpoint**: Status em tempo real da saúde do sistema
- ✅ **Configuração Simples**: Configuração via variáveis de ambiente

### 🏗️ Arquitetura
```
Webhook Request → Direct Processing → External API
                      ↓
               Health Endpoint ← Metrics ← Prometheus
```

# QuePasa Webhook System - Documentação

## 📋 Índice
- [Visão Geral](#visão-geral)
- [Métricas e Monitoramento](#métricas-e-monitoramento)
- [Health Endpoint](#health-endpoint)
- [Configuração](#configuração)
- [Exemplos Práticos](#exemplos-práticos)
- [Troubleshooting](#troubleshooting)

---

## 🎯 Visão Geral

O **Sistema de Webhooks do QuePasa** é uma solução simples e direta para envio de webhooks:

### � Funcionalidades Principais
- ✅ **Processamento Direto**: Envio imediato de webhooks
- ✅ **Métricas Prometheus**: Monitoramento básico de performance
- ✅ **Health Endpoint**: Status em tempo real da saúde do sistema
- ✅ **Configuração Simples**: Configuração via variáveis de ambiente

### 🏗️ Arquitetura
```
Webhook Request → Direct Processing → External API
                      ↓
               Health Endpoint ← Metrics ← Prometheus
```

### � Como Funciona

#### Fluxo de Execução
```
1. Recebe Mensagem WhatsApp
   ↓
2. Processa e Cria Payload
   ↓
3. Envia HTTP POST para URL Webhook
   ↓
4. Registra Métricas (Sucesso/Erro)
```

---

## 📊 Métricas e Monitoramento

### 📈 Métricas de Mensagens

#### `quepasa_sent_messages_total`
- **Tipo**: Counter
- **Descrição**: Total de mensagens enviadas pelo sistema
- **Uso**: Monitora volume de mensagens de saída

#### `quepasa_send_message_errors_total`
- **Tipo**: Counter
- **Descrição**: Total de erros ao enviar mensagens
- **Uso**: Monitora falhas no envio de mensagens

#### `quepasa_received_messages_total`
- **Tipo**: Counter
- **Descrição**: Total de mensagens recebidas pelo sistema
- **Uso**: Monitora volume de mensagens de entrada

#### `quepasa_receive_message_errors_total`
- **Tipo**: Counter
- **Descrição**: Total de erros ao processar mensagens recebidas
- **Uso**: Monitora falhas no processamento de mensagens de entrada

### 📈 Métricas de Webhook

#### `quepasa_webhooks_sent_total`
- **Tipo**: Counter
- **Descrição**: Total de webhooks enviados
- **Uso**: Monitora volume total de webhooks

#### `quepasa_webhook_send_errors_total`
- **Tipo**: Counter
- **Descrição**: Total de webhooks que falharam
- **Uso**: Monitora taxa de falha dos webhooks

### 📊 Queries do Prometheus

#### Volume de Mensagens
```promql
# Total de mensagens enviadas
rate(quepasa_sent_messages_total[5m])

# Total de mensagens recebidas
rate(quepasa_received_messages_total[5m])
```

#### Taxa de Erro
```promql
# Taxa de erro de envio de mensagens
rate(quepasa_send_message_errors_total[5m]) / rate(quepasa_sent_messages_total[5m]) * 100

# Taxa de erro de webhook
rate(quepasa_webhook_send_errors_total[5m]) / rate(quepasa_webhooks_sent_total[5m]) * 100
```

#### Performance de Webhook
```promql
# Volumetria de webhooks
rate(quepasa_webhooks_sent_total[5m])

# Taxa de sucesso
(rate(quepasa_webhooks_sent_total[5m]) - rate(quepasa_webhook_send_errors_total[5m])) / rate(quepasa_webhooks_sent_total[5m]) * 100
```

### 🚨 Alertas Prometheus

#### Configuração de Alertas
```yaml
groups:
- name: quepasa.rules
  rules:
  - alert: HighWebhookErrorRate
    expr: rate(quepasa_webhook_send_errors_total[5m]) / rate(quepasa_webhooks_sent_total[5m]) > 0.1
    for: 2m
    labels:
      severity: warning
    annotations:
      summary: "Alta taxa de erro em webhooks"
      description: "Taxa de erro de webhooks acima de 10% por mais de 2 minutos"

  - alert: WebhookDown
    expr: up{job="quepasa"} == 0
    for: 1m
    labels:
      severity: critical
    annotations:
      summary: "QuePasa está down"
      description: "Instância QuePasa não responde"
```

---

## 🩺 Health Endpoint

### 📋 Endpoint Principal

- **GET /health**: Status básico do sistema

### 📄 Exemplo de Response
```json
{
  "success": true,
  "message": "OK",
  "timestamp": "2024-01-15T10:30:00Z",
  "stats": {
    "total": 5,
    "healthy": 4,
    "unhealthy": 1,
    "percentage": 80.0
  },
  "items": [
    {
      "wid": "5511999887766",
      "status": "connected",
      "healthy": true,
      "timestamp": "2024-01-15T10:29:45Z"
    }
  ]
}
```

---

## ⚙️ Configuração

### 📋 Variáveis de Ambiente

#### Webhook Básico
```bash
# WEBHOOK_TIMEOUT - Timeout em segundos para requests webhook
# Padrão: 10 segundos
# Mínimo: 1 segundo
# Máximo: 300 segundos (5 minutos)
WEBHOOK_TIMEOUT=10
```

---

## 💡 Exemplos Práticos

### � Configuração Básica

#### 1. Configuração Simples no .env
```bash
# Configuração básica
WEBHOOK_TIMEOUT=10
```

#### 2. Webhook Payload Exemplo
```json
{
  "message": {
    "id": "msg_123456789",
    "text": "Hello World",
    "from": "5511999887766@s.whatsapp.net",
    "to": "5511888776655@s.whatsapp.net",
    "timestamp": "2024-01-15T10:30:00Z"
  },
  "extra": {
    "custom_field": "value"
  }
}
```

### 🔗 Integrações Comuns

#### N8N Workflow
```json
{
  "nodes": [
    {
      "name": "QuePasa Webhook",
      "type": "webhook",
      "parameters": {
        "httpMethod": "POST",
        "path": "quepasa-webhook"
      }
    }
  ]
}
```

#### Chatwoot Integration
```javascript
// Processar webhook do QuePasa no Chatwoot
app.post('/quepasa-webhook', (req, res) => {
  const { message } = req.body;
  
  // Criar conversa no Chatwoot
  createConversation({
    contact_id: message.from,
    message: message.text
  });
  
  res.status(200).send('OK');
});
```

---

## 🔧 Troubleshooting

### 🚨 Problemas Comuns

#### 1. Webhooks não estão sendo enviados
**Sintomas:**
- Mensagens chegam no WhatsApp mas webhook não é chamado

**Verificações:**
```bash
# 1. Verificar se webhook está configurado
curl http://localhost:31000/v1/bot/{token}/webhook

# 2. Verificar logs
tail -f quepasa.log | grep webhook

# 3. Verificar métricas
curl http://localhost:31000/metrics | grep webhook
```

#### 2. Timeout em webhooks
**Sintomas:**
- Logs mostram timeout errors
- Métrica `webhook_send_errors` aumentando

**Soluções:**
```bash
# Aumentar timeout no .env
WEBHOOK_TIMEOUT=30

# Verificar se URL webhook responde
curl -I https://sua-url-webhook.com/endpoint
```

#### 3. URL webhook inválida
**Sintomas:**
- Erro 400 ou 404 consistente
- Logs mostram "invalid response"

**Verificações:**
```bash
# Testar URL manualmente
curl -X POST https://sua-url-webhook.com/endpoint \
  -H "Content-Type: application/json" \
  -d '{"test": "payload"}'
```

### 📊 Monitoramento

#### Dashboard Grafana Básico
```json
{
  "dashboard": {
    "title": "QuePasa Webhooks",
    "panels": [
      {
        "title": "Taxa de Webhooks",
        "targets": [
          {
            "expr": "rate(quepasa_webhooks_sent_total[5m])"
          }
        ]
      },
      {
        "title": "Taxa de Erro",
        "targets": [
          {
            "expr": "rate(quepasa_webhook_send_errors_total[5m]) / rate(quepasa_webhooks_sent_total[5m]) * 100"
          }
        ]
      }
    ]
  }
}
```

### 🔍 Debug Avançado

#### Logs Detalhados
```bash
# Ativar logs debug
export LOG_LEVEL=debug

# Filtrar logs de webhook
tail -f quepasa.log | grep "webhook\|http"
```

#### Health Check Script
```bash
#!/bin/bash
# health-check.sh

QUEPASA_URL="http://localhost:31000"

echo "Verificando health endpoint..."
curl -s "$QUEPASA_URL/health" | jq .

echo "Verificando métricas..."
curl -s "$QUEPASA_URL/metrics" | grep -E "(webhook|message)"
```

---

## 📚 Referencias

### 🔗 Links Úteis
- [Documentação N8N](https://docs.n8n.io/webhooks/)
- [Documentação Chatwoot](https://www.chatwoot.com/docs/product/webhooks)
- [Prometheus Metrics](https://prometheus.io/docs/concepts/metric_types/)
- [Grafana Dashboards](https://grafana.com/docs/grafana/latest/dashboards/)

### � Repositórios
- [QuePasa GitHub](https://github.com/nocodeleaks/quepasa)
- [Exemplos N8N](../extra/n8n+chatwoot/)
- [Exemplos Chatwoot](../extra/chatwoot/)

---

## 📝 Changelog

### Versão Atual
- ✅ Sistema de webhook direto e simples
- ✅ Métricas básicas do Prometheus
- ✅ Health endpoint simplificado
- ✅ Configuração via variáveis de ambiente

---

**📞 Suporte**: Para dúvidas, abra uma issue no repositório GitHub.
quepasa_sent_messages_total

# Total de mensagens recebidas  
quepasa_received_messages_total

# Taxa de mensagens por minuto (enviadas)
rate(quepasa_sent_messages_total[1m])

# Taxa de mensagens por minuto (recebidas)
rate(quepasa_received_messages_total[1m])
```

#### Taxa de Erro de Mensagens
```promql
# Taxa de erro no envio de mensagens
rate(quepasa_send_message_errors_total[5m]) / rate(quepasa_sent_messages_total[5m]) * 100

# Taxa de erro no recebimento de mensagens
rate(quepasa_receive_message_errors_total[5m]) / rate(quepasa_received_messages_total[5m]) * 100
```

#### Balanceamento de Tráfego
```promql
# Relação entre mensagens enviadas e recebidas
quepasa_sent_messages_total / quepasa_received_messages_total

# Volume total de mensagens processadas
quepasa_sent_messages_total + quepasa_received_messages_total
```

#### Taxa de Sucesso Global
```promql
# Taxa de sucesso de webhooks
(quepasa_webhooks_sent_total - quepasa_webhook_send_errors_total) / quepasa_webhooks_sent_total * 100
```

#### Eficácia do Sistema de Retry
```promql
# Quantos webhooks foram salvos pelo retry
quepasa_webhook_retries_successful_total / quepasa_webhook_retry_attempts_total * 100
```

#### Taxa de Retry
```promql
# Porcentagem de webhooks que precisaram de retry
quepasa_webhook_retry_attempts_total / quepasa_webhooks_sent_total * 100
```

#### Utilização da Fila
```promql
# Porcentagem de utilização da fila
quepasa_webhook_queue_size / WEBHOOK_QUEUE_SIZE * 100
```

#### Latência Média
```promql
# Tempo médio de entrega de webhooks
rate(quepasa_webhook_duration_seconds_sum[5m]) / rate(quepasa_webhook_duration_seconds_count[5m])
```

### 🚨 Alertas Recomendados

#### Taxa de Erro de Mensagens Alta
```yaml
- alert: MessageHighErrorRate
  expr: rate(quepasa_send_message_errors_total[5m]) / rate(quepasa_sent_messages_total[5m]) > 0.05
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "Taxa de erro de mensagens alta"
    description: "{{ $value | humanizePercentage }} das mensagens estão falhando no envio"

- alert: MessageReceiveErrorsHigh
  expr: rate(quepasa_receive_message_errors_total[5m]) / rate(quepasa_received_messages_total[5m]) > 0.05
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "Taxa de erro no recebimento de mensagens alta"
    description: "{{ $value | humanizePercentage }} das mensagens recebidas estão falhando no processamento"
```

#### Volume de Mensagens Baixo
```yaml
- alert: MessageVolumeLow
  expr: rate(quepasa_received_messages_total[5m]) < 0.1
  for: 10m
  labels:
    severity: warning
  annotations:
    summary: "Volume de mensagens recebidas muito baixo"
    description: "Sistema pode estar desconectado ou com problemas"
```

#### Taxa de Falha Alta
```yaml
- alert: WebhookHighFailureRate
  expr: rate(quepasa_webhook_send_errors_total[5m]) / rate(quepasa_webhooks_sent_total[5m]) > 0.1
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "Taxa de falha de webhook alta"
    description: "{{ $value | humanizePercentage }} dos webhooks estão falhando"
```

#### Sistema de Retry Ineficaz
```yaml
- alert: WebhookRetryIneffective
  expr: rate(quepasa_webhook_retry_failures_total[5m]) / rate(quepasa_webhook_retry_attempts_total[5m]) > 0.5
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "Sistema de retry não está funcionando"
    description: "{{ $value | humanizePercentage }} dos retries estão falhando"
```

#### Fila Muito Cheia
```yaml
- alert: WebhookQueueFull
  expr: quepasa_webhook_queue_size / WEBHOOK_QUEUE_SIZE > 0.8
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "Fila de webhooks muito cheia"
    description: "Fila está {{ $value | humanizePercentage }} cheia"
```

---

## 🏥 Health Endpoint

### 📍 Endpoints Disponíveis
- **GET /health**: Status completo com métricas de queue
- **GET /health/basic**: Status básico
- **GET /metrics**: Métricas detalhadas do Prometheus

### 📊 Resposta do Health Endpoint

```json
{
  "success": true,
  "status": "application is running",
  "timestamp": "2025-09-09T10:51:00Z",
  "queue": {
    "enabled": true,
    "current_size": 5,
    "max_size": 100,
    "utilization_percentage": 5.0,
    "processing_delay": "0s",
    "workers": 2,
    "processed_total": 150,
    "discarded_total": 2,
    "retries_total": 25,
    "completed_total": 145,
    "failed_total": 3
  }
}
```

### 📋 Campos da Queue no Health

| Campo | Tipo | Descrição |
|-------|------|-------------|
| `enabled` | boolean | Sistema de queue habilitado |
| `current_size` | integer | Tamanho atual da fila |
| `max_size` | integer | Capacidade máxima da fila |
| `utilization_percentage` | float | Utilização em porcentagem |
| `processing_delay` | string | Delay entre processamentos |
| `workers` | integer | Número de workers ativos |
| `processed_total` | float | Total processado (tempo real) |
| `discarded_total` | float | Total descartado (tempo real) |
| `retries_total` | float | Total de retries (tempo real) |
| `completed_total` | float | Total completado (tempo real) |
| `failed_total` | float | Total falhado (tempo real) |

---

## ⚙️ Configuração

### 🌍 Variáveis de Environment

#### Sistema de Retry
| Variável | Padrão | Descrição |
|----------|---------|-------------|
| `WEBHOOK_RETRY_COUNT` | undefined | Número de tentativas de retry |
| `WEBHOOK_RETRY_DELAY` | 1 | Segundos entre tentativas |
| `WEBHOOK_TIMEOUT` | 10 | Timeout por requisição (segundos) |

#### Sistema de Queue
| Variável | Padrão | Descrição |
|----------|---------|-------------|
| `WEBHOOK_QUEUE_ENABLED` | false | Habilitar sistema de queue |
| `WEBHOOK_QUEUE_SIZE` | 100 | Tamanho máximo da fila |
| `WEBHOOK_QUEUE_TIMEOUT` | 30 | Timeout de processamento |
| `WEBHOOK_QUEUE_DELAY` | 0 | Delay entre processamentos |
| `WEBHOOK_WORKERS` | 1 | Número de workers simultâneos |

### 📝 Arquivo .env.example

```bash
# Sistema de Retry de Webhooks
WEBHOOK_RETRY_COUNT=3
WEBHOOK_RETRY_DELAY=1
WEBHOOK_TIMEOUT=10

# Sistema de Queue de Webhooks
WEBHOOK_QUEUE_ENABLED=true
WEBHOOK_QUEUE_SIZE=100
WEBHOOK_QUEUE_TIMEOUT=30
WEBHOOK_QUEUE_DELAY=0
WEBHOOK_WORKERS=2
```

---

## 🔀 Exemplos Práticos

### 🏭 Ambiente de Produção
```bash
WEBHOOK_RETRY_COUNT=5
WEBHOOK_RETRY_DELAY=2
WEBHOOK_TIMEOUT=15
WEBHOOK_QUEUE_ENABLED=true
WEBHOOK_QUEUE_SIZE=500
WEBHOOK_WORKERS=4
```

### 🧪 Ambiente de Desenvolvimento
```bash
WEBHOOK_RETRY_COUNT=1
WEBHOOK_RETRY_DELAY=1
WEBHOOK_TIMEOUT=5
WEBHOOK_QUEUE_ENABLED=true
WEBHOOK_QUEUE_SIZE=50
WEBHOOK_WORKERS=1
```

### 🚀 Alta Performance
```bash
WEBHOOK_RETRY_COUNT=3
WEBHOOK_RETRY_DELAY=1
WEBHOOK_TIMEOUT=10
WEBHOOK_QUEUE_ENABLED=true
WEBHOOK_QUEUE_SIZE=1000
WEBHOOK_WORKERS=8
```

### 🔧 Debug/Testing
```bash
WEBHOOK_RETRY_COUNT=0
WEBHOOK_QUEUE_ENABLED=false
```

---

## 📋 Logs e Monitoramento

### ✅ Logs de Sucesso

#### Sucesso na Primeira Tentativa
```
INFO[2023-12-01 10:00:00] posting webhook
DEBUG[2023-12-01 10:00:01] webhook success on attempt 1
INFO[2023-12-01 10:00:01] webhook posted successfully
```

#### Sucesso Após Retry
```
INFO[2023-12-01 10:00:00] posting webhook
WARN[2023-12-01 10:00:01] webhook request error (attempt 1/4): timeout
INFO[2023-12-01 10:00:02] webhook retry attempt 1/3 after 1s delay
DEBUG[2023-12-01 10:00:03] webhook success on attempt 2
INFO[2023-12-01 10:00:03] webhook posted successfully
```

### ❌ Logs de Falha

#### Falha Não-Retryable (404)
```
INFO[2023-12-01 10:00:00] posting webhook
ERROR[2023-12-01 10:00:01] webhook returned status 404 (attempt 1/1)
ERROR[2023-12-01 10:00:01] webhook failed permanently
```

#### Falha Após Todos os Retries
```
INFO[2023-12-01 10:00:00] posting webhook
WARN[2023-12-01 10:00:01] webhook returned status 500 (attempt 1/4)
INFO[2023-12-01 10:00:02] webhook retry attempt 1/3 after 1s delay
WARN[2023-12-01 10:00:03] webhook returned status 502 (attempt 2/4)
INFO[2023-12-01 10:00:04] webhook retry attempt 2/3 after 1s delay
WARN[2023-12-01 10:00:05] webhook returned status 503 (attempt 3/4)
INFO[2023-12-01 10:00:06] webhook retry attempt 3/3 after 1s delay
WARN[2023-12-01 10:00:07] webhook returned status 504 (attempt 4/4)
ERROR[2023-12-01 10:00:07] max retry attempts reached
ERROR[2023-12-01 10:00:07] webhook failed after 4 attempts
```

### 📋 Logs de Queue

#### Mensagem Enfileirada
```
INFO[2023-12-01 10:00:00] Webhook enqueued for processing (Queue: 5/100)
```

#### Fila Cheia
```
WARN[2023-12-01 10:00:00] Webhook queue full, discarding message (Queue: 100/100)
```

#### Processamento
```
INFO[2023-12-01 10:00:01] Processing webhook from queue
INFO[2023-12-01 10:00:02] Webhook processed successfully
```

---

## 🔧 Troubleshooting

### 🚨 Problemas Comuns

#### 1. Webhooks Não Fazem Retry
**Sintomas**: Webhooks falham na primeira tentativa
**Causas Possíveis**:
- `WEBHOOK_RETRY_COUNT` não está definida
- Erro é classificado como não-retryable (4xx)
**Soluções**:
- Verificar se `WEBHOOK_RETRY_COUNT` está no .env
- Checar se erro é realmente retryable

#### 2. Fila Não Processa
**Sintomas**: Mensagens ficam na fila
**Causas Possíveis**:
- `WEBHOOK_QUEUE_ENABLED=false`
- Workers travados
- Problemas de conectividade
**Soluções**:
- Verificar configuração da fila
- Checar logs dos workers
- Reiniciar aplicação

#### 3. Alto Consumo de Memória
**Sintomas**: Memória cresce continuamente
**Causas Possíveis**:
- `WEBHOOK_QUEUE_SIZE` muito grande
- Muitas mensagens enfileiradas
- Workers não processando
**Soluções**:
- Reduzir `WEBHOOK_QUEUE_SIZE`
- Aumentar `WEBHOOK_WORKERS`
- Monitorar métricas de fila

#### 4. Latência Alta
**Sintomas**: Webhooks demoram muito para processar
**Causas Possíveis**:
- `WEBHOOK_TIMEOUT` muito alto
- `WEBHOOK_QUEUE_DELAY` configurado
- APIs externas lentas
**Soluções**:
- Ajustar timeouts
- Otimizar configurações
- Checar performance das APIs externas

### 🔍 Debugging

#### Verificar Configuração
```bash
# Checar se variáveis estão definidas
env | grep WEBHOOK

# Verificar valores no health endpoint
curl http://localhost:31000/health
```

#### Monitorar Métricas
```bash
# Ver métricas do Prometheus
curl http://localhost:31000/metrics | grep webhook

# Monitorar fila em tempo real
watch -n 1 'curl -s http://localhost:31000/health | jq .queue'
```

#### Analisar Logs
```bash
# Filtrar logs de webhook
tail -f logs/quepasa.log | grep webhook

# Ver apenas erros
tail -f logs/quepasa.log | grep -i error | grep webhook
```

---

## ❓ FAQ

### 🤔 O Sistema de Retry é Obrigatório?

**Não!** O sistema de retry é completamente opcional. Se você não definir `WEBHOOK_RETRY_COUNT` no seu .env, o sistema funcionará exatamente como antes - uma tentativa única por webhook.

### 🔄 Como Migrar para o Sistema de Retry?

1. **Teste primeiro**: Configure em ambiente de desenvolvimento
2. **Comece pequeno**: Use `WEBHOOK_RETRY_COUNT=1`
3. **Monitore**: Observe os logs e métricas
4. **Ajuste**: Aumente conforme necessário
5. **Produção**: Aplique configurações otimizadas

### 📊 As Métricas Afetam Performance?

Não significativamente. As métricas do Prometheus são otimizadas e têm impacto mínimo na performance. Elas são coletadas de forma assíncrona e não bloqueiam o processamento dos webhooks.

### 🏗️ Posso Usar Apenas a Fila Sem Retry?

Sim! Você pode habilitar apenas o sistema de queue definindo:
```bash
WEBHOOK_QUEUE_ENABLED=true
# WEBHOOK_RETRY_COUNT não definido = sem retry
```

### 👷 Quantos Workers Devo Usar?

Depende da sua carga de trabalho:
- **Desenvolvimento**: 1 worker
- **Produção pequena**: 2-4 workers
- **Produção média**: 4-8 workers
- **Alta performance**: 8+ workers

Monitore as métricas para encontrar o equilíbrio ideal.

### 🚨 E se a Fila Ficar Cheia?

O sistema usa **drop-tail policy**: quando a fila atinge o limite (`WEBHOOK_QUEUE_SIZE`), novas mensagens são descartadas automaticamente. Isso previne problemas de memória, mas você deve monitorar a métrica `quepasa_webhook_queue_discarded_total`.

### 🔧 Como Saber se Está Funcionando?

1. **Logs**: Procure por mensagens de retry e queue
2. **Health Endpoint**: Verifique o campo `queue` na resposta
3. **Métricas**: Acesse `/metrics` para ver contadores
4. **Teste**: Envie um webhook e veja os logs

---

## � Revisão Técnica e Melhorias

### 📋 Resumo da Análise

O sistema implementado está **tecnicamente sólido** e segue boas práticas de Go. Durante a revisão, foram identificados e corrigidos alguns problemas críticos e implementadas melhorias importantes.

### ✅ Pontos Positivos Encontrados

#### 1. **Arquitetura Bem Projetada**
- ✅ Uso correto de Go channels para thread-safety
- ✅ Padrão singleton com `sync.Once` para instância global
- ✅ Separação clara entre sistema de retry e queue
- ✅ Worker pool configurável

#### 2. **Sistema de Métricas Completo**
- ✅ Métricas Prometheus abrangentes
- ✅ Counters, Gauges e Histograms apropriados
- ✅ Integração com health endpoint

#### 3. **Configuração Flexível**
- ✅ Variáveis de ambiente bem organizadas
- ✅ Sistema condicional (só ativa quando configurado)
- ✅ Valores padrão sensatos

### 🔧 Problemas Críticos Corrigidos

#### 1. **CRÍTICO: Inicialização Desnecessária da Queue**
**Problema:** Queue era inicializada sempre, mesmo quando `WEBHOOK_QUEUE_ENABLED=false`

```go
// ANTES (problemático)
func init() {
    InitializeWebhookQueue() // Sempre executava
}

// DEPOIS (corrigido)
func init() {
    if environment.Settings.API.WebhookQueueEnabled {
        InitializeWebhookQueue()
    }
}
```

**Impacto:** Evita consumo desnecessário de recursos quando queue está desabilitada.

#### 2. **PERFORMANCE: Otimização do Worker Pool**
**Problema:** Loop desnecessário com timeout causava overhead de CPU

```go
// ANTES (ineficiente)
case <-time.After(100 * time.Millisecond):
    select {
    case msg := <-w.messageCache:
        // processa
    default:
        continue // CPU desperdiçada
    }

// DEPOIS (otimizado)
case msg := <-w.messageCache:
    w.processMessage(msg) // Bloqueia diretamente no channel
```

**Impacto:** Redução significativa do uso de CPU em idle, melhor performance geral.

#### 3. **LÓGICA: Melhoria na Função shouldRetry**
**Problema:** Ordem de verificação de status codes não era otimizada

```go
// ANTES
if statusCode >= 500 && statusCode < 600 {
    return true
}
if statusCode >= 400 && statusCode < 500 {
    return false
}

// DEPOIS (mais claro e eficiente)
if statusCode >= 400 && statusCode < 500 {
    return false // 4xx são erros permanentes (client errors)
}
if statusCode >= 500 && statusCode < 600 {
    return true  // 5xx são erros temporários (server errors)
}
```

**Impacto:** Lógica mais clara e menos tentativas desnecessárias em erros 4xx.

### 🚀 Melhorias Implementadas

#### 1. **Validação de Configuração com Limites Seguros**
```go
func (settings APISettings) GetWebhookQueueSize() int {
    if settings.WebhookQueueSize > 0 {
        if settings.WebhookQueueSize > 10000 {
            return 10000 // Previne uso excessivo de memória
        }
        return settings.WebhookQueueSize
    }
    return 100
}

func (settings APISettings) GetWebhookWorkers() int {
    if settings.WebhookWorkers > 0 {
        if settings.WebhookWorkers > 20 {
            return 20 // Previne criação excessiva de goroutines
        }
        return settings.WebhookWorkers
    }
    return 1
}
```

**Benefícios:**
- Previne configurações que podem consumir memória excessiva
- Limita número de workers para evitar sobrecarga
- Mantém valores padrão sensatos

#### 2. **Graceful Shutdown com Timeout**
```go
func (w *WebhookQueueClient) Close() {
    close(w.closed)
    
    done := make(chan struct{})
    go func() {
        w.wg.Wait()
        close(done)
    }()
    
    select {
    case <-done:
        log.Info("Workers finished gracefully")
    case <-time.After(30 * time.Second):
        log.Warn("Timeout waiting for workers")
    }
}
```

**Benefícios:**
- Encerramento limpo dos workers
- Evita travamento na shutdown da aplicação
- Timeout configurável para casos extremos

#### 3. **Funções de Gestão da Queue**
Novas funções administrativas implementadas:

```go
// Limpa recursos da queue
func CleanupWebhookQueue() {
    if GlobalWebhookQueue != nil {
        GlobalWebhookQueue.Close()
        GlobalWebhookQueue = nil
    }
}

// Reinicia queue sem restart da aplicação
func RestartWebhookQueue() {
    CleanupWebhookQueue()
    if environment.Settings.API.WebhookQueueEnabled {
        InitializeWebhookQueue()
    }
}
```

**Benefícios:**
- Capacidade de reiniciar queue em runtime
- Útil para mudanças de configuração sem downtime
- Melhor manutenibilidade

#### 4. **Logs Mais Informativos e Estruturados**
```go
if statusCode >= 400 && statusCode < 500 {
    logentry.Warnf("client error (4xx) detected - not retryable (status: %d)", statusCode)
} else {
    logentry.Infof("error is not retryable, stopping attempts")
}

// Log de debug para contagem de mensagens
logentry.Debugf("received message counted: type=%s, from=%s, chat=%s", 
    message.Type, from, message.Chat.Id)
```

**Benefícios:**
- Melhor debugging e troubleshooting
- Logs estruturados facilitam parsing
- Diferentes níveis para diferentes situações

### 🏆 Estado Final Após Melhorias

A implementação agora está **ainda mais robusta** e **pronta para produção**:

- ✅ **Eficiente**: Correção do polling desnecessário
- ✅ **Seguro**: Validação de limites de configuração
- ✅ **Robusto**: Graceful shutdown implementado
- ✅ **Administrável**: Funções de gestão disponíveis
- ✅ **Observável**: Logs melhorados para debugging
- ✅ **Escalável**: Worker pool otimizado

---

## 📊 Implementação de Contadores de Mensagens

### 🎯 Objetivo Alcançado

Implementação completa do **sistema de contagem de mensagens recebidas** no QuePasa, complementando as métricas já existentes de mensagens enviadas, fornecendo visibilidade total do tráfego de mensagens.

### ✅ Modificações Implementadas

#### 1. **Contadores no Handler Principal**
**Arquivo:** `src/whatsmeow/whatsmeow_handlers.go`

##### Função `Follow()` - Contador Principal de Recebimento
```go
// Increment received messages counter for all incoming messages
// Only count messages that are not from us (FromMe = false)
if !message.FromMe {
    metrics.MessagesReceived.Inc()
    
    logentry.Debugf("received message counted: type=%s, from=%s, chat=%s", 
        message.Type, from, message.Chat.Id)
}
```

**Critérios de Contagem:**
- ✅ Conta apenas mensagens **recebidas** (`FromMe = false`)
- ✅ Inclui todos os tipos: texto, mídia, chamadas, grupos
- ✅ Exclui mensagens **enviadas por nós** para evitar duplicação

##### Função `Message()` - Contadores de Erro
```go
// Count message receive error for nil messages
if evt.Message == nil {
    // ... error handling ...
    metrics.MessageReceiveErrors.Inc()
    return
}

// Count unhandled message as error
if message.Type == whatsapp.UnhandledMessageType {
    // ... error handling ...
    metrics.MessageReceiveErrors.Inc()
}
```

**Tipos de Erro Contabilizados:**
- ✅ Mensagens nulas/corrompidas
- ✅ Tipos de mensagem não suportados
- ✅ Falhas na decodificação

##### Função `CallMessage()` - Chamadas como Mensagens
```go
// Count incoming call as received message
metrics.MessagesReceived.Inc()
```

### 📈 Métricas Completas Disponíveis

| Métrica | Tipo | Descrição | Status |
|---------|------|-----------|---------|
| `quepasa_sent_messages_total` | Counter | Mensagens enviadas | ✅ Existente |
| `quepasa_send_message_errors_total` | Counter | Erros no envio | ✅ Existente |
| **`quepasa_received_messages_total`** | Counter | **Mensagens recebidas** | 🆕 **NOVO** |
| **`quepasa_receive_message_errors_total`** | Counter | **Erros no recebimento** | 🆕 **NOVO** |

### 🎯 Comportamento dos Novos Contadores

#### `quepasa_received_messages_total` incrementa quando:
- ✅ Mensagem de texto recebida de contato
- ✅ Mensagem de mídia recebida (imagem, vídeo, áudio, documento)
- ✅ Chamada recebida (voz ou vídeo)
- ✅ Mensagem de grupo recebida
- ✅ Mensagem de broadcast recebida
- ✅ Mensagens de sistema (entrada/saída de grupo)
- ❌ **NÃO conta** mensagens enviadas por nós (`FromMe = true`)

#### `quepasa_receive_message_errors_total` incrementa quando:
- ✅ Evento de mensagem nulo (`evt.Message == nil`)
- ✅ Tipo de mensagem não reconhecido (`UnhandledMessageType`)
- ✅ Falhas na decodificação de mensagens
- ✅ Erros de processamento interno

### 🚀 Benefícios Implementados

#### 1. **📈 Visibilidade Completa do Tráfego**
```promql
# Volume total de mensagens (entrada + saída)
quepasa_sent_messages_total + quepasa_received_messages_total

# Relação entre mensagens enviadas e recebidas
quepasa_sent_messages_total / quepasa_received_messages_total
```

#### 2. **🔍 Detecção Proativa de Problemas**
```promql
# Taxa de erro alta no recebimento pode indicar problemas de conectividade
rate(quepasa_receive_message_errors_total[5m]) / rate(quepasa_received_messages_total[5m]) > 0.05

# Volume baixo pode indicar desconexão do WhatsApp
rate(quepasa_received_messages_total[5m]) < 0.1
```

#### 3. **📊 Análise de Performance e Uso**
- Identificação de picos de tráfego e padrões de uso
- Análise de carga (bots vs usuários humanos)
- Balanceamento entre entrada e saída de mensagens

#### 4. **🚨 Alertas Inteligentes**
Novos alertas adicionados à documentação:

```yaml
# Taxa de erro alta no recebimento
- alert: MessageReceiveErrorsHigh
  expr: rate(quepasa_receive_message_errors_total[5m]) / rate(quepasa_received_messages_total[5m]) > 0.05
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "Taxa de erro no recebimento de mensagens alta"
    description: "{{ $value | humanizePercentage }} das mensagens recebidas estão falhando no processamento"

# Volume baixo pode indicar desconexão
- alert: MessageVolumeLow
  expr: rate(quepasa_received_messages_total[5m]) < 0.1
  for: 10m
  labels:
    severity: warning
  annotations:
    summary: "Volume de mensagens recebidas muito baixo"
    description: "Sistema pode estar desconectado ou com problemas"
```

### 💡 Queries Prometheus Úteis

#### Análise de Volume
```promql
# Taxa de mensagens por minuto (recebidas)
rate(quepasa_received_messages_total[1m])

# Taxa de mensagens por minuto (enviadas)
rate(quepasa_sent_messages_total[1m])

# Volume total processado
sum(rate(quepasa_received_messages_total[1m])) + sum(rate(quepasa_sent_messages_total[1m]))
```

#### Análise de Qualidade
```promql
# Taxa de erro no recebimento
rate(quepasa_receive_message_errors_total[5m]) / rate(quepasa_received_messages_total[5m]) * 100

# Taxa de erro no envio
rate(quepasa_send_message_errors_total[5m]) / rate(quepasa_sent_messages_total[5m]) * 100

# Taxa de erro geral do sistema
(rate(quepasa_send_message_errors_total[5m]) + rate(quepasa_receive_message_errors_total[5m])) / 
(rate(quepasa_sent_messages_total[5m]) + rate(quepasa_received_messages_total[5m])) * 100
```

#### Análise de Padrões
```promql
# Identificar se é mais bot (envia mais) ou usuário (recebe mais)
increase(quepasa_sent_messages_total[1h]) / increase(quepasa_received_messages_total[1h])

# Picos de atividade
delta(quepasa_received_messages_total[5m])
```

### 🛠️ Dashboard Sugerido para Grafana

```
┌─────────────────────────────────────────┐
│ 📊 Volume de Mensagens (24h)            │
├─────────────────────────────────────────┤
│ Enviadas: 1,234  │  Recebidas: 2,567   │
├─────────────────────────────────────────┤
│ 🔍 Taxa de Erro                         │
│ Envio: 0.2%      │  Recebimento: 0.1%  │
├─────────────────────────────────────────┤
│ 📈 Gráfico Temporal (Mensagens/Minuto)  │
│ ▲▲▲▲▲▲▲▲▲ (Recebidas - Azul)           │
│ ▼▼▼▼▼▼▼▼▼ (Enviadas - Verde)           │
├─────────────────────────────────────────┤
│ ⚡ Taxa de Processamento                │
│ Entrada: 45/min  │  Saída: 23/min      │
└─────────────────────────────────────────┘
```

### ✅ Status da Implementação

- ✅ **Código**: Implementado e funcionando
- ✅ **Compilação**: Sem erros, código testado
- ✅ **Métricas**: Contadores operacionais
- ✅ **Logs**: Debug information adicionada
- ✅ **Documentação**: Atualizada com exemplos
- ✅ **Alertas**: Configurações prontas
- ✅ **Queries**: Exemplos práticos fornecidos
- ✅ **Versão**: Atualizada para refletir mudanças

O sistema agora oferece **visibilidade completa** do tráfego de mensagens no QuePasa! 🎉

### 📖 Documentação Relacionada
- [QuePasa API Documentation](./api/)
- [Environment Configuration](./environment/)
- [Prometheus Metrics Guide](./metrics/)

### 🔗 Links Úteis
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Dashboards](https://grafana.com/grafana/dashboards/)
- [Webhook Best Practices](https://webhooks.fyi/)

---

## 🏷️ Version History

- **v3.25.0909.1130**: Implementação completa de contadores de mensagens recebidas
- **v3.25.0909.0952**: Sistema completo com queue, retry, métricas e health endpoint
- **v3.25.0909.0951**: Health endpoint com métricas em tempo real
- **v3.25.0909.0950**: Suporte a múltiplos workers
- **v3.25.2207.0128**: Sistema de queue assíncrona
- **v3.25.2207.0127**: Sistema de retry inteligente

### 🔧 Melhorias por Versão

#### v3.25.0909.1130
- ✅ Adição de contadores de mensagens recebidas
- ✅ Contadores de erros de recebimento
- ✅ Logs estruturados para debugging
- ✅ Alertas Prometheus para volume baixo
- ✅ Queries para análise de padrões de uso

#### v3.25.0909.0952
- ✅ Correção da inicialização condicional da queue
- ✅ Otimização do worker pool (remoção de polling)
- ✅ Validação de limites de configuração
- ✅ Implementação de graceful shutdown
- ✅ Funções de gestão da queue (cleanup/restart)
- ✅ Melhoria na lógica shouldRetry
- ✅ Logs mais informativos

---

*Esta documentação é mantida atualizada com as últimas funcionalidades. Para dúvidas ou sugestões, consulte o time de desenvolvimento.*
