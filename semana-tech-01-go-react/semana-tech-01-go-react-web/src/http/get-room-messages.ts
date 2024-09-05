interface GetRoomMessagesRequest {
  roomId: string
}

export interface GetRoomMessagesResponse {
  messages: {
    id: string
    text: string
    amountOfReactions: number
    answered: boolean
  }[]
}

// API will return the room ID created.
export async function getRoomMessages({
  roomId,
}: GetRoomMessagesRequest): Promise<GetRoomMessagesResponse> {
  // Request
  const response = await fetch(
    `${import.meta.env.VITE_APP_API_URL}/rooms/${roomId}/messages`
  )

  // Response then defines the json response array
  const data: Array<{
    id: string
    room_id: string
    message: string
    reaction_count: number
    answered: boolean
  }> = await response.json()

  // Returning with a map to rename some itens used in the frontend
  return {
    messages: data.map(item => {
      return {
        id: item.id,
        text: item.message,
        amountOfReactions: item.reaction_count,
        answered: item.answered,
      }
    }),
  }
}
