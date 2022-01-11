package tools

// Worker task priorities
const (
  PriorityImmediate = -100 // Run task before any others
  PriorityApi       = 20   // API generation
  PriorityIndices   = 25   // 6502 api indices
  PriorityAutodoc   = 30   // Autodoc generation
  PriorityChip      = 50   // Chip reference tables
  PrioritySVG       = 70   // SVG generation
  PriorityHugo      = 100  // Hugo page generation
  PriorityPDF       = 400  // PDF Generation
  PriorityExcel     = 500  // Priority for Excel generation
)
