import { useParams } from "react-router-dom"
import { Message } from "./message"
import { getRoomMessages, type GetRoomMessagesResponse } from "../http/get-room-messages"
import { useQueryClient, useSuspenseQuery } from "@tanstack/react-query"
import { useEffect } from "react"


export function Messages() {
  const queryClient = useQueryClient()
  const { roomId } = useParams()

  if (!roomId) {
    throw new Error("Messages component must be used within room page");
  }

  const { data } = useSuspenseQuery({
    queryKey: ["messages", roomId],
    queryFn: () => getRoomMessages({ roomId }),
  })

  // Real Time part using Websockts
  // useEffect: Se nenhuma variável mudar eu não executo o código dentro de useEffect
  // Caso o ID da sala mude eu refaço a conexão
  // Esse código irá mudar o já gerado no código acima e armazenado em data 
  // caso receba do servidor alguma mensagem de mudança
  useEffect(() => {
    const ws = new WebSocket(`ws://localhost:8080/subscribe/${roomId}`)
    ws.onopen = () => {
      console.log('WebSocket connected!')
    }

    ws.onclose = () => {
      console.log('WebSocket closed!')
    }

    ws.onmessage = (event) => {
      const data: {
        // filtros
        kind: 'message_created' | 'message_answered' | 'message_reaction_increased' | 'message_reaction_decreased',
        value: any
      } = JSON.parse(event.data)

      console.log(data)

      switch (data.kind) {
        case 'message_created':
          queryClient.setQueryData<GetRoomMessagesResponse>(['messages', roomId], state => {

            console.log(state)

            return {
              messages: [
                // Se houver estado anterior eu retorno senão retona vazio
                // Para ter estado anterior a useSuspenseQuery deve ter executado antes
                ...(state?.messages ?? []),
                {
                  id: data.value.id,
                  text: data.value.message,
                  amountOfReactions: 0,
                  answered: false,
                }
              ],

            }
          })
          break;

        case 'message_answered':
          queryClient.setQueryData<GetRoomMessagesResponse>(['messages', roomId], state => {
            if (!state) {
              return undefined
            }

            return {
              messages: state.messages.map(item => {
                if (item.id == data.value.id) {
                  return { ...item, answered: true }
                }
                return item
              })
            }

          })
          break;

        case 'message_reaction_increased':
        case 'message_reaction_decreased':
          queryClient.setQueryData<GetRoomMessagesResponse>(['messages', roomId], state => {
            if (!state) {
              return undefined
            }

            return {
              messages: state.messages.map(item => {
                if (item.id == data.value.id) {
                  return { ...item, amountOfReactions: data.value.count }
                }
                return item
              })
            }

          })
          break;
      }

    }

    return () => {
      ws.close()
    }
  }, [roomId, queryClient])

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
