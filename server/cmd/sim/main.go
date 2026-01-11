package main

import (
    "fmt"
    "sync"
    "darkscar/internal/game/engine"
    "darkscar/internal/game/entities"
    "darkscar/internal/game/combat"
)

func main() {
    fmt.Println("--- STARTING ARCHETYPE SIMULATION (Level 10) ---")
    var wg sync.WaitGroup

    // 1. Level 10 Dreadnought (Tank)
    tank := entities.NewCharacter("Tank_Dread", "dreadnought", 10, "Players")
    
    // 2. Level 10 Dervish (DPS)
    dps := entities.NewCharacter("DPS_Derv", "dervish", 10, "Players")

    // 3. Level 10 Martyr (Bruiser)
    bruiser := entities.NewCharacter("Bruiser_Mart", "martyr", 10, "Players")

    printStats(tank)
    printStats(dps)
    printStats(bruiser)

    // Party Up
    party := []combat.Combatant{tank, dps, bruiser}
    sess := engine.NewSession("User_Archetypes", party)
    engine.GlobalManager.AddSession(sess)

    // Enemy: A Giant Golem
    boss := entities.NewCharacter("Boss_Golem", "martyr", 20, "Enemies") 
    boss.TargetMode = combat.ModeHighestThreat

    enemies := []combat.Combatant{boss}

    sess.StartCombat(enemies, &wg)
    wg.Wait()
}

func printStats(c *entities.Character) {
    fmt.Printf("\n[%s] Stats:\n", c.GetID())
    fmt.Printf(" HP: %.0f | DEF: %.0f\n", c.GetMaxHP(), c.GetDefense(0))
    fmt.Printf(" ATK: %.0f | SPD: %.1f\n", c.GetAttack(0), c.GetAttackSpeed())
    fmt.Printf(" ACC: %.0f | EVA: %.0f\n", c.GetAccuracy(), c.GetEvasion())
    fmt.Printf(" CRIT: %.1f%% | CDMG: %.1fx\n", c.GetCritChance()*100, c.GetCritDamage())
}
