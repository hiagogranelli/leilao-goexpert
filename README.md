# Leilão GoExpert

## Subindo com Docker Compose

No diretório do projeto, execute:

```
docker compose up -d --build
```

- A API ficará disponível em http://localhost:8080
- Para acompanhar os logs:
```
docker compose logs -f app
```
- Para parar os serviços:
```
docker compose down
```

Se precisar recriar o volume de dados do Mongo (limpando o banco):
```
docker compose down -v
```

## Testes rápidos com curl

A seguir estão exemplos mínimos para exercitar a API. Os textos foram atualizados e o produto de exemplo foi alterado.

1) Criar um novo leilão

```
curl localhost:8080/auction \
  -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "product_name": "Cafeteira Espresso Automática",
    "category": "Eletrodomésticos",
    "description": "Cafeteira 110V com reservatório de 1,8L, vaporizador e timer programável.",
    "condition": 1
  }'
```

2) Consultar leilões abertos

```
curl -X GET "localhost:8080/auction?status=0"
```

3) Após o intervalo definido em `AUCTION_INTERVAL` (por exemplo, 10 segundos), consultar leilões finalizados

```
curl -X GET "localhost:8080/auction?status=1"
```
