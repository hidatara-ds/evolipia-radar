import { NextRequest } from "next/server";

export async function GET(_request: NextRequest) {
  const responseStream = new ReadableStream({
    start(controller) {
      const encoder = new TextEncoder();
      const initialEvent = `event: progress\ndata: ${JSON.stringify({
        step: 6,
        message: "Standing by for next crawl",
        progress: 100,
        is_complete: true,
        has_error: false,
      })}\n\n`;
      controller.enqueue(encoder.encode(initialEvent));
      controller.close();
    },
  });

  return new Response(responseStream, {
    headers: {
      "Content-Type": "text/event-stream",
      "Cache-Control": "no-cache, no-transform",
      Connection: "keep-alive",
    },
  });
}
