package gear

// This is mostly WIP, since I'm not handling gear currently (just stat totals)
// In the future, this will be used to build different gear sets for a BiS solver

type Gear struct {
	MainHand  Item
	OffHand   Item
	Head      Item
	Body      Item
	Hands     Item
	Waist     Item
	Legs      Item
	Feet      Item
	Earrings  Item
	Necklace  Item
	Bracelets Item
	Ring1     Item
	Ring2     Item
	Food      Item
	Potion    Item
}
