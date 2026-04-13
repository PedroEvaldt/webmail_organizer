# webmail_organizer

Monitor de email que verifica mensagens não lidas via IMAP e envia notificações no Discord via webhook. Pensado para rodar periodicamente num Raspberry Pi.

## Como funciona

1. Conecta ao servidor IMAP via TLS e busca emails não lidos na INBOX
2. Compara com os UIDs já processados (`data/seen_uids.txt`) para evitar notificações duplicadas
3. Se houver emails novos, envia um embed formatado no Discord via webhook
4. Salva os novos UIDs localmente

## Pré-requisitos

- Go 1.21+
- Acesso a um servidor IMAP com TLS
- Um webhook do Discord configurado

## Configuração

Crie um arquivo `.env` na raiz do projeto:

```env
IMAP_USERNAME=seu_usuario
PASSWORD=sua_senha
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
```

**Como criar o webhook no Discord:**
1. Configurações do canal → Integrações → Webhooks → Criar Webhook
2. Copie a URL e cole no `.env`

Crie o diretório de dados:

```bash
mkdir -p data
```

## Instalação e execução

```bash
# Clonar e entrar no projeto
git clone https://github.com/seu-usuario/webmail_organizer
cd webmail_organizer

# Baixar dependências
go mod download

# Compilar
go build -o webmail_organizer .

# Executar
./webmail_organizer
```

## Deploy no Raspberry Pi

**Cross-compilar a partir do seu PC:**

```bash
# Raspberry Pi 4/5 (ARM64)
GOOS=linux GOARCH=arm64 go build -o webmail_organizer_arm64 .

# Raspberry Pi 3 ou Zero (ARM 32-bit)
GOOS=linux GOARCH=arm GOARM=7 go build -o webmail_organizer_arm .
```

**Copiar para o Raspberry Pi:**

```bash
scp webmail_organizer_arm64 pi@<IP_DO_PI>:/home/pi/webmail_organizer/
scp .env pi@<IP_DO_PI>:/home/pi/webmail_organizer/
```

**Agendar com cron (executar a cada 2 horas):**

```bash
crontab -e
```

```
0 */2 * * * /home/pi/webmail_organizer/webmail_organizer >> /home/pi/webmail_organizer/logs/run.log 2>&1
```

## Estrutura do projeto

```
webmail_organizer/
├── main.go                        # Orquestra o fluxo principal
├── internal/
│   ├── imapclient/
│   │   └── imapclient.go          # Conexão IMAP e fetch de emails
│   ├── model/
│   │   └── email.go               # Struct Email
│   ├── storage/
│   │   ├── seen_uids.go           # Persistência de UIDs já processados
│   │   └── txt_writer.go          # Salva emails em .txt (debug)
│   └── notifier/
│       └── discord.go             # Envio de notificação via Discord webhook
└── data/
    └── seen_uids.txt              # Gerado automaticamente na primeira execução
```

## Dependências

| Biblioteca | Propósito |
|---|---|
| `github.com/emersion/go-imap/v2` | Cliente IMAP |
| `github.com/joho/godotenv` | Carrega variáveis do `.env` |
| `net/http` + `encoding/json` | Discord webhook (stdlib) |
