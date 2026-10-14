Notification Runtime

Um serviço de notificações escalável, em tempo real e de alto desempenho construído em Go.

Este projeto foi desenhado utilizando princípios de Clean Architecture, garantindo alta coesão e baixo acoplamento. Ele suporta envio de notificações em tempo real via WebSockets, processamento assíncrono utilizando RabbitMQ, e persistência de dados segura com PostgreSQL.

Tecnologias e Ferramentas

Linguagem: Go (Golang)

Banco de Dados: PostgreSQL

Mensageria / Filas: RabbitMQ

Tempo Real: WebSockets (hub e notifier)

Autenticação: JWT (JSON Web Tokens)

Observabilidade: Métricas e instrumentação nativa

Infraestrutura e Deploy: Docker, Docker Compose e Terraform

Arquitetura: Hexagonal

Arquitetura e Estrutura do Projeto

O projeto está dividido em camadas bem definidas para facilitar a manutenção e escalabilidade:

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


Funcionalidades Principais

Notificações em Tempo Real: Conexão WebSocket bidirecional contínua entre o cliente e o servidor.

Processamento em Background: Desacoplamento do envio de notificações utilizando RabbitMQ, garantindo que a API não sofra gargalos durante picos de uso.

Segurança: Rotas protegidas por autenticação JWT (auth_middleware).

Métricas Integradas: Endpoint dedicado para coleta de métricas de saúde e performance do sistema.

Infraestrutura Pronta: Arquivos do Terraform já configurados para provisionamento em nuvem, juntamente com o docker-compose para ambiente local.

Como Executar Localmente

Pré-requisitos

Docker e Docker Compose

Go 1.20+ (caso deseje rodar fora dos containers)

Passo a Passo

Clone o repositório:

git clone https://github.com/seu-usuario/Notification_RunTime.git
cd Notification_RunTime


Suba a infraestrutura (PostgreSQL, RabbitMQ):

docker-compose up -d


Nota: O banco de dados será inicializado automaticamente através do script docker/postgres/init.sql.

Inicie a API:

go run cmd/api/main.go


Testes

A aplicação possui uma suíte de testes robusta cobrindo integrações, casos de uso e workers. Para executar todos os testes:

go test -v ./tests/...


Para rodar testes com cobertura de código:

go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out


Infraestrutura como Código (Terraform)

O diretório /terraform contém os scripts necessários para provisionar a infraestrutura em nuvem (ex: AWS).

cd terraform
terraform init
terraform plan
terraform apply
