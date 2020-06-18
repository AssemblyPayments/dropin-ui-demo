package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func main() {
	http.HandleFunc("/dropin.html", handleHosted)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleHosted(w http.ResponseWriter, req *http.Request) {
	uid, err := user()
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(failed, 500, err, "save card user")))
		return
	}
	token, err := token(uid)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(failed, 500, err, "API token acquisition")))
		return
	}
	w.Write([]byte(fmt.Sprintf(payload, token)))
}

func user() (string, error) {
	auth, err := readSecret()
	if err != nil {
		return "", err
	}
	uid := uuid.New().String()
	reader := bytes.NewReader(newUser(uid))
	req, err := http.NewRequest("POST", "https://test.api.promisepay.com/users", reader)
	if err != nil {
		return "", err
	}
	req.Header["Authorization"] = []string{"Basic " + auth}
	req.Header["Content-Type"] = []string{"application/json"}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 201 {
		var msg string
		if len(body) > 0 {
			msg = string(body)
		} else {
			msg = fmt.Sprintf("user create got HTTP %d", res.StatusCode)
		}
		return "", errors.New(msg)
	}
	return uid, err
}

func readSecret() (string, error) {
	b, err := ioutil.ReadFile("secret")
	if err != nil {
		return "", err
	}
	clean := bytes.TrimSpace(b)
	return base64.StdEncoding.EncodeToString(clean), nil
}

func newUser(uid string) []byte {
	user := userFields{
		ID:        uid,
		FirstName: "Alice",
		Email:     fmt.Sprintf("%s@test.com", uid),
		LastName:  "Baker",
		Country:   "AUS",
	}
	slice, _ := json.Marshal(user)
	return slice
}

type userFields struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	Email     string `json:"email"`
	LastName  string `json:"last_name"`
	Country   string `json:"country"`
}

const failed = `<html>
<head>
	<title>fail</title>
</head>
<body>
<h2>HTTP %d %s on %s</h2>
</body>
</html>`

func token(user string) (string, error) {
	auth, err := readSecret()
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(
		"POST",
		"https://test.api.promisepay.com/token_auths?token_type=card&user_id="+user,
		nil,
	)
	if err != nil {
		return "", err
	}
	req.Header["Authorization"] = []string{"Basic " + auth}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if res.StatusCode > 201 {
		var msg string
		if len(body) > 0 {
			msg = string(body)
		} else {
			msg = fmt.Sprintf("token acquisition got HTTP %d", res.StatusCode)
		}
		return "", errors.New(msg)
	}
	tr := tokenResponse{}
	err = json.Unmarshal(body, &tr)
	return tr.Auth.Token, err
}

type tokenResponse struct {
	Auth tokenAuth `json:"token_auth"`
}

type tokenAuth struct {
	Name  string `json:"token_type"`
	Token string `json:"token"`
	User  string `json:"user_id"`
}

const payload = `<html>
<head>
	<title>Drop-In UI test page</title>
	<meta http-equiv="refresh" content="1800"/>
	<meta name="viewport" content="width=device-width, initial-scale=1" />
	<style type="text/css">	
		body {
			font-family: Arial, Helvetica, sans-serif;
			font-size: 13px;
		}
		#container {
			max-width: 400px;
			border: 0.5px solid #bdbdbd;
		}	
	</style>
</head>
<body>
	<div id="container"></div>
	<script src="https://hosted.assemblypay.com/assembly.js"></script>
	<script>    
		let dropinHandler = Dropin.create({        
			cardTokenAuth: '%s',
			environment: 'prelive',        
			targetElementId: '#container',
			cardAccountCreationCallback: function(cardAccountResult) { 
				alert(` + "`Card Account ID received: ${cardAccountResult.id}`" + `)
			}
		}, function (error, instance) {        
			if(error) { 
				alert(` + "`Error: ${error}`" + `);
			}    
		});
	</script>
</body>
</html>`
