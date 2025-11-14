# Korp_Teste_VitorMozer

Sistema de Emissão de Notas Fiscais com Arquitetura de Microsserviços

**GitHub:** [VitorMozer9](https://github.com/VitorMozer9)

## 🏗️ Arquitetura

O projeto segue os princípios de **Clean Architecture** e **SOLID**, dividido em dois microsserviços:

- **Stock Service** (Porta 8081): Gerencia produtos e saldos de estoque
- **Billing Service** (Porta 8082): Gerencia notas fiscais

## 🚀 Tecnologias

### Backend
- Go 1.25.1
- Chi Router (HTTP)
- Arquitetura Limpa (Domain, UseCase, Repository, Transport)

### Frontend
- Angular 17+
- RxJS para gerenciamento de estado reativo
- Material UI / PrimeNG para componentes visuais

## 📦 Estrutura do Projeto

```
Korp_Teste_SeuNome/
├── services/
│   ├── stock/       # Microsserviço de Estoque
│   └── billing/     # Microsserviço de Faturamento
├── pkg/             # Pacotes compartilhados
├── frontend/        # Aplicação Angular
└── docker-compose.yml
```

## 🏃 Como Executar

### Pré-requisitos
- Go 1.21+
- Docker e Docker Compose
- Node.js 18+ e npm

### Executar com Docker
```bash
docker-compose up --build
```

### Executar localmente

**Backend:**
```bash
# Terminal 1 - Stock Service
cd services/stock
go run cmd/stock/main.go

# Terminal 2 - Billing Service
cd services/billing
go run cmd/billing/main.go
```

**Frontend:**
```bash
cd frontend
npm install
ng serve
```

## 📋 Funcionalidades

### ✅ Implementadas
- [ ] Cadastro de Produtos
- [ ] Listagem de Produtos
- [ ] Cadastro de Notas Fiscais
- [ ] Impressão de Notas (com atualização de status e saldo)
- [ ] Tratamento de falhas entre microsserviços
- [ ] Validação de estoque antes de fechar nota

## 🛠️ Decisões Técnicas

### Backend (Go)

**Ciclos de Vida / Padrões:**
- Repository Pattern para abstração de dados
- Dependency Injection via construtores
- Interface segregation (ISP)

**Gerenciamento de Dependências:**
- Go Modules (go.mod)
- Go Workspace (go.work) para multi-módulos

**Tratamento de Erros:**
- Erros customizados por domínio
- Middleware de recuperação de panic
- Logs estruturados

**Frameworks:**
- Chi Router (HTTP routing)
- net/http padrão do Go

### Frontend (Angular)

**Ciclos de Vida:**
- ngOnInit para inicialização
- ngOnDestroy para limpeza de subscriptions

**RxJS:**
- Observables para chamadas HTTP
- BehaviorSubject para estado compartilhado
- Operators: map, catchError, switchMap

**Bibliotecas:**
- Angular Material / PrimeNG (componentes UI)
- RxJS (programação reativa)

## 📝 APIs

### Stock Service (http://localhost:8081)

```
GET    /api/products          # Lista produtos
POST   /api/products          # Cria produto
GET    /api/products/:id      # Busca produto
PUT    /api/products/:id      # Atualiza produto
POST   /api/products/reserve  # Reserva estoque
DELETE /api/products/:id      # Deleta estoque 
```

### Billing Service (http://localhost:8082)

```
GET    /api/invoices              # Lista notas
POST   /api/invoices              # Cria nota
GET    /api/invoices/:id          # Busca nota
POST   /api/invoices/:id/print    # Imprime (fecha) nota
```

## 👨‍💻 Autor

Vitor Mozer - [GitHub](https://github.com/VitorMozer9)

