package processes

import "github.com/spf13/cobra"

func newChatCmd() *cobra.Command {
	chatCmd := &cobra.Command{
		Use:   "chat",
		Short: "Business verbs for data chatbot workflows",
	}

	chatCmd.AddCommand(newChatAskCmd())
	chatCmd.AddCommand(newCollectionListCmd("list", "List chatbot queries", "data-chatbot", "queries"))

	return chatCmd
}

func newChatAskCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ask <domain> <question> [context] [channel] [sender_id] [origin_session_id]",
		Short: "Send a natural-language query to the data chatbot",
		Args:  cobra.RangeArgs(2, 6),
		RunE: func(cmd *cobra.Command, args []string) error {
			payload := map[string]any{
				"domain":   args[0],
				"question": args[1],
			}
			if len(args) >= 3 {
				payload["context"] = args[2]
			}
			if len(args) >= 4 {
				payload["channel"] = args[3]
			}
			if len(args) >= 5 {
				payload["sender_id"] = args[4]
			}
			if len(args) == 6 {
				payload["origin_session_id"] = args[5]
			}
			return runModuleAction(cmd, "data-chatbot", "send_chat_query", payload)
		},
	}
}
