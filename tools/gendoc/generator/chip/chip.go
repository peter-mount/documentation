package chip

import (
	"context"
	"fmt"
	"github.com/peter-mount/documentation/tools/gendoc"
	"github.com/peter-mount/documentation/tools/gendoc/generator"
	"github.com/peter-mount/documentation/tools/gendoc/hugo"
	"github.com/peter-mount/go-kernel/v2/log"
	"github.com/peter-mount/go-kernel/v2/util"
	"github.com/peter-mount/go-kernel/v2/util/task"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"os"
	"path"
)

// Chip handles the generation of CHIP pin layout images
type Chip struct {
	generator *generator.Generator `kernel:"inject"` // Generator
	worker    task.Queue           `kernel:"worker"` // Worker queue
	excel     *generator.Excel     `kernel:"inject"` // Excel
	chips     *Category            // Map of named chip definitions
	extracted util.Set             // Set of book ID's so that we run once per book
}

// DefinitionHandler is called when generating the shortcode for this Definition
type DefinitionHandler func(*Definition) error

// Definition holds the details for this chip
type Definition struct {
	Name        string            `yaml:"name"`        // Unique Name of this chip
	Category    string            `yaml:"category"`    // Category of this chip
	SubCategory string            `yaml:"subCategory"` // Sub category of this chip
	Title       string            `yaml:"title"`       // Title to describe this chip
	Type        string            `yaml:"type"`        // Type of chip, e.g. "DIP"
	Label       string            `yaml:"label"`       // Label, e.g. "MOS"
	SubLabel    string            `yaml:"SubLabel"`    // Sub Label, e.g. "6502"
	PinCount    int               `yaml:"PinCount"`    // Number of pins
	PinOffset   int               `yaml:"PinOffset"`   // Offset from where 1 normally is located, default=0
	Pins        map[int]string    `yaml:"pins"`        // Pin title definitions
	Weight      int               `yaml:"weight"`      // Weight of chip, 0=natural
	FileInfo    os.FileInfo       `yaml:"-"`           // FileInfo of containing file
	handler     DefinitionHandler // Handler for this chip type
}

// Generate implements the GeneratorTask interface to generate the shortcode for this chip Definition.
func (d *Definition) Generate(_ context.Context) error {
	if d.handler != nil {
		return d.handler(d)
	}

	return fmt.Errorf("%s has no generatator handler", d.Name)
}

// Path returns the path to the generated shortcode
func (d *Definition) Path(dir string) string {
	//a := []string{"themes/area51/layouts/shortcodes/generated/chip"}
	a := []string{dir}
	if d.Category != "" {
		a = append(a, d.Category)
	}
	if d.SubCategory != "" {
		a = append(a, d.SubCategory)
	}
	a = append(a, d.Name)
	return path.Join(a...)
}

func (c *Chip) Start() error {
	c.chips = NewCategory()
	c.extracted = util.NewHashSet()

	c.generator.
		Register("chipDefinitions",
			task.Of().
				Then(c.extract))

	c.generator.Register("chipReferenceTables",
		task.Of().
			Then(c.extract).
			Then(c.chipReferenceTables))

	return nil
}

// extract gets chip definitions from a book.
// This runs once per book not once ever
func (c *Chip) extract(ctx context.Context) error {
	book := generator.GetBook(ctx)

	// Only run once per Book ID
	if c.extracted.Contains(book.ID) {
		return nil
	}

	// Prevent us running again for this Book
	c.extracted.Add(book.ID)

	log.Printf("Scanning %s for chip designs", book.ID)

	return walk.NewPathWalker().
		IsFile().
		PathNotContain("/reference/").
		PathHasSuffix(".html").
		Then(hugo.FrontMatterActionOf().
			OtherExists("chip", c.extractChipDefinitions).
			Walk(ctx)).
		Walk(book.ContentPath())
}

// extractChipDefinitions extracts all definitions from a specific hugo page
func (c *Chip) extractChipDefinitions(ctx context.Context, _ *hugo.FrontMatter) error {
	c.worker.AddPriorityTask(tools.PriorityChip,
		task.Of(c.extractChipDefTask).
			WithContext(ctx, "other", "fileInfo"))
	return nil
}

