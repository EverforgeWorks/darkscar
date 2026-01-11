package engine

import (
    "darkscar/internal/game/combat"
    "math/rand"
)

// SelectTarget decides who 'actor' attacks from the list of 'enemies'
func SelectTarget(actor combat.Combatant, enemies []combat.Combatant) combat.Combatant {
    mode := actor.GetTargetMode()
    
    // 1. Threat Logic (Enemies usually do this)
    if mode == combat.ModeHighestThreat {
        topThreatID := actor.GetHighestThreatID()
        
        for _, e := range enemies {
            if e.GetID() == topThreatID && !e.IsDead() {
                return e
            }
        }
        // Fallback: If no threat yet, pick Random
        mode = combat.ModeRandom
    }

    var candidate combat.Combatant
    bestVal := 0.0

    switch mode {
    case combat.ModeLowestHP:
        bestVal = 9999999.0
        for _, e := range enemies {
            if !e.IsDead() && e.GetCurrentHP() < bestVal {
                bestVal = e.GetCurrentHP()
                candidate = e
            }
        }

    case combat.ModeHighestHP:
        bestVal = -1.0
        for _, e := range enemies {
            if !e.IsDead() && e.GetCurrentHP() > bestVal {
                bestVal = e.GetCurrentHP()
                candidate = e
            }
        }

    case combat.ModeRandom:
        var live []combat.Combatant
        for _, e := range enemies {
            if !e.IsDead() {
                live = append(live, e)
            }
        }
        if len(live) > 0 {
            candidate = live[rand.Intn(len(live))]
        }
    }

    return candidate
}
