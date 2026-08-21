# Sistema de Emissão de Notas Fiscais - Korp / Viasoft

Projeto desenvolvido para o teste prático de desenvolvimento, a aplicação permite o cadastro de produtos, gerenciamento de estoque e emissão de notas fiscais. O sistema foi desenhado sob uma arquitetura de microsserviços orientada a eventos. Ao fechar (imprimir) uma nota, o serviço de faturamento publica uma mensagem no RabbitMQ, e o serviço de estoque consome essa fila para realizar a baixa dos itens de forma assíncrona e resiliente.


## 🚀 Tecnologias

* **Frontend:** Angular 17+ 
* **Backend:** Go (Golang) com framework Gin
* **Banco de Dados:** PostgreSQL 
* **Mensageria:** RabbitMQ para comunicação assíncrona entre os microsserviços
* **Inteligência Artificial:** Google Generative AI SDK (Gemini API) 
* **Reatividade (RxJS):** Uso de Observables retornados pelo HttpClient do Angular
* **UI/UX:** Não foi utilizada biblioteca de componentes visuais (como Bootstrap ou Material). A interface foi construída com HTML e CSS próprios e componentizados.

## 🏗️ Estrutura e System Design

```text
Korp_Teste_MayconVyctor/
├── frontend/              # Aplicação Angular (porta 4200)
├── inventory-service/     # Microsserviço de Estoque (Go - porta 8081)
├── billing-service/       # Microsserviço de Faturamento (Go - porta 8082)
└── docker-compose.yml     # Infraestrutura (PostgreSQL e RabbitMQ)



                ┌──────────────────────────┐
                │   Frontend (Angular)     │
                │   http://localhost:4200  │
                └───────────┬──────────────┘
                            │ HTTP / JSON
            ┌───────────────┴───────────────┐
            ▼                               ▼
┌───────────────────────┐      ┌───────────────────────┐
│  Serviço de Estoque   │◄─────┤ Serviço de Faturamento│
│  Go  ·  porta 8081    │ HTTP │  Go  ·  porta 8082    │
│                       │(Sync)│                       │
│ - CRUD de produtos    │      │ - CRUD de notas       │
│ - Controle de Saldos  │      │ - Integração IA Gemini│
│ - Consome Fila Rabbit ├──────◄ - Publica Mensageria  │
└───────────┬───────────┘Async └───────────┬───────────┘
            │          (RabbitMQ)          │
            │                              │
            │      Banco de Dados          │
            └──────────────┬───────────────┘
                           ▼
                ┌──────────────────────┐
                │    PostgreSQL 16     │
                │      porta 5432      │
                └──────────────────────┘

## Funcionalidades

**Produtos (Estoque):**
* Cadastro com código, descrição e saldo;
* Listagem dos produtos cadastrados;
* Validação rigorosa de campos obrigatórios.

**Notas Fiscais (Faturamento):**
* Numeração sequencial;
* Status inicial como **Aberta**;
* Inclusão de múltiplos produtos;
* Padrão de **Snapshot** (a descrição do produto é salva na nota no instante da compra, garantindo imutabilidade);
* Impressão restrita a notas abertas;
* Alteração de status para **Fechada** e disparo de evento via RabbitMQ para baixa de estoque;
* Geração de insights estratégicos da nota utilizando IA (Google Gemini).

## Endpoints Principais

**Estoque (Inventory) - Porta 8081**
* `GET /products` - Lista todos os produtos
* `GET /products/:code` - Busca produto específico
* `POST /products` - Cadastra um novo produto

**Faturamento (Billing) - Porta 8082**
* `GET /invoices` - Lista o histórico de notas
* `POST /invoices` - Emite uma nova nota fiscal (com validação síncrona no estoque)
* `PUT /invoices/:id/print` - Imprime/fecha a nota (dispara evento para o RabbitMQ)
* `GET /invoices/:id/insights` - Gera análise de vendas com IA

## Demonstração de Falha 

O sistema foi preparado para lidar com indisponibilidade, se o serviço de estoque estiver fora do ar no momento da emissão de uma nota fiscal, a operação é bloqueada com segurança e um feedback é exibido, ou o servico de faturamento retorna o mesmo 

**Para demonstrar o cenário:**
1. Navegue até a tela de Faturamento (Invoices);
2. Derrube (pare a execução) do terminal do `inventory-service`;
3. Tente emitir uma nova nota fiscal;
4. A tela exibirá um alerta formatado informando que o serviço de estoque está indisponível (Erro 503);
5. Reinicie o `inventory-service`;
6. Repita a emissão e a operação será concluída com sucesso (mensagem verde).

## Detalhamento Tecnico

**Frontend (Angular)**
* **Ciclos de vida:** Utilizado o `ngOnInit` nas páginas de componentes para disparar o carregamento de dados assim que a view é inicializada.
* **RxJS e Reatividade:** O consumo das APIs é feito via `HttpClient`, retornando `Observables`. Utilizamos o `.subscribe()` para processar as respostas assíncronas e gerenciar o fluxo de dados.
* **Renderização:** Para garantir a exibição imediata das mensagens de sucesso/erro dinâmicas após o retorno das requisições, utilizamos o `ChangeDetectorRef` (`detectChanges()`).
* **Bibliotecas:** Foram utilizadas `@angular/forms` (Reactive Forms para validações) e `@angular/common/http`.

**Backend (Go)**
* **Gerenciamento de dependências:** Feito de forma nativa através do Go Modules (`go.mod` e `go.sum`).
* **Framework:** Utilizado o Gin Framework para o roteamento HTTP, validação de JSON e middlewares.
* **Persistência:** Driver `github.com/lib/pq` para comunicação com o PostgreSQL, as operações de salvamento de notas contam com transações explícitas (`Begin`, `Commit`, `Rollback`) para assegurar consistência.
* **Tratamento de Erros:** Exceções e falhas são capturadas nos handlers e retornam HTTP Status Codes padronizados (ex: `400 Bad Request`, `500 Internal Server Error`, `503 Service Unavailable`).

## Como Executar

**1. Infraestrutura (Bancos de dados e Mensageria)**
Na raiz do projeto, suba os containers via Docker:
```bash
docker-compose up -d

**2. Servico de estoque (em outro terminal)
cd inventory-service
go run main.go

O serviço inicia em http://localhost:8081.

**3. Servico de faturamento (em outro terminal)
cd billing-service
go run main.go

O serviço inicia em http://localhost:8082.

**4. frontend (em outro terminal)
cd frontend
npm install
npm start

O serviço inicia em http://localhost:4200.
