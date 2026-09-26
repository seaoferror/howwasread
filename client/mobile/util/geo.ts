import { GeoInfo } from "@/types/conversation";

const ALPHABET = "23456789CFGHJMPQRVWX";
const LAT_UNITS = 8000 * 3125; // units per degree at max precision (lat)
const LNG_UNITS = 8000 * 1024; // units per degree at max precision (lng)

function decodePlusCode(code: string): { lat: number; lng: number } {
  const s = String(code).trim().toUpperCase();
  if (s.indexOf("+") !== 8) throw new Error(`Not a full Plus Code: ${code}`);
  const digits = s.replace("+", "").slice(0, 15);

  let lat = 0;
  let lng = 0;
  let latPlace = 20 * LAT_UNITS;
  let lngPlace = 20 * LNG_UNITS;

  // Digits 1-10: lat/lng pairs, base 20.
  for (let i = 0; i < Math.min(10, digits.length); i += 2) {
    if (digits[i] === "0") break; // padding, e.g. "8Q980000+"
    const a = ALPHABET.indexOf(digits[i]);
    const b = ALPHABET.indexOf(digits[i + 1]);
    if (a < 0 || b < 0) throw new Error(`Invalid character in ${code}`);
    if (i > 0) {
      latPlace /= 20;
      lngPlace /= 20;
    }
    lat += a * latPlace;
    lng += b * lngPlace;
  }

  // Digits 11-15: 5 rows x 4 columns grid refinement.
  for (let i = 10; i < digits.length; i++) {
    const d = ALPHABET.indexOf(digits[i]);
    if (d < 0) throw new Error(`Invalid character in ${code}`);
    latPlace /= 5;
    lngPlace /= 4;
    lat += Math.floor(d / 4) * latPlace;
    lng += (d % 4) * lngPlace;
  }

  return {
    lat: (lat + latPlace / 2) / LAT_UNITS - 90,
    lng: (lng + lngPlace / 2) / LNG_UNITS - 180,
  };
}

function extractCoords(url: string): { lat: number; lng: number } | null {
  const pin = url.match(/!3d(-?\d+(?:\.\d+)?)!4d(-?\d+(?:\.\d+)?)/);
  if (pin) return { lat: parseFloat(pin[1]), lng: parseFloat(pin[2]) };

  const plus = url.match(/!20s([^!?&#]+)/);
  if (plus) {
    try {
      return decodePlusCode(decodeURIComponent(plus[1]));
    } catch {
      // Short/invalid code: fall through, keep waiting for a better URL.
    }
  }
  return null;
}

function extractPlaceName(url: string): string | null {
  const m = url.match(/\/maps\/place\/([^/@?]+)/);
  if (!m) return null;
  let name: string;
  try {
    name = decodeURIComponent(m[1].replace(/\+/g, " "));
  } catch {
    name = m[1];
  }
  const isCoords = /^-?\d+(?:\.\d+)?,\s*-?\d+(?:\.\d+)?$/.test(name);
  return isCoords ? "Dropped Pin" : name;
}

export function parseGoogleMapsUrl(url: string): GeoInfo | null {
  const placeName = extractPlaceName(url);
  const coords = extractCoords(url);
  if (!placeName || !coords) return null;
  return { ...coords, placeName };
}

