import { useRef, useState, useEffect } from "react";

export const Test = () => {
    const inputVal = useRef(null);
    const wsRef = useRef(null);

    const [serverData, setServerData] = useState("");

    useEffect(() => {
        const ws = new WebSocket("ws://localhost:8080/ws");

        wsRef.current = ws;

        ws.onopen = () => {
            console.log("Connected");
        };

        ws.onmessage = (event) => {
            console.log(event.data);
            setServerData(event.data);
        };

        ws.onclose = () => {
            console.log("Disconnected");
        };

        return () => ws.close();
    }, []);

    const sendMessage = (e) => {
        e.preventDefault();

        if (wsRef.current?.readyState === WebSocket.OPEN) {
            wsRef.current.send(inputVal.current.value);
            inputVal.current.value = "";
        }
    };

    return (
        <div>
            <h2>WebSocket Test</h2>

            <p>{serverData}</p>

            <form onSubmit={sendMessage}>
                <input type="text" ref={inputVal} />
                <button type="submit">Send</button>
            </form>
        </div>
    );
};