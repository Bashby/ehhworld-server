package object

// Rarity is an enum for the rarity of the object in a game
type Rarity int

//go:generate stringer -type=Rarity

// Rarity enum
const (
	RarityTrash     Rarity = iota // Grey
	RarityCommon                  // White
	RarityUncommon                // Green
	RarityRare                    // Blue
	RarityLegendary               // Purple
	RarityExotic                  // Orange
	RarityMythic                  // Pink
	RarityQuest                   // Brown
)

type ObjectYAMLDefinition struct {
	Id          int
	Name        string
	Description string
	Stackable   bool
	Stacksize   int `yaml:,omitempty`

	Tags       []string `yaml:,omitempty,flow`
	Rarity     Rarity
	Durability float32
}

// Object defines objects held in inventories, traded, bought, sold, dropped, crafted, etc.
type Object struct {
	Count      int
	definition ObjectYAMLDefinition
}

func (i *Object) IsStackable() bool {
	return i.definition.Stackable
}

// func LoadObjectDefinitions() *map[int]*ObjectDefinition {
// 	m := make(map[int]*ObjectDefinition)

// 	// Get all files in Object definition directory
// 	files, err := ioutil.ReadDir("./asset/Object")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, file := range files {
// 		// !file.IsDir()
// 		fmt.Printf("Parsing Object %v.", file.Name())

// 		iDef := ObjectDefinition{}

// 		err := yaml.Unmarshal([]byte(file), &iDef)
// 		if err != nil {
// 			log.Fatalf("error: %v", err)
// 		}

// 	}

// 	return &m
// }
