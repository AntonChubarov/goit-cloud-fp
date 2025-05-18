import { useState } from 'react'
import './App.css'

function App() {
    const [url, setUrl] = useState("");
    const [short, setShort] = useState("");

    async function shorten() {
        const res = await fetch("/api/links", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ url }),
        });
        const data = await res.json();
        setShort(`${window.location.origin}/${data.short}`);
    }

    return (
        <main className="min-h-screen flex flex-col items-center justify-center bg-gray-900 text-gray-200">
            <h1 className="text-3xl mb-6">URL Shortener</h1>
            <input
                className="p-2 w-96 text-black rounded-l"
                placeholder="https://example.com"
                value={url}
                onChange={e => setUrl(e.target.value)}
            />
            <button onClick={shorten} className="p-2 bg-purple-700 rounded-r">
                Shorten
            </button>

            {short && (
                <p className="mt-4">
                    Short link:{" "}
                    <a className="text-sky-400" href={short}>
                        {short}
                    </a>
                </p>
            )}
        </main>
    );
}

export default App
