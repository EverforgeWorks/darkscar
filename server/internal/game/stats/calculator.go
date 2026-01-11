package stats

import (
    "math"
    "darkscar/internal/game/data"
)

// Diminishing Returns Curve
// points: The raw stat value (from Attributes * Relations)
// min: Baseline value (e.g. 5.0)
// max: Hard cap (e.g. 90.0)
// targetPoints: Cost factor (Higher = slower curve)
func finalizeStat(points, min, max, targetPoints float64) float64 {
    if points <= 0 {
        return min
    }
    k := -math.Log(0.05) / targetPoints
    return min + (max - min) * (1 - math.Exp(-k * points))
}

func CalculateStats(archetype data.Archetype, level int) map[string]float64 {
    derived := make(map[string]float64)

    // 1. Calculate Raw Attributes (Base + Growth)
    lvlMod := float64(level - 1)
    
    str := float64(archetype.BaseAttributes.Str) + (float64(archetype.Growth.Str) * lvlMod)
    dex := float64(archetype.BaseAttributes.Dex) + (float64(archetype.Growth.Dex) * lvlMod)
    intl := float64(archetype.BaseAttributes.Int) + (float64(archetype.Growth.Int) * lvlMod)
    wis := float64(archetype.BaseAttributes.Wis) + (float64(archetype.Growth.Wis) * lvlMod)
    end := float64(archetype.BaseAttributes.End) + (float64(archetype.Growth.End) * lvlMod)

    primaries := map[string]float64{
        "str": str, "dex": dex, "int": intl, "wis": wis, "end": end,
    }

    // 2. Accumulate Raw Points from Relations
    for statName, value := range primaries {
        weights, exists := archetype.Relations[statName]
        if !exists { continue }

        derived["ATK"]       += value * weights.PAtk
        derived["DEF"]       += value * weights.PDef
        derived["MATK"]      += value * weights.MAtk
        derived["MDEF"]      += value * weights.MDef
        derived["HP"]        += value * weights.MaxHP
        derived["MP"]        += value * weights.MaxMP
        derived["ACC"]       += value * weights.Acc
        derived["EVA"]       += value * weights.Eva
        derived["CRIT_RATE"] += value * weights.CritHit // Accumulating POINTS here
        derived["CRIT_DMG"]  += value * weights.CritDmg // Accumulating POINTS here
    }

    // 3. Apply Diminishing Returns (The Fix)
    // Crit Rate: Min 5%, Max 90%, Cost 100k
    rawCritRate := derived["CRIT_RATE"]
    finalCritPercent := finalizeStat(rawCritRate, 5.0, 90.0, 100000.0)
    derived["CRIT_RATE"] = finalCritPercent / 100.0 // Convert 5.0 -> 0.05 for logic

    // Crit Dmg: Min 1.5x, Max 10.0x, Cost 100k
    rawCritDmg := derived["CRIT_DMG"]
    derived["CRIT_DMG"] = finalizeStat(rawCritDmg, 1.5, 10.0, 100000.0)

    // Threat Modifier (Linear)
    derived["THREAT_MOD"] = 0.5 + (archetype.ThreatRating * 0.25)

    return derived
}
