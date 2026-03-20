# DragonKi — Sistema de Desgaste e Reparo de Equipamentos
Documento de Caso de Uso Completo

---

# 🧭 1. Visão Geral

Equipamentos com nível mínimo acima de 15 passam a possuir durabilidade (desgaste), consumida a cada batalha.

Quando a durabilidade chega a zero, o equipamento **quebra de forma irreparável**.  
Por isso, o jogador deve procurar o **Ferreiro** e reparar o item **antes** que a durabilidade chegue a 0.

---

# 🎯 2. Objetivos

- Criar sink de gold
- Incentivar manutenção preventiva
- Adicionar estratégia ao uso de equipamentos
- Valorizar itens raros e avançados
- Dar importância ao módulo do Ferreiro

---

# ⚙️ 3. Equipamentos Afetados

A regra se aplica apenas a equipamentos com:

```text
required_level > 15
```

Ou seja:

- equipamentos iniciais não sofrem desgaste
- equipamentos intermediários e avançados sofrem desgaste

---

# 🔋 4. Durabilidade

Cada item elegível possui:

- `max_durability`
- `current_durability`

## Valor padrão
Valor inicial sugerido:

```text
30 batalhas
```

Esse valor deve ser configurável via base de dados.

---

# ⚔️ 5. Consumo de Durabilidade

A cada batalha:

- o item perde **1 ponto de durabilidade**

Esse consumo ocorre tanto em:

- vitória
- derrota

## Modos onde o desgaste se aplica

- PvE
- Dungeons
- Boss
- Câmara do Ki Ancestral

---

# 💥 6. Quebra do Equipamento

Quando:

```text
current_durability <= 0
```

o equipamento entra no estado:

```text
BROKEN
```

## Regra fundamental
Ao atingir 0, o item:

- quebra
- torna-se **irreparável** 
- deixa de conceder bônus
- não pode mais ser recuperado no Ferreiro
- não pode ser vendido nem trocado
- podendo apenas ser descartado (dever ter a opção de descarte) 
- devera ser armazenado numa tabela o histórico de descarte, contendo as info do ‘item’ e do ‘player’ que o descartou

Isso obriga o jogador a fazer manutenção **antes** do colapso total do ‘item’.

---

# 🚫 7. Efeito de Item Quebrado

Equipamento quebrado:

- não concede atributos
- não pode ser equipado
- se já estiver equipado, perde o efeito automaticamente
- não pode mais ser reparado

---

# ⚠️ 8. Avisos ao Jogador

Para evitar perda acidental, o sistema deve alertar quando a durabilidade estiver baixa.
Essa mensagem deve ser exibida apenas se o jogador estiver no modo de combate que consuma durabilidade.

## Limite de alerta
Exemplo sugerido:

```text
<= 5
```

## Mensagem
```text
⚠️ Seu equipamento está com baixa durabilidade.
Procure um ferreiro antes que ele quebre.
```

---

# 🔨 9. Ferreiro

A função do Ferreiro é permitir a manutenção preventiva dos equipamentos antes da quebra.

Acesso:

```text
Menu → Ferreiro
```

---

# 🔄 10. Reparo

## Fluxo

1. Jogador acessa o Ferreiro
2. Sistema lista os itens reparáveis
3. Jogador escolhe um ‘item’
4. Sistema mostra:
   - durabilidade atual
   - durabilidade máxima
   - custo do reparo
5. Jogador confirma
6. Sistema restaura a durabilidade

## Importante
Apenas equipamentos com durabilidade **acima de 0** podem ser reparados.

---

# 💰 11. Cálculo de Custo

Fórmula sugerida:

```text
repair_cost = (max_durability - current_durability) * repair_cost_per_point
```

## Exemplo

- max = 30
- atual = 18
- faltando = 12
- custo por ponto = 25 gold

Resultado:

```text
12 * 25 = 300 gold
```

---

# 🗄️ 12. Configuração via Banco

O sistema não deve ser hardcoded.

## Tabela sugerida: equipment_durability_config

Campos:

```text
id
is_enabled
min_required_level
default_max_durability
durability_loss_per_battle
low_durability_threshold
repair_cost_per_point
break_on_zero
updated_at
```

## Valores iniciais sugeridos

```text
is_enabled = true
min_required_level = 15
default_max_durability = 30
durability_loss_per_battle = 1
low_durability_threshold = 5
repair_cost_per_point = 25
break_on_zero = true
```

---

# 📦 13. Persistência por Item

A durabilidade deve ser salva por instância do item do jogador, e não no item base.

Campos sugeridos em `player_items`:

```text
max_durability
current_durability
is_broken
last_repaired_at
```

Isso permite que dois jogadores tenham o mesmo item com desgastes diferentes.

---

# 🔄 14. Fluxo Pós-Batalha

1. Finaliza combate
2. Sistema identifica equipamentos elegíveis
3. Aplica desgaste
4. Verifica se algum item chegou a 0
5. Marca item como quebrado
6. Remove bônus, se necessário
7. Exibe aviso ao jogador

---

# 🧠 15. Casos de Uso do Módulo

## ApplyDurabilityUseCase
Aplica desgaste após combate.

Responsável por:
- reduzir durabilidade
- quebrar item quando necessário
- atualizar status

## GetRepairableItemsUseCase
Lista equipamentos que ainda podem ser reparados.

## RepairEquipmentUseCase
Executa o reparo do item.

Responsável por:
- validar gold
- validar se o item ainda é reparável
- restaurar durabilidade

## GetDurabilityConfigUseCase
Busca configuração ativa do sistema.

---

# ⚠️ 16. Edge Cases

## Item abaixo do nível mínimo
Não sofre desgaste.

## Item sem desgaste
Não aparece no ferreiro.

## Jogador sem gold
Reparo é bloqueado.

## Item quebrado
Não pode ser reparado.

---

# 💬 17. Mensagens Sugeridas

## Baixa durabilidade
```text
⚠️ Sua Espada Vulcânica está com baixa durabilidade: 4/30.
Procure um ferreiro antes que ela quebre.
```

## Quebra
```text
💥 Sua Armadura Gélida quebrou e tornou-se irreparável.
Ela não concede mais bônus.
```

## Reparo
```text
🔨 Equipamento reparado com sucesso!

Espada Vulcânica
Durabilidade: 30/30
Custo: 300 gold
```

## Item quebrado no ferreiro
```text
❌ Este equipamento já quebrou e não pode mais ser reparado.
```

---

# 🏁 18. Conclusão

O sistema de desgaste e reparo adiciona:

- economia
- estratégia
- manutenção preventiva
- risco real no uso de equipamentos avançados

E reforça o papel do **Ferreiro** como parte importante do ciclo do jogador no DragonKi.
