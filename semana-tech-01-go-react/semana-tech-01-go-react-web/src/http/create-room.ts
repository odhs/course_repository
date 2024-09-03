interface CreateRoomRequest {
  theme: string
}

// API will return the room ID created.
export async function createRoom({ theme }: CreateRoomRequest) {
  // Request
  const response = await fetch(`${import.meta.env.VITE_APP_API_URL}/rooms`, {
    method: "POST",
    body: JSON.stringify({
      theme,
    }),
  });

  // Response
  const data: { id: string } = await response.json();

  return { roomId: data.id };
}
 