# DragonKi --- Caso de Uso

## Módulo: Player VIP

------------------------------------------------------------------------

# Descrição

O módulo **Player VIP** permite que o jogador adquira um pacote VIP com
duração de **7 dias corridos**, pagando **10 Ki Crystals**.

Ao adquirir o pacote, o jogador recebe benefícios temporários no jogo
como:

-   bônus de XP
-   bônus de Gold
-   melhor valor de venda de equipamentos
-   desconto na lojinha
-   mais ataques diários contra o Big Boss

Todos os benefícios devem ser **configuráveis via banco de dados**,
evitando valores hardcoded no código.

------------------------------------------------------------------------

# Objetivo do módulo

Criar um sistema premium que:

-   aumente retenção de jogadores
-   incentive monetização
-   permita balanceamento sem deploy
-   mantenha configuração dinâmica via banco

------------------------------------------------------------------------

# Atores

-   Player
-   Sistema DragonKi

------------------------------------------------------------------------

# Compra do VIP

O jogador pode adquirir:

VIP semanal

Custo:

10 Ki Crystals

Duração:

7 dias completos

------------------------------------------------------------------------

# Regra de ativação

A contagem inicia **exatamente no momento da compra**.

Exemplo:

Compra realizada:

2026‑03‑09 14:35

Expiração:

2026‑03‑16 14:35

------------------------------------------------------------------------

# Validação de VIP ativo

Um jogador é considerado VIP quando:

vip_expires_at \> now()

------------------------------------------------------------------------

# Renovação do VIP

Caso o jogador compre VIP novamente **enquanto já possui VIP ativo**, o
tempo restante será **acumulado**.

Regra aplicada:

    vip_expires_at = vip_expires_at + 7 dias

------------------------------------------------------------------------

# Benefícios do VIP

## Bônus de XP

Jogadores VIP recebem XP adicional.

    xp_final = xp_base + (xp_base * vip_xp_bonus_percent)

------------------------------------------------------------------------

## Bônus de Gold

Jogadores VIP recebem gold adicional.

    gold_final = gold_base + (gold_base * vip_gold_bonus_percent)

------------------------------------------------------------------------

## Melhor valor de venda de equipamentos

Ao vender equipamentos na loja:

    sell_value_final = sell_value_base * vip_sell_bonus_multiplier

------------------------------------------------------------------------

## Desconto em compras da lojinha

Jogadores VIP recebem desconto percentual.

    shop_price_final = shop_price_base - (shop_price_base * vip_shop_discount_percent)

------------------------------------------------------------------------

## Mais ataques diários no Big Boss

Jogadores VIP possuem limite maior de ataques diários.

Exemplo:

Normal: 1 ataque/dia

VIP: valor definido na configuração.

------------------------------------------------------------------------

# Configuração dinâmica

Todos os valores devem ser armazenados na tabela **vip_config**.

------------------------------------------------------------------------

# Estrutura de domínio sugerida

## PlayerVIP

``` go
type PlayerVIP struct {
    ID        int64
    PlayerID  int64
    StartedAt time.Time
    ExpiresAt time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

------------------------------------------------------------------------

## VIPConfig

``` go
type VIPConfig struct {
    ID                      int64
    VIPPriceKiCrystals      int
    DurationDays            int
    XPBonusPercent          float64
    GoldBonusPercent        float64
    SellBonusPercent        float64
    ShopDiscountPercent     float64
    BigBossDailyAttackLimit int
    UpdatedAt               time.Time
}
```

------------------------------------------------------------------------

# Valores iniciais sugeridos

    VIPPriceKiCrystals = 10
    DurationDays = 7
    XPBonusPercent = 0.15
    GoldBonusPercent = 0.20
    SellBonusPercent = 0.25
    ShopDiscountPercent = 0.10
    BigBossDailyAttackLimit = 3

------------------------------------------------------------------------

# Fluxo principal

## 1. Player abre menu VIP

Sistema apresenta:

👑 DragonKi VIP

Custo: 10 Ki Crystals\
Duração: 7 dias

Benefícios:

-   bônus de XP
-   bônus de Gold
-   melhor venda de equipamentos
-   desconto na lojinha
-   mais ataques no Big Boss

------------------------------------------------------------------------

## 2. Player confirma compra

Sistema valida:

-   saldo de Ki Crystals
-   configuração VIP existente

------------------------------------------------------------------------

## 3. Sistema desconta os Ki Crystals

    player.ki_crystals -= vip_price_ki_crystals

------------------------------------------------------------------------

## 4. Sistema cria ou estende o VIP

Se o player **não possui VIP ativo**:

    started_at = now()
    expires_at = now() + 7 dias

Se o player **já possui VIP ativo**:

    expires_at = expires_at + 7 dias

------------------------------------------------------------------------

## 5. Sistema retorna confirmação

Mensagem ao jogador:

👑 VIP ativado com sucesso!

Seu VIP ficará ativo até:

DATA E HORA DA EXPIRAÇÃO

------------------------------------------------------------------------

# Fluxos alternativos

## Saldo insuficiente

❌ Você não possui Ki Crystals suficientes.

------------------------------------------------------------------------

## Configuração VIP indisponível

❌ O sistema VIP não está disponível no momento.

------------------------------------------------------------------------

# Pseudocódigo principal

``` go
func PurchaseVIP(playerID int64) {

    player := playerRepo.FindByID(playerID)
    config := vipConfigRepo.Get()

    if player.KiCrystals < config.VIPPriceKiCrystals {
        return ErrInsufficientKiCrystals
    }

    now := time.Now()

    player.KiCrystals -= config.VIPPriceKiCrystals

    vip := playerVIPRepo.FindByPlayerID(playerID)

    if vip == nil {
        vip.StartedAt = now
        vip.ExpiresAt = now.AddDate(0,0,config.DurationDays)
    } else {

        if vip.ExpiresAt.After(now) {
            vip.ExpiresAt = vip.ExpiresAt.AddDate(0,0,config.DurationDays)
        } else {
            vip.StartedAt = now
            vip.ExpiresAt = now.AddDate(0,0,config.DurationDays)
        }
    }
}
```

------------------------------------------------------------------------

# Nome técnico sugerido

player_vip

vip_config

------------------------------------------------------------------------

DragonKi Game Design Document VIP Membership System
