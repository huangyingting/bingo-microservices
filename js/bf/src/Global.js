export const API_ENDPOINT = process.env.REACT_APP_ENV === "dev" ? "http://localhost:8080" : window.location.protocol + "//" + window.location.host
export const WS_ENDPOINT = process.env.REACT_APP_ENV === "dev" ? "ws://localhost:8080/ws" : "ws://" + window.location.host + "/ws"
export const BF_CLIENT_ID = "50493366-b44d-48ff-bbc4-dd1bc1bf4c56"
export const BF_SCOPES = ["https://aliasesbiz.onmicrosoft.com/0b794c85-05a3-4c5e-8150-51dd647363fc/API.Read", "https://aliasesbiz.onmicrosoft.com/0b794c85-05a3-4c5e-8150-51dd647363fc/API.Write"]
export const BF_AUTHORITY = "https://login.microsoftonline.com/common"
