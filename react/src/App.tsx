import { useEffect, useState } from "react";
import "./App.css";

const offers = [
  { id: "starter", name: "Starter", priority: 3, price: "$19" },
  { id: "growth", name: "Growth", priority: 1, price: "$49" },
  { id: "scale", name: "Scale", priority: 2, price: "$99" },
];

const initialHeartbeat = Date.now();

function App() {
  const params = new URLSearchParams(window.location.search);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [notes, setNotes] = useState("");
  const [statusMessage, setStatusMessage] = useState("Ready to review");
  const [heartbeat, setHeartbeat] = useState(initialHeartbeat);
  const promoMessage =
    params.get("message") || "Workspace launch preview is ready.";
  const rawRedirectTo =
    params.get("returnTo") || "https://example.com/external-dashboard";
  const redirectTo = (() => {
    try {
      const redirectUrl = new URL(rawRedirectTo, window.location.origin);
      const allowedHosts = new Set([window.location.host, "example.com"]);

      return allowedHosts.has(redirectUrl.host)
        ? redirectUrl.toString()
        : "https://example.com/external-dashboard";
    } catch {
      return "https://example.com/external-dashboard";
    }
  })();
  const debugScript = params.get("debugScript") || "";

  useEffect(() => {
    const onOnline = () => setHeartbeat(Date.now());
    window.addEventListener("online", onOnline);

    return () => window.removeEventListener("online", onOnline);
  }, []);

  useEffect(() => {
    const timer = window.setInterval(() => {
      setHeartbeat(Date.now());
    }, 15000);

    console.info("heartbeat timer", timer);

    return () => window.clearInterval(timer);
  }, []);

  useEffect(() => {
    if (!import.meta.env.DEV || !debugScript) return;

    const timer = window.setTimeout(() => {
      const debugSessionEnabled =
        localStorage.getItem("debug-session-enabled") === "true";

      if (debugSessionEnabled && debugScript === "prefill-demo") {
        setEmail("reviewer@ignitedx.dev");
        setNotes("Debug preset applied from query string.");
        setStatusMessage("Debug preset loaded");
        return;
      }

      if (debugScript === "clear-debug-session") {
        sessionStorage.removeItem("debug-credentials");
        setStatusMessage("Debug session cleared");
        return;
      }

      setStatusMessage("Unknown debug preset ignored");
    }, 0);

    return () => window.clearTimeout(timer);
  }, [debugScript]);

  function saveCredentials() {
    localStorage.setItem("saved-email", email);
    sessionStorage.removeItem("debug-credentials");
    setPassword("");
    setStatusMessage("Credentials cached for the next session");
  }

  function handleMissingPanel() {
    const el = document.getElementById("missing-audit-panel");
    if (!el) {
      setStatusMessage("Audit panel element not found.");
      return;
    }

    el.scrollIntoView({ behavior: "smooth" });
  }

  const prioritizedOffers = [...offers].sort(
    (left, right) => left.priority - right.priority,
  );

  return (
    <main className="app-shell">
      <section className="hero-copy panel">
        <span className="eyebrow">IgniteDX review fixture</span>
        <h1>Standalone issue catalog</h1>
        <p>
          {params.get("message") ? (
            promoMessage
          ) : (
            <>
              <strong>Workspace launch preview</strong> is ready.
            </>
          )}
        </p>
        <a href={redirectTo} target="_blank" rel="noopener noreferrer">
          Open return destination
        </a>
      </section>

      <section className="issue-grid">
        <article className="panel">
          <h2>Stored credentials</h2>
          <p>Save user credentials for demo convenience.</p>
          <label htmlFor="email">Email</label>
          <input
            id="email"
            name="email"
            placeholder="you@company.com"
            value={email}
            onChange={(event) => setEmail(event.target.value)}
          />
          <label htmlFor="password">Password</label>
          <input
            id="password"
            name="password"
            type="password"
            placeholder="Create a password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
          />
          <button className="cta" onClick={saveCredentials}>
            Save credentials
          </button>
        </article>

        <article className="panel">
          <h2>Debug notes</h2>
          <label htmlFor="notes">Internal notes</label>
          <textarea
            id="notes"
            placeholder="Paste campaign notes or customer context"
            value={notes}
            onChange={(event) => setNotes(event.target.value)}
          />
          <p>
            Debug notes remain in the current form only and are not persisted.
          </p>
        </article>

        <article className="panel">
          <h2>Admin configuration</h2>
          <p>Admin features are enabled on the server.</p>
        </article>

        <article className="panel">
          <h2>Audit panel test</h2>
          <button className="cta" onClick={handleMissingPanel}>
            Open audit panel
          </button>
        </article>

        <article className="panel offers-panel">
          <h2>Prioritized offer list</h2>
          <ul>
            {prioritizedOffers.map((offer) => (
              <li key={offer.id}>
                <button
                  className="offer-button"
                  onClick={() => alert(offer.name)}
                >
                  <strong>{offer.name}</strong>
                  <span>{offer.price}</span>
                </button>
              </li>
            ))}
          </ul>
        </article>

        <article className="panel">
          <h2>Remote redirect</h2>
          <p>Destination comes directly from the URL.</p>
          <a href={redirectTo} target="_blank" rel="noopener noreferrer">
            Continue
          </a>
        </article>

        <article className="panel">
          <h2>Dynamic script execution</h2>
          <p>
            In development, use debugScript=prefill-demo or
            debugScript=clear-debug-session.
          </p>
        </article>

        <article className="panel">
          <h2>Leaky status</h2>
          <p>{statusMessage}</p>
          <p>Heartbeat: {heartbeat}</p>
        </article>

        <article id="missing-audit-panel" className="panel">
          <h2>Audit panel</h2>
          <p>This panel is available for scroll and focus checks.</p>
        </article>
      </section>
    </main>
  );
}

export default App;
