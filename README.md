🛠️ SocialDev - Backend

SocialDev é uma rede social para desenvolvedores, criada para ser um espaço colaborativo onde profissionais de tecnologia possam compartilhar conhecimento, interagir em comunidades específicas e criar conexões com outros devs.

Este repositório contém apenas o backend, desenvolvido em Golang.

🚀 Funcionalidades Principais:

🔐 Sistema de Autenticação
Registro e login para usuários, moderadores e administradores.
Alteração de dados de conta e permissões por nível de acesso.

👥 Sistema de Grupos
Comunidades automáticas baseadas em tecnologias escolhidas (ex.: Go, JavaScript, Python).
Participação dinâmica em grupos com base nas preferências do usuário.

📝 Postagens Interativas
Sistema de likes, reações, comentários e compartilhamentos.

➕ Sistema de Seguidores
Siga e seja seguido por outros desenvolvedores.

💬 Mensageria em Tempo Real
Chat privado com criptografia ponta a ponta.
Mensagens em tempo real.

🔔 Notificações
Sistema de notificações em tempo real no aplicativo.
Envio de e-mails para notificações importantes.

📊 Dashboard Administrativo (Web)
Gerenciamento de usuários e conteúdo.
Sistema de filtragem avançada para moderação.

⚙️ Tecnologias Utilizadas:
Golang - Linguagem principal do backend
gRPC - Comunicação eficiente entre serviços
SQLC - Geração de código SQL segura e rápida
WebSockets - Comunicação em tempo real para o chat
JWT - Autenticação baseada em tokens
PostgreSQL - Banco de dados relacional
Redis - Cache e mensagens em tempo real

🔥 Como Rodar o Projeto Localmente 🔥

Clone o repositório:
```
git clone https://github.com/seu-usuario/devconnect-backend.git
cd devconnect-backend
```

Configure as variáveis de ambiente

Crie um arquivo .env com as seguintes variáveis:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=senha
DB_NAME=devconnect
JWT_SECRET=sua_chave_secreta

```

Execute o projeto:
```
go run cmd/main.go
```

Rodar as migrações (se necessário):
```
make migrate
```
