package main

import (
	"time"

	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/bedrock-gophers/inv/inv"
	"github.com/bedrock-gophers/knockback/knockback"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sirupsen/logrus"
)

func main() {
	err := knockback.Load("assets/knockback.json")
	if err != nil {
		panic(err)
	}
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	chat.Global.Subscribe(chat.StdoutSubscriber{})
	c := server.DefaultConfig()
	c.Players.SaveData = false

	conf, err := c.Config(log)
	if err != nil {
		log.Fatalln(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()
	cmd.Register(cmd.New("kb", "", nil, knockback.Menu{}))

	srv.Listen()
	inv.PlaceFakeContainer(srv.World(), cube.Pos{0, 255, 0})

	for srv.Accept(func(p *player.Player) {
		intercept.Intercept(p)
		p.Handle(handler{})
		p.SetGameMode(world.GameModeSurvival)
	}) {

	}
}

type handler struct {
	player.NopHandler
}

func (handler) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource) {
	knockback.ApplyHitDelay(attackImmunity)
}

func (handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	knockback.ApplyForce(force)
	knockback.ApplyHeight(height)
}
