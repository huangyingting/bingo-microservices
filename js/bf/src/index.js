import React from 'react';
import ReactDOM from 'react-dom';
import './scss/Customization.scss';
import App from './App';
import { MsalProvider } from "@azure/msal-react";
import { MSAL_INSTANCE} from "./AAD"

ReactDOM.render(
  <React.StrictMode>
    <MsalProvider instance={MSAL_INSTANCE}>
      <App />
    </MsalProvider>
  </React.StrictMode>,
  document.getElementById('root')
);
