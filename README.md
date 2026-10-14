# Notification Runtime

Um serviço de notificações escalável, em tempo real e de alto desempenho construído em **Go**.

Este projeto foi desenhado utilizando princípios de **Clean Architecture**, garantindo alta coesão e baixo acoplamento. Ele suporta envio de notificações em tempo real via **WebSockets**, processamento assíncrono utilizando **RabbitMQ**, e persistência de dados segura com **PostgreSQL**.

## Tecnologias e Ferramentas

* **Linguagem:** [Go (Golang)](https://go.dev/)
* **Banco de Dados:** [PostgreSQL](https://www.postgresql.org/)
* **Mensageria / Filas:** [RabbitMQ](https://www.rabbitmq.com/)
* **Tempo Real:** WebSockets (`hub` e `notifier`)
* **Autenticação:** JWT (JSON Web Tokens)
* **Observabilidade:** Métricas e instrumentação nativa
* **Infraestrutura e Deploy:** Docker, Docker Compose e Terraform
* **Arquitetura:** Hexagonal

## Arquitetura e Estrutura do Projeto

O projeto está dividido em camadas bem definidas para facilitar a manutenção e escalabilidade:

```
├── cmd/api/                  # Ponto de entrada da aplicação (main.go)
├── docker/                   # Arquivos de configuração Docker (ex: init do Postgres)
├── internal/                 # Código principal da aplicação (não exportável)
│   ├── entity/               # Entidades de domínio (Regras de negócio core)
│   ├── usecase/              # Casos de uso (Ex: SendNotificationUseCase)
│   └── infra/                # Implementações técnicas e adaptadores
│       ├── auth/             # Middleware e geração de JWT
│       ├── delivery/http/    # Handlers HTTP, Middlewares, Rotas e WebSockets
│       ├── messaging/        # Implementação de filas (RabbitMQ, Local Queue)
│       ├── observability/    # Coleta de métricas da aplicação
│       ├── realtime/         # Gerenciamento de WebSockets (Hub, Notifier)
│       ├── repository/       # Acesso ao banco de dados Postgres
│       └── worker/           # Consumidores em background (Notification Worker)
├── terraform/                # Infraestrutura como código (IaaC - Provider, SG, Vars)
└── tests/                    # Testes automatizados (Integração, Casos de Uso, Workers)
```

## Funcionalidades Principais

1. **Notificações em Tempo Real:** Conexão WebSocket bidirecional contínua entre o cliente e o servidor.
2. **Processamento em Background:** Desacoplamento do envio de notificações utilizando RabbitMQ, garantindo que a API não sofra gargalos durante picos de uso.
3. **Segurança:** Rotas protegidas por autenticação JWT (`auth_middleware`).
4. **Métricas Integradas:** Endpoint dedicado para coleta de métricas de saúde e performance do sistema.
5. **Infraestrutura Pronta:** Arquivos do Terraform já configurados para provisionamento em nuvem, juntamente com o `docker-compose` para ambiente local.

## Como Executar Localmente

### Pré-requisitos

* [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/)
* [Go 1.20+](https://go.dev/dl/) (caso deseje rodar fora dos containers)

### Passo a Passo

1. **Clone o repositório:**
   ```bash
   git clone https://github.com/seu-usuario/Notification_RunTime.git
   cd Notification_RunTime
   ```

2. **Suba a infraestrutura (PostgreSQL, RabbitMQ):**
   ```bash
   docker-compose up -d
   ```
   *Nota: O banco de dados será inicializado automaticamente através do script `docker/postgres/init.sql`.*

3. **Inicie a API:**
   ```bash
   go run cmd/api/main.go
   ```

## Testes

A aplicação possui uma suíte de testes robusta cobrindo integrações, casos de uso e workers. Para executar todos os testes:

```bash
go test -v ./tests/...
```

Para rodar testes com cobertura de código:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Infraestrutura como Código (Terraform)

O diretório `/terraform` contém os scripts necessários para provisionar a infraestrutura em nuvem (ex: AWS).

```bash
cd terraform
terraform init
terraform plan
terraform apply
```
