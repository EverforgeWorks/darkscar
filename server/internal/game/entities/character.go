package entities

import (
    "darkscar/internal/game/combat"
    "darkscar/internal/game/data"
    "darkscar/internal/game/stats"
    "darkscar/internal/game/skills" 
)

type Item struct {
    Name      string
    Slot      string
    StatBonus map[string]float64
}

type Character struct {
    ID, Name, ClassID, Team string
    Level, Incarnation      int
    XP                      float64
    
    HP, MaxHP, MP, MaxMP float64
    BaseStats map[string]float64
    Equipment []Item
    
    TargetMode  combat.TargetMode
    ThreatTable map[string]float64 
    
    Buffs          []skills.Effect
    SkillCooldowns map[string]float64
    ActiveSkillID, PassiveSkillID string
}

func NewCharacter(name string, classID string, level int, team string) *Character {
    arch, ok := data.GetArchetype(classID)
    if !ok { arch = data.Registry["dreadnought"] }

    finalStats := stats.CalculateStats(arch, level)

    c := &Character{
        Name: name, ClassID: classID, Level: level, Team: team,
        Incarnation: 1, XP: 0,
        BaseStats: finalStats,
        HP: finalStats["HP"], MaxHP: finalStats["HP"],
        MP: finalStats["MP"], MaxMP: finalStats["MP"],
        ThreatTable: make(map[string]float64),
        TargetMode: combat.ModeRandom,
        Equipment: make([]Item, 0),
        Buffs: make([]skills.Effect, 0),
        SkillCooldowns: make(map[string]float64),
    }
    
    if classID == "dreadnought" { c.ActiveSkillID = "dread_active"; c.PassiveSkillID = "dread_passive" }
    if classID == "dervish"     { c.ActiveSkillID = "dervish_active"; c.PassiveSkillID = "dervish_passive" }
    if classID == "martyr"      { c.ActiveSkillID = "martyr_active"; c.PassiveSkillID = "martyr_passive" }
    if classID == "saint"       { c.ActiveSkillID = "saint_active"; c.PassiveSkillID = "saint_passive" }
    
    // --- SYNTAX FIX HERE ---
    if arch.ThreatRating >= 4 { 
        c.TargetMode = combat.ModeHighestHP 
    } else if arch.ThreatRating <= 2 { 
        c.TargetMode = combat.ModeLowestHP 
    }

    return c
}

func (c *Character) AddXP(amount float64) bool {
    c.XP += amount
    leveledUp := false
    for {
        req := data.GetXPToNextLevel(c.Level, c.Incarnation)
        if c.XP >= req {
            c.XP -= req
            c.LevelUp()
            leveledUp = true
        } else { break }
    }
    return leveledUp
}

func (c *Character) LevelUp() {
    c.Level++
    arch, _ := data.GetArchetype(c.ClassID)
    newStats := stats.CalculateStats(arch, c.Level)
    c.BaseStats = newStats
    c.MaxHP = newStats["HP"]; c.HP = c.MaxHP
    c.MaxMP = newStats["MP"]; c.MP = c.MaxMP
}

func (c *Character) GetID() string { return c.Name }
func (c *Character) GetTeamID() string { return c.Team }
func (c *Character) IsDead() bool { return c.HP <= 0 }
func (c *Character) GetCurrentHP() float64 { return c.HP }
func (c *Character) GetMaxHP() float64 { return c.MaxHP }
func (c *Character) GetCurrentMP() float64 { return c.MP }
func (c *Character) GetMaxMP() float64 { return c.MaxMP }
func (c *Character) Heal(amount float64) { c.HP += amount; if c.HP > c.MaxHP { c.HP = c.MaxHP } }
func (c *Character) AddBuff(e interface{}) { if effect, ok := e.(skills.Effect); ok { c.Buffs = append(c.Buffs, effect) } }

func (c *Character) GetAttack(dt combat.DamageType) float64 {
    val := c.BaseStats["ATK"]
    for _, item := range c.Equipment { val += item.StatBonus["ATK"] }
    pctMod := 0.0
    for _, b := range c.Buffs { if b.StatMod != nil { val += b.StatMod["ATK"]; pctMod += b.StatMod["ATK_PCT"]; pctMod += b.StatMod["DMG_MULT"] } }
    if pctMod != 0 { val *= (1 + pctMod) }
    return val
}
func (c *Character) GetDefense(dt combat.DamageType) float64 {
    val := c.BaseStats["DEF"]
    for _, item := range c.Equipment { val += item.StatBonus["DEF"] }
    for _, b := range c.Buffs { if b.StatMod != nil { val += b.StatMod["DEF"]; if pct, ok := b.StatMod["PDEF_PCT"]; ok { val *= (1 + pct) } } }
    return val
}
func (c *Character) GetAttackSpeed() float64 {
    val := 2.0
    if v, ok := c.BaseStats["SPD"]; ok { val = v }
    for _, b := range c.Buffs { if b.StatMod != nil { if pct, ok := b.StatMod["SPD_PCT"]; ok { val -= (val * pct) } } }
    if val < 0.2 { val = 0.2 }
    return val
}
func (c *Character) GetAccuracy() float64 { val := c.BaseStats["ACC"]; for _, item := range c.Equipment { val += item.StatBonus["ACC"] }; if val < 1 { val = 1 }; return val }
func (c *Character) GetEvasion() float64 { val := c.BaseStats["EVA"]; for _, item := range c.Equipment { val += item.StatBonus["EVA"] }; for _, b := range c.Buffs { if b.StatMod != nil { val += b.StatMod["EVA"]; if pct, ok := b.StatMod["EVA_PCT"]; ok { val *= (1 + pct) } } }; return val }
func (c *Character) GetCritChance() float64 { val := c.BaseStats["CRIT_RATE"]; for _, item := range c.Equipment { val += item.StatBonus["CRIT_RATE"] }; for _, b := range c.Buffs { if b.StatMod != nil { val += b.StatMod["CRIT_RATE"] } }; return val }
func (c *Character) GetCritDamage() float64 { val := c.BaseStats["CRIT_DMG"]; for _, item := range c.Equipment { val += item.StatBonus["CRIT_DMG"] }; if val == 0 { val = 1.5 }; return val }
func (c *Character) TakeDamage(amount float64, source combat.Combatant) { c.HP -= amount; c.AddThreat(source.GetID(), amount); for _, b := range c.Buffs { if b.OnGetHit != nil { b.OnGetHit(c, source, amount) } } }
func (c *Character) AddThreat(sourceID string, amount float64) { c.ThreatTable[sourceID] += amount }
func (c *Character) GetHighestThreatID() string { highestVal := -1.0; highestID := ""; for id, val := range c.ThreatTable { if val > highestVal { highestVal = val; highestID = id } }; return highestID }
func (c *Character) GetTargetMode() combat.TargetMode { return c.TargetMode }
