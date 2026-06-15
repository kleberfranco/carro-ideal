# Diagramas UML - Carro Ideal

Este documento resume a arquitetura do projeto `carro-ideal` a partir do código em `app/`, `config/`, `web/` e `migrations/`.

## Diagrama de Classes / Domínio

```mermaid
classDiagram
    direction LR

    class User {
        +int64 ID
        +string Name
        +string Email
        +string PasswordHash
        +string Role
        +bool Active
        +time CreatedAt
        +time UpdatedAt
    }

    class VehicleCategory {
        +int64 ID
        +string Name
        +string Description
        +bool Active
        +int VehicleCount
    }

    class Vehicle {
        +int64 ID
        +int64 CategoryID
        +string Brand
        +string Model
        +string Version
        +int Year
        +string Condition
        +string FuelType
        +string Transmission
        +float64 PriceMin
        +float64 PriceMax
        +int Seats
        +int TrunkCapacity
        +float64 ConsumptionCity
        +float64 ConsumptionHighway
        +map MatchProfile
        +bool Active
    }

    class Question {
        +int64 ID
        +string Text
        +string Type
        +float64 Weight
        +int DisplayOrder
        +bool Active
    }

    class AnswerOption {
        +int64 ID
        +int64 QuestionID
        +string Text
        +map ScoreProfile
        +int DisplayOrder
        +bool Active
    }

    class SubmittedAnswer {
        +int64 QuestionID
        +int64 AnswerOptionID
    }

    class Recommendation {
        +int64 ID
        +int64 UserID
        +string Summary
        +string AISummary
        +int ItemCount
        +time CreatedAt
    }

    class RecommendationItem {
        +int64 ID
        +int64 RecommendationID
        +int Rank
        +float64 Score
        +string Reason
        +string[] MatchedCriteria
    }

    class AdminStats {
        +int Users
        +int Vehicles
        +int Recommendations
        +int Questions
        +int ActiveUsersWeek
        +int NewUsersWeek
    }

    User "1" --> "0..*" Recommendation : recebe
    User "1" --> "0..*" SubmittedAnswer : envia
    VehicleCategory "1" --> "0..*" Vehicle : categoriza
    Question "1" --> "1..*" AnswerOption : possui
    Question "1" --> "0..*" SubmittedAnswer : respondida por
    AnswerOption "1" --> "0..*" SubmittedAnswer : selecionada em
    Recommendation "1" --> "1..*" RecommendationItem : contem
    RecommendationItem "1" --> "1" Vehicle : recomenda
```

## Diagrama de Componentes / Camadas

```mermaid
flowchart TB
    Browser[Browser / Cliente Web]
    AdminUI[Painel Admin]

    subgraph HTTP[Camada HTTP - app/internal]
        WebHandler[web.Handler]
        APIHandler[api.Handler]
        AdminHandler[admin.Handler]
        HealthHandler[health.Handler]
        Middleware[platform middleware<br/>CORS, CSRF, rate limit, logs, recovery]
    end

    subgraph Services[Camada de Servicos - app/service]
        UserService[UserService]
        AuthService[AuthService]
        QuestionnaireService[QuestionnaireService]
        VehicleService[VehicleService]
        RecommendationService[RecommendationService]
        AdminService[AdminService]
        AIService[AIService]
        CatalogCache[CatalogCache]
    end

    subgraph Repositories[Camada de Persistencia - app/repository]
        UserRepository[UserRepository]
        SessionRepository[SessionRepository]
        QuestionRepository[QuestionRepository]
        VehicleRepository[VehicleRepository]
        RecommendationRepository[RecommendationRepository]
        AdminRepository[AdminRepository]
    end

    subgraph Infra[Infraestrutura]
        Config[config.Config]
        DB[(PostgreSQL)]
        Migrations[golang-migrate<br/>migrations/*.sql]
        OpenAIClient[clients.OpenAIClient]
        OpenAI[(OpenAI Chat Completions)]
    end

    Browser --> Middleware
    AdminUI --> Middleware
    Middleware --> WebHandler
    Middleware --> APIHandler
    Middleware --> AdminHandler
    Middleware --> HealthHandler

    WebHandler --> UserService
    WebHandler --> AuthService
    APIHandler --> UserService
    APIHandler --> AuthService
    APIHandler --> QuestionnaireService
    APIHandler --> VehicleService
    APIHandler --> RecommendationService
    AdminHandler --> UserService
    AdminHandler --> AuthService
    AdminHandler --> AdminService
    HealthHandler --> DB

    UserService --> UserRepository
    AuthService --> SessionRepository
    QuestionnaireService --> QuestionRepository
    QuestionnaireService --> CatalogCache
    VehicleService --> VehicleRepository
    VehicleService --> CatalogCache
    RecommendationService --> QuestionnaireService
    RecommendationService --> VehicleService
    RecommendationService --> RecommendationRepository
    RecommendationService --> AIService
    AdminService --> AdminRepository
    AdminService --> CatalogCache
    AIService --> OpenAIClient

    UserRepository --> DB
    SessionRepository --> DB
    QuestionRepository --> DB
    VehicleRepository --> DB
    RecommendationRepository --> DB
    AdminRepository --> DB
    Migrations --> DB
    Config --> DB
    OpenAIClient --> OpenAI
```

