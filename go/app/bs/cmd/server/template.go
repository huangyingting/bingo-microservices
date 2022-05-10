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
			var c = new WebSocket("ws://{{.Host}}/ws");
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
        console.log(response.status, response.statusText);
        response.json().then((json) => {
          console.log(json);
          error.innerText = json.error;
        })
      })
    })
  });    
</script>

</html>
`))
