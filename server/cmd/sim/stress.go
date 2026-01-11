package main

import (
    "fmt"
    "sync"
    "time"
    // "darkscar/internal/game/engine" <-- Removed unused import
    "darkscar/internal/game/entities"
    "darkscar/internal/game/combat"
)

func main() {
    count := 10000 // Ten Thousand Fights
    fmt.Printf("--- STARTING STRESS TEST: %d Concurrent Fights ---\n", count)
    
    var wg sync.WaitGroup
    start := time.Now()

    for i := 0; i < count; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // Build the heavy structs to prove RAM usage
            tank := entities.NewCharacter("Tank", "dreadnought", 10, "A")
            dps := entities.NewCharacter("DPS", "dervish", 10, "A")
            boss := entities.NewCharacter("Boss", "martyr", 20, "B")
            
            // Simulate the math load of 100 ticks (10 seconds of combat)
            for tick := 0; tick < 100; tick++ {
                // 1. Calc Stats
                _ = tank.GetAttack(combat.Physical)
                
                // 2. Run Formula
                res := combat.CalculateHit(dps, boss, combat.Physical, 50)
                boss.TakeDamage(res.Damage, dps)
                
                // 3. Threat Logic
                boss.GetHighestThreatID()
            }
        }(i)
    }

    wg.Wait()
    duration := time.Since(start)
    
    fmt.Printf("--- COMPLETED ---\n")
    fmt.Printf("Process: %d Fights x 100 Ticks (1 Million Combat Rounds)\n", count)
    fmt.Printf("Time Taken: %s\n", duration)
    fmt.Printf("Throughput: %.0f Rounds/Sec\n", 1000000/duration.Seconds())
}
