package ability

// Context that will be sent back to us in the next request
type Context struct {
	LastAbility string      `json:"last_ability,omitempty"`
	SlotFilling interface{} `json:"slot_filling,omitempty"`
}