package data

var XPTable = map[int]map[int]float64{
    1: { 
        2: 250, 3: 650, 4: 1500, 5: 3200, 6: 6400, 
        7: 12000, 8: 22200, 9: 40550, 10: 73000, 
    },
    // Add other incarnations as needed...
}

func GetXPToNextLevel(level, incarnation int) float64 {
    if level >= 30 { return 999999999 }
    
    incMap, ok := XPTable[incarnation]
    if !ok { incMap = XPTable[1] } 
    
    req, ok := incMap[level+1]
    if !ok { return 5000 } // Fallback
    
    return req
}
