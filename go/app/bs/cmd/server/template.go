package main

import "html/template"

var WS_DEBUG_TEMPLATE = template.Must(template.New("WsDebug").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
	<title>WebSocket Client</title>
  </head>
  <body>
    <h3>WebSocket Client</h3>
    <pre id="output"></pre>
    <script type="text/javascript">
		(function() {
			var data = document.getElementById("output");
      var schema = window.location.protocol === "https:" ? "wss://": "ws://"
			var c = new WebSocket(schema+"{{.Host}}/ws");
			c.onclose = function(msg) {
					data.append((new Date().toUTCString())+" <== Connection closed")
			}
			c.onmessage = function(msg) {
					console.log('On message: '+msg.data);
					data.append((new Date().toUTCString())+" <== "+msg.data+"\n")
			}
		})();
    </script>
  </body>
</html>
`))

var CAPTCHA_TEMPLATE = template.Must(template.New("Captcha").Parse(`
<html>

<head>
  <meta name="viewport" content="width=device-width, height=device-height">
  <title>Bingo Short Url Service</title>
  <meta name="referrer" content="{{.Referrer}}" />
  <style>
  @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@500&display=swap');

  body {
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    margin: 0;
    background: #a18cd1;
    color: #fff;
    text-align: center;
    font-family: 'Poppins', sans-serif;
    /* Chrome 10-25, Safari 5.1-6 */
    background: -webkit-linear-gradient(to right, rgba(161, 140, 209, 0.5), rgba(251, 194, 235, 0.5));
    /* W3C, IE 10+/ Edge, Firefox 16+, Chrome 26+, Opera 12+, Safari 7+ */
    background: linear-gradient(to right, rgba(161, 140, 209, 0.5), rgba(251, 194, 235, 0.5));
  }

  #loading {
    display: inline-block;
    width: 64px;
    height: 63px;
    border: 4px solid rgba(255, 255, 255, .3);
    border-radius: 50%;
    border-top-color: #fff;
    animation: spin 3s ease-in-out infinite;
    -webkit-animation: spin 3s ease-in-out infinite;
  }

  @keyframes spin {
    to {
      -webkit-transform: rotate(360deg);
    }
  }

  @-webkit-keyframes spin {
    to {
      -webkit-transform: rotate(360deg);
    }
  }

  .grecaptcha-badge {
    visibility: hidden;
  }
    .grecaptcha-badge {
      visibility: hidden;
    }
  </style>
</head>

<body>
  <h1>Bingo Short Url Service</h1>
  <div id="loading"></div>
  <h3 id="title">Redirecting...</h3>
  <h6 id="error"></h6>
</body>
<script src="https://www.google.com/recaptcha/api.js?render={{.RecaptchaSiteKey}}"></script>
<script>
  function redirect(url) {
    window.location.href = url
  }
  const site_key = '{{.RecaptchaSiteKey}}';
  grecaptcha.ready(function () {
    grecaptcha.execute(site_key, { action: 'http_ok_redirect' }).then(function (token) {
      fetch("/v1/captcha/verify", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ "token": token, "alias": {{.Alias}} }),
      }).then((response) => {
        if (response.ok) {
          return response.json()
        }
        return Promise.reject(response);
      }).then((data) => {
        if (data.score < 0.5 && true) {
          title.innerText = "Fraud Detected";
          loading.style.display = "none";
        } else {
          redirect(data.url);
        }
      }).catch((response) => {
        title.innerText = "Failed to redirect";
        loading.style.display = "none";
        response.json().then((e) => {
          console.log(e);
          error.innerText = e.message;
        })
      })
    })
  });    
</script>

</html>
`))

var EXPAND_TEMPLATE = template.Must(template.New("Expand").Parse(`
<!DOCTYPE html>
<html>

<head>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Expanded URL Information</title>
  <style>
    @import url('https://fonts.googleapis.com/css2?family=Poppins:wght@500&display=swap');
    h1, h2, h3, h4, p, body, a {
      font-family: 'Poppins', sans-serif;
    }
    body {
      padding: 0px 12px;
      background: #a18cd1;
      background: #a18cd1;
      color: #fff;
      background: -webkit-linear-gradient(to right, rgba(161, 140, 209, 0.5), rgba(251, 194, 235, 0.5));
      background: linear-gradient(to right, rgba(161, 140, 209, 0.5), rgba(251, 194, 235, 0.5));
    }
    .w-25 {
      width: 25%;
      height: auto;
      overflow-wrap: break-word;
    }
    .w-75 {
      width: 75%;
      height: auto;
      overflow-wrap: break-word;
    }
    .flex {
      display: flex;
      flex-direction: row;
    }
    .left {
      float: left;
      width: 50%;
    }
    .right {
      float: right;
      width: 50%;
    }
    img {
      max-width: 100%;
      height: auto;
    }
    .group:after {
      content: "";
      display: table;
      clear: both;
    }
    @media screen and (max-width: 768px) {
      .left,
      .right {
        float: none;
        width: auto;
      }
    }
  </style>
</head>
<body>
  <h1>Expand Short URL</h1>
  <h2>Allows you to retrieve the original URL from a shortened link before clicking on it and visiting the destination.
    We provide furthermore information about unshortened URL such as title, description, keywords and summary of the
    page.</h2>
  <div class="group">
    <div class="left">
      <h2>Information</h2>
      <div class="flex">
        <h3 class="w-25">Short URL Alias:</h3>
        <p class="w-75">{{.Alias}}</p>
      </div>
      <div class="flex">
        <h3 class="w-25">Original URL:</h3>
        <p class="w-75">{{.Url}}</p>
      </div>
      <div class="flex">
        <h3 class="w-25">Title:</h3>
        <p class="w-75">{{.Title}}</p>
      </div>
      <div class="flex">
        <h3 class="w-25">Keywords:</h3>
        <p class="w-75">{{.Keywords}}</p>
      </div>
      <div class="flex">
        <h3 class="w-25">Summary:</h3>
        <p class="w-75">{{.Summary}}</p>
      </div>
    </div>
    <div class="right">
      <h2>URL Snapshot</h2>
      <img
        src="{{.Snapshot}}" loading="lazy" alt="snapshot"/>
    </div>
  </div>
</body>
</html>
`))
