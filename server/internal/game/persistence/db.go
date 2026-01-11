package persistence

import (
    "encoding/json"
    
    "darkscar/internal/game/entities"
    "github.com/pocketbase/pocketbase/core"
)

type Store struct {
    app core.App
}

func NewStore(app core.App) *Store {
    return &Store{app: app}
}

// --- SAVE CHARACTER ---
func (s *Store) SaveCharacter(c *entities.Character, userID string) error {
    collection, err := s.app.FindCollectionByNameOrId("characters")
    if err != nil {
        return err
    }

    // 1. Serialize Complex Data to JSON
    equipJSON, _ := json.Marshal(c.Equipment)
    
    // 2. Find or Create Record
    var record *core.Record
    
    // If we have an ID that looks valid (15 chars), try to find it
    if c.ID != "" && len(c.ID) == 15 { 
        record, err = s.app.FindRecordById("characters", c.ID)
    }
    
    // If not found or new, create one
    if record == nil {
        record = core.NewRecord(collection)
        record.Set("user_id", userID)
    }

    // 3. Map Fields (Go -> DB)
    record.Set("name", c.Name)
    record.Set("class_id", c.ClassID)
    record.Set("level", c.Level)
    record.Set("incarnation", c.Incarnation)
    record.Set("xp", c.XP)
    record.Set("current_hp", c.HP)
    
    record.Set("active_skill", c.ActiveSkillID)
    record.Set("passive_skill", c.PassiveSkillID)
    
    record.Set("equipment", string(equipJSON))

    // 4. Commit to DB
    return s.app.Save(record)
}

// --- LOAD CHARACTER ---
func (s *Store) LoadCharacter(recordID string) (*entities.Character, error) {
    record, err := s.app.FindRecordById("characters", recordID)
    if err != nil {
        return nil, err
    }

    // 1. Extract Core Fields
    name := record.GetString("name")
    classID := record.GetString("class_id")
    level := record.GetInt("level")
    
    // 2. Re-Hydrate the Entity
    char := entities.NewCharacter(name, classID, level, "Players")
    
    // 3. Restore State
    char.ID = record.Id
    char.Incarnation = record.GetInt("incarnation")
    char.XP = record.GetFloat("xp")
    char.HP = record.GetFloat("current_hp")
    
    char.ActiveSkillID = record.GetString("active_skill")
    char.PassiveSkillID = record.GetString("passive_skill")

    // 4. Restore Complex Data
    equipJSON := record.GetString("equipment")
    if equipJSON != "" {
        var items []entities.Item
        if err := json.Unmarshal([]byte(equipJSON), &items); err == nil {
            char.Equipment = items
        }
    }

    return char, nil
}

// --- FIND USER CHARACTERS ---
func (s *Store) GetCharactersByUser(userID string) ([]*entities.Character, error) {
    records, err := s.app.FindRecordsByFilter("characters", "user_id = {:uid}", "-updated", 50, 0, map[string]interface{}{"uid": userID})
    if err != nil {
        return nil, err
    }
    
    var chars []*entities.Character
    for _, r := range records {
        c, _ := s.LoadCharacter(r.Id)
        chars = append(chars, c)
    }
    return chars, nil
}