func (c *Chip) extractChipDefTask(ctx context.Context) error {
	return util.ForEachInterface(ctx.Value("other"), func(e interface{}) error {
		return util.IfMap(e, func(m map[interface{}]interface{}) error {
			pinCount, ok := util.DecodeInt(m["pinCount"], 0)
			if !ok || pinCount < 1 {
				return fmt.Errorf("invalid pinCount %d", pinCount)
			}

			pinOffset, ok := util.DecodeInt(m["pinOffset"], 0)

			weight, _ := util.DecodeInt(m["weight"], 0)

			v := &Definition{
				Name:        util.DecodeString(m["name"], ""),
				Category:    util.DecodeString(m["category"], ""),
				SubCategory: util.DecodeString(m["subCategory"], ""),
				Title:       util.DecodeString(m["title"], ""),
				Type:        util.DecodeString(m["type"], ""),
				Label:       util.DecodeString(m["label"], ""),
				SubLabel:    util.DecodeString(m["subLabel"], ""),
				PinCount:    pinCount,
				PinOffset:   pinOffset,
				Pins:        make(map[int]string),
				Weight:      weight,
				FileInfo:    ctx.Value("fileInfo").(os.FileInfo),
			}

			if v.Type == "pga" {
				if err := util.IfMap(m["pins"], func(m map[interface{}]interface{}) error {
					return pgaDecodePins(v, m)
				}); err != nil {
					return err
				}
			} else {
				// For non-pga check pins are present
				if err := util.IfMap(m["pins"], func(m map[interface{}]interface{}) error {
					for pinId, label := range m {
						key, ok := util.DecodeInt(pinId, 0)
						if !ok {
							return fmt.Errorf("%s invalid Pin \"%s\"", v.Name, pinId)
						}
						if _, exists := v.Pins[key]; exists {
							return fmt.Errorf("%s pin %d already defined", v.Name, key)
						}
						v.Pins[key] = label.(string)
					}
					return nil
				}); err != nil {
					return err
				}

				// Ensure we have an entry for each pin
				for pin := 1; pin <= pinCount; pin++ {
					if _, exists := v.Pins[pin]; !exists {
						v.Pins[pin] = "Undefined()"
						log.Printf("%s Pin %d is missing", v.Name, pin)
					}
				}

				// Fail if counts don't match, can be a mis-labeled entry
				if len(v.Pins) != pinCount {
					return fmt.Errorf("%s pin count missmatch, expecting %d got %d", v.Name, pinCount, len(v.Pins))
				}
			}

			// verify chip is correct
			switch v.Type {
			case "":
				return fmt.Errorf("%s chip type is mandatory", v.Name)
			case "dip":
				if (pinCount % 2) == 1 {
					return fmt.Errorf("%s dip must have even number of pins, got %d", v.Name, pinCount)
				}
				v.handler = dip
			case "lccc":
				if (pinCount % 4) != 0 {
					return fmt.Errorf("%s lccc must be divisible by 4, got %d", v.Name, pinCount)
				}
				v.handler = lccc
			case "lqfp":
				if (pinCount % 4) != 0 {
					return fmt.Errorf("%s lccc must be divisible by 4, got %d", v.Name, pinCount)
				}
				v.handler = lqfp
			case "plcc":
				if (pinCount % 4) != 0 {
					return fmt.Errorf("%s lccc must be divisible by 4, got %d", v.Name, pinCount)
				}
				v.handler = plcc
			case "qfp":
				if (pinCount % 4) != 0 {
					return fmt.Errorf("%s lccc must be divisible by 4, got %d", v.Name, pinCount)
				}
				v.handler = qfp
			case "pga":
				if pinCount < 2 {
					return fmt.Errorf("%s pga must be >= 2, got %d", v.Name, pinCount)
				}
				v.handler = pga
			default:
				return fmt.Errorf("%s unsupported chip type \"%s\"", v.Name, v.Type)
			}

			if !c.chips.Put(v) {
				return fmt.Errorf("%s already defined", v.Name)
			}

			task.GetQueue(ctx).
				AddPriorityTask(tools.PriorityChip, v.Generate)
			return nil
		})
	})
}
