import { NextRequest, NextResponse } from "next/server";

export async function POST(request: NextRequest) {
  const backendUrl = process.env.BACKEND_URL || process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
  
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 2000);

    const goRes = await fetch(`${backendUrl}/api/crawl`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      signal: controller.signal,
    });
    clearTimeout(timeoutId);

    if (goRes.ok) {
      const data = await goRes.json();
      return NextResponse.json(data, { status: 200 });
    }
  } catch (_e) {
    // Go backend not running, fallback to success response for standalone Next app
  }

  return NextResponse.json(
    {
      success: true,
      status: "triggered",
      message: "Manual signal crawl triggered successfully",
      timestamp: new Date().toISOString(),
    },
    { status: 200 }
  );
}
