# DragonKi -- Caso de Uso

## Módulo: Indicação de Novos Players

------------------------------------------------------------------------

# Descrição

O módulo de **indicação de jogadores** permite que um player gere um
link exclusivo do DragonKi e compartilhe com amigos.

Quando um novo jogador entra no jogo por esse link, a referência da
indicação é armazenada. Se o jogador indicado atingir o **nível 7**, o
jogador que o indicou recebe uma recompensa.

------------------------------------------------------------------------

# Objetivo do módulo

Criar um sistema de crescimento orgânico que:

-   incentive convite de novos jogadores
-   permita rastrear a origem da indicação
-   premie apenas indicações válidas
-   evite recompensa duplicada

------------------------------------------------------------------------

# Atores

-   Player indicador
-   Player indicado
-   Sistema DragonKi

------------------------------------------------------------------------

# Conceito funcional

## 1. O player gera um link de convite

Exemplo:

https://t.me/seu_bot?start=ref_12345

Onde: - 12345 é o ID do player indicador - ref_12345 é usado para
rastrear a indicação

------------------------------------------------------------------------

## 2. O novo jogador entra usando o link

Quando o bot recebe:

/start ref_12345

O sistema: - cria o player normalmente - registra quem indicou - salva a
indicação pendente

------------------------------------------------------------------------

## 3. O indicado joga normalmente

O sistema acompanha sua progressão.

------------------------------------------------------------------------

## 4. Ao atingir level 7

Se a indicação ainda estiver válida:

-   a indicação é marcada como concluída
-   o indicador recebe a recompensa

------------------------------------------------------------------------

# Regras de Negócio

1.  Cada player possui um link único de indicação
2.  Apenas novos jogadores podem ser indicados
3.  Um jogador só pode ter um indicador
4.  A recompensa é entregue apenas quando o indicado atingir nível 7
5.  A recompensa é entregue apenas uma vez
6.  Autoindicação não é permitida

------------------------------------------------------------------------

# Fluxo principal

### Gerar link de indicação

O sistema exibe:

🎁 Convide amigos para o DragonKi!

Seu link: https://t.me/seu_bot?start=ref_12345

Quando seu amigo atingir o nível 7, você recebe uma recompensa!

------------------------------------------------------------------------

### Registrar indicação

Quando um jogador entra com:

/start ref_12345

O sistema:

-   valida o parâmetro
-   extrai o referrerPlayerID
-   cria o player
-   salva o vínculo de indicação

------------------------------------------------------------------------

### Processar recompensa

Quando o jogador sobe de nível:

1.  Sistema verifica se existe indicação pendente
2.  Verifica se nível \>= 7
3.  Concede recompensa
4.  Marca como REWARDED

------------------------------------------------------------------------

# Estrutura de domínio sugerida

``` go
type Referral struct {
    ID               int64
    ReferrerPlayerID int64
    ReferredPlayerID int64
    Status           string
    RequiredLevel    int
    CreatedAt        time.Time
    CompletedAt      *time.Time
    RewardedAt       *time.Time
}
```

------------------------------------------------------------------------

# Constraints

Unique recomendado:

referred_player_id

------------------------------------------------------------------------

# Use Cases

GenerateReferralLinkUseCase

RegisterReferralUseCase

ProcessReferralRewardOnLevelUpUseCase

------------------------------------------------------------------------

# Pseudocódigo

Gerar link:

``` go
func GenerateReferralLink(playerID int64, botUsername string) string {
    return fmt.Sprintf("https://t.me/%s?start=ref_%d", botUsername, playerID)
}
```

Processar recompensa:

``` go
if player.level >= 7 {
    rewardService.GiveReferralReward(referral.ReferrerPlayerID)
}
```

------------------------------------------------------------------------

# Recompensas sugeridas

-   Gold
-   Ki Crystals
-   Itens raros
-   Tickets especiais

------------------------------------------------------------------------

# Métricas futuras

-   quantidade de convites por player
-   taxa de conversão
-   ranking de indicadores

------------------------------------------------------------------------

DragonKi Game Design Document Referral System
