package api

import (
	"fmt"
	"os"
	"time"

	"github.com/Not-Cyrus/GoGuardian/utils"
	"github.com/valyala/fastjson"

	"github.com/Not-Cyrus/GoGuardian/commands"
	"github.com/Not-Cyrus/GoGuardian/handlers"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) Setup() {
	token := getToken()
	if len(token) == 0 {
		fmt.Print("Enter your token: ")
		fmt.Scan(&token)
	}

	b.DS, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("Couldn't use that token: %s", err.Error())
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	b.BU, err = b.DS.User("@me")
	if err != nil {
		fmt.Printf("Couldn't get a local User: %s", err.Error())
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	b.DS.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMembers)

	handlerNames := []interface{}{handlers.BanHandler, handlers.ChannelCreate, handlers.ChannelRemove, handlers.KickHandler, handlers.MemberAdded, handlers.ReadyHandler, handlers.MemberRoleUpdate, handlers.RoleCreate, handlers.RoleRemove, handlers.RoleUpdate}
	for _, handler := range handlerNames {
		b.DS.AddHandler(handler)
	}

	b.DS.AddHandler(route.MessageCreate)

	// yeah lets just hide the sheer UGLINESS of this V on one line and forget about it ok? thanks
	route.Add("config", "- Turns on/off a protection setting that you specify. (type help config for more info)", []string{"antiadminrole\n(Needs audit,ban,manage roles perms)\nblocks people/bots from making roles have Administator permissions", "antiban\n(needs audit,ban perms)\n blocks wizz/nuke bots from banning everyone from your guild", "antibots\n(audits,ban perms)\nbans any bot that gets invited (and the person who invited it)", "antichannelspam\n(audits,ban perms)\nstops people/bots from spam creating channels (mainly bots)", "antichannelnuke\n(audits,ban perms)\nstops people/bots from deleting lots of channels (mainly bots)", "antihijack\n(audit perms)\nif GoGuardian ever gets comprimised and starts \"nuking\" the server it will leave.", "antikick\n(audits,ban perms)\nblocks wizz/nuke/destroyer bots from kicking everyone", "antimemberadmin\n(audit,manage roles,ban perms)\nblocks people/bots from giving a member a role that has Administator permissions", "antirolespam\n(audits, ban perms)\nstops people/bots from spamm creating roles (mainly bots)", "antirolenuke\n(audits,ban perms)\nstops people/bots from deleting lots of roles (mainly bots)", "seconds\nhow many seconds it'll take before that moderation action is no longer classified as \"malicious\" by the bot.", "threshold\nchanges the amount of times someone can do a certain moderation action in x seconds (see above)"}, route.Config, true, true)

	route.Add("help", "Literally just help info what else do you expect a cookie?", []string{"WOW, You got a cookie!"}, route.Help, false, false)

	route.Add("invite", "Gives you an invite link to invite GoGuardian to your server.", []string{"ANOTHER COOKIE!!! (I'm running out of ideas for what to put here)"}, route.Invite, false, false)

	route.Add("whitelist", "whitelists a user so that they do not get affected by any protection settings.", []string{"whitelist @user"}, route.AddWhitelist, true, true)

	route.Add("unwhitelist", "unwhitelists a user so that they are affected by any toggled protection settings (default for all users minus the Guild owner.)", []string{"unwhitelist @user"}, route.RemoveWhitelist, true, true)
}

func (b *Bot) Run() error {
	err := b.DS.Open()
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) Stop() {
	b.DS.Close()
}

func getToken() string {
	fileContents := utils.ReadFile("Config.json")
	parsed, err := parser.Parse(fileContents)
	if err != nil {
		fmt.Printf("Couldn't parse Config.json to get your Token: %s", err.Error())
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}
	if fastjson.Exists([]byte(fileContents), "Token") {
		return string(parsed.GetStringBytes("Token"))
	}
	return ""
}

type (
	Bot struct {
		DS *discordgo.Session
		BU *discordgo.User
	}
)

var (
	err    error
	token  string
	parser fastjson.Parser
	route  = commands.New()
)
