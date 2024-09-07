import { useParams } from "react-router-dom"
import { Message } from "./message"
import { getRoomMessages } from "../http/get-room-messages"
import { useMessagesWebSockets } from "../hooks/use-messages-web-sockets"
import { useSuspenseQuery } from "@tanstack/react-query"


export function Messages() {
  const { roomId } = useParams()

  if (!roomId) {
    throw new Error("Messages component must be used within room page");
  }

  const { data } = useSuspenseQuery({
    queryKey: ["messages", roomId],
    queryFn: () => getRoomMessages({ roomId }),
  })

  useMessagesWebSockets({ roomId })

  // Sort messages
  const sortedMessages = data.messages.sort((a, b) => {
    return b.amountOfReactions - a.amountOfReactions
  })

  return (
    <ol className="list-decimal list-outside px-3 space-y-8">
      {sortedMessages.map((message) => {
        return (
          <Message
            key={message.id}
            id={message.id}
            text={message.text}
            amountOfReactions={message.amountOfReactions}
            answered={message.answered}
          />
        )
      })}
    </ol>
  );
}
