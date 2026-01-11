package main

import (
    "log"
    "fmt"
    "sync" // Added missing import
    
    "darkscar/internal/game/persistence"
    "darkscar/internal/game/entities"
    "darkscar/internal/game/engine"
    "darkscar/internal/game/combat"

    "github.com/pocketbase/pocketbase"
    "github.com/pocketbase/pocketbase/core"
    "github.com/pocketbase/pocketbase/plugins/migratecmd"
)

func main() {
    app := pocketbase.New()

    // 1. Initialize Persistence Store
    var store *persistence.Store
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        store = persistence.NewStore(app)
        ensureSchema(app)
        return e.Next()
    })

    // 2. Register Custom Game API Routes
    app.OnServe().BindFunc(func(e *core.ServeEvent) error {
        
        // GROUP: /api/game
        g := e.Router.Group("/api/game")
        
        // GET /api/game/list - List my characters
        g.GET("/list", func(e *core.RequestEvent) error {
            userID := "test_user" 
            
            chars, err := store.GetCharactersByUser(userID)
            if err != nil {
                return e.Error(500, "Failed to fetch characters", err)
            }
            return e.JSON(200, chars)
        })

        // POST /api/game/create - Create a new character
        g.POST("/create", func(e *core.RequestEvent) error {
            data := struct {
                Name    string `json:"name"`
                ClassID string `json:"class_id"`
            }{}
            if err := e.BindBody(&data); err != nil {
                return e.Error(400, "Invalid JSON", err)
            }
            
            // Logic: Create Level 1 Hero
            newChar := entities.NewCharacter(data.Name, data.ClassID, 1, "Players")
            
            // Save to DB
            if err := store.SaveCharacter(newChar, "test_user"); err != nil {
                return e.Error(500, "Failed to save character", err)
            }
            
            return e.JSON(200, newChar)
        })

        // POST /api/game/start - Launch Simulation for a Character
        g.POST("/start", func(e *core.RequestEvent) error {
            data := struct {
                CharacterID string `json:"character_id"`
            }{}
            if err := e.BindBody(&data); err != nil {
                return e.Error(400, "Invalid Data", err)
            }
            
            // 1. Load Character from DB
            hero, err := store.LoadCharacter(data.CharacterID)
            if err != nil {
                return e.Error(404, "Character not found", err)
            }
            
            // 2. Setup a Dummy Fight (1v1 vs Rat)
            enemy := entities.NewCharacter("Sewer_Rat", "martyr", 1, "Enemies")
            
            // 3. Start Session (Async)
            party := []combat.Combatant{hero}
            session := engine.NewSession("API_Session", party)
            
            var wg sync.WaitGroup
            go func() {
                 session.StartCombat([]combat.Combatant{enemy}, &wg)
                 wg.Wait()
                 // After fight, save state to persist XP/HP changes
                 store.SaveCharacter(hero, "test_user")
                 fmt.Println(">> Fight Finished. Progress Saved.")
            }()
            
            return e.JSON(200, map[string]string{"status": "Combat Started", "session_id": session.UserID})
        })

        return e.Next()
    })

    migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
        Automigrate: true,
    })

    if err := app.Start(); err != nil {
        log.Fatal(err)
    }
}

// --- HELPER: Schema Setup ---
func ensureSchema(app core.App) {
    _, err := app.FindCollectionByNameOrId("characters")
    if err == nil { return } 

    fmt.Println(">> INITIALIZING DATABASE SCHEMA...")
    col := core.NewBaseCollection("characters")
    
    col.Fields.Add(&core.TextField{Name: "user_id", Required: true})
    col.Fields.Add(&core.TextField{Name: "name", Required: true})
    col.Fields.Add(&core.TextField{Name: "class_id", Required: true})
    col.Fields.Add(&core.NumberField{Name: "level", Required: true})
    col.Fields.Add(&core.NumberField{Name: "incarnation"})
    col.Fields.Add(&core.NumberField{Name: "xp"})
    col.Fields.Add(&core.NumberField{Name: "current_hp"})
    col.Fields.Add(&core.JSONField{Name: "equipment"})
    col.Fields.Add(&core.TextField{Name: "active_skill"})
    col.Fields.Add(&core.TextField{Name: "passive_skill"})
    
    if err := app.Save(col); err != nil {
        log.Fatal("Failed to create table: ", err)
    }
    fmt.Println(">> SCHEMA CREATED SUCCESSFULLY.")
}
