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
			} else if messageArgs[0] == Prefix + "sortinghat" && len(sortingQueue) >= 1 {
				s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+", someone is being sorted, please wait!")
				fmt.Println(m.Author, "tried to sort but the queue is full.")
			} else if messageArgs[0] == Prefix + "changehouse" && len(messageArgs) == 2 {
				userRoles := [] string{}
				getMember, err := s.State.Member(serverID, m.Author.ID)
				if err != nil{
					return
				}else{
					for _, n := range getMember.Roles {
						role, err := s.State.Role(serverID, n)
						if err != nil {
							return
						}
						userRoles = append(userRoles, role.Name)
					}
					if checkUserRole(userRoles) == true{
						roleIDMap := map[string]string{"Slytherin": slytherinID, "Gryffindor": gryffindorID, "Hufflepuff": hufflepuffID, "Ravenclaw": ravenclawID}
						userCurrentHouse := checkRoleToRemove(userRoles)
						userFutureHouse := messageArgs[1]
						if strings.Title(getRoleName(userCurrentHouse)) == strings.Title(userFutureHouse){
							s.ChannelMessageSend(m.ChannelID, "You are already in that house, " +
								"please do not waste my time.")
						}else{
							if checkIfInHouses(userFutureHouse) == true{
								s.GuildMemberRoleRemove(serverID, m.Author.ID, userCurrentHouse)
								s.GuildMemberRoleAdd(serverID, m.Author.ID, roleIDMap[strings.Title(userFutureHouse)])
								s.ChannelMessageSend(m.ChannelID, m.Author.Mention() + ", as you wish. I have " +
									"removed you from " + strings.Title(getRoleName(userCurrentHouse)) +
									" and put you in " + strings.Title(userFutureHouse) + ".")
							}else{
								s.ChannelMessageSend(m.ChannelID, userFutureHouse + " is not a house!")
							}
						}
					}else{
						s.ChannelMessageSend(m.ChannelID, "I have not assigned you to a house.\n" +
							"Please use `*sortinghat`, and then if you are unhappy, run `*changehouse <house>`.")
					}
				}

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

func checkRoleToRemove(userRoles [] string) string{
	for _, n := range userRoles{
		if n == "Gryffindor"{
			return gryffindorID
		}else if n == "Slytherin"{
			return slytherinID
		}else if n == "Hufflepuff"{
			return hufflepuffID
		}else if n == "Ravenclaw"{
			return ravenclawID
		}
	}
	return "None"
}

func checkIfInHouses(house string) bool {
	house = strings.ToLower(house)
	if house == "slytherin" || house == "hufflepuff" || house == "ravenclaw" || house == "gryffindor"{
		return true
	}
	return false
}

func getRoleName(roleID string) string{
	if roleID == slytherinID{
		return "Slytherin"
	}else if roleID == hufflepuffID{
		return "Hufflepuff"
	}else if roleID == ravenclawID{
		return "Ravenclaw"
	}else if roleID == gryffindorID{
		return "Gryffindor"
	}
	return "None"
}
