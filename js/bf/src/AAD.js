import { PublicClientApplication, LogLevel, InteractionRequiredAuthError } from "@azure/msal-browser";

export const loginRequest = {
  scopes: ["User.Read"]
};

const msalConfig = {
  auth: {
    clientId: "eea9dc4e-6439-4800-b3b5-8d8fcf926369",
    authority: "https://login.microsoftonline.com/736e8d18-4edf-4080-96b9-69afaec94892",
    redirectUri: window.location.protocol + '//' + window.location.host + '/'
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
        scopes: ["api://52d63002-b52b-4233-a733-092896d960b6/API.Read", "api://52d63002-b52b-4233-a733-092896d960b6/API.Write"],
        account: account
      });
    }
    catch (error) {
      if (error instanceof InteractionRequiredAuthError) {
        try {
          token = await MSAL_INSTANCE.acquireTokenPopup({
            scopes: ["api://52d63002-b52b-4233-a733-092896d960b6/API.Read", "api://52d63002-b52b-4233-a733-092896d960b6/API.Write"],
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