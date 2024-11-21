**Abstract**

The website monitoring and alerting system checks and verifies if the specific 
websites are up and working, so that the site visitors can use the site as expected. 
This tool monitors the given websites using a specific status code (a status code is a 
message a website 's server sends to the browser to indicate whether or not that 
request can be fulfilled). 

Our project checks the status code for the following websites in periodic intervals : 
- www.facebook.com
- www.twitter.com
- www.google.com
- http://localhost
   
If the response code does not match with the expected status code ( here, 200), 
connection refused from server, then it triggers an e-mail and noti es the user. 

GoLang Packages Used : 
- net/http
- time
- net/stmp
- os
