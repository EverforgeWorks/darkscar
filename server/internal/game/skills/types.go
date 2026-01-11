package skills

import "darkscar/internal/game/combat"

type SkillDef struct {
    ID          string
    Name        string
    ManaCost    float64
    HPCostPct   float64 // 0.10 = 10% Max HP
    Cooldown    float64 // Seconds
    IsPassive   bool
    
    // The Logic: What happens when triggered?
    OnCast func(caster combat.Combatant, allies, enemies []combat.Combatant)
}

type Effect struct {
    ID        string
    Name      string
    Duration  float64 // Seconds remaining
    
    // Stat Modifiers
    StatMod   map[string]float64 
    
    // Hooks
    OnTick    func(owner combat.Combatant, dt float64)
    OnHit     func(owner, target combat.Combatant, damage float64) 
    OnGetHit  func(owner, source combat.Combatant, damage float64)
}
