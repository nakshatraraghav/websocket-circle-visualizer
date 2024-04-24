import React, { useEffect, useRef } from "react";

interface PointData {
  sin: number | null;
  cos: number | null;
}

const Circle: React.FC = () => {
  // const [radius, setRadius] = useState<number>(100);
  const point = useRef<PointData>({ sin: null, cos: null });
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const ws = useRef<WebSocket | null>(null);

  const radius = 100;

  useEffect(() => {
    async function HandleRenderingAndSocket() {
      const canvas = canvasRef.current;
      if (!canvas) return;

      const ctx = canvas.getContext("2d");
      if (!ctx) return;

      // Define circle properties
      const centerX = canvas.width / 2;
      const centerY = canvas.height / 2;

      // Clear canvas
      ctx.clearRect(0, 0, canvas.width, canvas.height);

      // Draw circle
      drawCircle(ctx, centerX, centerY, radius);

      // Connect to WebSocket
      ws.current = new WebSocket("ws://localhost:3000/ws");

      ws.current.onopen = () => {
        console.log("Connected to WebSocket");
      };

      // Listen for WebSocket messages
      ws.current.onmessage = (event) => {
        const msg = event.data as string;
        console.log(msg);

        // Parse the message
        const [key, value] = msg.split(":");
        if (key === "cos") {
          point.current.cos = parseFloat(value);
        } else if (key === "sin") {
          point.current.sin = parseFloat(value);
        } else if (key == "change") {
          // setRadius(+value);
        }

        // Render the point on the circle
        if (point.current.cos !== null && point.current.sin !== null) {
          ctx.clearRect(0, 0, canvas.width, canvas.height);
          drawCircle(ctx, centerX, centerY, radius);

          const projectionX =
            centerX +
            radius * Math.cos(Math.atan2(point.current.sin, point.current.cos));
          const projectionY =
            centerY +
            radius * Math.sin(Math.atan2(point.current.sin, point.current.cos));
          drawPoint(ctx, projectionX, projectionY);
        }
      };
    }

    HandleRenderingAndSocket();

    return () => {
      ws.current?.close();
      alert("connection closed");
    };
  }, []);

  const drawCircle = (
    ctx: CanvasRenderingContext2D,
    centerX: number,
    centerY: number,
    radius: number
  ) => {
    ctx.beginPath();
    ctx.arc(centerX, centerY, radius, 0, 2 * Math.PI);
    ctx.stroke();
  };

  const drawPoint = (ctx: CanvasRenderingContext2D, x: number, y: number) => {
    ctx.beginPath();
    ctx.arc(x, y, 3, 0, 2 * Math.PI);
    ctx.fillStyle = "blue";
    ctx.fill();
  };

  // async function UpdateRadius(
  //   ev: React.MouseEvent<HTMLButtonElement, MouseEvent>
  // ) {
  //   ev.preventDefault();
  //   const req = await fetch("http://localhost:3000/update_radius", {
  //     method: "POST",
  //     body: JSON.stringify(radius),
  //   });

  //   if (!req.ok) {
  //     alert("failed to update the radius");
  //   }

  //   ws.current?.send(JSON.stringify(radius));
  // }

  return (
    <div className="w-screen h-screen flex items-center justify-center space-x-4">
      <canvas ref={canvasRef} width={400} height={400} />;
    </div>
  );
};

export default Circle;
