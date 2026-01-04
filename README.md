# Desafio Pós Go Expert - Concorrência em GO (aplicação de Leilão)

> Este projeto contém a solução para o desafio de concorrência em GO (aplicação de Leilão) da pós-graduação Go Expert da FullCycle.

(...) Adicionar uma nova funcionalidade ao projeto já existente para o leilão fechar automaticamente a partir de um tempo definido.

Clone o seguinte repositório: [clique para acessar o repositório](https://github.com/devfullcycle/labs-auction-goexpert).

Toda rotina de criação do leilão e lances já está desenvolvida, entretanto, o projeto clonado necessita de melhoria: adicionar a rotina de fechamento automático a partir de um tempo.

Para essa tarefa, você utilizará o go routines e deverá se concentrar no processo de criação de leilão (auction). A validação do leilão (auction) estar fechado ou aberto na rotina de novos lançes (bid) já está implementado.

Você deverá desenvolver:

- Uma função que irá calcular o tempo do leilão, baseado em parâmetros previamente definidos em variáveis de ambiente;
- Uma nova go routine que validará a existência de um leilão (auction) vencido (que o tempo já se esgotou) e que deverá realizar o update, fechando o leilão (auction);
- Um teste para validar se o fechamento está acontecendo de forma automatizada;

Dicas

- Concentre-se na no arquivo internal/infra/database/auction/create_auction.go, você deverá implementar a solução nesse arquivo;
Lembre-se que estamos trabalhando com concorrência, implemente uma solução que solucione isso;
- Verifique como o cálculo de intervalo para checar se o leilão (auction) ainda é válido está sendo realizado na rotina de criação de bid;

Entrega

- O código-fonte completo da implementação.
- Documentação explicando como rodar o projeto em ambiente dev.
- Utilize docker/docker-compose para podermos realizar os testes de sua aplicação.

# Executando a aplicação

### Pré-requisitos

- Ajuste as variáveis de ambiente no arquivo `cmd/auction/.env`

### Iniciando os serviços

Inicie os containers através do docker compose:

```sh
docker-compose up -d
```

A aplicação irá rodar na porta `:8080`.

### Endpoints da API

A aplicação expõe os seguintes endpoints:

- `GET /auction` Lista todos os leilões
- `GET /auction/:auctionId` Busca um leilão pelo ID
- `POST /auction` Cria um novo leilão
- `GET /auction/winner/:auctionId` Busca o lance vencedor de um leilão pelo ID
- `POST /bid` Cria um novo lance
- `GET /bid/:auctionId` Lista todos os lances de um leilão
- `GET /user/:userId` Busca um usuário pelo ID

Utilize a extensão REST Client, do vscode, e execute as requisições do arquivo `api/auction.http`.

### Rodando os testes

```sh
go test -v internal/infra/database/auction/create_auction_test.go 
```