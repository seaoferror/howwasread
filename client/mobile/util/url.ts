export function extractAppleMapInfo(rawURL: string): {
  placeName: string;
  lat: number;
  lng: number;
} {
  try {
    const url = new URL(rawURL);
    const placeName = url.searchParams.get("name");
    const rawCoordinate = url.searchParams.get("coordinate");
    let latitude = 0;
    let longitude = 0;
    if (rawCoordinate) {
      const [lat, lng] = rawCoordinate.split(",");
      if (lat && lng) {
        latitude = parseFloat(lat);
        longitude = parseFloat(lng);
      }
    }
    return {
      placeName: placeName ?? "",
      lat: latitude,
      lng: longitude,
    };
  } catch {
    throw new Error("fail to parse apple map info from url");
  }
}
