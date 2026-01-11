package data

type Archetype struct {
    ID             string                 
    Name           string                 
    ThreatRating   float64                
    BaseAttributes Attributes             
    Growth         Attributes             
    Relations      map[string]StatWeights 
}

type Attributes struct {
    Str int 
    Dex int 
    Int int 
    Wis int 
    End int 
}

type StatWeights struct {
    PAtk     float64 
    PDef     float64 
    MAtk     float64 
    MDef     float64 
    MaxHP    float64 
    MaxMP    float64 
    HPRegen  float64 
    MPRegen  float64 
    Acc      float64 
    Eva      float64 
    CritHit  float64 
    CritDmg  float64 
}
