export const API_ENDPOINT = process.env.REACT_APP_ENV === "dev" ? "http://localhost:8080" : window.location.protocol + "//" + window.location.host
export const WS_ENDPOINT = process.env.REACT_APP_ENV === "dev" ? "ws://localhost:8080/ws" : "ws://" + window.location.host + "/ws"
