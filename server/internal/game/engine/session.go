package engine

import (
    "fmt"
    "time"
    "sync"
    "darkscar/internal/game/combat"
    "darkscar/internal/game/skills"
    "darkscar/internal/game/entities"
    "darkscar/internal/game/data"
)

type Session struct {
    UserID string
    PlayerParty []combat.Combatant
    EnemyParty  []combat.Combatant
    Cooldowns map[string]float64 
    RegenTimer float64
    quit      chan bool
}

func NewSession(uid string, players []combat.Combatant) *Session {
    return &Session{
        UserID: uid, PlayerParty: players, EnemyParty: make([]combat.Combatant, 0),
        Cooldowns: make(map[string]float64), quit: make(chan bool),
    }
}

func (s *Session) StartCombat(enemies []combat.Combatant, wg *sync.WaitGroup) {
    s.EnemyParty = enemies
    wg.Add(1)
    s.applyPassives(s.PlayerParty, s.EnemyParty)
    s.applyPassives(s.EnemyParty, s.PlayerParty)
    for _, c := range s.PlayerParty { s.Cooldowns[c.GetID()] = 0.5 + (0.1 * c.GetAttackSpeed()) }
    for _, c := range s.EnemyParty  { s.Cooldowns[c.GetID()] = 1.0 + (0.1 * c.GetAttackSpeed()) }

    go func() {
        defer wg.Done()
        dt := 0.1
        ticker := time.NewTicker(100 * time.Millisecond)
        defer ticker.Stop()
        fmt.Printf("[Session %s] ⚔️ Combat Started (%d vs %d)\n", s.UserID, len(s.PlayerParty), len(s.EnemyParty))
        for {
            select {
            case <-ticker.C: if !s.tick(dt) { return }
            case <-s.quit: return
            }
        }
    }()
}

func (s *Session) applyPassives(team, opponents []combat.Combatant) {
    for _, member := range team {
        if char, ok := member.(*entities.Character); ok {
            if def, exists := skills.Registry[char.PassiveSkillID]; exists {
                def.OnCast(char, team, opponents)
                fmt.Printf("[%s] activates Passive: %s\n", char.GetID(), def.Name)
            }
        }
    }
}

func (s *Session) tick(dt float64) bool {
    activePlayers := s.processPartyTurn(s.PlayerParty, s.EnemyParty, dt)
    activeEnemies := s.processPartyTurn(s.EnemyParty, s.PlayerParty, dt)

    s.RegenTimer += dt
    if s.RegenTimer >= 5.0 {
        s.RegenTimer = 0
        s.applyRegen(s.PlayerParty)
        s.applyRegen(s.EnemyParty)
    }

    if activePlayers == 0 { fmt.Printf(">> [%s] DEFEAT.\n", s.UserID); return false }
    if activeEnemies == 0 { fmt.Printf(">> [%s] VICTORY.\n", s.UserID); return false }
    return true
}

func (s *Session) applyRegen(team []combat.Combatant) {
    for _, c := range team { if !c.IsDead() { c.Heal(c.GetMaxHP() * 0.10) } }
}

func (s *Session) processPartyTurn(allies, enemies []combat.Combatant, dt float64) int {
    aliveCount := 0
    for _, actor := range allies {
        if actor.IsDead() { continue }
        aliveCount++
        char, isChar := actor.(*entities.Character)
        if isChar {
            activeBuffs := char.Buffs[:0]
            for _, b := range char.Buffs {
                b.Duration -= dt
                if b.OnTick != nil { b.OnTick(char, dt) }
                if b.Duration > 0 { activeBuffs = append(activeBuffs, b) }
            }
            char.Buffs = activeBuffs
            s.tryCastSkill(char, allies, enemies, dt)
        }
        id := actor.GetID()
        s.Cooldowns[id] -= dt
        if s.Cooldowns[id] <= 0 {
            target := SelectTarget(actor, enemies)
            if target != nil {
                s.performAttack(actor, target)
                s.Cooldowns[id] = actor.GetAttackSpeed()
            }
        }
    }
    return aliveCount
}

func (s *Session) tryCastSkill(c *entities.Character, allies, enemies []combat.Combatant, dt float64) {
    def, exists := skills.Registry[c.ActiveSkillID]
    if !exists { return }
    if c.SkillCooldowns[def.ID] > 0 { c.SkillCooldowns[def.ID] -= dt; return }
    if c.MP < def.ManaCost { return }
    hpCost := c.MaxHP * def.HPCostPct
    if c.HP <= hpCost { return }
    c.MP -= def.ManaCost; c.HP -= hpCost; c.SkillCooldowns[def.ID] = def.Cooldown
    def.OnCast(c, allies, enemies)
}

func (s *Session) performAttack(attacker, defender combat.Combatant) {
    result := combat.CalculateHit(attacker, defender, combat.Physical, 50) 
    if result.IsMiss { fmt.Printf("%s attacks %s -> MISS\n", attacker.GetID(), defender.GetID()); return }
    defender.TakeDamage(result.Damage, attacker)
    if result.IsCrit { fmt.Printf("%s CRITS %s! (%.0f)\n", attacker.GetID(), defender.GetID(), result.Damage) }
    
    if defender.IsDead() {
        fmt.Printf(">> %s KILLED %s!\n", attacker.GetID(), defender.GetID())
        _, isEnemy := defender.(*entities.Character)
        if isEnemy && attacker.GetTeamID() == "Players" {
            s.distributeXP(defender)
        }
    } 
}

func (s *Session) distributeXP(victim combat.Combatant) {
    mob, ok := victim.(*entities.Character)
    if !ok { return }
    
    // 1. Calculate Base XP (Level * 50 * Threat)
    // e.g. Lvl 20 Mob (Threat 5) = 5000 XP
    arch, _ := data.GetArchetype(mob.ClassID)
    baseXP := float64(mob.Level) * 50.0 * arch.ThreatRating
    
    // 2. Calculate Party Bonus Multiplier
    // Formula: 1.0 + (ExtraMembers * 0.10)
    // Party of 1: 1.0
    // Party of 3: 1.2
    // Party of 5: 1.4
    partySize := float64(len(s.PlayerParty))
    partyBonus := 1.0 + ((partySize - 1.0) * 0.10)
    
    // 3. Final Amount per person (No Splitting!)
    xpAmount := baseXP * partyBonus

    fmt.Printf("   [XP Event] Mob Value: %.0f | Party Bonus: %.0f%% | Total: %.0f\n", baseXP, (partyBonus-1)*100, xpAmount)
    
    // 4. Distribute to Survivors
    for _, member := range s.PlayerParty {
        char, ok := member.(*entities.Character)
        if ok && !char.IsDead() {
            // Note: If you add 'XP_BONUS' gear later, multiply it here!
            // e.g. xpAmount *= char.GetStat("XP_BONUS")
            
            leveled := char.AddXP(xpAmount)
            fmt.Printf("   + %.0f XP -> %s", xpAmount, char.GetID())
            if leveled { fmt.Printf(" (LEVEL UP! %d)\n", char.Level) } else { fmt.Printf("\n") }
        }
    }
}
