interface CreateMessageRequest {
  roomId: string
  message: string
}

// API will return the room ID created.
export async function createMessage({ roomId, message }: CreateMessageRequest) {
  // Request
  const response = await fetch(`${import.meta.env.VITE_APP_API_URL}/rooms/${roomId}/messages`, {
    method: "POST",
    body: JSON.stringify({
      message,
    }),
  });

  // Response
  const data: { id: string } = await response.json();

  return { messageId: data.id };
}
 