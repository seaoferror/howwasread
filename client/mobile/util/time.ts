export function getTimestamp(id: string): string {
  if (id.length < 36 || id.charAt(14) !== "7") {
    throw new Error("Invalid UUIDv7 string provided.");
  }

  const hexTimestamp = id.substring(0, 8) + id.substring(9, 13);

  const timestampMs = parseInt(hexTimestamp, 16);

  return new Date(timestampMs).toISOString();
}
