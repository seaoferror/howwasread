import { GeoCoordinates } from "@/types/conversation";

export function extractPlaceName(
  longURL: string,
): string {
  const namePattern = /\/maps\/(?:place|search)\/([^/@]+)/;
  const match = longURL.match(namePattern);
  if (match && match[1]) {
    const spaceCleaned = match[1].replace(/\+/g, " ");
    const decodedName = decodeURIComponent(spaceCleaned);
    const isRawCoordinateString = /[\d°'"]+[NSEW]/.test(decodedName);
    if (isRawCoordinateString) {
      return "Dropped Pin";
    }
    return decodedName;
  }
  return "";
}

export function extractCoords(
  longURL: string,
): GeoCoordinates | null {
  const atPattern = /@(-?\d+\.\d+),(-?\d+\.\d+)/;
  const atMatch = longURL.match(atPattern);

  if (atMatch && atMatch[1] && atMatch[2]) {
    return {
      lat: parseFloat(atMatch[1]),
      lng: parseFloat(atMatch[2]),
    };
  }
  return null;
}
