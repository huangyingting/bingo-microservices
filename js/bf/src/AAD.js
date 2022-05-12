import { PublicClientApplication, LogLevel, InteractionRequiredAuthError } from "@azure/msal-browser";
import { BF_CLIENT_ID, BF_SCOPES, BF_AUTHORITY } from "./Global";
export const loginRequest = {
  scopes: ["User.Read"]
};

const msalConfig = {
  auth: {
    clientId: BF_CLIENT_ID,
    authority: BF_AUTHORITY,
    redirectUri: window.location.protocol + '//' + window.location.host + '/blank.html'
  },
  cache: {
    cacheLocation: "sessionStorage", // This configures where your cache will be stored
    storeAuthStateInCookie: false, // Set this to "true" if you are having issues on IE11 or Edge
  },
  system: {
    loggerOptions: {
      loggerCallback: (level, message, containsPii) => {
        if (containsPii) {
          return;
        }
        switch (level) {
          case LogLevel.Error:
            console.error(message);
            return;
          case LogLevel.Warning:
            console.warn(message);
            return;
          //case LogLevel.Info:
          //  console.info(message);
          //  return;
          //case LogLevel.Verbose:
          //  console.debug(message);
          //  return;
          default:
            return;
        }
      }
    }
  }
};

export const MSAL_INSTANCE = new PublicClientApplication(msalConfig);

export const GetAccessToken = async () => {
  const account = MSAL_INSTANCE.getAllAccounts()[0];
  var token = null
  if (account) {
    try {
      token = await MSAL_INSTANCE.acquireTokenSilent({
        scopes: BF_SCOPES,
        account: account
      });
    }
    catch (error) {
      if (error instanceof InteractionRequiredAuthError) {
        try {
          token = await MSAL_INSTANCE.acquireTokenPopup({
            scopes: BF_SCOPES,
            account: account
          });
        }
        catch (error) {
          console.log(error);
        }
      } else {
        console.log(error);
      }
    }
  }
  return token.accessToken
}