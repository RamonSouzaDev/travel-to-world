<h1 align="center">Hello 👋, I'm Ramon Mendes - Software Developer</h1>
<h3 align="center">A back-end developer passionate about technology</h3>

- 🔭 I am currently working on Software Architecture and Engineering Projects 

- 🌱 I'm currently learning **Software Architecture and Engineering**

- 📫 How to reach me **dwmom@hotmail.com**


   I ended up getting excited and developing, even after the delivery date.
  
<h3 align="left"> Let's network :</h3>

<p align="left">
<a href="https://linkedin.com/in/https://www.linkedin.com/in/ramon-mendes-b44456164/" target="blank"><img align="center" src="https://raw.githubusercontent.com/rahuldkjain/github-profile-readme-generator/master/src/images/icons/Social/linked-in-alt.svg" alt="https://www.linkedin.com/in/ramon-mendes-b44456164/" height="30" width="40" /></a>
</p>

# 🚀 Travel Requests - Sistema de Gerenciamento de Viagens Corporativas

[![Go](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![Vue.js](https://img.shields.io/badge/Vue.js-3.0+-green.svg)](https://vuejs.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](https://docker.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://postgresql.org/)

## 📋 Sobre o Projeto

Sistema completo para gerenciamento de pedidos de viagem corporativa, desenvolvido como **case técnico** demonstrando expertise em **Go (Golang)** e **Vue.js**.

### ✨ Funcionalidades Implementadas

- ✅ **Autenticação JWT** completa com registro e login
- ✅ **CRUD de Pedidos de Viagem** com validações robustas
- ✅ **Sistema de Status** (solicitado, aprovado, cancelado)
- ✅ **Filtros Avançados** por status, destino, período e datas de criação
- ✅ **Regras de Negócio** - usuário criador não pode alterar próprio pedido
- ✅ **Notificações** simuladas via logs estruturados
- ✅ **Interface Responsiva** com Tailwind CSS
- ✅ **Testes Automatizados** para backend
- ✅ **Dockerização Completa** com Docker Compose

## 🏗️ Arquitetura

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │    Backend      │    │   PostgreSQL    │
│   Vue.js 3      │◄──►│   Go + Gin      │◄──►│   Database      │
│   Port: 3000    │    │   Port: 8080    │    │   Port: 5432    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## 🚀 Instalação e Execução

### Pré-requisitos
- **Docker** 20.10+
- **Docker Compose** 2.0+

### Execução Rápida

```bash
# 1. Clone o repositório
git clone <seu-repositorio>
cd travel-requests

# 2. Inicie o sistema
./start.sh

# 3. Acesse as aplicações
# Frontend: http://localhost:3000
# Backend:  http://localhost:8080
```

### Comandos Disponíveis

```bash
./start.sh   # Iniciar sistema completo
./stop.sh    # Parar todos os serviços
./logs.sh    # Ver logs em tempo real
./clean.sh   # Limpeza completa
```

## 📚 Documentação da API

### 🔐 Autenticação

#### Registrar Usuário
```http
POST /api/auth/register
Content-Type: application/json

{
  "name": "João Silva",
  "email": "joao@empresa.com",
  "password": "senha123"
}
```

#### Login
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "joao@empresa.com",
  "password": "senha123"
}
```

### ✈️ Pedidos de Viagem

#### Criar Pedido
```http
POST /api/travel-requests
Authorization: Bearer {token}
Content-Type: application/json

{
  "requester_name": "João Silva",
  "destination": "São Paulo",
  "departure_date": "2025-08-15",
  "return_date": "2025-08-20"
}
```

#### Listar Pedidos (com filtros avançados)
```http
GET /api/travel-requests?status=aprovado&destination=São Paulo&start_date=2025-08-01&end_date=2025-08-31&created_after=2025-07-01&created_before=2025-07-31
Authorization: Bearer {token}
```

**Filtros Disponíveis:**
- `status`: solicitado, aprovado, cancelado
- `destination`: busca parcial por destino
- `start_date`: pedidos com ida após esta data
- `end_date`: pedidos com volta antes desta data
- `created_after`: pedidos criados após esta data
- `created_before`: pedidos criados antes desta data

#### Consultar Pedido por ID
```http
GET /api/travel-requests/1
Authorization: Bearer {token}
```

#### Atualizar Status (apenas se não foi o criador)
```http
PUT /api/travel-requests/1/status
Authorization: Bearer {token}
Content-Type: application/json

{
  "status": "aprovado"
}
```

#### Cancelar Pedido
```http
DELETE /api/travel-requests/1
Authorization: Bearer {token}
```

## 🧪 Testes Automatizados

### Executar Testes

```bash
# Parar containers se estiverem rodando
./stop.sh

# Executar testes do backend
cd backend
go test -v ./...

# Ou via Docker
docker-compose exec backend go test -v ./...
```

### Cobertura dos Testes

- ✅ **Health Check** - Verificação de saúde da API
- ✅ **Registro de Usuário** - Validação de criação
- ✅ **Login** - Autenticação e geração de token
- ✅ **Criação de Pedidos** - CRUD completo
- ✅ **Listagem** - Filtros e paginação
- ✅ **Validações** - Casos de erro e edge cases
- ✅ **Autenticação** - Proteção de rotas

## 🔧 Configuração do Ambiente

### Variáveis de Ambiente

Arquivo `.env` (criado automaticamente):
```env
# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=travel_requests
DB_USER=postgres
DB_PASSWORD=postgres

# JWT
JWT_SECRET=travel-requests-jwt-secret-2024-secure-key

# Server
PORT=8080
ENV=development
```

### Estrutura do Projeto

```
travel-requests/
├── backend/                    # API Go + Gin Framework
│   ├── main.go                # Entry point com toda lógica
│   ├── main_test.go           # Testes automatizados
│   ├── go.mod                 # Dependências Go
│   └── Dockerfile             # Container do backend
├── frontend/                   # Aplicação Vue.js
│   ├── src/
│   │   ├── Login.vue          # Tela de login/registro
│   │   ├── Dashboard.vue      # Dashboard principal
│   │   ├── FilterComponent.vue # Componente de filtros
│   │   └── main.js            # Entry point Vue
│   ├── Dockerfile             # Container do frontend
│   ├── package.json           # Dependências npm
│   └── vite.config.js         # Configuração Vite
├── docker-compose.yml         # Orquestração completa
├── start.sh                   # Script de inicialização
├── stop.sh                    # Script para parar
├── logs.sh                    # Script para logs
├── clean.sh                   # Script de limpeza
└── README.md                  # Esta documentação
```

## 🎯 Regras de Negócio Implementadas

### 1. **Autenticação e Autorização**
- ✅ JWT obrigatório para rotas protegidas
- ✅ Usuários só veem seus próprios pedidos
- ✅ Tokens com expiração de 24 horas

### 2. **Criação de Pedidos**
- ✅ Validação de datas (volta > ida)
- ✅ Campos obrigatórios validados
- ✅ Pedidos vinculados ao usuário criador

### 3. **Alteração de Status**
- ✅ **REGRA PRINCIPAL**: Usuário que criou o pedido **NÃO pode** alterar o status
- ✅ Apenas outros usuários podem aprovar/cancelar
- ✅ Status válidos: solicitado, aprovado, cancelado

### 4. **Cancelamento**
- ✅ Permite cancelar pedidos aprovados
- ✅ Pedidos já cancelados não podem ser alterados

### 5. **Notificações**
- ✅ Logs estruturados para todas as mudanças
- ✅ Notificações simuladas no console
- ✅ Formato: `[NOTIFICATION] Ação realizada`

## 🌐 Funcionalidades do Frontend

### 🔐 **Autenticação**
- **Login/Registro** em tela única com toggle
- **Gerenciamento de token** automático
- **Redirecionamento** baseado em autenticação

### 📊 **Dashboard**
- **Listagem** de pedidos com paginação
- **Filtros avançados** por múltiplos critérios
- **Cards responsivos** com informações completas
- **Ações contextuais** baseadas em regras de negócio

### 🎨 **Interface**
- **Design moderno** com Tailwind CSS
- **Responsivo** para mobile e desktop
- **Estados de loading** e feedback visual
- **Validações** em tempo real

## 🚨 Validação dos Requisitos

### ✅ **Funcionalidades Obrigatórias**
- [x] Criar pedido de viagem (ID, nome, destino, datas, status)
- [x] Atualizar status (com regra: criador não pode alterar)
- [x] Consultar pedido por ID
- [x] Listar pedidos com filtros (status, período, destino)
- [x] Cancelar pedido aprovado
- [x] Notificações (logs estruturados)

### ✅ **Requisitos Técnicos - Backend**
- [x] Go versão estável (1.21)
- [x] Princípios de microsserviços
- [x] PostgreSQL como banco relacional
- [x] Migrations automáticas via GORM
- [x] Validação de dados robusta
- [x] Tratamento de erros adequado
- [x] Autenticação JWT completa
- [x] Isolamento por usuário
- [x] Testes automatizados implementados
- [x] Docker para execução

### ✅ **Requisitos Técnicos - Frontend**
- [x] Vue.js 3 com Composition API
- [x] Autenticação via JWT
- [x] Telas de login e cadastro
- [x] Criação de pedidos
- [x] Listagem com filtros
- [x] Visualização de detalhes
- [x] Atualização de status
- [x] Estilização com Tailwind CSS
- [x] Comunicação via Axios

## 🔍 Como Testar

### 1. **Teste Manual Completo**
```bash
# 1. Acesse http://localhost:3000
# 2. Registre dois usuários diferentes
# 3. Com usuário A: crie um pedido
# 4. Com usuário B: aprove o pedido do usuário A
# 5. Tente alterar com usuário A: deve dar erro
# 6. Teste filtros avançados
# 7. Teste cancelamento de pedido aprovado
```

### 2. **Teste de API**
```bash
# Health check
curl http://localhost:8080/health

# Registrar usuário
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"123456"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"123456"}'
```

## 📈 Melhorias Futuras

### Funcionalidades
- [ ] Dashboard administrativo
- [ ] Relatórios e estatísticas
- [ ] Aprovação multinível
- [ ] Integração com sistemas de viagem
- [ ] Notificações por email/SMS

### Técnicas
- [ ] Cache com Redis
- [ ] Rate limiting
- [ ] Observabilidade com Prometheus
- [ ] CI/CD pipeline
- [ ] Testes E2E com Cypress

## 🤝 Desenvolvimento

### Estrutura de Código
- **Backend**: Arquitetura limpa com separação de responsabilidades
- **Frontend**: Componentes reutilizáveis e estado reativo
- **Testes**: Cobertura de casos críticos e edge cases
- **Docker**: Multi-stage builds otimizados

### Práticas Implementadas
- ✅ **Código limpo** e bem documentado
- ✅ **Tratamento de erros** consistente
- ✅ **Validações** em frontend e backend
- ✅ **Logs estruturados** para debugging
- ✅ **Segurança** com JWT e sanitização
- ✅ **Performance** com queries otimizadas

## 📞 Suporte

Para dúvidas sobre implementação ou execução:

1. **Logs do sistema**: `./logs.sh`
2. **Status dos containers**: `docker-compose ps`
3. **Health check**: `curl http://localhost:8080/health`

---

## 🎉 Resumo da Implementação

Este projeto demonstra:

- **Expertise em Go** com APIs REST robustas
- **Proficiência em Vue.js** com interfaces modernas
- **Arquitetura de microsserviços** bem estruturada
- **Testes automatizados** com boa cobertura
- **DevOps** com Docker e containerização
- **Boas práticas** de desenvolvimento

**Todos os requisitos do case técnico foram implementados com qualidade de código profissional! 🚀**
