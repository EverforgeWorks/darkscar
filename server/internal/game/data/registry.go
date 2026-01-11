package data

import "strings"

var Registry = map[string]Archetype{
    "martyr": {
        ID: "martyr", Name: "Martyr", ThreatRating: 4,
        BaseAttributes: Attributes{Str: 120, Dex: 100, Int: 75, Wis: 90, End: 120},
        Growth:         Attributes{Str: 5, Dex: 3, Int: 1, Wis: 1, End: 4},
        Relations: map[string]StatWeights{
            "str": {PAtk: 3.5, PDef: 1, MaxHP: 2.5, Acc: 1.5},
            "dex": {PAtk: 0.5, PDef: 0.5, Eva: 1, Acc: 1.5, CritHit: 0.5, CritDmg: 0.8},
            "int": {MAtk: 0.2, MDef: 0.5, Acc: 1, CritHit: 0.1, CritDmg: 0.1},
            "wis": {MAtk: 0.5, MDef: 1.5, MaxMP: 5.5, MPRegen: 0.05},
            "end": {PDef: 2.5, MDef: 1.5, MaxHP: 3.5, HPRegen: 0.15},
        },
    },
    "dreadnought": {
        ID: "dreadnought", Name: "Dreadnought", ThreatRating: 5,
        BaseAttributes: Attributes{Str: 115, Dex: 80, Int: 70, Wis: 110, End: 130},
        Growth:         Attributes{Str: 4, Dex: 1, Int: 1, Wis: 3, End: 5},
        Relations: map[string]StatWeights{
            "str": {PAtk: 2.5, PDef: 4, MaxHP: 1.5, Acc: 0.5},
            "dex": {PAtk: 0.2, PDef: 1, Eva: 0.5, Acc: 1, CritHit: 0.1, CritDmg: 0.2},
            "int": {MAtk: 0, MDef: 1, Acc: 0.5, CritHit: 0, CritDmg: 0},
            "wis": {MAtk: 0.2, MDef: 2.5, MaxMP: 5, MPRegen: 0.04},
            "end": {PDef: 5.5, MDef: 3.5, MaxHP: 4.5, HPRegen: 0.22},
        },
    },
    "dervish": {
        ID: "dervish", Name: "Dervish", ThreatRating: 3,
        BaseAttributes: Attributes{Str: 80, Dex: 130, Int: 105, Wis: 80, End: 100},
        Growth:         Attributes{Str: 3, Dex: 5, Int: 2, Wis: 1, End: 3},
        Relations: map[string]StatWeights{
            "str": {PAtk: 2.2, PDef: 0.5, MaxHP: 1.5, Acc: 1},
            "dex": {PAtk: 1.8, PDef: 0.5, Eva: 4, Acc: 3.5, CritHit: 4, CritDmg: 4.8},
            "int": {MAtk: 0.5, MDef: 0.5, Acc: 2, CritHit: 3.5, CritDmg: 1.5},
            "wis": {MAtk: 0.5, MDef: 0.5, MaxMP: 5, MPRegen: 0.05},
            "end": {PDef: 1, MDef: 1, MaxHP: 3.5, HPRegen: 0.12},
        },
    },
    "dancer": {
        ID: "dancer", Name: "Dancer", ThreatRating: 1,
        BaseAttributes: Attributes{Str: 90, Dex: 125, Int: 110, Wis: 105, End: 75},
        Growth:         Attributes{Str: 2, Dex: 5, Int: 3, Wis: 3, End: 3},
        Relations: map[string]StatWeights{
            "str": {PAtk: 1.5, PDef: 1.5, MaxHP: 2.2, Acc: 1},
            "dex": {PAtk: 1.2, PDef: 0.8, Eva: 5, Acc: 3, CritHit: 4.4, CritDmg: 2},
            "int": {MAtk: 1.5, MDef: 2.5, Acc: 2.5, CritHit: 3, CritDmg: 3},
            "wis": {MAtk: 1.5, MDef: 3.5, MaxMP: 6.5, MPRegen: 0.08},
            "end": {PDef: 1.5, MDef: 2, MaxHP: 5.5, HPRegen: 0.18},
        },
    },
    "evoker": {
        ID: "evoker", Name: "Evoker", ThreatRating: 3,
        BaseAttributes: Attributes{Str: 80, Dex: 110, Int: 130, Wis: 110, End: 75},
        Growth:         Attributes{Str: 1, Dex: 2, Int: 6, Wis: 3, End: 2},
        Relations: map[string]StatWeights{
            "str": {PAtk: 0.5, PDef: 0.5, MaxHP: 2, Acc: 0.5},
            "dex": {PAtk: 0.8, PDef: 0.5, Eva: 1.5, Acc: 2.5, CritHit: 0.8, CritDmg: 1.2},
            "int": {MAtk: 5.5, MDef: 2, Acc: 4, CritHit: 1, CritDmg: 1.5},
            "wis": {MAtk: 2, MDef: 2, MaxMP: 5, MPRegen: 0.12},
            "end": {PDef: 0.5, MDef: 3, MaxHP: 4.5, HPRegen: 0.15},
        },
    },
    "seer": {
        ID: "seer", Name: "Seer", ThreatRating: 2,
        BaseAttributes: Attributes{Str: 75, Dex: 90, Int: 125, Wis: 115, End: 100},
        Growth:         Attributes{Str: 1, Dex: 2, Int: 5, Wis: 3, End: 3},
        Relations: map[string]StatWeights{
            "str": {PAtk: 0.8, PDef: 1, MaxHP: 2, Acc: 1},
            "dex": {PAtk: 0.2, PDef: 0.5, Eva: 2, Acc: 2, CritHit: 0.8, CritDmg: 1},
            "int": {MAtk: 3, MDef: 4.5, Acc: 3.5, CritHit: 0.4, CritDmg: 0.5},
            "wis": {MAtk: 2, MDef: 5.5, MaxMP: 5.2, MPRegen: 0.14},
            "end": {PDef: 3, MDef: 6, MaxHP: 4.5, HPRegen: 0.12},
        },
    },
    "warden": {
        ID: "warden", Name: "Warden", ThreatRating: 4,
        BaseAttributes: Attributes{Str: 115, Dex: 75, Int: 90, Wis: 110, End: 115},
        Growth:         Attributes{Str: 3, Dex: 1, Int: 2, Wis: 4, End: 5},
        Relations: map[string]StatWeights{
            "str": {PAtk: 2.2, PDef: 2.5, MaxHP: 2, Acc: 1.5},
            "dex": {PAtk: 0.5, PDef: 1, Eva: 1, Acc: 1.5, CritHit: 0.2, CritDmg: 0.3},
            "int": {MAtk: 1, MDef: 1.5, Acc: 1.5, CritHit: 0.1, CritDmg: 0.1},
            "wis": {MAtk: 3.5, MDef: 4, MaxMP: 6.5, MPRegen: 0.12},
            "end": {PDef: 3.5, MDef: 3, MaxHP: 4.5, HPRegen: 0.18},
        },
    },
    "saint": {
        ID: "saint", Name: "Saint", ThreatRating: 1,
        BaseAttributes: Attributes{Str: 75, Dex: 75, Int: 115, Wis: 125, End: 115},
        Growth:         Attributes{Str: 1, Dex: 1, Int: 3, Wis: 5, End: 4},
        Relations: map[string]StatWeights{
            "str": {PAtk: 0.5, PDef: 2, MaxHP: 2, Acc: 0.5},
            "dex": {PAtk: 0.1, PDef: 0.5, Eva: 0.5, Acc: 1, CritHit: 0.1, CritDmg: 0.1},
            "int": {MAtk: 1.5, MDef: 2.5, Acc: 1, CritHit: 0.2, CritDmg: 0.2},
            "wis": {MAtk: 2.5, MDef: 6, MaxMP: 6.5, MPRegen: 0.16},
            "end": {PDef: 4, MDef: 5, MaxHP: 4, HPRegen: 0.2},
        },
    },
}

func GetArchetype(id string) (Archetype, bool) {
    val, ok := Registry[strings.ToLower(id)]
    return val, ok
}
