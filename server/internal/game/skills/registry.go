package skills

import (
    "darkscar/internal/game/combat"
    "fmt"
)

var Registry = map[string]SkillDef{
    // --- MARTYR ---
    "martyr_active": {
        ID: "martyr_active", Name: "Blood Sacrifice", Cooldown: 10, HPCostPct: 0.10,
        OnCast: func(caster combat.Combatant, allies, enemies []combat.Combatant) {
            caster.AddBuff(Effect{
                ID: "blood_sac", Name: "Blood Sacrifice", Duration: 5.0,
                StatMod: map[string]float64{"DMG_MULT": 0.40},
            })
            fmt.Printf("[%s] uses Blood Sacrifice (HP -> Dmg)!\n", caster.GetID())
        },
    },
    "martyr_passive": {
        ID: "martyr_passive", Name: "Vampirism", IsPassive: true,
        OnCast: func(caster combat.Combatant, allies, enemies []combat.Combatant) {
            caster.AddBuff(Effect{
                ID: "vampirism", Duration: 9999,
                OnHit: func(owner, target combat.Combatant, dmg float64) {
                    owner.Heal(dmg * 0.01) // 1% Life Steal
                },
            })
        },
    },

    // --- DREADNOUGHT ---
    "dread_active": {
        ID: "dread_active", Name: "Roar", Cooldown: 15, ManaCost: 20,
        OnCast: func(caster combat.Combatant, allies, enemies []combat.Combatant) {
            fmt.Printf("[%s] ROARS! Taunting enemies!\n", caster.GetID())
            for _, e := range enemies {
                e.AddThreat(caster.GetID(), 2000) 
            }
        },
    },
    "dread_passive": {
        ID: "dread_passive", Name: "Iron Skin", IsPassive: true,
        OnCast: func(caster combat.Combatant, allies, enemies []combat.Combatant) {
            caster.AddBuff(Effect{
                ID: "iron_skin", Duration: 9999,
                StatMod: map[string]float64{"PDEF_PCT": 0.10, "MDEF_PCT": 0.10}, // +5-10% Def
            })
        },
    },

    // --- DERVISH ---
    "dervish_active": {
        ID: "dervish_active", Name: "Lacerate", Cooldown: 8, ManaCost: 15,
        OnCast: func(caster combat.Combatant, allies, enemies []combat.Combatant) {
            if len(enemies) == 0 { return }
            target := enemies[0] 
            dmg := caster.GetAttack(combat.Physical) * 0.10
            
            target.AddBuff(Effect{
                ID: "bleed", Name: "Bleed", Duration: 5.0,
                OnTick: func(owner combat.Combatant, dt float64) {
                    owner.TakeDamage(dmg * dt, caster) 
                },
            })
            fmt.Printf("[%s] applies BLEED to [%s]\n", caster.GetID(), target.GetID())
        },
    },
    "dervish_passive": {
        ID: "dervish_passive", Name: "Adrenaline", IsPassive: true,
        OnCast: func(caster combat.Combatant, allies, enemies []combat.Combatant) {
            caster.AddBuff(Effect{
                ID: "adrenaline_tracker", Duration: 9999,
                OnGetHit: func(owner, source combat.Combatant, dmg float64) {
                    // Trigger buff when hit
                    owner.AddBuff(Effect{
                        ID: "adrenaline", Name: "Adrenaline", Duration: 3.0,
                        StatMod: map[string]float64{"SPD_PCT": 0.10, "EVA_PCT": 0.10},
                    })
                },
            })
        },
    },

    // --- SAINT ---
    "saint_active": {
        ID: "saint_active", Name: "Flash Heal", Cooldown: 6, ManaCost: 30,
        OnCast: func(caster combat.Combatant, allies, enemies []combat.Combatant) {
            // Find lowest HP ally
            var target combat.Combatant
            lowest := 1.1 // Percent
            
            for _, a := range allies {
                pct := a.GetCurrentHP() / a.GetMaxHP()
                if pct < lowest && !a.IsDead() {
                    lowest = pct
                    target = a
                }
            }
            
            if target != nil {
                amt := target.GetMaxHP() * 0.15
                target.Heal(amt)
                fmt.Printf("[%s] HEALS [%s] for %.0f HP!\n", caster.GetID(), target.GetID(), amt)
            }
        },
    },
    "saint_passive": {
        ID: "saint_passive", Name: "Aura of Light", IsPassive: true,
        OnCast: func(caster combat.Combatant, allies, enemies []combat.Combatant) {
            // Apply to ALL allies
            for _, a := range allies {
                a.AddBuff(Effect{
                    ID: "saint_aura", Name: "Aura", Duration: 9999,
                    StatMod: map[string]float64{"PDEF_PCT": 0.08, "MDEF_PCT": 0.08},
                })
            }
        },
    },
}
