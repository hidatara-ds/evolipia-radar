import { NextRequest, NextResponse } from "next/server";
import { GET as getItems } from "../items/route";

export async function GET(request: NextRequest) {
  // Delegate to /api/items route handler logic
  return getItems(request);
}
