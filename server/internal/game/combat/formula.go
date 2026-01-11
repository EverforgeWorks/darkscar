package combat

import (
    "math"
    "math/rand"
)

// HitResult explains WHAT happened, not just how much it hurt.
type HitResult struct {
    Damage  float64
    IsMiss  bool
    IsCrit  bool
    Chance  float64 // Debug info: what was the hit chance?
}

func CalculateHit(attacker, defender Combatant, damageType DamageType, basePower float64) HitResult {
    // 1. Get Stats
    acc := attacker.GetAccuracy()
    eva := defender.GetEvasion()
    
    // 2. Calculate Hit Chance (Your Formula)
    // Formula: (Acc * 1.22) / ((Acc * 1.22) + Eva)
    c := 1.22
    numerator := acc * c
    hitChance := numerator / (numerator + eva)
    
    // 3. Roll for Hit
    if rand.Float64() > hitChance {
        return HitResult{IsMiss: true, Chance: hitChance}
    }

    // --- PHASE 1: Base Damage (Existing Logic) ---
    atk := attacker.GetAttack(damageType)
    def := defender.GetDefense(damageType)

    mitigatedDmg := (atk / ((def * 1.25) + 10.0)) * (basePower / 2.0)
    damageCap := 2.5 * basePower
    dCalc := math.Min(damageCap, mitigatedDmg)

    // --- PHASE 2: Crit & Variance ---
    isCrit := false
    critMult := 1.0
    
    // Roll for Crit
    if rand.Float64() < attacker.GetCritChance() {
        isCrit = true
        critMult = attacker.GetCritDamage()
    }

    // Variance (0.9 to 1.1)
    variance := 0.9 + (rand.Float64() * 0.2)

    finalDamage := math.Round(dCalc * critMult * variance)

    return HitResult{
        Damage: finalDamage,
        IsMiss: false,
        IsCrit: isCrit,
        Chance: hitChance,
    }
}
