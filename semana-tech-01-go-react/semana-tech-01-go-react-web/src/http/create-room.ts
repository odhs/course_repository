interface CreateRoomRequest {
  theme: string;
}

// API will return the room ID created.
export async function createRoom({ theme }: CreateRoomRequest) {
  // Request
  const response = await fetch("http://localhost:8080/api/", {
    method: "POST",
    body: JSON.stringify({
      theme,
    }),
  });

  // Response in Type Script
  const data: { id: string } = await response.json();
  // Response in Javascript
  //const data: = await response.json();
  
  return { roomId: data.id };
}