## Sequencia: Gerar Recomendacao

```mermaid
sequenceDiagram
    autonumber
    actor Usuario
    participant API as api.Handler
    participant Auth as AuthService
    participant Rec as RecommendationService
    participant Quiz as QuestionnaireService
    participant Vehicles as VehicleService
    participant AI as AIService
    participant Repo as RecommendationRepository
    participant DB as PostgreSQL
    participant OpenAI as OpenAI API

    Usuario->>API: POST /api/recommendations/generate
    API->>Auth: Authenticate(session token)
    Auth->>DB: consulta sessao ativa
    DB-->>Auth: userID
    Auth-->>API: usuario autenticado

    API->>Rec: Generate(userID, answers)
    Rec->>Quiz: BuildProfile(answers)
    Quiz->>DB: carrega perguntas e opcoes ativas
    DB-->>Quiz: Question[]
    Quiz-->>Rec: perfil ponderado

    Rec->>Vehicles: GetActive(categoryID=0)
    Vehicles->>DB: carrega catalogo ativo
    DB-->>Vehicles: Vehicle[]
    Vehicles-->>Rec: veiculos disponiveis

    alt OPENAI_API_KEY configurada e chamada bem-sucedida
        Rec->>AI: Recommend(answers, questions, vehicles)
        AI->>OpenAI: ChatComplete(systemPrompt, userPrompt)
        OpenAI-->>AI: JSON com ranking
        AI-->>Rec: AIRecommendation
    else IA ausente ou erro
        Rec->>Rec: recommendWithScoring(perfil, veiculos)
    end

    Rec->>Repo: Create(recommendation, answers)
    Repo->>DB: salva recommendations, items e respostas
    DB-->>Repo: ids persistidos
    Repo-->>Rec: OK
    Rec-->>API: Recommendation
    API-->>Usuario: 201 Created + JSON
```

## Rotas Principais

```mermaid
flowchart LR
    Root[/ /] --> Home[web.HomeHandler]
    Login[/login, /web/login/] --> LoginHandler[web.LoginHandler]
    Register[/register, /web/register/] --> RegisterHandler[web.RegisterHandler]
    Recommend[/recommend, /web/recommend/] --> RecommendHandler[web.RecommendHandler]
    AdminPage[/admin/] --> AdminPageHandler[admin.Page]

    AuthAPI[/api/auth/*/] --> UserAuth[api.Register/Login/Logout/Me]
    QuestionsAPI[/api/questions/] --> Questions[api.Questions]
    VehiclesAPI[/api/vehicles/] --> Vehicles[api.Vehicles]
    RecommendationsAPI[/api/recommendations/*/] --> Recommendations[api.Generate/History/Details]
    AdminAPI[/api/admin/*/] --> AdminCRUD[admin dashboard e CRUD]

    QuestionsAPI -. requer sessao .-> AuthMW[RequireAuth]
    RecommendationsAPI -. requer sessao .-> AuthMW
    AdminAPI -. requer sessao e role admin .-> AdminMW[RequireAuth + RequireAdmin]
```
