import { useState } from 'react';
import './App.css';

function App() {
    const [url, setUrl] = useState('');
    const [short, setShort] = useState('');
    const [copied, setCopied] = useState(false);

    async function shorten() {
        const res = await fetch('/api/links', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url }),
        });
        const data = await res.json();
        const fullShort = `${window.location.origin}/r/${data.short}`;
        setShort(fullShort);
        setCopied(false);
    }

    const copyToClipboard = async () => {
        await navigator.clipboard.writeText(short);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    const clearForm = () => {
        setUrl('');
        setShort('');
        setCopied(false);
    };

    return (
        <main className="container">
            <h1 className="title">🔗 URL Shortener</h1>

            <div className="form">
                <input
                    className="input"
                    placeholder="https://example.com"
                    value={url}
                    onChange={e => setUrl(e.target.value)}
                />
                <button onClick={shorten} className="button" disabled={!url.trim()}>
                    Shorten
                </button>
                {(url || short) && (
                    <button onClick={clearForm} className="clear-button">
                        🗑 Clear
                    </button>
                )}
            </div>

            {short && (
                <div className="result-box">
                    <p className="result">
                        Short link:{' '}
                        <a className="link" href={short} target="_blank" rel="noopener noreferrer">
                            {short}
                        </a>
                    </p>
                    <button onClick={copyToClipboard} className="copy-button">
                        {copied ? '✅ Copied!' : '📋 Copy'}
                    </button>
                </div>
            )}
        </main>
    );
}

export default App;
