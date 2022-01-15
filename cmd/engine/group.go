package main

import (
	"fmt"
	msdk "github.com/fox-one/mixin-sdk-go"
	"github.com/spf13/cobra"
	"log"
	"option-dance/cmd/config"
	"option-dance/pkg/mixin"
)

var (
	groupAct string
	value    string
)

const (
	ActionAddMember    = "add"
	ActionRemoveMember = "remove"
	ActionRenameGroup  = "rename"
	ActionAnnouncement = "announcement"
)

func init() {
	GroupManageCmd.PersistentFlags().StringVarP(&groupAct, "act", "a", "", "actions: add / remove / rename / announcement")
	GroupManageCmd.PersistentFlags().StringVarP(&value, "value", "v", "", "action value")
}

var GroupManageCmd = &cobra.Command{
	Use:   "group",
	Short: "GroupManageCmd",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := mixin.Client()
		if err != nil {
			panic(err)
		}
		dapp := config.Cfg.DApp
		switch groupAct {
		case ActionAddMember:
			_, err := client.AddParticipants(cmd.Context(), dapp.GroupConversationId, value)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("add user %s to group %s success", value, dapp.GroupConversationId)
			break
		case ActionRemoveMember:
			_, err := client.RemoveParticipants(cmd.Context(), dapp.GroupConversationId, value)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("remove user %s from group %s success", value, dapp.GroupConversationId)
			break
		case ActionRenameGroup:
			_, err = client.UpdateConversation(cmd.Context(), dapp.GroupConversationId, msdk.ConversationUpdate{
				Name: value,
			})
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("update group name to '%s' success", value)
			break
		case ActionAnnouncement:
			_, err = client.UpdateConversation(cmd.Context(), dapp.GroupConversationId, msdk.ConversationUpdate{
				Announcement: value,
			})
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Printf("publish announcement '%s' success", value)
			break
		default:
			break
		}
	},
}
