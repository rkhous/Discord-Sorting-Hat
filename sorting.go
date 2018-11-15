package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

//CONFIG
var Token string //Bot Token will go under the main function!
var Prefix = "*" //Prefix for the bot. Limited to one char.
var sortingQueue = [] string {} //Leave untouched! This is for queues.
var serverID = "" //Your servers ID
var slytherinID = "" //ID for slytherin role
var gryffindorID = "" //ID for gryffindor role
var hufflepuffID = "" //ID for hufflepuff role
var ravenclawID = "" //ID for ravenclaw role

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	Token := "" //Add bot token here
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is connected to discord. Ctrl+C to quit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	if len(m.Content) == 0 {
		return
	} else {
		var getPrefix string = strings.Split(m.Content, "")[0]
		if getPrefix == Prefix {
			messageArgs := strings.Split(m.Content, " ")
			if messageArgs[0] == Prefix + "sortinghat" && len(sortingQueue) == 0 {
				getMember, err := s.State.Member(serverID, m.Author.ID)
				if err != nil {
					return
				} else {
					userRoles := [] string{}
					for _, n := range getMember.Roles {
						role, err := s.State.Role(serverID, n)
						if err != nil {
							return
						}
						userRoles = append(userRoles, role.Name)
					}
					if checkUserRole(userRoles) == true {
						s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+", I have already assigned you to a house, do you mistake me for a fool?")
					} else {
						roleIDMap := map[string]string{"Slytherin": slytherinID, "Gryffindor": gryffindorID, "Hufflepuff": hufflepuffID, "Ravenclaw": ravenclawID}
						sortingQueue = append(sortingQueue, m.Author.ID)
						sortingChoices := [] string{"Slytherin", "Gryffindor", "Hufflepuff", "Ravenclaw"}
						rand.Seed(time.Now().Unix())
						explainHouses := [] string{"You might belong in Gryffindor", "Where dwell the brave at heart",
							"Their daring, nerve, and chivalry", "Set Gryffindors apart;", "You might belong in Hufflepuff",
							"Where they are just and loyal", "Those patient Hufflepuffs are true", "And unafraid of toil;",
							"Or yet in wise old Ravenclaw", "If you've a ready mind", "Where those of wit and learning",
							"Will always find their kind;", "Or perhaps in Slytherin", "You'll make your real friends",
							"These cunning folks use any means", "To achieve their ends."}
						for _, n := range explainHouses {
							s.ChannelMessageSend(m.ChannelID, n)
							time.Sleep(time.Second * 2)
						}
						s.ChannelMessageSend(m.ChannelID, "So where shall I send you, "+m.Author.Mention()+"?")
						time.Sleep(time.Second * 2)
						rand.Seed(time.Now().Unix())
						sortingHatChoiceMade := sortingChoices[rand.Intn(len(sortingChoices))]
						sortingQueue = [] string{}
						s.GuildMemberRoleAdd(serverID, m.Author.ID, roleIDMap[sortingHatChoiceMade])
						embed := &discordgo.MessageEmbed{
							Author:      &discordgo.MessageEmbedAuthor{},
							Color:       0x00ff00,
							Description: "**" + m.Author.Username + ",\n Welcome to " + sortingHatChoiceMade + "**",
							Thumbnail: &discordgo.MessageEmbedThumbnail{
								URL: "https://raw.githubusercontent.com/rkhous/Discord-Sorting-Hat/master/" + strings.ToLower(sortingHatChoiceMade) + ".png",
							},
							Footer:	   &discordgo.MessageEmbedFooter{
										Text:"Created by github.com/rkhous",
										IconURL:"https://d1q6f0aelx0por.cloudfront.net/product-logos/81630ec2-d253-4eb2-b36c-eb54072cb8d6-golang.png"},
							Title:     "The Sorting Hat Has Spoken!",
						}
						s.ChannelMessageSendEmbed(m.ChannelID, embed)
				}
			}
			} else if messageArgs[0] == Prefix+"sortinghat" && len(sortingQueue) >= 1 {
				s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+", someone is being sorted, please wait!")
				fmt.Println(m.Author, "tried to sort but the queue is full.")
			} else {
				return
			}
		}

	}
}

func checkUserRole(userRoles [] string) bool{
	for _, n := range userRoles{
		if n == "Gryffindor" || n == "Slytherin" || n == "Hufflepuff" || n == "Ravenclaw"{
			return true
		}
	}
	return false
}
