package knockback

import (
	"fmt"
	"github.com/bedrock-gophers/inv/inv"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strconv"
	"strings"
	"time"
)

// Menu is the command that opens the knockback menu.
type Menu struct {
	Allower cmd.Allower
}

// Run ...
func (m Menu) Run(src cmd.Source, _ *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	sendKnockBackMenu(p)
}

// Allow ...
func (m Menu) Allow(src cmd.Source) bool {
	if m.Allower != nil {
		return m.Allower.Allow(src)
	}
	_, ok := src.(*player.Player)
	return ok
}

// selections are the selections that can be made in the knockback menu.
var selections = []string{
	"Force",
	"Height",
	"Hit Delay",
}

// menuSubmittable is a submittable that is used to handle the knockback menu.
type menuSubmittable struct {
	// close is a channel that is closed when the menu should be closed.
	close chan struct{}
	// selectionIndex is the index of the currently selected option.
	selectionIndex int
}

// sendKnockBackMenu sends the knockback menu to a player.
func sendKnockBackMenu(p *player.Player) {
	s := &menuSubmittable{
		close: make(chan struct{}),
	}
	resendMenu(s, p)
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				resendMenu(s, p)
			case <-s.close:
				return
			}
		}
	}()
}

// Submit ...
func (s *menuSubmittable) Submit(p *player.Player, it item.Stack) {
	_, ok := it.Item().(item.Sword)
	if ok {
		s.selectionIndex = (s.selectionIndex + 1) % len(selections)
		resendMenu(s, p)
		return
	}
	change, ok := it.Value("change")
	if !ok {
		return
	}

	switch s.selectionIndex {
	case 0:
		force += change.(float64)
	case 1:
		height += change.(float64)
	case 2:
		hitDelay += time.Duration(change.(float64) * float64(time.Millisecond))
	}

	save()
	resendMenu(s, p)
}

// Close ...
func (s *menuSubmittable) Close() {
	close(s.close)
}

// resendMenu resends the knockback menu to a player.
func resendMenu(s *menuSubmittable, p *player.Player) {
	menu := inv.NewMenu(s, "KnockBack Settings", inv.ContainerChest{})
	stacks := make([]item.Stack, 27)

	var selectedValue any
	switch s.selectionIndex {
	case 0:
		selectedValue = force
	case 1:
		selectedValue = height
	case 2:
		selectedValue = hitDelay
	}

	var amounts = []float64{-0.025, -0.01, -0.001, 0.001, 0.01, 0.025}
	if s.selectionIndex == 2 {
		amounts = []float64{-25, -10, -1, 1, 10, 25}
	}

	stacks[10] = s.formatStack(item.NewStack(block.StainedGlassPane{Colour: item.ColourRed()}, 25), amounts[0], selectedValue)
	stacks[11] = s.formatStack(item.NewStack(block.StainedGlassPane{Colour: item.ColourRed()}, 10), amounts[1], selectedValue)
	stacks[12] = s.formatStack(item.NewStack(block.StainedGlassPane{Colour: item.ColourRed()}, 1), amounts[2], selectedValue)

	var selectionsFormat []string
	for i, selection := range selections {
		col := "grey"
		if i == s.selectionIndex {
			col = "green"
		}
		selectionsFormat = append(selectionsFormat, text.Colourf("<%s>%s</%s>", col, selection, col))
	}

	stacks[13] = item.NewStack(item.Sword{Tier: item.ToolTierGold}, 1).
		WithCustomName(text.Colourf("<red>KnockBack Configuration</red>")).
		WithLore(
			text.Colourf("%s <yellow>%s</yellow>", selectionsFormat[0], formatFloat(force, 4)),
			text.Colourf("%s <yellow>%s</yellow>", selectionsFormat[1], formatFloat(height, 4)),
			text.Colourf("%s <yellow>%dms</yellow> <grey>(%.2f ticks)</grey>", selectionsFormat[2], hitDelay.Milliseconds(), hitDelay.Seconds()*20),
			"",
			text.Colourf("<grey>Click to change selection</grey>"),
		)

	stacks[14] = s.formatStack(item.NewStack(block.StainedGlassPane{Colour: item.ColourGreen()}, 1), amounts[3], selectedValue)
	stacks[15] = s.formatStack(item.NewStack(block.StainedGlassPane{Colour: item.ColourGreen()}, 10), amounts[4], selectedValue)
	stacks[16] = s.formatStack(item.NewStack(block.StainedGlassPane{Colour: item.ColourGreen()}, 25), amounts[5], selectedValue)

	inv.SendMenu(p, menu.WithStacks(stacks...))
}

// formatStack formats a stack with the given amount and value.
func (s *menuSubmittable) formatStack(it item.Stack, amt float64, value any) item.Stack {
	it = it.WithValue("change", amt)
	symbol := text.Colourf("<green>+</green>")
	if amt < 0 {
		symbol = text.Colourf("<red>-</red>")
		amt = -amt
	}
	var currentFormat string
	switch v := value.(type) {
	case float64:
		currentFormat = text.Colourf("<grey>Value: <yellow>%s</yellow></grey>", formatFloat(v, 4))
	case time.Duration:
		currentFormat = text.Colourf("<grey>Value: <yellow>%dms</yellow> <grey>(%.2f ticks)</grey>", v.Milliseconds(), v.Seconds()*20)
	default:
		panic("should never happen")
	}
	return it.
		WithCustomName(text.Colourf("%s %s", symbol, formatFloat(amt, 4))).
		WithLore(
			text.Colourf("<grey>Selected: <yellow>%s</yellow></grey>", selections[s.selectionIndex]),
			currentFormat,
		)
}

// formatFloat formats a float to a string with the given precision.
func formatFloat(num float64, prc int) string {
	var (
		zero, dot = "0", "."

		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)

	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}
