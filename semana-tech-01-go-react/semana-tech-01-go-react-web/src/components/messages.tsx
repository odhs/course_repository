import { useParams } from "react-router-dom"
import { Message } from "./message"
import { getRoomMessages } from "../http/get-room-messages"
import { useSuspenseQuery } from "@tanstack/react-query"
import { useEffect } from "react"


export function Messages() {
  const { roomId } = useParams()

  if (!roomId) {
    throw new Error("Messages component must be used within room page");
  }

  const { data } = useSuspenseQuery({
    queryKey: ["messages", roomId],
    queryFn: () => getRoomMessages({ roomId }),
  })

  // Real Time part
  // useEffect: Se nenhuma variável mudar eu não executo o código dentro de useEffect
  // Caso o ID da sala mude eu refaço a conexão
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/subscribe/${roomId}`)
    ws.onopen = () => {
      console.log('WebSocket connected!')
    }

    ws.onclose = () => {
      console.log('WebSocket closed!')
    }

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data)
      console.log(data)
    }

    return () => {
      ws.close()
    }
  }, [roomId])

  return (
    <ol className="list-decimal list-outside px-3 space-y-8">
      {data.messages.map((message) => {
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
