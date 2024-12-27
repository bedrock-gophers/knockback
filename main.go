package main

import (
	"log/slog"
	"time"

	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/bedrock-gophers/knockback/knockback"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
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

	chat.Global.Subscribe(chat.StdoutSubscriber{})
	c := server.DefaultConfig()
	c.Players.SaveData = false

	conf, err := c.Config(slog.Default())
	if err != nil {
		logrus.Fatalln(err)
	}

	srv := conf.New()
	srv.CloseOnProgramEnd()
	cmd.Register(cmd.New("kb", "", nil, knockback.Menu{}))

	srv.Listen()

	for p := range srv.Accept() {
		intercept.Intercept(p)
		p.Handle(handler{})
		p.SetGameMode(world.GameModeSurvival)
	}
}

type handler struct {
	player.NopHandler
}

func (handler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	if immune {
		return
	}

	knockback.ApplyHitDelay(attackImmunity)
}

func (handler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	knockback.ApplyForce(force)
	knockback.ApplyHeight(height)
}
