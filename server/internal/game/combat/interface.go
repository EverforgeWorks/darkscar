package combat

type DamageType int
const (
    Physical DamageType = iota
    Magical
)

type TargetMode int
const (
    ModeManual TargetMode = iota
    ModeLowestHP
    ModeHighestHP
    ModeHighestThreat
    ModeRandom
)

type Combatant interface {
    GetID() string
    GetTeamID() string
    
    // Stats
    GetAttack(dt DamageType) float64
    GetDefense(dt DamageType) float64
    GetAttackSpeed() float64
    GetAccuracy() float64
    GetEvasion() float64
    GetCritChance() float64 
    GetCritDamage() float64
    
    // Resources
    GetCurrentHP() float64
    GetMaxHP() float64
    GetCurrentMP() float64
    GetMaxMP() float64
    
    // Actions
    TakeDamage(amount float64, source Combatant)
    Heal(amount float64)
    IsDead() bool
    
    // Threat & AI
    AddThreat(sourceID string, amount float64)
    GetHighestThreatID() string
    GetTargetMode() TargetMode
    
    // Buff Support (Use empty interface to avoid circular deps)
    AddBuff(e interface{}) 
}
