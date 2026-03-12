# DragonKi - Guia do Jogador (Acoes e Regras)

Este guia explica, em linguagem de jogador, o que cada botao faz e quais regras o jogo aplica por tras.

## Menu Principal

No menu principal voce encontra:
- `Aventura`
- `Personagem`
- `Inventario`
- `Economia`
- `Convites`
- `Eventos`
- `Rankings`
- `Entrar na Comunidade`

---

## 1) Aventura

### `Cacar`
O que faz:
- Encontra um inimigo da sua regiao atual.
- Abre o combate com opcoes de `Atacar` ou `Escapar`.

Regra do jogo:
- So funciona se voce tiver energia.
- Os monstros dependem da regiao e do peso de spawn.

### `Explorar`
O que faz:
- Mostra as regioes disponiveis para seu nivel.

Regra do jogo:
- Ao escolher uma regiao, seu personagem muda de local.
- Cada regiao tem nivel minimo e penalidades diferentes.

### `Treinamento`
O que faz:
- Abre a area de treino/transformacao.

Regra do jogo:
- Transformacao exige requisitos (nivel e energia).

---

## 2) Personagem

### `Status`
O que faz:
- Mostra seus atributos atuais (HP, Energia, ATK, DEF, nivel, gold etc.).

Regra do jogo:
- Valores podem incluir efeitos de equipamentos e estado atual.

### `Transformar`
O que faz:
- Tenta ativar sua transformacao.

Regra do jogo:
- Consome energia.
- Exige requisitos minimos.

### `Pousada`
O que faz:
- Abre opcoes de descanso.

Opcoes:
- `Dormir`: maior recuperacao (custo em gold).
- `Descansar`: recuperacao parcial (custo menor).

---

## 3) Inventario

### `Abrir Inventario`
O que faz:
- Abre suas categorias de itens.

### `Pocoes`
O que faz:
- Lista pocoes que voce possui.

Regra do jogo:
- Ao clicar em uma pocao, ela e consumida (`Usar item`).

### `Equipamentos`
O que faz:
- Lista equipamentos no inventario por tipo (arma, armadura, luvas etc.).

Regra do jogo:
- Ao clicar no item, ele e equipado.
- Itens equipados nao aparecem na lista de venda.

### `Equipados`
O que faz:
- Mostra os itens atualmente equipados.

Opcoes:
- Desequipar por slot.
- Desequipar tudo.

---

## 4) Economia

### `Lojinha`
Dentro da lojinha:
- `Pocoes`
- `Equipamentos`

#### Comprar itens
O que faz:
- Compra o item selecionado.

Regra do jogo:
- Verifica gold suficiente.
- Equipamentos seguem requisitos de nivel e classe.

#### Venda de equipamentos
O que faz:
- Lista equipamentos do player (nao equipados), paginados.

Regra do jogo:
- Valor de venda = `40%` do preco do item.

### `Cambio`
O que faz:
- Converte Ki Crystals em Gold.

### `Comprar Ki`
O que faz:
- Inicia compra de Ki Crystals.

Regra do jogo:
- So aparece quando a feature de pagamento esta habilitada.

---

## 5) Eventos

## Big Boss Semanal

### `Abrir Big Boss`
O que mostra:
- Nome do boss da semana.
- HP global atual.
- Tempo restante.
- Seu dano acumulado na semana.
- Status do ataque diario.

### `Atacar`
O que faz:
- Executa um combate especial contra o Big Boss.

Regras:
- Cada player pode atacar `1 vez por dia`.
- Custo por ataque: `5 de energia` e `200 de gold`.
- Ao final do duelo, o player termina com `HP = 1`.
- O dano causado reduz o HP global do boss.

### `Ranking`
O que faz:
- Mostra os maiores contribuidores de dano na semana.

### Derrota do Big Boss (HP global = 0)
O que acontece:
- Evento muda para derrotado.
- Premiacao Top 3 e processada.

Regra de premiacao:
- 1 lugar: gold + itens (conforme configuracao do servidor).
- 2 lugar: gold + itens.
- 3 lugar: gold.

Observacao:
- A premiacao e processada por fila (worker) para garantir estabilidade.

---

## Camara do Ki Ancestral

### `Abrir Camara`
O que mostra:
- custo de entrada (`1 Energia` + `150 Gold`)
- maior nivel de camara ja desbloqueado
- status da run ativa (se houver)

### `Iniciar Nivel`
O que faz:
- inicia uma run no nivel escolhido
- consome o custo de entrada uma unica vez

Regras:
- nivel 1 sempre liberado
- nivel N so libera ao concluir nivel N-1
- se ja existir run ativa, nao inicia outra

### `Atacar Batalha Atual`
O que faz:
- executa a batalha atual da run
- usa dificuldade escalada pelo nivel da camara e do player

Regras:
- cada run tem `5 batalhas`
- ao vencer, avanca para a proxima batalha
- ao perder, a run falha e encerra

### `Entre batalhas`
Apos cada vitoria, voce pode:
- `Continuar`
- `Usar Pocoes`
- `Transformar`
- `Sair da Camara`

### `Concluir nivel`
Quando vence as `5/5` batalhas:
- nivel da camara e concluido
- proximo nivel e desbloqueado

### `Sair da Camara`
Se sair antes de concluir 5/5:
- progresso da run atual e perdido
- nova tentativa volta para batalha 1

---

## 6) Convites

### `🎁 Convites`
O que faz:
- Mostra seu link de indicacao.
- Mostra quantas indicacoes estao pendentes e quantas ja foram recompensadas.

Como funciona:
- Compartilhe seu link com amigos.
- O amigo precisa entrar no bot usando o link (`/start ref_<seu_player_id>`).
- Quando o indicado atingir o nivel minimo configurado no servidor, voce recebe recompensa.

Regras:
- Nao vale autoindicacao.
- Cada jogador pode ser vinculado a apenas um convidador.
- A recompensa de cada indicacao e paga apenas uma vez.
- Valores de recompensa e nivel minimo podem variar por configuracao do servidor.

---

## 7) Rankings

### `Top Batalhas`
Mostra quem mais batalhou na semana.

### `Top Vitorias`
Mostra quem mais venceu na semana.

### `Top Derrotas`
Mostra quem mais perdeu na semana.

### `Top Escapes`
Mostra quem mais fugiu na semana.

---

## 8) Regras gerais importantes

- Energia limita acoes de combate.
- Gold e necessario para compras, descanso e ataque do Big Boss.
- Equipamentos alteram atributos do personagem.
- Algumas funcoes podem ser ligadas/desligadas pelo servidor (feature flags).

---

## 9) Dicas para evoluir mais rapido

- Mantenha energia sempre em ciclo (cacar + descanso).
- Invista em equipamentos com bonus para seu estilo.
- Participe do Big Boss todo dia para subir no ranking semanal.
- Use o inventario com frequencia para nao carregar item parado.
