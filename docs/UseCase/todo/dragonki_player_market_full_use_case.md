
# DragonKi — Sistema de Loja entre Jogadores (Player Market)
Documento de Caso de Uso Completo

---

# 🧭 Visão Geral

O sistema de **Player Market** permite que jogadores anunciem e comprem equipamentos entre si utilizando:

- Gold
- Ki Crystals

Este sistema cria uma **economia interna viva**, incentivando farm, trade e builds estratégicas.

---

# 🎯 Objetivos

- Criar economia entre jogadores
- Dar valor a itens raros
- Aumentar retenção
- Incentivar builds e sets

---

# 🛒 1. Acesso ao Sistema

Menu → Loja → Mercado de Jogadores

---

# 📦 2. Criação de Anúncio

## Fluxo

1. Player abre inventário
2. Seleciona item
3. Clica em "Vender"
4. Define:
   - Preço
   - Moeda (Gold / Ki Crystal)
5. Confirma anúncio

## Regras

- Item não pode estar equipado
- Item não pode estar em dungeon ativa
- Item fica **bloqueado**
- Limite de anúncios: 5 por jogador

---

# 🏪 3. Visualização do Mercado

Cada anúncio exibe:

- Nome do item
- Tipo (arma, armadura, acessório)
- Atributos
- Raridade
- Vendedor
- Preço
- Moeda

---

# 🔎 4. Filtros e Ordenação

## Filtros

- Por classe
- Por tipo de item
- Por raridade
- Por faixa de level

## Ordenação

- Menor preço
- Maior atributo
- Melhor custo-benefício

---

# 💰 5. Compra de Item

## Fluxo

1. Player seleciona anúncio
2. Sistema valida saldo
3. Remove moeda do comprador
4. Credita vendedor
5. Transfere item

---

# ✏️ 6. Edição de Anúncios

O vendedor pode:

- Alterar preço
- Alterar moeda
- Cancelar anúncio

---

# ❌ 7. Cancelamento

- Item retorna ao inventário
- Anúncio removido

---

# ⏳ 8. Expiração

- Tempo: 48 horas
- Após isso:
  - Item retorna ao jogador
  - Anúncio expirado

---

# ⚖️ 9. Taxa de Venda (Gold Sink)

Recomendado:

- 5% de taxa sobre vendas em Gold
- 0–2% em Ki Crystal

Objetivo:
- Controlar inflação
- Criar sink econômico

---

# 📊 10. Sistema de Valor de Mercado (Avançado)

- Histórico de preços
- Média de venda por item
- Sugestão automática de preço

---

# 🧩 11. Integração com Builds e Sets

## Ideia chave

Itens fazem parte de **sets**

Exemplo:

Set do Vento:
- +10% Speed (2 peças)
- +20% Dodge (4 peças)

Impacto:
- Jogadores buscarão peças específicas
- Aumenta uso do mercado

---

# 🧠 12. Segurança

Validações:

- Saldo suficiente
- Item disponível
- Anúncio ativo
- Player não pode comprar próprio item

---

# 🗄️ 13. Persistência (Banco de Dados)

## market_listings

id  
seller_id  
item_id  
price  
currency (gold | ki_crystal)  
status (active, sold, expired, canceled)  
expires_at  
created_at  
updated_at  

---

## market_transactions

id  
listing_id  
buyer_id  
seller_id  
price  
currency  
tax  
created_at  

---

# 🔄 14. Fluxo Completo

Inventário  
↓  
Criar anúncio  
↓  
Item listado  
↓  
Outro player compra  
↓  
Transferência  
↓  
Registro da transação  

---

# 🚀 15. Evoluções Futuras

- Sistema de leilão
- Destaque pago de anúncios
- Ranking de vendedores
- Sistema de reputação
- Notificação de venda

---

# 🏁 Conclusão

O Player Market:

- Cria economia viva
- Incentiva progressão
- Aumenta retenção
- Gera valor para itens

Sistema essencial para o crescimento do DragonKi.
