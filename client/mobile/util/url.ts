export function extractPlaceNameFromPathVariable(longURL: string): string {
  const namePattern = /\/maps\/(?:place|search)\/([^/@]+)/;
  const match = longURL.match(namePattern);
  if (match && match[1]) {
    const spaceCleaned = match[1].replace(/\+/g, " ");
    const decodedName = decodeURIComponent(spaceCleaned);
    const isRawCoordinateString = /^[-]?\d+\.\d+,\s*[-]?\d+\.\d+$/.test(
      decodedName,
    );
    if (isRawCoordinateString) {
      return "Dropped Pin";
    }
    return decodedName;
  }
  return "";
}

export function extractPlaceNameFromQueryParam(longUrl: string): string {
  const match = longUrl.match(/[?&]q=([^&]+)/i);

  if (match && match[1]) {
    const withSpaces = match[1].replace(/\+/g, " ");
    const decodedName = decodeURIComponent(withSpaces);
    const isRawCoordinateString = /^[-]?\d+\.\d+,\s*[-]?\d+\.\d+$/.test(
      decodedName,
    );
    if (isRawCoordinateString) {
      return "Dropped Pin";
    }
    return decodedName.split(",")[0].trim();
  }
  return "";
}

export function extractCoordsFromURL(
  longURL: string,
): { lat: number; lng: number } | null {
  const pinPattern = /!3d(-?\d+\.\d+)!4d(-?\d+\.\d+)/;
  const pinMatch = longURL.match(pinPattern);

  if (pinMatch && pinMatch[1] && pinMatch[2]) {
    return {
      lat: parseFloat(pinMatch[1]),
      lng: parseFloat(pinMatch[2]),
    };
  }
  return null;
}
